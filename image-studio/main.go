package main

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"image-studio/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	wailsmac "github.com/wailsapp/wails/v2/pkg/options/mac"
	wailswindows "github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// WebView2 用户数据(localStorage / IndexedDB / cookies / cache)的稳定根目录,
// 跟 exe 文件名解绑。v0.1.6 及之前没显式设 WebviewUserDataPath,Wails 退到
// go-webview2 的默认 %APPDATA%/<exe文件名>/EBWebView —— 用户改 exe 名 / 下载
// 副本(Windows 自动加 (1) 后缀)就会找不到老数据。
//
// 路径选 %LOCALAPPDATA% 而不是 %APPDATA%(Roaming),理由跟 Chromium / VSCode 一致:
//  1. WebView2 的 IndexedDB + cache 可以几百 MB,Roaming 域同步会卡死 OneDrive
//     /企业 SSO 桌面流;
//  2. 部分企业策略禁止 Roaming 写入,但允许本地 AppData;
//  3. 它本就是"机器本地"的状态,跨机漫游也没意义。
const stableWebView2Subdir = "image-studio"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	svc := backend.NewService()
	appOptions := &options.App{
		Title:     "Image Studio",
		Width:     1440,
		Height:    980,
		MinWidth:  1100,
		MinHeight: 780,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 20, B: 26, A: 1},
		OnStartup:        svc.Startup,
		Bind: []interface{}{
			svc,
		},
	}

	if runtime.GOOS == "darwin" {
		appOptions.Mac = &wailsmac.Options{
			Appearance:           wailsmac.DefaultAppearance,
			TitleBar:             wailsmac.TitleBarHiddenInset(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		}
	}
	if runtime.GOOS == "windows" {
		// 决定稳定的 WebView2 user data path 并迁移老数据。MUST 在 wails.Run 前,
		// 不然 WebView2 先初始化了空的新目录,migration 检查「新路径无数据」
		// 就会跳过,反而把空目录当成 already-migrated。
		var stableUserDataPath string
		localAppData := os.Getenv("LOCALAPPDATA")
		appData := os.Getenv("APPDATA")
		if localAppData != "" {
			stableUserDataPath = filepath.Join(localAppData, stableWebView2Subdir, "WebView2")
			migrateWebView2DataIfNeeded(appData, localAppData, stableUserDataPath)
		}
		appOptions.Windows = &wailswindows.Options{
			Theme:                wailswindows.SystemDefault,
			BackdropType:         wailswindows.Mica,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
			WebviewUserDataPath:  stableUserDataPath,
			CustomTheme: &wailswindows.ThemeSettings{
				DarkModeTitleBar:           wailswindows.RGB(32, 32, 32),
				DarkModeTitleBarInactive:   wailswindows.RGB(38, 38, 38),
				DarkModeTitleText:          wailswindows.RGB(245, 245, 245),
				DarkModeTitleTextInactive:  wailswindows.RGB(200, 200, 200),
				DarkModeBorder:             wailswindows.RGB(54, 54, 54),
				DarkModeBorderInactive:     wailswindows.RGB(45, 45, 45),
				LightModeTitleBar:          wailswindows.RGB(243, 243, 243),
				LightModeTitleBarInactive:  wailswindows.RGB(237, 237, 237),
				LightModeTitleText:         wailswindows.RGB(31, 31, 31),
				LightModeTitleTextInactive: wailswindows.RGB(96, 96, 96),
				LightModeBorder:            wailswindows.RGB(219, 219, 219),
				LightModeBorderInactive:    wailswindows.RGB(226, 226, 226),
			},
		}
	}

	err := wails.Run(appOptions)

	if err != nil {
		println("Error:", err.Error())
	}
}

