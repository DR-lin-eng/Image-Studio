package ui

import (
	"fmt"
	"strings"

	"image-studio/gio-client/internal/kernel"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/yuanhua/image-gptcodex/pkg/client"
)

func (a *App) layoutControls(gtx layout.Context) layout.Dimensions {
	for a.composeToggleButton.Clicked(gtx) {
		a.composeOpen = !a.composeOpen
	}
	for a.advancedToggleButton.Clicked(gtx) {
		a.advancedOpen = !a.advancedOpen
	}
	for a.manageUpstreamButton.Clicked(gtx) {
		a.advancedOpen = true
	}

	return a.borderedSurface(gtx, fluent.sidebar, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 12, Bottom: 12, Left: 12, Right: 12}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return a.controlsList.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
					layout.Rigid(a.layoutWorkbenchCard),
					layout.Rigid(a.layoutModeCard),
					layout.Rigid(a.layoutPromptCard),
					layout.Rigid(a.layoutComposeCard),
					layout.Rigid(a.layoutAdvancedCard),
					layout.Rigid(a.layoutActions),
				)
			})
		})
	})
}

func (a *App) layoutWorkbenchCard(gtx layout.Context) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "图像工作台", unit.Sp(18), fluent.text, font.SemiBold)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "保持界面简洁，把注意力留给 prompt、参考图和结果。", unit.Sp(12), fluent.textMuted, font.Normal)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.badge(gtx, a.modeLabel(), fluent.accentSoft, fluent.accent)
			}),
		)
	})
}

func (a *App) layoutModeCard(gtx layout.Context) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.sectionEyebrow(gtx, "模式")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.segmented(gtx, modeChoices, a.mode, a.modeButtons, func(value string) { a.mode = value })
			}),
		)
	})
}

func (a *App) layoutPromptCard(gtx layout.Context) layout.Dimensions {
	promptLen := len([]rune(strings.TrimSpace(a.promptInput.Text())))
	footer := "模板 / 历史在 WebView 版更完整，Gio 版先对齐工作流主干。"
	if a.mode == string(client.ModeEdit) {
		count := len(kernel.ParseSourcePaths(a.sourcePathsInput.Text()))
		footer = fmt.Sprintf("图生图模式 · 当前 %d 张参考图", count)
	}

	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, "PROMPT 提示词", unit.Sp(11), fluent.textMuted, font.Bold)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, fmt.Sprintf("%d", promptLen), unit.Sp(11), fluent.textDim, font.Medium)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return fixedHeight(gtx, unit.Dp(136), func(gtx layout.Context) layout.Dimensions {
					return a.borderedSurface(gtx, fluent.surface, unit.Dp(4), fluent.border2, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 10, Bottom: 10, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.editorText(gtx, &a.promptInput, "输入生成或编辑要求", unit.Sp(13))
						})
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.label(gtx, footer, unit.Sp(10), fluent.textDim, font.Normal)
			}),
		)
	})
}

