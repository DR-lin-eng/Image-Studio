package ui

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func (a *App) field(gtx layout.Context, title string, editor *widget.Editor, hint string, height unit.Dp) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, title, unit.Sp(12), fluent.textMuted, font.Medium)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return fixedHeight(gtx, height, func(gtx layout.Context) layout.Dimensions {
				return a.borderedSurface(gtx, fluent.surface, unit.Dp(4), fluent.border2, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: 9, Bottom: 9, Left: 10, Right: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						style := material.Editor(a.th, editor, hint)
						style.Color = fluent.text
						style.HintColor = fluent.textDim
						style.SelectionColor = rgba(0x005fb8, 0x3d)
						style.TextSize = unit.Sp(13)
						return style.Layout(gtx)
					})
				})
			})
		}),
	)
}

func (a *App) segmentedWithTitle(gtx layout.Context, title string, options []choice, selected string, buttons []widget.Clickable, set func(string)) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, title, unit.Sp(12), fluent.textMuted, font.Medium)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.segmented(gtx, options, selected, buttons, set)
		}),
	)
}

func (a *App) segmentedGridWithTitle(gtx layout.Context, title string, options []choice, selected string, buttons []widget.Clickable, columns int, set func(string)) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.label(gtx, title, unit.Sp(12), fluent.textMuted, font.Medium)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.segmentedGrid(gtx, options, selected, buttons, columns, set)
		}),
	)
}

func (a *App) segmented(gtx layout.Context, options []choice, selected string, buttons []widget.Clickable, set func(string)) layout.Dimensions {
	children := make([]layout.FlexChild, 0, len(options))
	for i := range options {
		i := i
		for buttons[i].Clicked(gtx) {
			set(options[i].Value)
		}
		children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			bg := fluent.surface
			fg := fluent.textMuted
			if options[i].Value == selected {
				bg = fluent.accentSoft
				fg = fluent.accent
			}
			return a.button(gtx, &buttons[i], options[i].Label, bg, fg)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx, children...)
}

func (a *App) segmentedGrid(gtx layout.Context, options []choice, selected string, buttons []widget.Clickable, columns int, set func(string)) layout.Dimensions {
	if columns <= 0 {
		columns = 2
	}
	rows := (len(options) + columns - 1) / columns
	children := make([]layout.FlexChild, 0, rows)
	for row := 0; row < rows; row++ {
		row := row
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			cellChildren := make([]layout.FlexChild, 0, columns)
			for col := 0; col < columns; col++ {
				idx := row*columns + col
				if idx >= len(options) {
					cellChildren = append(cellChildren, layout.Flexed(1, layout.Spacer{}.Layout))
					continue
				}
				for buttons[idx].Clicked(gtx) {
					set(options[idx].Value)
				}
				cellChildren = append(cellChildren, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					bg := fluent.surface
					fg := fluent.textMuted
					if options[idx].Value == selected {
						bg = fluent.accentSoft
						fg = fluent.accent
					}
					return a.button(gtx, &buttons[idx], options[idx].Label, bg, fg)
				}))
			}
			return layout.Inset{Bottom: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Gap: gtx.Dp(unit.Dp(6))}.Layout(gtx, cellChildren...)
			})
		}))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
