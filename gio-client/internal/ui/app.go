package ui

import (
	"context"
	"image"
	"strconv"
	"strings"
	"sync"

	gioCompat "image-studio/gio-client/internal/compat"
	"image-studio/gio-client/internal/kernel"
	sharedCompat "image-studio/shared/compat"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type resultState struct {
	Image         image.Image
	SavedPath     string
	RawPath       string
	RevisedPrompt string
	SourceEvent   string
	Item          sharedCompat.HistoryItem
	HasItem       bool
	Rev           int
}

type snapshot struct {
	Running           bool
	Status            string
	Logs              []string
	History           []sharedCompat.HistoryItem
	Profiles          []sharedCompat.UpstreamProfile
	ActiveProfileID   string
	SelectedHistoryID string
	Result            resultState
	SavePromptVisible bool
}

type cachedImage struct {
	Image  image.Image
	Failed bool
}

type App struct {
	th     *material.Theme
	runner kernel.Runner

	controlsList widget.List
	logList      widget.List
	historyList  widget.List

	apiKeyInput         widget.Editor
	baseURLInput        widget.Editor
	textModelInput      widget.Editor
	imageModelInput     widget.Editor
	promptInput         widget.Editor
	sourcePathsInput    widget.Editor
	outputDirInput      widget.Editor
	seedInput           widget.Editor
	negativePromptInput widget.Editor
	partialImagesInput  widget.Editor
	proxyURLInput       widget.Editor
	savePromptPathInput widget.Editor
	historyQueryInput   widget.Editor

	mode    string
	api     string
	size    string
	quality string
	format  string
	policy  string
	proxy   string

	modeButtons           []widget.Clickable
	apiButtons            []widget.Clickable
	sizeButtons           []widget.Clickable
	qualityButtons        []widget.Clickable
	formatButtons         []widget.Clickable
	policyButtons         []widget.Clickable
	proxyButtons          []widget.Clickable
	historyModeButtons    []widget.Clickable
	historyDateButtons    []widget.Clickable
	runButton             widget.Clickable
	cancelButton          widget.Clickable
	clearLogButton        widget.Clickable
	composeToggleButton   widget.Clickable
	advancedToggleButton  widget.Clickable
	profilePickerButton   widget.Clickable
	manageUpstreamButton  widget.Clickable
	historyCollapseButton widget.Clickable
	savePromptSaveButton  widget.Clickable
	savePromptSkipButton  widget.Clickable
	savePromptNeverAsk    widget.Bool

	mu                sync.Mutex
	running           bool
	cancel            context.CancelFunc
	status            string
	logs              []string
	history           []sharedCompat.HistoryItem
	profiles          []sharedCompat.UpstreamProfile
	activeProfileID   string
	selectedHistoryID string
	result            resultState
	imageOp           paint.ImageOp
	imageOpRev        int
	imageCache        map[string]cachedImage

	savePromptVisible    bool
	savePromptSuppressed bool
	savePromptSourcePath string
	composeOpen          bool
	advancedOpen         bool
	profilePickerOpen    bool
	historyRailCollapsed bool
	historyModeFilter    string
	historyDateFilter    string
	profileButtons       map[string]*widget.Clickable
	historyButtons       map[string]*widget.Clickable

	invalidate func()
}

func New() *App {
	cfg := kernel.DefaultConfig()
	compatState, compatPath, compatErr := gioCompat.LoadState()
	if compatErr == nil {
		cfg = gioCompat.ConfigFromState(cfg, compatState)
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette = material.Palette{
		Bg:         fluent.bg,
		Fg:         fluent.text,
		ContrastBg: fluent.accent,
		ContrastFg: fluent.white,
	}
	th.TextSize = unit.Sp(14)
	a := &App{
		th:                   th,
		runner:               kernel.Runner{},
		mode:                 string(cfg.Mode),
		api:                  string(cfg.APIMode),
		size:                 cfg.Size,
		quality:              cfg.Quality,
		format:               cfg.OutputFormat,
		policy:               string(cfg.RequestPolicy),
		proxy:                cfg.ProxyMode,
		modeButtons:          make([]widget.Clickable, len(modeChoices)),
		apiButtons:           make([]widget.Clickable, len(apiChoices)),
		sizeButtons:          make([]widget.Clickable, len(sizeChoices)),
		qualityButtons:       make([]widget.Clickable, len(qualityChoices)),
		formatButtons:        make([]widget.Clickable, len(formatChoices)),
		policyButtons:        make([]widget.Clickable, len(policyChoices)),
		proxyButtons:         make([]widget.Clickable, len(proxyChoices)),
		historyModeButtons:   make([]widget.Clickable, 3),
		historyDateButtons:   make([]widget.Clickable, 3),
		status:               "Gio 原生客户端就绪",
		logs:                 []string{"独立 Gio 高性能测试客户端已启动。"},
		history:              append([]sharedCompat.HistoryItem(nil), compatState.History...),
		profiles:             append([]sharedCompat.UpstreamProfile(nil), compatState.Profiles...),
		savePromptSuppressed: gioCompat.SavePromptSuppressed(compatState),
		imageCache:           map[string]cachedImage{},
		composeOpen:          false,
		advancedOpen:         false,
		profilePickerOpen:    false,
		historyRailCollapsed: false,
		historyModeFilter:    "all",
		historyDateFilter:    "all",
		profileButtons:       map[string]*widget.Clickable{},
		historyButtons:       map[string]*widget.Clickable{},
	}
	if profile, ok := gioCompat.ActiveProfile(compatState); ok {
		a.activeProfileID = profile.ID
	}
	a.savePromptNeverAsk.Value = a.savePromptSuppressed
	if compatPath != "" {
		a.logs = append(a.logs, "兼容状态文件: "+compatPath)
	}
	if compatErr != nil {
		a.logs = append(a.logs, "兼容状态读取失败: "+compatErr.Error())
	}
	a.controlsList.List.Axis = layout.Vertical
	a.logList.List.Axis = layout.Vertical
	a.historyList.List.Axis = layout.Vertical
	a.configureEditors(cfg)
	a.historyQueryInput.SingleLine = true
	if latest, ok := newestHistoryItem(a.history); ok {
		if err := a.loadHistoryPreview(latest, false); err != nil && !isMissingPreview(err) {
			a.logs = appendBounded(a.logs, "载入最近历史失败: "+err.Error())
		}
	}
	return a
}

func (a *App) configureEditors(cfg kernel.Config) {
	singleLine := []*widget.Editor{
		&a.apiKeyInput,
		&a.baseURLInput,
		&a.textModelInput,
		&a.imageModelInput,
		&a.outputDirInput,
		&a.seedInput,
		&a.partialImagesInput,
		&a.proxyURLInput,
		&a.savePromptPathInput,
		&a.historyQueryInput,
	}
	for _, editor := range singleLine {
		editor.SingleLine = true
	}
	a.apiKeyInput.Mask = '*'
	a.seedInput.Filter = "0123456789"
	a.partialImagesInput.Filter = "0123456789"
	a.apiKeyInput.SetText(cfg.APIKey)
	a.baseURLInput.SetText(cfg.BaseURL)
	a.textModelInput.SetText(cfg.TextModelID)
	a.imageModelInput.SetText(cfg.ImageModelID)
	a.outputDirInput.SetText(cfg.OutputDir)
	a.partialImagesInput.SetText(strconv.Itoa(cfg.PartialImages))
	a.proxyURLInput.SetText(cfg.ProxyURL)
	a.promptInput.SetText("")
}

func (a *App) applyRuntimeConfig(cfg kernel.Config) {
	if strings.TrimSpace(cfg.APIKey) != "" || strings.TrimSpace(a.apiKeyInput.Text()) != "" {
		a.apiKeyInput.SetText(cfg.APIKey)
	}
	a.baseURLInput.SetText(cfg.BaseURL)
	a.textModelInput.SetText(cfg.TextModelID)
	a.imageModelInput.SetText(cfg.ImageModelID)
	a.proxy = cfg.ProxyMode
	a.proxyURLInput.SetText(cfg.ProxyURL)
	a.outputDirInput.SetText(cfg.OutputDir)
	if strings.TrimSpace(cfg.OutputFormat) != "" {
		a.format = cfg.OutputFormat
	}
	if strings.TrimSpace(string(cfg.APIMode)) != "" {
		a.api = string(cfg.APIMode)
	}
	if strings.TrimSpace(string(cfg.RequestPolicy)) != "" {
		a.policy = string(cfg.RequestPolicy)
	}
}

func (a *App) Run(w *app.Window) error {
	a.invalidate = w.Invalidate
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			a.saveCurrentConfig()
			a.cancelRun()
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
