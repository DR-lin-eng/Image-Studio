package ui

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sharedCompat "image-studio/shared/compat"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
)

func (a *App) layoutHistoryAndLogs(gtx layout.Context) layout.Dimensions {
	for a.profilePickerButton.Clicked(gtx) {
		a.profilePickerOpen = !a.profilePickerOpen
	}
	for a.historyCollapseButton.Clicked(gtx) {
		a.historyRailCollapsed = !a.historyRailCollapsed
	}
	for idx, value := range []string{"all", "generate", "edit"} {
		for a.historyModeButtons[idx].Clicked(gtx) {
			a.historyModeFilter = value
		}
	}
	for idx, value := range []string{"all", "today", "week"} {
		for a.historyDateButtons[idx].Clicked(gtx) {
			a.historyDateFilter = value
		}
	}

	snap := a.readSnapshot()
	filtered := a.filteredHistory(snap.History)
	generateCount, editCount := historyCounts(snap.History)
	latest, hasLatest := newestHistoryItem(filtered)
	visible := filtered
	if len(visible) > 18 {
		visible = visible[:18]
	}

	for _, profile := range snap.Profiles {
		button := a.profileButton("profile:" + profile.ID)
		for button.Clicked(gtx) {
			a.switchActiveProfile(profile.ID)
		}
	}
	for _, item := range visible {
		button := a.historyButton("row:" + item.ID)
		for button.Clicked(gtx) {
			if err := a.loadHistoryPreview(item, true); err != nil && !isMissingPreview(err) {
				a.appendLog("载入历史结果失败: " + err.Error())
			}
		}
	}
	if hasLatest {
		button := a.historyButton("feature:" + latest.ID)
		for button.Clicked(gtx) {
			if err := a.loadHistoryPreview(latest, true); err != nil && !isMissingPreview(err) {
				a.appendLog("载入最近作品失败: " + err.Error())
			}
		}
	}

	return a.borderedSurface(gtx, fluent.inspector, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			children := []layout.FlexChild{
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutUpstreamCard(gtx, snap)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.layoutHistorySummaryCard(gtx, snap, filtered, generateCount, editCount)
				}),
			}

			if !a.historyRailCollapsed && hasLatest {
				children = append(children,
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutLatestHistoryCard(gtx, latest, snap.SelectedHistoryID == latest.ID)
					}),
				)
			}

			children = append(children,
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if a.historyRailCollapsed {
						return a.layoutLogsCard(gtx, snap)
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return a.layoutHistoryResultsCard(gtx, snap, filtered, visible)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return fixedHeight(gtx, unit.Dp(196), func(gtx layout.Context) layout.Dimensions {
								return a.layoutLogsCard(gtx, snap)
							})
						}),
					)
				}),
			)
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
	})
}

