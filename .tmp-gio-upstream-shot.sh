#!/usr/bin/env bash
set -euo pipefail

export PATH=/usr/local/go/bin:$PATH
export DEBIAN_FRONTEND=noninteractive

apt-get update >/tmp/apt-update.log
apt-get install -y --no-install-recommends \
  xvfb \
  x11-utils \
  xdotool \
  imagemagick \
  pkg-config \
  gcc \
  g++ \
  libc6-dev \
  libx11-dev \
  libx11-xcb-dev \
  libxkbcommon-x11-dev \
  libxcursor-dev \
  libxfixes-dev \
  libxi-dev \
  libxinerama-dev \
  libxrandr-dev \
  libxxf86vm-dev \
  libwayland-dev \
  libegl1-mesa-dev \
  libgles2-mesa-dev \
  libvulkan-dev >/tmp/apt-install.log

mkdir -p /tmp/xdg /tmp/home/.config/image-studio/compat "/tmp/home/Pictures/Image Studio"
chmod 700 /tmp/xdg

convert -size 768x768 xc:'#90b4ef' "/tmp/home/Pictures/Image Studio/demo-a.png"

cat > /tmp/home/.config/image-studio/compat/state.json <<'EOF'
{
  "schemaVersion": 1,
  "client": "gio",
  "updatedAt": 1760000000000,
  "settings": {
    "background": "auto",
    "outputCompression": 100,
    "inputFidelity": "auto",
    "moderation": "low",
    "partialImages": 1
  },
  "profiles": [
    {
      "id": "demo-profile",
      "name": "配置1",
      "apiMode": "responses",
      "requestPolicy": "openai",
      "baseURL": "https://api.example.com",
      "textModelID": "gpt-5.5",
      "imageModelID": "gpt-image-2",
      "concurrencyLimit": 0,
      "createdAt": 1760000000000,
      "lastUsedAt": 1760000000000
    }
  ],
  "activeProfileId": "demo-profile",
  "history": [
    {
      "id": "hist-1",
      "prompt": "一只橘猫坐在雨夜窗边，电影级侧逆光，50mm，浅景深，写实摄影",
      "revisedPrompt": "一只橘猫坐在雨夜窗边，电影感侧逆光，50mm，浅景深，写实摄影，细节丰富",
      "mode": "generate",
      "size": "1536x1024",
      "quality": "medium",
      "outputFormat": "png",
      "createdAt": 1760000000000,
      "styleTag": "cyberpunk",
      "batchIndex": 0,
      "elapsedSec": 18.2,
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-a.png"
    }
  ]
}
EOF

export HOME=/tmp/home
export GOPATH=/src/.tmp-gio-docker-gopath
export GOMODCACHE=/src/.tmp-gio-docker-gomodcache
export GOCACHE=/src/.tmp-gio-docker-gocache
export CGO_ENABLED=1
export LIBGL_ALWAYS_SOFTWARE=1
export DISPLAY=:99
export XDG_RUNTIME_DIR=/tmp/xdg
export XDG_SESSION_TYPE=x11
export WAYLAND_DISPLAY=

Xvfb :99 -screen 0 1440x1024x24 >/tmp/xvfb.log 2>&1 &
xvfb_pid=$!
trap 'kill "$xvfb_pid" >/dev/null 2>&1 || true' EXIT

go run ./gio-client/cmd/image-studio-gio >/tmp/gio.log 2>&1 &
app_pid=$!
sleep 20

# Click the right-rail "上游配置" button.
xdotool mousemove --sync 1222 155 click 1
sleep 2

xwininfo -display :99 -root -tree > /src/.tmp-gio-upstream-xwininfo.txt || true
window_id=$(awk '/"Image Studio Gio"/ {print $1; exit}' /src/.tmp-gio-upstream-xwininfo.txt)
if [ -n "$window_id" ]; then
  import -display :99 -window "$window_id" /src/.tmp-gio-upstream-shot.png
fi

kill "$app_pid" >/dev/null 2>&1 || true
wait "$app_pid" || true
