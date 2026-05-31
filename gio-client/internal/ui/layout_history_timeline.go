package ui

import (
	"image"
	"strconv"
	"strings"

	sharedCompat "image-studio/shared/compat"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func (a *App) layoutHistoryTimelineModal(gtx layout.Context) layout.Dimensions {
	for a.closeHistoryTimelineButton.Clicked(gtx) {
		a.closeHistoryTimeline()
	}
	for idx, value := range []string{"all", "generate", "edit"} {
		for a.historyTimelineModeButtons[idx].Clicked(gtx) {
			a.historyTimelineModeFilter = value
		}
	}
	for idx, value := range []string{"all", "today", "week"} {
		for a.historyTimelineDateButtons[idx].Clicked(gtx) {
			a.historyTimelineDateFilter = value
		}
	}

	snap := a.readSnapshot()
	if !snap.HistoryTimelineOpen {
		return layout.Dimensions{}
	}

	filtered := a.filteredTimelineHistory(snap.History)
	dayGroups := buildHistoryDayGroups(filtered)
	for _, dayGroup := range dayGroups {
		for _, entry := range dayGroup.Entries {
			if entry.Kind == "group" {
				summaryBtn := a.historyButton("timeline-group:" + entry.Group.Key)
				for summaryBtn.Clicked(gtx) {
					if err := a.loadHistoryPreview(entry.Group.Representative, true); err != nil && !isMissingPreview(err) {
						a.appendLog("载入历史结果失败: " + err.Error())
					} else {
						a.closeHistoryTimeline()
					}
				}
				openBtn := a.historyActionButton("timeline-group-open:" + entry.Group.Key)
				for openBtn.Clicked(gtx) {
					a.openPromptGroup(entry.Group)
					a.closeHistoryTimeline()
				}
				continue
			}

			item := entry.Item
			rowBtn := a.historyButton("timeline-row:" + item.ID)
			for rowBtn.Clicked(gtx) {
				if err := a.loadHistoryPreview(item, true); err != nil && !isMissingPreview(err) {
					a.appendLog("载入历史结果失败: " + err.Error())
				} else {
					a.closeHistoryTimeline()
				}
			}
			detailBtn := a.historyActionButton("timeline-detail:" + item.ID)
			for detailBtn.Clicked(gtx) {
				a.openResultDetail(item)
				a.closeHistoryTimeline()
			}
			reuseBtn := a.historyActionButton("timeline-reuse:" + item.ID)
			for reuseBtn.Clicked(gtx) {
				a.reuseHistoryItemAsSource(item)
				a.closeHistoryTimeline()
			}
			deleteBtn := a.historyActionButton("timeline-delete:" + item.ID)
			for deleteBtn.Clicked(gtx) {
				a.deleteHistoryItem(item.ID)
			}
		}
	}

	paint.FillShape(gtx.Ops, rgba(0x000000, 0x52), clip.Rect{Max: gtx.Constraints.Max}.Op())
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		return fixedWidth(gtx, unit.Dp(920), func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(660), func(gtx layout.Context) layout.Dimensions {
				return a.borderedSurface(gtx, fluent.surface, unit.Dp(8), fluent.border, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(12))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return a.label(gtx, "完整历史", unit.Sp(18), fluent.text, font.SemiBold)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												countText := strconv.Itoa(len(filtered))
												if len(filtered) != len(snap.History) {
													countText += " / " + strconv.Itoa(len(snap.History))
												}
												return a.label(gtx, countText+" 项", unit.Sp(11), fluent.textMuted, font.Normal)
											}),
										)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return fixedWidth(gtx, unit.Dp(104), func(gtx layout.Context) layout.Dimensions {
											return a.compactIconTextButton(gtx, &a.closeHistoryTimelineButton, uiIconClose, "关闭", false)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.searchField(gtx, &a.historyTimelineQueryInput, "搜索 prompt / revised prompt...")
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineModeButtons[0], uiIconList, "全部模式", a.historyTimelineModeFilter == "all")
											}),
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineModeButtons[1], uiIconPlay, "文生图", a.historyTimelineModeFilter == "generate")
											}),
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineModeButtons[2], uiIconEdit, "图生图", a.historyTimelineModeFilter == "edit")
											}),
										)
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineDateButtons[0], uiIconList, "全部日期", a.historyTimelineDateFilter == "all")
											}),
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineDateButtons[1], uiIconHistory, "今天", a.historyTimelineDateFilter == "today")
											}),
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												return a.pillIconTextButton(gtx, &a.historyTimelineDateButtons[2], uiIconHistory, "近 7 天", a.historyTimelineDateFilter == "week")
											}),
										)
									}),
								)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if len(dayGroups) == 0 {
									return a.emptyPanel(gtx, "没有匹配的历史记录")
								}
								return a.historyTimelineList.Layout(gtx, len(dayGroups), func(gtx layout.Context, index int) layout.Dimensions {
									return layout.Inset{Bottom: unit.Dp(14)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return a.layoutHistoryTimelineDayGroup(gtx, dayGroups[index], snap.SelectedHistoryID)
									})
								})
							}),
						)
					})
				})
			})
		})
	})
}

