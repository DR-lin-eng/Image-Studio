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
	sourcePaths := kernel.ParseSourcePaths(a.sourcePathsInput.Text())
	showSourceStrip := a.mode == "edit" || len(sourcePaths) > 0

	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(48), func(gtx layout.Context) layout.Dimensions {
				return a.canvasToolbar(gtx, snap)
			})
		}),
	}
	if showSourceStrip {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(64), func(gtx layout.Context) layout.Dimensions {
				return a.sourceStrip(gtx, sourcePaths)
			})
		}))
	}
	children = append(children,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.resultSurface(gtx, snap)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.canvasStatusBar(gtx, snap)
		}),
	)

	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (a *App) canvasToolbar(gtx layout.Context, snap snapshot) layout.Dimensions {
	sizeLabel := sizeChoiceLabel(a.size)
	qualityLabel := qualityChoiceLabel(a.quality)
	if snap.Result.HasItem {
		if strings.TrimSpace(snap.Result.Item.Size) != "" {
			sizeLabel = snap.Result.Item.Size
		}
		if strings.TrimSpace(snap.Result.Item.Quality) != "" {
			qualityLabel = snap.Result.Item.Quality
		}
	}
	trailing := strings.TrimSpace(snap.Result.SourceEvent)
	if trailing == "" {
		trailing = snap.Status
	}

	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, "画布", true)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, a.modeLabel(), false)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, sizeLabel, false)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, qualityLabel, false)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.toolPill(gtx, strings.ToUpper(a.format), false)
				}),
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.badge(gtx, trailing, fluent.surface, fluent.textMuted)
				}),
			)
		})
	})
}

func (a *App) sourceStrip(gtx layout.Context, sourcePaths []string) layout.Dimensions {
	label := "参考图 0 张"
	if len(sourcePaths) > 0 {
		label = "参考图 " + strconv.Itoa(len(sourcePaths)) + " 张"
	}

	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 8, Bottom: 8, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			children := []layout.FlexChild{
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, label, unit.Sp(11), fluent.textMuted, font.Medium)
				}),
			}
			if len(sourcePaths) == 0 {
				children = append(children,
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "图生图时会显示解析后的本地参考图。", unit.Sp(11), fluent.textDim, font.Normal)
					}),
				)
			} else {
				for _, path := range sourcePaths {
					path := path
					children = append(children,
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							img, _ := a.imageForPath(path)
							return a.imageThumb(gtx, img, unit.Dp(48), unit.Dp(48), unit.Dp(4))
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return fixedWidth(gtx, unit.Dp(132), func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, strings.TrimSpace(path), unit.Sp(10), fluent.textDim, font.Normal)
							})
						}),
					)
				}
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
		})
	})
}

func (a *App) resultSurface(gtx layout.Context, snap snapshot) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return a.surface(gtx, fluent.canvasBg, unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.UniformInset(unit.Dp(18)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return a.borderedSurface(gtx, fluent.surface, unit.Dp(8), fluent.border, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = gtx.Constraints.Max
				if snap.Result.Image == nil {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, "等待生成结果", unit.Sp(18), fluent.text, font.Medium)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, "当前布局已经对齐 Windows 三栏结构，生成完成后会在这里常驻显示结果。", unit.Sp(13), fluent.textMuted, font.Normal)
							}),
						)
					})
				}
				if snap.Result.Rev != a.imageOpRev {
					a.imageOp = paint.NewImageOp(snap.Result.Image)
					a.imageOpRev = snap.Result.Rev
				}
				return layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						img := widget.Image{
							Src:      a.imageOp,
							Fit:      widget.Contain,
							Position: layout.Center,
						}
						return img.Layout(gtx)
					})
				})
			})
		})
	})
}

func (a *App) canvasStatusBar(gtx layout.Context, snap snapshot) layout.Dimensions {
	lastLog := ""
	if len(snap.Logs) > 0 {
		lastLog = snap.Logs[len(snap.Logs)-1]
	}

	return a.borderedSurface(gtx, fluent.panel2, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: 9, Bottom: 9, Left: 14, Right: 14}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if snap.Running {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, snap.Status, unit.Sp(11), fluent.accent, font.Medium)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if strings.TrimSpace(lastLog) == "" {
							return layout.Dimensions{}
						}
						return a.label(gtx, lastLog, unit.Sp(11), fluent.textMuted, font.Normal)
					}),
				)
			}

			if snap.Result.HasItem {
				headline := "生成结果"
				if snap.Result.Item.Mode == "edit" {
					headline = "编辑结果"
				}
				meta := historyMetaText(snap.Result.Item)
				revised := strings.TrimSpace(snap.Result.RevisedPrompt)
				if revised == "" {
					revised = "暂无修订提示词"
				}
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, headline, unit.Sp(11), fluent.accent, font.Medium)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, meta, unit.Sp(11), fluent.textMuted, font.Normal)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, revised, unit.Sp(11), fluent.textDim, font.Normal)
					}),
				)
			}

			return a.label(gtx, "准备就绪", unit.Sp(11), fluent.textMuted, font.Normal)
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
