package ui

import "github.com/yuanhua/image-gptcodex/pkg/client"

type choice struct {
	Label string
	Value string
}

var (
	modeChoices = []choice{
		{"文生图", string(client.ModeGenerate)},
		{"图生图", string(client.ModeEdit)},
	}
	apiChoices = []choice{
		{"Responses", string(client.APIModeResponses)},
		{"Images", string(client.APIModeImages)},
	}
	sizeChoices = []choice{
		{"自适应 auto", "auto"},
		{"方形 1024x1024", "1024x1024"},
		{"横版 1536x1024", "1536x1024"},
		{"竖版 1024x1536", "1024x1536"},
		{"2K 方形 2048x2048", "2048x2048"},
		{"2K 横版 2048x1360", "2048x1360"},
		{"2K 竖版 1360x2048", "1360x2048"},
		{"2K 横版 2048x1152", "2048x1152"},
		{"2K 竖版 1152x2048", "1152x2048"},
		{"4K 方形 2880x2880", "2880x2880"},
		{"4K 横版 3456x2304", "3456x2304"},
		{"4K 竖版 2304x3456", "2304x3456"},
		{"4K 横版 3840x2160", "3840x2160"},
		{"4K 竖版 2160x3840", "2160x3840"},
	}
	qualityChoices = []choice{
		{"自适应 auto", "auto"},
		{"高质量 high", "high"},
		{"中等 medium", "medium"},
		{"快速草稿 low", "low"},
	}
	formatChoices = []choice{
		{"PNG", "png"},
		{"JPEG", "jpeg"},
		{"WebP", "webp"},
	}
	policyChoices = []choice{
		{"OpenAI 标准", string(client.RequestPolicyOpenAI)},
		{"兼容中转扩展", string(client.RequestPolicyCompat)},
	}
	proxyChoices = []choice{
		{"系统配置", client.ProxyModeSystem},
		{"不使用", client.ProxyModeNone},
		{"自定义", client.ProxyModeCustom},
	}
)

func (a *App) modeLabel() string {
	if a.mode == string(client.ModeEdit) {
		return "图生图"
	}
	return "文生图"
}

func sizeChoiceLabel(value string) string {
	return choiceLabel(sizeChoices, value)
}

func qualityChoiceLabel(value string) string {
	return choiceLabel(qualityChoices, value)
}

func choiceLabel(choices []choice, value string) string {
	for _, item := range choices {
		if item.Value == value {
			return item.Label
		}
	}
	return value
}
