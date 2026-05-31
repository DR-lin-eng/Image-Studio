package ui

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	gioCompat "image-studio/gio-client/internal/compat"
	"image-studio/gio-client/internal/kernel"
	sharedCompat "image-studio/shared/compat"

	"github.com/yuanhua/image-gptcodex/pkg/client"
)

func (a *App) startRun() {
	if a.isRunning() {
		return
	}
	cfg := a.currentConfig()
	if err := gioCompat.SaveConfig(cfg); err != nil {
		a.appendLog("兼容配置保存失败: " + err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.mu.Lock()
	a.running = true
	a.cancel = cancel
	a.status = "正在提交请求"
	a.logs = appendBounded(a.logs, "开始任务: "+shortPrompt(cfg.Prompt))
	a.mu.Unlock()
	a.invalidateNow()

	go func() {
		started := time.Now()
		res, err := a.runner.Run(ctx, cfg, kernel.Callbacks{
			Log: func(line string) {
				a.appendLog(line)
			},
			Progress: func(stage string, elapsed int, bytes int64) {
				a.setStatus(fmt.Sprintf("%s - %ds - %s", stage, elapsed, client.FormatBytes(bytes)))
			},
			Partial: func(partial client.PartialImage) {
				a.setStatus(fmt.Sprintf("收到流式预览 #%d", partial.PartialImageIndex))
			},
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				a.finishCancelled()
				return
			}
			a.finishWithError(err, res.RawPath)
			return
		}
		img, err := decodeImageB64(res.ImageB64)
		if err != nil {
			a.finishWithError(err, res.RawPath)
			return
		}
		elapsedSec := time.Since(started).Seconds()
		if err := gioCompat.SaveConfigAndHistory(cfg, res, elapsedSec); err != nil {
			a.appendLog("兼容历史保存失败: " + err.Error())
		}
		compatState, _, _ := gioCompat.LoadState()
		compatState = sharedCompat.Normalize(compatState)
		selectedItem, hasSelected := historyItemBySavedPath(compatState.History, res.SavedPath)
		if !hasSelected {
			selectedItem, hasSelected = newestHistoryItem(compatState.History)
		}
		activeProfileID := ""
		if profile, ok := gioCompat.ActiveProfile(compatState); ok {
			activeProfileID = profile.ID
		}
		a.mu.Lock()
		a.running = false
		a.cancel = nil
		a.status = fmt.Sprintf("完成 - %.1fs", elapsedSec)
		a.result = resultState{
			Image:         img,
			SavedPath:     res.SavedPath,
			RawPath:       res.RawPath,
			RevisedPrompt: res.RevisedPrompt,
			SourceEvent:   res.SourceEvent,
			Item:          selectedItem,
			HasItem:       hasSelected,
			Rev:           a.result.Rev + 1,
		}
		a.history = append([]sharedCompat.HistoryItem(nil), compatState.History...)
		a.profiles = append([]sharedCompat.UpstreamProfile(nil), compatState.Profiles...)
		a.activeProfileID = activeProfileID
		if hasSelected {
			a.selectedHistoryID = selectedItem.ID
		}
		if !a.savePromptSuppressed && res.SavedPath != "" {
			a.savePromptVisible = true
			a.savePromptSourcePath = res.SavedPath
			a.savePromptPathInput.SetText(res.SavedPath)
		}
		a.logs = appendBounded(a.logs, "生成完成: "+res.SavedPath)
		a.mu.Unlock()
		a.invalidateNow()
	}()
}

func (a *App) currentConfig() kernel.Config {
	seed, _ := strconv.ParseInt(strings.TrimSpace(a.seedInput.Text()), 10, 64)
	partial, _ := strconv.Atoi(strings.TrimSpace(a.partialImagesInput.Text()))
	return kernel.Config{
		APIKey:         a.apiKeyInput.Text(),
		BaseURL:        a.baseURLInput.Text(),
		TextModelID:    a.textModelInput.Text(),
		ImageModelID:   a.imageModelInput.Text(),
		Prompt:         a.promptInput.Text(),
		Mode:           client.Mode(a.mode),
		APIMode:        client.APIMode(a.api),
		RequestPolicy:  client.RequestPolicy(a.policy),
		Size:           a.size,
		Quality:        a.quality,
		OutputFormat:   a.format,
		ProxyMode:      a.proxy,
		ProxyURL:       a.proxyURLInput.Text(),
		SourcePaths:    kernel.ParseSourcePaths(a.sourcePathsInput.Text()),
		OutputDir:      a.outputDirInput.Text(),
		Seed:           seed,
		NegativePrompt: a.negativePromptInput.Text(),
		PartialImages:  partial,
	}
}
