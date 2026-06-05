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

convert -size 768x768 gradient:'#49b2e8-#4b27d2' \
  -fill 'rgba(255,255,255,0.18)' -draw 'circle 560,180 560,20' \
  -fill 'rgba(12,20,90,0.26)' -draw 'circle 180,560 180,360' \
  -fill 'rgba(26,10,68,0.58)' -draw 'roundrectangle 120,470 640,590 30,30' \
  "/tmp/home/Pictures/Image Studio/demo-a.png"
convert -size 768x768 xc:'#d7a3a3' "/tmp/home/Pictures/Image Studio/demo-b.png"

cat > /tmp/home/.config/image-studio/compat/state.json <<'EOF'
{
  "schemaVersion": 1,
  "client": "gio",
  "updatedAt": 1760000000000,
  "settings": {
    "background": "auto",
    "outputCompression": 100,
    "inputFidelity": "auto",
    "imageStyle": "default",
    "moderation": "low",
    "partialImages": 1
  },
  "profiles": [
    {
      "id": "demo-profile",
      "name": "Preview Responses",
      "apiMode": "responses",
      "requestPolicy": "openai",
      "baseURL": "https://code1.linzefeng.top",
      "textModelID": "gpt-4.1-mini",
      "imageModelID": "gpt-image-1",
      "concurrencyLimit": 1,
      "createdAt": 1760000000000,
      "lastUsedAt": 1760000000000
    }
  ],
  "activeProfileId": "demo-profile",
  "history": [
    {
      "id": "hist-1",
      "prompt": "赛博雨夜角色海报，湿地街道反光，红青霓虹边缘光，35mm，电影感，超细节",
      "revisedPrompt": "同一提示词批量结果 1，强化雨滴高光、轮廓光与街面反射",
      "mode": "edit",
      "size": "2880x2880",
      "quality": "high",
      "outputFormat": "png",
      "background": "opaque",
      "outputCompression": 100,
      "inputFidelity": "auto",
      "imageStyle": "default",
      "moderation": "low",
      "createdAt": 1760000000000,
      "styleTag": "电影海报",
      "batchIndex": 0,
      "seed": 3200,
      "elapsedSec": 7,
      "thumbPath": "/tmp/home/Pictures/Image Studio/demo-a.png",
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

# Expand the compose accordion in the left rail.
xdotool mousemove --sync 180 520 click 1
sleep 2

xwininfo -display :99 -root -tree > /src/.tmp-gio-compose-xwininfo.txt || true
window_id=$(awk '/"Image Studio Gio"/ {print $1; exit}' /src/.tmp-gio-compose-xwininfo.txt)
if [ -n "$window_id" ]; then
  import -display :99 -window "$window_id" /src/.tmp-gio-compose-shot.png
fi

kill "$app_pid" >/dev/null 2>&1 || true
wait "$app_pid" || true
