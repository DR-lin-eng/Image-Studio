package ui

import (
	"strconv"
	"strings"

	sharedCompat "image-studio/shared/compat"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/yuanhua/image-gptcodex/pkg/client"
)

func (a *App) layoutHistoryAndLogs(gtx layout.Context) layout.Dimensions {
	return a.borderedSurface(gtx, fluent.inspector, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
				layout.Rigid(a.layoutUpstreamCard),
				layout.Flexed(0.56, a.layoutHistory),
				layout.Flexed(0.44, a.layoutLogs),
			)
		})
	})
}

func (a *App) layoutUpstreamCard(gtx layout.Context) layout.Dimensions {
	ready := strings.TrimSpace(a.apiKeyInput.Text()) != "" && strings.TrimSpace(a.baseURLInput.Text()) != ""
	status := "未配置"
	statusColor := fluent.danger
	dotColor := fluent.danger
	if ready {
		status = "已配置"
		statusColor = fluent.accent
		dotColor = fluent.accent
	}
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
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
						return a.label(gtx, status, unit.Sp(11), statusColor, font.Medium)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				apiMode := "Responses API"
				if a.api == string(client.APIModeImages) {
					apiMode = "Images API"
				}
				return a.label(gtx, apiMode+" - "+strings.TrimSpace(a.baseURLInput.Text()), unit.Sp(11), fluent.textMuted, font.Normal)
			}),
		)
	})
}

func (a *App) layoutHistory(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "历史")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.badge(gtx, strconv.Itoa(len(snap.History))+" 项", fluent.surface2, fluent.textMuted)
					}),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(snap.History) == 0 {
					return a.emptyPanel(gtx, "还没有结果")
				}
				return a.historyList.Layout(gtx, len(snap.History), func(gtx layout.Context, i int) layout.Dimensions {
					return layout.Inset{Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.historyRow(gtx, snap.History[i])
					})
				})
			}),
		)
	})
}

func (a *App) historyRow(gtx layout.Context, item sharedCompat.HistoryItem) layout.Dimensions {
	return a.borderedSurface(gtx, fluent.surface, unit.Dp(6), fluent.border, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			mode := "文生图"
			if item.Mode == string(client.ModeEdit) {
				mode = "图生图"
			}
			meta := strings.Join(compactNonEmpty([]string{mode, item.Size, item.Quality, item.OutputFormat}), " · ")
			saved := item.SavedPath
			if saved == "" {
				saved = "未登记保存路径"
			}
			return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(5))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, shortPrompt(item.Prompt), unit.Sp(12), fluent.text, font.Medium)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, meta, unit.Sp(11), fluent.textMuted, font.Normal)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, saved, unit.Sp(10), fluent.textDim, font.Normal)
				}),
			)
		})
	})
}

func (a *App) layoutLogs(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "运行日志")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return fixedWidth(gtx, unit.Dp(68), func(gtx layout.Context) layout.Dimensions {
							return a.button(gtx, &a.clearLogButton, "清空", fluent.surface2, fluent.text)
						})
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
						return a.borderedSurface(gtx, fluent.surface, unit.Dp(6), fluent.border, func(gtx layout.Context) layout.Dimensions {
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
				return a.label(gtx, raw, unit.Sp(11), fluent.textDim, font.Normal)
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
