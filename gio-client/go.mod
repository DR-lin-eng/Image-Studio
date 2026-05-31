module image-studio/gio-client

go 1.25.5

toolchain go1.26.3

require (
	gioui.org v0.10.0
	github.com/yuanhua/image-gptcodex v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.41.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	github.com/go-text/typesetting v0.3.4 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/yuanhua/image-gptcodex => ../go-cli
