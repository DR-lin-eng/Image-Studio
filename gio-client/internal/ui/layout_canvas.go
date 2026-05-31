package ui

import (
	"image"
	"strconv"
	"strings"

	"image-studio/gio-client/internal/kernel"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (a *App) layoutCanvas(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return fixedHeight(gtx, unit.Dp(48), func(gtx layout.Context) layout.Dimensions {
					return a.canvasToolbar(gtx, snap)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return fixedHeight(gtx, unit.Dp(56), func(gtx layout.Context) layout.Dimensions {
					return a.sourceStrip(gtx)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return a.resultSurface(gtx, snap)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.canvasStatusBar(gtx, snap)
			}),
		)
	})
}

func (a *App) canvasToolbar(gtx layout.Context, snap snapshot) layout.Dimensions {
	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		source := snap.Result.SourceEvent
		if source == "" {
			source = "idle"
		}
		return layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, "画布", true)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, a.modeLabel(), false)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, strings.ToUpper(a.format), false)
				}),
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.badge(gtx, source, fluent.surface, fluent.textMuted)
				}),
			)
		})
	})
}

func (a *App) sourceStrip(gtx layout.Context) layout.Dimensions {
	paths := len(kernel.ParseSourcePaths(a.sourcePathsInput.Text()))
	label := "源图片: 未添加"
	if paths > 0 {
		label = "源图片: " + strconv.Itoa(paths) + " 张"
	}
	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, label, unit.Sp(11), fluent.textMuted, font.Medium)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return fixedWidth(gtx, unit.Dp(56), func(gtx layout.Context) layout.Dimensions {
						return fixedHeight(gtx, unit.Dp(38), func(gtx layout.Context) layout.Dimensions {
							return a.surface(gtx, fluent.surface, unit.Dp(4), layout.Spacer{}.Layout)
						})
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, strings.TrimSpace(a.sourcePathsInput.Text()), unit.Sp(11), fluent.textDim, font.Normal)
				}),
			)
		})
	})
}

func (a *App) resultSurface(gtx layout.Context, snap snapshot) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return a.surface(gtx, fluent.canvasBg, unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		if snap.Result.Image == nil {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "等待生成结果", unit.Sp(18), fluent.text, font.Medium)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "Gio 原生画布 - 不经过 WebView2/WebKit", unit.Sp(13), fluent.textMuted, font.Normal)
					}),
				)
			})
		}
		if snap.Result.Rev != a.imageOpRev {
			a.imageOp = paint.NewImageOp(snap.Result.Image)
			a.imageOpRev = snap.Result.Rev
		}
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			img := widget.Image{
				Src:      a.imageOp,
				Fit:      widget.Contain,
				Position: layout.Center,
			}
			return img.Layout(gtx)
		})
	})
}

func (a *App) canvasStatusBar(gtx layout.Context, snap snapshot) layout.Dimensions {
	revised := strings.TrimSpace(snap.Result.RevisedPrompt)
	if revised == "" {
		revised = "暂无修订提示词"
	}
	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: 9, Bottom: 9, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, revised, unit.Sp(12), fluent.textMuted, font.Normal)
		})
	})
}

func (a *App) layoutSavePrompt(gtx layout.Context) layout.Dimensions {
	if a.savePromptNeverAsk.Update(gtx) {
		a.setSavePromptSuppressed(a.savePromptNeverAsk.Value)
	}
	for a.savePromptSkipButton.Clicked(gtx) {
		a.closeSavePrompt()
	}
	for a.savePromptSaveButton.Clicked(gtx) {
		a.savePromptCopy()
	}

	paint.FillShape(gtx.Ops, rgba(0x000000, 0x52), clip.Rect{Max: gtx.Constraints.Max}.Op())
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		return fixedWidth(gtx, unit.Dp(520), func(gtx layout.Context) layout.Dimensions {
			return a.borderedSurface(gtx, fluent.surface, unit.Dp(8), fluent.border, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(12))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, "图片已生成, 是否另存到指定位置?", unit.Sp(18), fluent.text, font.SemiBold)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, "默认目录已保存一份。需要放到项目、相册或其他目录时, 可以现在填写目标路径再保存副本。", unit.Sp(13), fluent.textMuted, font.Normal)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.field(gtx, "保存到", &a.savePromptPathInput, "输入完整文件路径或目录", unit.Dp(48))
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							style := material.CheckBox(a.th, &a.savePromptNeverAsk, "以后不再提示")
							style.Color = fluent.text
							style.IconColor = fluent.accent
							return style.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return a.button(gtx, &a.savePromptSkipButton, "稍后", fluent.surface2, fluent.text)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return a.button(gtx, &a.savePromptSaveButton, "保存副本", fluent.accent, fluent.white)
								}),
							)
						}),
					)
				})
			})
		})
	})
}
