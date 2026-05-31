package ui

import "image/color"

type fluentColors struct {
	accent     color.NRGBA
	accent2    color.NRGBA
	accentSoft color.NRGBA
	bg         color.NRGBA
	bg2        color.NRGBA
	panel      color.NRGBA
	panel2     color.NRGBA
	surface    color.NRGBA
	surface2   color.NRGBA
	sidebar    color.NRGBA
	inspector  color.NRGBA
	toolbar    color.NRGBA
	border     color.NRGBA
	border2    color.NRGBA
	text       color.NRGBA
	textMuted  color.NRGBA
	textDim    color.NRGBA
	canvasBg   color.NRGBA
	canvasTile color.NRGBA
	success    color.NRGBA
	danger     color.NRGBA
	white      color.NRGBA
}

var fluent = fluentColors{
	accent:     rgb(0x005fb8),
	accent2:    rgb(0x0a6fcb),
	accentSoft: rgba(0x005fb8, 0x1f),
	bg:         rgb(0xf3f3f3),
	bg2:        rgb(0xe9e9e9),
	panel:      rgb(0xfbfbfb),
	panel2:     rgb(0xf6f6f6),
	surface:    rgb(0xffffff),
	surface2:   rgb(0xf2f2f2),
	sidebar:    rgb(0xfbfbfb),
	inspector:  rgb(0xf8f8f8),
	toolbar:    rgb(0xf7f7f7),
	border:     rgba(0x000000, 0x14),
	border2:    rgba(0x000000, 0x24),
	text:       rgb(0x1f1f1f),
	textMuted:  rgb(0x5f6368),
	textDim:    rgb(0x8a8f98),
	canvasBg:   rgb(0xeeeeee),
	canvasTile: rgb(0xdedede),
	success:    rgb(0x0f7b0f),
	danger:     rgb(0xc42b1c),
	white:      rgb(0xffffff),
}

func rgb(v uint32) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}

func rgba(v uint32, alpha uint8) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: alpha}
}
