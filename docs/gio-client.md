# Gio 高性能测试客户端

Gio 客户端是 Windows / Linux 的独立测试版本，目录为 `gio-client/`。它的目标是验证不经过 WebView2/WebKitGTK 的原生 Gio 渲染路径，同时保持当前请求内核不变。

## 边界

- 不改 `image-studio/main.go`，不影响 Wails / WebView2 桌面端。
- 不改 `image-studio/frontend/` 的 React 视觉实现。
- 复用 `go-cli/pkg/client` 的 Responses API、Images API、SSE、retry、proxy、模型默认值和请求字段策略。
- Gio 前端为新的 immediate-mode 架构，UI 结构沿用桌面端的控制面板、画布和运行日志布局。
- GUI 入口仅面向 Windows / Linux；其他平台只编译 unsupported stub，避免误判为 macOS 支持。

## 构建

```bash
cd gio-client
go test ./...
go build -o /tmp/image-studio-gio ./cmd/image-studio-gio
```

Linux 需要 Gio 原生依赖:

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

Release workflow 会额外产出 `image-studio-gio-*` artifacts。它们与现有 `image-studio-*` Wails artifacts 分开上传。