func (a *App) layoutComposeCard(gtx layout.Context) layout.Dimensions {
	sourceLabel := "文生图"
	if a.mode == string(client.ModeEdit) {
		count := len(kernel.ParseSourcePaths(a.sourcePathsInput.Text()))
		if count > 0 {
			sourceLabel = fmt.Sprintf("%d 张源图", count)
		} else {
			sourceLabel = "未添加源图"
		}
	}
	summary := strings.Join(compactNonEmpty([]string{
		sizeChoiceLabel(a.size),
		qualityChoiceLabel(a.quality),
		strings.ToUpper(a.format),
		sourceLabel,
	}), " · ")

	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.layoutDisclosureHeader(gtx, &a.composeToggleButton, "创作参数", summary, a.composeOpen)
			}),
		}
		if a.composeOpen {
			children = append(children,
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "参考图路径", &a.sourcePathsInput, "图生图时填写本地图片路径，多张可换行", unit.Dp(72))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedWithTitle(gtx, "API 形态", apiChoices, a.api, a.apiButtons, func(value string) { a.api = value })
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedGridWithTitle(gtx, "尺寸", sizeChoices, a.size, a.sizeButtons, 2, func(value string) { a.size = value })
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedWithTitle(gtx, "质量", qualityChoices, a.quality, a.qualityButtons, func(value string) { a.quality = value })
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedWithTitle(gtx, "格式", formatChoices, a.format, a.formatButtons, func(value string) { a.format = value })
				}),
			)
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (a *App) layoutAdvancedCard(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	profileName := strings.TrimSpace(activeProfileName(snap.Profiles, snap.ActiveProfileID))
	if profileName == "" {
		profileName = "未命名配置"
	}
	summary := strings.Join(compactNonEmpty([]string{
		profileName,
		apiChoiceLabel(a.api),
		policyChoiceLabel(a.policy),
		proxyChoiceLabel(a.proxy),
	}), " · ")

	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.layoutDisclosureHeader(gtx, &a.advancedToggleButton, "高级参数", summary, a.advancedOpen)
			}),
		}
		if a.advancedOpen {
			children = append(children,
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "BASE_URL", &a.baseURLInput, "https://example.com", unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "API Key", &a.apiKeyInput, "sk-...", unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "文本模型", &a.textModelInput, client.TextModel, unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "图像模型", &a.imageModelInput, client.ImageModel, unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedWithTitle(gtx, "请求字段", policyChoices, a.policy, a.policyButtons, func(value string) { a.policy = value })
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.segmentedWithTitle(gtx, "代理", proxyChoices, a.proxy, a.proxyButtons, func(value string) { a.proxy = value })
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "自定义代理 URL", &a.proxyURLInput, "http://127.0.0.1:7890", unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "输出目录", &a.outputDirInput, "生成图片保存目录", unit.Dp(44))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return a.field(gtx, "Seed", &a.seedInput, "0", unit.Dp(44))
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return a.field(gtx, "Partial", &a.partialImagesInput, "1", unit.Dp(44))
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.field(gtx, "负向提示词", &a.negativePromptInput, "兼容模式可发给部分上游", unit.Dp(64))
				}),
			)
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (a *App) layoutActions(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	ready := strings.TrimSpace(a.apiKeyInput.Text()) != "" && strings.TrimSpace(a.baseURLInput.Text()) != ""

	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		children := make([]layout.FlexChild, 0, 4)
		if !ready {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.borderedSurface(gtx, fluent.accentSoft, unit.Dp(6), rgba(0x005fb8, 0x28), func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, "还没有可用上游配置", unit.Sp(12), fluent.accent, font.SemiBold)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.label(gtx, "先配置 BASE_URL 和 API Key，才能开始测试连接或生成。", unit.Sp(11), fluent.textMuted, font.Normal)
							}),
						)
					})
				})
			}))
			children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout))
		}

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flexChildren := []layout.FlexChild{
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if !ready {
						return a.button(gtx, &a.manageUpstreamButton, "配置上游", fluent.accent, fluent.white)
					}
					label := "生成"
					bg := fluent.accent
					if a.mode == string(client.ModeEdit) {
						label = "编辑"
					}
					if snap.Running {
						label = "运行中"
						bg = fluent.textMuted
					}
					return a.button(gtx, &a.runButton, label, bg, fluent.white)
				}),
			}
			if snap.Running {
				flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return fixedWidth(gtx, unit.Dp(92), func(gtx layout.Context) layout.Dimensions {
						return a.button(gtx, &a.cancelButton, "取消", fluent.surface2, fluent.text)
					})
				}))
			} else if ready {
				flexChildren = append(flexChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return fixedWidth(gtx, unit.Dp(92), func(gtx layout.Context) layout.Dimensions {
						return a.button(gtx, &a.manageUpstreamButton, "上游", fluent.surface2, fluent.text)
					})
				}))
			}
			return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx, flexChildren...)
		}))
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (a *App) layoutDisclosureHeader(gtx layout.Context, btn *widget.Clickable, title string, summary string, open bool) layout.Dimensions {
	stateText := "展开"
	if open {
		stateText = "收起"
	}
	return a.surfaceButton(
		gtx,
		btn,
		chooseColor(open, fluent.surface2, fluent.surface),
		fluent.surface2,
		fluent.border,
		unit.Dp(6),
		layout.Inset{Top: 10, Bottom: 10, Left: 10, Right: 10},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(4))}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, title, unit.Sp(11), fluent.text, font.Bold)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return a.label(gtx, summary, unit.Sp(11), fluent.textMuted, font.Normal)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.label(gtx, stateText, unit.Sp(11), fluent.textDim, font.Medium)
				}),
			)
		},
	)
}

func (a *App) editorText(gtx layout.Context, editor *widget.Editor, hint string, size unit.Sp) layout.Dimensions {
	style := material.Editor(a.th, editor, hint)
	style.Color = fluent.text
	style.HintColor = fluent.textDim
	style.SelectionColor = rgba(0x005fb8, 0x3d)
	style.TextSize = size
	return style.Layout(gtx)
}

func apiChoiceLabel(value string) string {
	return choiceLabel(apiChoices, value)
}

func policyChoiceLabel(value string) string {
	return choiceLabel(policyChoices, value)
}

func proxyChoiceLabel(value string) string {
	return choiceLabel(proxyChoices, value)
}
