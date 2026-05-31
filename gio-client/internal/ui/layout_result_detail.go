package ui

import (
	"fmt"
	"image"
	"io"
	"path/filepath"
	"strings"

	sharedCompat "image-studio/shared/compat"

	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func (a *App) layoutResultDetailModal(gtx layout.Context) layout.Dimensions {
	for a.closeResultDetailButton.Clicked(gtx) {
		a.closeResultDetail()
	}
	snap := a.readSnapshot()
	item := snap.ActiveResultDetail
	if item.ID == "" && strings.TrimSpace(item.SavedPath) == "" {
		return layout.Dimensions{}
	}
	for a.resultDetailUsePromptButton.Clicked(gtx) {
		a.useResultPrompt(item.Prompt)
	}
	for a.resultDetailUseRevisedButton.Clicked(gtx) {
		a.useResultPrompt(item.RevisedPrompt)
	}
	for a.resultDetailCopyPromptButton.Clicked(gtx) {
		copyResultDetailText(gtx, item.Prompt)
		a.appendLog("已复制原始提示词")
	}
	for a.resultDetailCopyRevisedButton.Clicked(gtx) {
		copyResultDetailText(gtx, item.RevisedPrompt)
		a.appendLog("已复制优化后提示词")
	}
	for a.resultDetailOpenPathButton.Clicked(gtx) {
		path := strings.TrimSpace(item.SavedPath)
		if path == "" {
			continue
		}
		if err := openPath(filepath.Dir(path)); err != nil {
			a.appendLog("打开文件夹失败: " + err.Error())
		}
	}
	for a.resultDetailCopyPathButton.Clicked(gtx) {
		copyResultDetailText(gtx, item.SavedPath)
		a.appendLog("已复制文件路径")
	}
	paint.FillShape(gtx.Ops, rgba(0x000000, 0x52), clip.Rect{Max: gtx.Constraints.Max}.Op())
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		return fixedWidth(gtx, unit.Dp(760), func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, unit.Dp(620), func(gtx layout.Context) layout.Dimensions {
				return a.borderedSurface(gtx, fluent.surface, unit.Dp(8), fluent.border, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(12))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return a.label(gtx, "结果详情", unit.Sp(18), fluent.text, font.SemiBold)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return a.label(gtx, detailHeadline(item), unit.Sp(11), fluent.textMuted, font.Normal)
											}),
										)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return fixedWidth(gtx, unit.Dp(104), func(gtx layout.Context) layout.Dimensions {
											return a.compactIconTextButton(gtx, &a.closeResultDetailButton, uiIconClose, "关闭", false)
										})
									}),
								)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(16))}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return fixedWidth(gtx, unit.Dp(280), func(gtx layout.Context) layout.Dimensions {
											return a.layoutResultDetailPreview(gtx, item)
										})
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return a.layoutResultDetailSections(gtx, item)
									}),
								)
							}),
						)
					})
				})
			})
		})
	})
}

func (a *App) layoutResultDetailPreview(gtx layout.Context, item sharedCompat.HistoryItem) layout.Dimensions {
	img, _ := a.imageForHistoryItem(item)
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.imageThumb(gtx, img, unit.Dp(244), unit.Dp(244), unit.Dp(6))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if strings.TrimSpace(item.SavedPath) == "" {
					return layout.Dimensions{}
				}
				return a.label(gtx, historyPathText(item.SavedPath), unit.Sp(10), fluent.textDim, font.Normal)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if strings.TrimSpace(item.SavedPath) == "" {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, &a.resultDetailOpenPathButton, uiIconFolder, "打开文件夹", false)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, &a.resultDetailCopyPathButton, uiIconCopy, "复制路径", false)
					}),
				)
			}),
		)
	})
}