func (a *App) layoutUpstreamCard(gtx layout.Context, snap snapshot) layout.Dimensions {
	activeName := strings.TrimSpace(activeProfileName(snap.Profiles, snap.ActiveProfileID))
	if activeName == "" {
		activeName = "还没有上游配置"
	}
	apiModeLabel := "Responses API"
	if a.api == "images" {
		apiModeLabel = "Images API"
	}
	ready := strings.TrimSpace(a.apiKeyInput.Text()) != "" && strings.TrimSpace(a.baseURLInput.Text()) != ""
	statusLabel := "未配置"
	statusColor := fluent.danger
	dotColor := fluent.danger
	if ready {
		statusLabel = "已配置"
		statusColor = fluent.accent
		dotColor = fluent.accent
	}

	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "上游")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return fixedWidth(gtx, unit.Dp(7), func(gtx layout.Context) layout.Dimensions {
							return fixedHeight(gtx, unit.Dp(7), func(gtx layout.Context) layout.Dimensions {
								return a.surface(gtx, dotColor, unit.Dp(4), layout.Spacer{}.Layout)
							})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, statusLabel, unit.Sp(11), statusColor, font.Medium)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		}

		if len(snap.Profiles) == 0 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.label(gtx, "还没有上游配置，先在左侧高级参数里补上 BASE_URL 和 API Key。", unit.Sp(11), fluent.textMuted, font.Normal)
			}))
		} else {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.surfaceButton(
						gtx,
						&a.profilePickerButton,
						chooseColor(a.profilePickerOpen, fluent.surface2, fluent.surface),
						fluent.surface2,
						fluent.border,
						unit.Dp(4),
						layout.Inset{Top: 9, Bottom: 9, Left: 10, Right: 10},
						func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(3))}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return a.label(gtx, activeName, unit.Sp(12), fluent.text, font.Medium)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return a.label(gtx, apiModeLabel, unit.Sp(11), fluent.textMuted, font.Normal)
										}),
									)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									state := "展开"
									if a.profilePickerOpen {
										state = "收起"
									}
									return a.label(gtx, state, unit.Sp(11), fluent.textDim, font.Medium)
								}),
							)
						},
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			)
			if a.profilePickerOpen {
				for _, profile := range snap.Profiles {
					profile := profile
					children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutProfileOption(gtx, profile, profile.ID == snap.ActiveProfileID)
					}))
					children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout))
				}
			}
		}

		if strings.TrimSpace(a.baseURLInput.Text()) != "" {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.label(gtx, apiModeLabel+" · "+strings.TrimSpace(a.baseURLInput.Text()), unit.Sp(11), fluent.textDim, font.Normal)
			}))
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (a *App) layoutProfileOption(gtx layout.Context, profile sharedCompat.UpstreamProfile, active bool) layout.Dimensions {
	btn := a.profileButton("profile:" + profile.ID)
	return a.surfaceButton(
		gtx,
		btn,
		chooseColor(active, fluent.accentSoft, fluent.surface),
		chooseColor(active, rgba(0x005fb8, 0x28), fluent.surface2),
		fluent.border,
		unit.Dp(4),
		layout.Inset{Top: 9, Bottom: 9, Left: 10, Right: 10},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(3))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, strings.TrimSpace(profile.Name), unit.Sp(12), chooseColor(active, fluent.accent, fluent.text), font.Medium)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, strings.ToUpper(strings.TrimSpace(profile.APIMode)), unit.Sp(10), fluent.textDim, font.Normal)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if !active {
						return layout.Dimensions{}
					}
					return a.badge(gtx, "当前", fluent.accentSoft, fluent.accent)
				}),
			)
		},
	)
}

func (a *App) layoutHistorySummaryCard(
	gtx layout.Context,
	snap snapshot,
	filtered []sharedCompat.HistoryItem,
	generateCount int,
	editCount int,
) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(3))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.sectionEyebrow(gtx, "历史")
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
						return a.compactButton(gtx, &a.historyCollapseButton, chooseHistoryCollapseLabel(a.historyRailCollapsed), a.historyRailCollapsed)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyModeButtons[0], "全部 "+strconv.Itoa(len(snap.History)), a.historyModeFilter == "all")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyModeButtons[1], "文生图 "+strconv.Itoa(generateCount), a.historyModeFilter == "generate")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyModeButtons[2], "图生图 "+strconv.Itoa(editCount), a.historyModeFilter == "edit")
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.searchField(gtx, &a.historyQueryInput, "搜索 prompt...")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyDateButtons[0], "全部", a.historyDateFilter == "all")
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyDateButtons[1], "今天", a.historyDateFilter == "today")
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.pillButton(gtx, &a.historyDateButtons[2], "本周", a.historyDateFilter == "week")
					}),
				)
			}),
		)
	})
}

func (a *App) layoutLatestHistoryCard(gtx layout.Context, item sharedCompat.HistoryItem, active bool) layout.Dimensions {
	btn := a.historyButton("feature:" + item.ID)
	return a.surfaceButton(
		gtx,
		btn,
		chooseColor(active, fluent.surface2, fluent.surface),
		fluent.surface2,
		chooseColor(active, rgba(0x005fb8, 0x48), fluent.border),
		unit.Dp(6),
		layout.Inset{Top: 10, Bottom: 10, Left: 10, Right: 10},
		func(gtx layout.Context) layout.Dimensions {
			img, _ := a.imageForHistoryItem(item)
			return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return a.sectionEyebrow(gtx, "最近作品")
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, formatHistoryClock(item.CreatedAt), unit.Sp(11), fluent.textDim, font.Normal)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.imageThumb(gtx, img, unit.Dp(88), unit.Dp(88), unit.Dp(4))
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(5))}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.label(gtx, shortPrompt(item.Prompt), unit.Sp(12), fluent.text, font.Medium)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.label(gtx, historyMetaText(item), unit.Sp(11), fluent.textMuted, font.Normal)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return a.label(gtx, historyPathText(item.SavedPath), unit.Sp(10), fluent.textDim, font.Normal)
								}),
							)
						}),
					)
				}),
			)
		},
	)
}

