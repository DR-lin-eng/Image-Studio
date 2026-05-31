# Image Studio Gio Client

`gio-client/` is an independent high-performance desktop test client for Windows and Linux. It uses Gio for the native GUI and reuses the existing Go request kernel from `go-cli/pkg/client`.

It does not embed the React frontend, Wails, WebView2, or WebKitGTK. The current Wails desktop app remains in `image-studio/` and continues to build through the existing WebView2/WebKit path.

## Architecture

```text
gio-client/
├── cmd/image-studio-gio/      # Gio app entrypoint
├── internal/ui/               # Gio immediate-mode frontend
└── internal/kernel/           # adapter around go-cli/pkg/client
```

The UI keeps the current desktop control-panel / canvas / log-rail structure, but its frontend architecture is native Gio instead of React/CSS. Request payload construction, retry behavior, SSE parsing, Images API support, proxy handling, and default model constants remain owned by `go-cli/pkg/client`.

The GUI entrypoint is built only for Windows and Linux. Other platforms compile a small unsupported stub so accidental local launches do not imply macOS support for this test client.

## Local Build

```bash
cd gio-client
go test ./...
go build -o /tmp/image-studio-gio ./cmd/image-studio-gio
```

Linux requires Gio's native build libraries:

```bash
sudo apt-get update
sudo apt-get install -y \
  pkg-config \
  libegl1-mesa-dev \
  libvulkan-dev \
  libwayland-dev \
  libx11-dev \
  libx11-xcb-dev \
  libxcursor-dev \
  libxfixes-dev \
  libxkbcommon-dev \
  libxkbcommon-x11-dev
```

Generated images default to `~/Pictures/Image Studio Gio`, unless the output directory field is changed.