func (a *App) layoutResultDetailSections(gtx layout.Context, item sharedCompat.HistoryItem) layout.Dimensions {
	sections := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions { return a.layoutResultDetailMeta(gtx, item) },
		func(gtx layout.Context) layout.Dimensions {
			return a.layoutResultDetailTextSection(gtx, "原始提示词", item.Prompt)
		},
	}
	if strings.TrimSpace(item.RevisedPrompt) != "" {
		sections = append(sections, func(gtx layout.Context) layout.Dimensions {
			return a.layoutResultDetailTextSection(gtx, "优化后提示词", item.RevisedPrompt)
		})
	}
	if strings.TrimSpace(item.NegativePrompt) != "" {
		sections = append(sections, func(gtx layout.Context) layout.Dimensions {
			return a.layoutResultDetailTextSection(gtx, "负向提示词", item.NegativePrompt)
		})
	}
	if strings.TrimSpace(item.SavedPath) != "" {
		sections = append(sections, func(gtx layout.Context) layout.Dimensions {
			return a.layoutResultDetailFileSection(gtx, item)
		})
	}
	return a.settingsList.Layout(gtx, len(sections), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, sections[index])
	})
}

func (a *App) layoutResultDetailMeta(gtx layout.Context, item sharedCompat.HistoryItem) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		rows := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.sectionEyebrow(gtx, "参数")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "模式", chooseModeLabel(item.Mode))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "尺寸", item.Size)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "质量", item.Quality)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "格式", strings.ToUpper(strings.TrimSpace(item.OutputFormat)))
			}),
		}
		if item.Seed != 0 {
			rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "Seed", detailValue(item.Seed))
			}))
		}
		if strings.TrimSpace(item.StyleTag) != "" {
			rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "风格", "#"+styleChoiceLabel(item.StyleTag))
			}))
		}
		rows = append(rows,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.detailKV(gtx, "创建时间", formatHistoryClock(item.CreatedAt))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if item.ElapsedSec <= 0 {
					return layout.Dimensions{}
				}
				return a.detailKV(gtx, "耗时", detailValue(item.ElapsedSec)+"s")
			}),
		)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
	})
}

func (a *App) layoutResultDetailTextSection(gtx layout.Context, title string, text string) layout.Dimensions {
	actionBtn := &a.resultDetailUsePromptButton
	copyBtn := &a.resultDetailCopyPromptButton
	actionLabel := "用作下次提示词"
	actionAccent := false
	if strings.Contains(title, "优化后") {
		actionBtn = &a.resultDetailUseRevisedButton
		copyBtn = &a.resultDetailCopyRevisedButton
		actionAccent = true
	}
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.sectionEyebrow(gtx, title)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				content := strings.TrimSpace(text)
				if content == "" {
					content = "(空)"
				}
				return a.borderedSurface(gtx, fluent.surface2, unit.Dp(6), fluent.border, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, content, unit.Sp(11), fluent.textMuted, font.Normal)
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if strings.TrimSpace(text) == "" {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, copyBtn, uiIconCopy, "复制", false)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, actionBtn, uiIconRefresh, actionLabel, actionAccent)
					}),
				)
			}),
		)
	})
}

func (a *App) layoutResultDetailFileSection(gtx layout.Context, item sharedCompat.HistoryItem) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.sectionEyebrow(gtx, "文件")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.borderedSurface(gtx, fluent.surface2, unit.Dp(6), fluent.border, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, strings.TrimSpace(item.SavedPath), unit.Sp(11), fluent.textMuted, font.Normal)
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, &a.resultDetailOpenPathButton, uiIconFolder, "打开文件夹", false)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.compactIconTextButton(gtx, &a.resultDetailCopyPathButton, uiIconCopy, "复制路径", false)
					}),
				)
			}),
		)
	})
}

func copyResultDetailText(gtx layout.Context, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	gtx.Execute(clipboard.WriteCmd{Type: "application/text", Data: io.NopCloser(strings.NewReader(text))})
}

func (a *App) detailKV(gtx layout.Context, label string, value string) layout.Dimensions {
	value = strings.TrimSpace(value)
	if value == "" {
		return layout.Dimensions{}
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedWidth(gtx, unit.Dp(68), func(gtx layout.Context) layout.Dimensions {
				return a.label(gtx, label, unit.Sp(10), fluent.textDim, font.Medium)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, value, unit.Sp(11), fluent.text, font.Normal)
		}),
	)
}

func detailHeadline(item sharedCompat.HistoryItem) string {
	return chooseModeLabel(item.Mode) + " · " + historyMetaText(item)
}

func chooseModeLabel(mode string) string {
	if mode == "edit" {
		return "图生图"
	}
	return "文生图"
}

func detailValue[T any](value T) string {
	return strings.TrimSpace(fmt.Sprint(value))
}
