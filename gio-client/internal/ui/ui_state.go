package ui

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	gioCompat "image-studio/gio-client/internal/compat"
	"image-studio/gio-client/internal/kernel"
	sharedCompat "image-studio/shared/compat"

	"gioui.org/widget"
)

var errMissingPreview = errors.New("missing preview image")

func isMissingPreview(err error) bool {
	return errors.Is(err, errMissingPreview)
}

func newestHistoryItem(items []sharedCompat.HistoryItem) (sharedCompat.HistoryItem, bool) {
	if len(items) == 0 {
		return sharedCompat.HistoryItem{}, false
	}
	return items[0], true
}

func historyItemBySavedPath(items []sharedCompat.HistoryItem, savedPath string) (sharedCompat.HistoryItem, bool) {
	savedPath = strings.TrimSpace(savedPath)
	if savedPath == "" {
		return sharedCompat.HistoryItem{}, false
	}
	for _, item := range items {
		if strings.TrimSpace(item.SavedPath) == savedPath {
			return item, true
		}
	}
	return sharedCompat.HistoryItem{}, false
}

func historyItemByID(items []sharedCompat.HistoryItem, id string) (sharedCompat.HistoryItem, bool) {
	id = strings.TrimSpace(id)
	if id == "" {
		return sharedCompat.HistoryItem{}, false
	}
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return sharedCompat.HistoryItem{}, false
}

func historyCounts(items []sharedCompat.HistoryItem) (generate int, edit int) {
	for _, item := range items {
		if item.Mode == "edit" {
			edit++
			continue
		}
		generate++
	}
	return generate, edit
}

func todayHistoryCount(items []sharedCompat.HistoryItem, now time.Time) int {
	start := localDayStart(now)
	count := 0
	for _, item := range items {
		if item.CreatedAt >= start.UnixMilli() {
			count++
		}
	}
	return count
}

func (a *App) profileButton(id string) *widget.Clickable {
	if a.profileButtons == nil {
		a.profileButtons = map[string]*widget.Clickable{}
	}
	if btn, ok := a.profileButtons[id]; ok {
		return btn
	}
	btn := new(widget.Clickable)
	a.profileButtons[id] = btn
	return btn
}

func (a *App) historyButton(id string) *widget.Clickable {
	if a.historyButtons == nil {
		a.historyButtons = map[string]*widget.Clickable{}
	}
	if btn, ok := a.historyButtons[id]; ok {
		return btn
	}
	btn := new(widget.Clickable)
	a.historyButtons[id] = btn
	return btn
}

func (a *App) filteredHistory(items []sharedCompat.HistoryItem) []sharedCompat.HistoryItem {
	query := strings.TrimSpace(strings.ToLower(a.historyQueryInput.Text()))
	modeFilter := strings.TrimSpace(a.historyModeFilter)
	dateFilter := strings.TrimSpace(a.historyDateFilter)
	if query == "" && modeFilter == "all" && dateFilter == "all" {
		return items
	}
	filtered := make([]sharedCompat.HistoryItem, 0, len(items))
	for _, item := range items {
		if modeFilter != "" && modeFilter != "all" && item.Mode != modeFilter {
			continue
		}
		if !matchHistoryDate(item.CreatedAt, dateFilter, time.Now()) {
			continue
		}
		if !matchHistoryQuery(item, query) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func (a *App) loadHistoryPreview(item sharedCompat.HistoryItem, addLog bool) error {
	img, err := a.imageForHistoryItem(item)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.result = resultState{
		Image:         img,
		SavedPath:     item.SavedPath,
		RawPath:       item.RawPath,
		RevisedPrompt: item.RevisedPrompt,
		SourceEvent:   "history",
		Item:          item,
		HasItem:       item.ID != "",
		Rev:           a.result.Rev + 1,
	}
	a.selectedHistoryID = item.ID
	a.status = "已载入历史结果"
	if addLog {
		a.logs = appendBounded(a.logs, "载入历史结果: "+shortPrompt(item.Prompt))
	}
	a.mu.Unlock()
	a.invalidateNow()
	return nil
}

func (a *App) imageForHistoryItem(item sharedCompat.HistoryItem) (image.Image, error) {
	cacheKey := strings.TrimSpace(item.SavedPath)
	if cacheKey == "" {
		cacheKey = "history:" + item.ID
	}
	if cached, ok := a.imageCache[cacheKey]; ok {
		if cached.Failed {
			return nil, errMissingPreview
		}
		return cached.Image, nil
	}

	load := func() (image.Image, error) {
		if strings.TrimSpace(item.ImageB64) != "" {
			return decodeImageB64(item.ImageB64)
		}
		if strings.TrimSpace(item.SavedPath) == "" {
			return nil, errMissingPreview
		}
		return a.imageForPath(item.SavedPath)
	}

	img, err := load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, errMissingPreview) {
			a.imageCache[cacheKey] = cachedImage{Failed: true}
			return nil, fmt.Errorf("%w: %v", errMissingPreview, err)
		}
		a.imageCache[cacheKey] = cachedImage{Failed: true}
		return nil, err
	}
	a.imageCache[cacheKey] = cachedImage{Image: img}
	return img, nil
}

func (a *App) imageForPath(path string) (image.Image, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errMissingPreview
	}
	cacheKey := "path:" + path
	if cached, ok := a.imageCache[cacheKey]; ok {
		if cached.Failed {
			return nil, errMissingPreview
		}
		return cached.Image, nil
	}
	img, err := decodeImageFile(path)
	if err != nil {
		a.imageCache[cacheKey] = cachedImage{Failed: true}
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %v", errMissingPreview, err)
		}
		return nil, err
	}
	a.imageCache[cacheKey] = cachedImage{Image: img}
	return img, nil
}

func (a *App) switchActiveProfile(profileID string) {
	profileID = strings.TrimSpace(profileID)
	if profileID == "" {
		return
	}
	state, _, err := gioCompat.LoadState()
	if err != nil {
		a.appendLog("读取上游配置失败: " + err.Error())
		return
	}
	state = sharedCompat.Normalize(state)
	if _, ok := historyItemByID(state.History, a.selectedHistoryID); !ok && len(state.History) > 0 {
		a.selectedHistoryID = state.History[0].ID
	}
	found := false
	for _, profile := range state.Profiles {
		if profile.ID == profileID {
			found = true
			break
		}
	}
	if !found {
		return
	}
	state.ActiveProfile = profileID
	state.UpdatedAt = time.Now().UnixMilli()
	if err := gioCompat.SaveState(state); err != nil {
		a.appendLog("切换上游失败: " + err.Error())
		return
	}
	cfg := gioCompat.ConfigFromState(kernel.DefaultConfig(), state)
	a.applyRuntimeConfig(cfg)
	a.mu.Lock()
	a.profiles = append([]sharedCompat.UpstreamProfile(nil), state.Profiles...)
	a.activeProfileID = profileID
	a.status = "已切换上游: " + activeProfileName(state.Profiles, profileID)
	a.logs = appendBounded(a.logs, "切换上游配置: "+activeProfileName(state.Profiles, profileID))
	a.mu.Unlock()
	a.profilePickerOpen = false
	a.invalidateNow()
}

func activeProfileName(profiles []sharedCompat.UpstreamProfile, profileID string) string {
	for _, profile := range profiles {
		if profile.ID == profileID {
			return strings.TrimSpace(profile.Name)
		}
	}
	return ""
}