// migrateWebView2DataIfNeeded 在 WebView2 初始化前把老版本残留的用户数据搬到
// 稳定路径。Only Windows。完整候选树:
//
//	%APPDATA%/image-studio-windows-amd64.exe/EBWebView   ← 老 release 默认位置
//	%APPDATA%/image-studio.exe/EBWebView                 ← 用户改名 / Win 自动 (1)
//	%LOCALAPPDATA%/image-studio/WebView2/EBWebView       ← v0.1.7 起的稳定位置
//
// 策略:
//  1. 如果稳定路径已经有 WebView2 数据(检查 Default/Local Storage 标志),说明
//     已迁移过或本来就是新装,什么都不做。
//  2. 否则扫 %APPDATA% 和 %LOCALAPPDATA% 下所有 "image-studio" 前缀 + 含
//     EBWebView/Default/Local Storage 签名的目录,按 Default/ 的 mtime 选最近
//     用过的一个 —— 用户可能多次改名留下多份数据,挑最近的最准。
//  3. 迁移用 copy → 写到 EBWebView.partial → 全部成功后 atomic rename 到
//     EBWebView。中途崩了 .partial 留在那里,下次启动检查"新路径无数据"还会
//     重新尝试。
//  4. 老目录不删,保留作为备份。代价是磁盘空间(MB~GB 量级,但只在升级当次
//     占用)。
func migrateWebView2DataIfNeeded(appData, localAppData, newDataPath string) {
	newEBWebView := filepath.Join(newDataPath, "EBWebView")
	if webView2HasData(newEBWebView) {
		return
	}

	candidates := make([]webViewCandidate, 0, 4)
	candidates = append(candidates, findWebView2Candidates(appData, newEBWebView)...)
	candidates = append(candidates, findWebView2Candidates(localAppData, newEBWebView)...)
	if len(candidates) == 0 {
		return
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].mtime.After(candidates[j].mtime)
	})
	best := candidates[0]

	if err := os.MkdirAll(newDataPath, 0o755); err != nil {
		return
	}
	partial := filepath.Join(newDataPath, "EBWebView.partial")
	_ = os.RemoveAll(partial) // 清理上次崩溃留下的半成品
	if err := copyDirAll(best.ebWebView, partial); err != nil {
		_ = os.RemoveAll(partial)
		return
	}
	if err := os.Rename(partial, newEBWebView); err != nil {
		// 同卷下 rename 不该失败;保险起见兜底
		_ = os.RemoveAll(partial)
	}
	// 不删 best.ebWebView,留作备份。用户要清理硬盘自己手动来。
}

type webViewCandidate struct {
	ebWebView string    // 源 EBWebView 绝对路径
	mtime     time.Time // Default/ 子目录的 mtime,代表"最近用过"
}

// findWebView2Candidates 扫 root 下一层目录,挑出含 WebView2 数据 + 名字像
// image-studio 的。filename 前缀过滤是防御:避免错搬其他 Wails 应用的数据
// (一台机上可能装了多个 Wails app,每个都按默认在 %APPDATA% 下留 <exe>/EBWebView,
// 仅靠签名检测会跨应用拽过来)。
//
// excludePath 用来排除当前的稳定目录自身,避免自我搬迁。
func findWebView2Candidates(root, excludePath string) []webViewCandidate {
	if root == "" {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	excludeAbs, _ := filepath.Abs(excludePath)
	out := make([]webViewCandidate, 0, 4)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		nameLower := strings.ToLower(e.Name())
		if !strings.HasPrefix(nameLower, "image-studio") {
			continue
		}
		// 注意:稳定目录 image-studio(无 .exe)被前缀匹配了,这是预期 ——
		// 它的 EBWebView 还在 WebView2/ 子目录下,我们这里只看 <parent>/EBWebView,
		// 所以新稳定路径不会被自己当成迁移源。
		eb := filepath.Join(root, e.Name(), "EBWebView")
		ebAbs, _ := filepath.Abs(eb)
		if ebAbs != "" && ebAbs == excludeAbs {
			continue
		}
		if !webView2HasData(eb) {
			continue
		}
		// 用 Default/ 的 mtime 判断"最近用过"
		mtime := time.Time{}
		if info, err := os.Stat(filepath.Join(eb, "Default")); err == nil {
			mtime = info.ModTime()
		}
		out = append(out, webViewCandidate{ebWebView: eb, mtime: mtime})
	}
	return out
}

// webView2HasData 判断 EBWebView 目录里是否有「真实数据」,而不是空壳。
// localStorage / IndexedDB 的实际文件都在 Default/ 下,挑几个标志性子项检查。
func webView2HasData(ebWebViewDir string) bool {
	for _, sub := range []string{
		"Default/Local Storage",
		"Default/Local Storage/leveldb",
		"Default/IndexedDB",
	} {
		if _, err := os.Stat(filepath.Join(ebWebViewDir, sub)); err == nil {
			return true
		}
	}
	return false
}

// copyDirAll 递归拷贝目录(rename 跨卷失败时的兜底)。
func copyDirAll(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