func (a *App) layoutHistoryTimelineDayGroup(gtx layout.Context, dayGroup historyDayGroup, selectedHistoryID string) layout.Dimensions {
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, dayGroup.Label, unit.Sp(13), fluent.text, font.SemiBold)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, strconv.Itoa(len(dayGroup.Entries))+" 组", unit.Sp(11), fluent.textDim, font.Normal)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
	}
	for idx, entry := range dayGroup.Entries {
		entry := entry
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutHistoryTimelineEntry(gtx, entry, selectedHistoryID)
		}))
		if idx != len(dayGroup.Entries)-1 {
			children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout))
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (a *App) layoutHistoryTimelineEntry(gtx layout.Context, entry historyPromptEntry, selectedHistoryID string) layout.Dimensions {
	if entry.Kind == "group" {
		return a.layoutHistoryTimelineGroupRow(gtx, entry.Group, selectedHistoryID)
	}
	return a.layoutHistoryTimelineRow(gtx, entry.Item, entry.Item.ID == selectedHistoryID)
}

func (a *App) layoutHistoryTimelineGroupRow(gtx layout.Context, group historyPromptGroup, selectedHistoryID string) layout.Dimensions {
	active := historyPromptGroupContains(group, selectedHistoryID)
	summaryBtn := a.historyButton("timeline-group:" + group.Key)
	openBtn := a.historyActionButton("timeline-group-open:" + group.Key)
	prompt := group.Prompt
	if prompt == "" {
		prompt = "(无 prompt)"
	}
	meta := strconv.Itoa(len(group.Items)) + " 张 · " + historyMetaText(group.Representative)

	return a.borderedSurface(gtx, chooseColor(active, fluent.surface2, fluent.surface), unit.Dp(6), chooseColor(active, rgba(0x005fb8, 0x48), fluent.border), func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return summaryBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.layoutHistoryGroupPile(gtx, group)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return summaryBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, shortPrompt(prompt), unit.Sp(12), fluent.text, font.Medium)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, meta, unit.Sp(10), fluent.textMuted, font.Normal)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.compactIconTextButton(gtx, openBtn, uiIconGrid, "同提示词", false)
				}),
			)
		})
	})
}

func (a *App) layoutHistoryTimelineRow(gtx layout.Context, item sharedCompat.HistoryItem, active bool) layout.Dimensions {
	rowBtn := a.historyButton("timeline-row:" + item.ID)
	detailBtn := a.historyActionButton("timeline-detail:" + item.ID)
	reuseBtn := a.historyActionButton("timeline-reuse:" + item.ID)
	deleteBtn := a.historyActionButton("timeline-delete:" + item.ID)
	return a.surfaceButton(
		gtx,
		rowBtn,
		chooseColor(active, fluent.surface2, fluent.surface),
		fluent.surface2,
		chooseColor(active, rgba(0x005fb8, 0x48), fluent.border),
		unit.Dp(6),
		layout.Inset{Top: 8, Bottom: 8, Left: 8, Right: 8},
		func(gtx layout.Context) layout.Dimensions {
			img, _ := a.imageForHistoryItem(item)
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutHistoryModeThumb(gtx, img, item.Mode, unit.Dp(72), unit.Dp(72))
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(5))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, shortPrompt(item.Prompt), unit.Sp(12), fluent.text, font.Medium)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, historyMetaText(item), unit.Sp(10), fluent.textMuted, font.Normal)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, strings.Join(compactNonEmpty([]string{
								historyPathText(item.SavedPath),
								formatHistoryClock(item.CreatedAt),
							}), " · "), unit.Sp(10), fluent.textDim, font.Normal)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.End, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if active {
								return a.badge(gtx, "当前", fluent.accentSoft, fluent.accent)
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.compactIconTextButton(gtx, detailBtn, uiIconInfo, "详情", false)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.compactIconTextButton(gtx, reuseBtn, uiIconSource, "设为源图", false)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.compactIconTextButton(gtx, deleteBtn, uiIconDelete, "删除", false)
								}),
							)
						}),
					)
				}),
			)
		},
	)
}