func (a *App) layoutHistoryResultsCard(
	gtx layout.Context,
	snap snapshot,
	filtered []sharedCompat.HistoryItem,
	visible []sharedCompat.HistoryItem,
) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "结果")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := strconv.Itoa(len(visible))
						if len(filtered) > len(visible) {
							label += " / " + strconv.Itoa(len(filtered))
						}
						return a.label(gtx, label, unit.Sp(11), fluent.textMuted, font.Normal)
					}),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(visible) == 0 {
					text := "还没有结果"
					if len(filtered) == 0 && len(snap.History) > 0 {
						text = "没有匹配项"
					}
					return a.emptyPanel(gtx, text)
				}
				return a.historyList.Layout(gtx, len(visible), func(gtx layout.Context, i int) layout.Dimensions {
					item := visible[i]
					return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.layoutHistoryRow(gtx, item, item.ID == snap.SelectedHistoryID)
					})
				})
			}),
		)
	})
}

func (a *App) layoutHistoryRow(gtx layout.Context, item sharedCompat.HistoryItem, active bool) layout.Dimensions {
	btn := a.historyButton("row:" + item.ID)
	return a.surfaceButton(
		gtx,
		btn,
		chooseColor(active, fluent.surface2, fluent.surface),
		fluent.surface2,
		chooseColor(active, rgba(0x005fb8, 0x48), fluent.border),
		unit.Dp(6),
		layout.Inset{Top: 7, Bottom: 7, Left: 7, Right: 7},
		func(gtx layout.Context) layout.Dimensions {
			img, _ := a.imageForHistoryItem(item)
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.imageThumb(gtx, img, unit.Dp(48), unit.Dp(48), unit.Dp(4))
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, shortPrompt(item.Prompt), unit.Sp(12), fluent.text, font.Medium)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, historyMetaText(item), unit.Sp(10), fluent.textMuted, font.Normal)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, historyPathText(item.SavedPath), unit.Sp(10), fluent.textDim, font.Normal)
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
							return a.label(gtx, formatHistoryClock(item.CreatedAt), unit.Sp(10), fluent.textDim, font.Normal)
						}),
					)
				}),
			)
		},
	)
}

func (a *App) layoutLogsCard(gtx layout.Context, snap snapshot) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "运行日志")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactButton(gtx, &a.clearLogButton, "清空", false)
					}),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(snap.Logs) == 0 {
					return a.emptyPanel(gtx, "暂无日志")
				}
				return a.logList.Layout(gtx, len(snap.Logs), func(gtx layout.Context, i int) layout.Dimensions {
					idx := len(snap.Logs) - 1 - i
					line := snap.Logs[idx]
					return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.borderedSurface(gtx, fluent.surface, unit.Dp(4), fluent.border, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, line, unit.Sp(11), fluent.textMuted, font.Normal)
							})
						})
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				raw := strings.TrimSpace(snap.Result.RawPath)
				if raw == "" {
					raw = "Raw response: 暂无"
				} else {
					raw = "Raw response: " + raw
				}
				return a.label(gtx, raw, unit.Sp(10), fluent.textDim, font.Normal)
			}),
		)
	})
}

func (a *App) emptyPanel(gtx layout.Context, text string) layout.Dimensions {
	return a.surface(gtx, fluent.surface2, unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, text, unit.Sp(12), fluent.textMuted, font.Normal)
		})
	})
}

func chooseHistoryCollapseLabel(collapsed bool) string {
	if collapsed {
		return "展开"
	}
	return "折叠"
}

func historyMetaText(item sharedCompat.HistoryItem) string {
	mode := "文生图"
	if item.Mode == "edit" {
		mode = "图生图"
	}
	format := strings.ToUpper(strings.TrimSpace(item.OutputFormat))
	return strings.Join(compactNonEmpty([]string{mode, item.Size, item.Quality, format}), " · ")
}

func historyPathText(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "未登记保存路径"
	}
	return filepath.Base(path)
}

func formatHistoryClock(createdAt int64) string {
	if createdAt <= 0 {
		return ""
	}
	return time.UnixMilli(createdAt).Format("15:04")
}
