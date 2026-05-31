package ui

import (
	"fmt"
	"strings"

	"image-studio/gio-client/internal/kernel"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/yuanhua/image-gptcodex/pkg/client"
)

func (a *App) layoutControls(gtx layout.Context) layout.Dimensions {
	return a.borderedSurface(gtx, fluent.sidebar, unit.Dp(0), fluent.border, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = gtx.Constraints.Max
		return layout.Inset{Top: 12, Bottom: 12, Left: 16, Right: 16}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return a.field(gtx, "提示词", &a.promptInput, "输入生成或编辑要求", unit.Dp(132))
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
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(5))}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.sectionEyebrow(gtx, "创作参数")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.label(gtx, summary, unit.Sp(12), fluent.textMuted, font.Normal)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.field(gtx, "参考图路径", &a.sourcePathsInput, "图生图时填写本地图片路径, 多张可换行", unit.Dp(76))
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
	})
}

func (a *App) layoutAdvancedCard(gtx layout.Context) layout.Dimensions {
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(10))}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.sectionEyebrow(gtx, "高级参数")
			}),
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
	})
}

func (a *App) layoutActions(gtx layout.Context) layout.Dimensions {
	snap := a.readSnapshot()
	return a.card(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(8))}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				txt := "生成"
				bg := fluent.accent
				if snap.Running {
					txt = "运行中"
					bg = fluent.textMuted
				}
				return a.button(gtx, &a.runButton, txt, bg, fluent.white)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return fixedWidth(gtx, unit.Dp(88), func(gtx layout.Context) layout.Dimensions {
					return a.button(gtx, &a.cancelButton, "取消", fluent.surface2, fluent.text)
				})
			}),
		)
	})
}
