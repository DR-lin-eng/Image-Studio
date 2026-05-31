package ui

import (
	"fmt"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	for a.runButton.Clicked(gtx) {
		a.startRun()
	}
	for a.cancelButton.Clicked(gtx) {
		a.cancelRun()
	}
	for a.clearLogButton.Clicked(gtx) {
		a.clearLogs()
	}

	paint.FillShape(gtx.Ops, fluent.bg, clip.Rect{Max: gtx.Constraints.Max}.Op())
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(48), a.layoutHeader)
		}),
		layout.Flexed(1, a.layoutBody),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(36), a.layoutFooter)
		}),
	)
	if a.readSnapshot().SavePromptVisible {
		a.layoutSavePrompt(gtx)
	}
	return dims
}

func (a *App) layoutHeader(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	return a.borderedSurface(gtx, fluent.toolbar, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return fixedWidth(gtx, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
						return fixedHeight(gtx, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
							return a.surface(gtx, fluent.accent, unit.Dp(5), layout.Spacer{}.Layout)
						})
					})
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(1))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, "Image Studio", unit.Sp(14), fluent.text, font.SemiBold)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, "Gio - Windows / Linux 原生客户端 · 对齐 Windows Fluent 工作流", unit.Sp(11), fluent.textMuted, font.Normal)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.badge(gtx, a.modeLabel(), fluent.accentSoft, fluent.accent)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := snap.Status
					if snap.Running {
						label = "运行中 - " + label
					}
					return a.badge(gtx, label, fluent.surface, fluent.textMuted)
				}),
			)
		})
	})
}

func (a *App) layoutFooter(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	state := "就绪"
	dot := fluent.textDim
	if snap.Running {
		state = "运行中"
		dot = fluent.accent
	}
	todayCount := todayHistoryCount(snap.History, time.Now())
	return a.borderedSurface(gtx, fluent.toolbar, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 8, Bottom: 8, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, "输出目录: "+strings.TrimSpace(a.outputDirInput.Text()), unit.Sp(11), fluent.textMuted, font.Normal)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, fmt.Sprintf("今日已生图: %d · 总生图: %d", todayCount, len(snap.History)), unit.Sp(11), fluent.textMuted, font.Medium)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(14)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return fixedWidth(gtx, unit.Dp(7), func(gtx layout.Context) layout.Dimensions {
								return fixedHeight(gtx, unit.Dp(7), func(gtx layout.Context) layout.Dimensions {
									return a.surface(gtx, dot, unit.Dp(4), layout.Spacer{}.Layout)
								})
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, state, unit.Sp(11), fluent.textMuted, font.Medium)
						}),
					)
				}),
			)
		})
	})
}

func (a *App) layoutBody(gtx layout.Context) layout.Dimensions {
	width := gtx.Constraints.Max.X
	rightWidth := unit.Dp(320)
	leftWidth := unit.Dp(372)
	if width < gtx.Dp(unit.Dp(1180)) {
		rightWidth = unit.Dp(300)
		leftWidth = unit.Dp(336)
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedWidth(gtx, leftWidth, a.layoutControls)
		}),
		layout.Flexed(1, a.layoutCanvas),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedWidth(gtx, rightWidth, a.layoutHistoryAndLogs)
		}),
	)
}
