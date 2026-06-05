#!/usr/bin/env bash
set -euo pipefail

export PATH=/usr/local/go/bin:$PATH
export DEBIAN_FRONTEND=noninteractive

apt-get update >/tmp/apt-update.log
apt-get install -y --no-install-recommends \
  xvfb \
  x11-utils \
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
convert -size 768x768 xc:'#d7a3a3' "/tmp/home/Pictures/Image Studio/demo-b.png"
convert -size 768x768 gradient:'#49b2e8-#4b27d2' \
  -fill 'rgba(255,255,255,0.18)' -draw 'circle 560,180 560,20' \
  -fill 'rgba(12,20,90,0.26)' -draw 'circle 180,560 180,360' \
  -fill 'rgba(26,10,68,0.58)' -draw 'roundrectangle 120,470 640,590 30,30' \
  "/tmp/home/Pictures/Image Studio/thumb-a.png"
convert -size 768x768 gradient:'#df5e7f-#6f1fd0' \
  -fill 'rgba(255,255,255,0.18)' -draw 'circle 560,180 560,20' \
  -fill 'rgba(66,12,76,0.24)' -draw 'circle 180,560 180,360' \
  -fill 'rgba(56,8,44,0.58)' -draw 'roundrectangle 120,470 640,590 30,30' \
  "/tmp/home/Pictures/Image Studio/thumb-b.png"
convert -size 768x768 gradient:'#7ab54f-#3473d1' \
  -fill 'rgba(255,255,255,0.16)' -draw 'circle 560,180 560,20' \
  -fill 'rgba(18,55,72,0.22)' -draw 'circle 180,560 180,360' \
  -fill 'rgba(24,40,62,0.58)' -draw 'roundrectangle 120,470 640,590 30,30' \
  "/tmp/home/Pictures/Image Studio/thumb-c.png"

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
      "sourcePaths": [
        "/tmp/home/Pictures/Image Studio/thumb-a.png",
        "/tmp/home/Pictures/Image Studio/thumb-b.png"
      ],
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-a.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-a.png"
    },
    {
      "id": "hist-2",
      "prompt": "赛博雨夜角色海报，湿地街道反光，红青霓虹边缘光，35mm，电影感，超细节",
      "revisedPrompt": "同一提示词批量结果 2，强化雨滴高光、轮廓光与街面反射",
      "mode": "generate",
      "size": "2880x2880",
      "quality": "high",
      "outputFormat": "png",
      "background": "auto",
      "outputCompression": 100,
      "inputFidelity": "high",
      "imageStyle": "natural",
      "moderation": "auto",
      "createdAt": 1759996700000,
      "styleTag": "胶片人像",
      "batchIndex": 1,
      "seed": 3201,
      "elapsedSec": 8,
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-b.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-b.png"
    },
    {
      "id": "hist-3",
      "prompt": "赛博雨夜角色海报，湿地街道反光，红青霓虹边缘光，35mm，电影感，超细节",
      "revisedPrompt": "同一提示词批量结果 3，强化雨滴高光、轮廓光与街面反射",
      "mode": "generate",
      "size": "2048x2048",
      "quality": "medium",
      "outputFormat": "png",
      "background": "auto",
      "outputCompression": 100,
      "inputFidelity": "auto",
      "imageStyle": "default",
      "moderation": "low",
      "createdAt": 1759993400000,
      "styleTag": "电影海报",
      "batchIndex": 2,
      "seed": 3202,
      "elapsedSec": 9,
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-c.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-a.png"
    },
    {
      "id": "hist-4",
      "prompt": "复古未来主义列车站台，雨雾、钠灯、金属结构、低角度广角",
      "revisedPrompt": "高对比、冷暖霓虹、边缘轮廓光、海报构图、主体占中",
      "mode": "edit",
      "size": "2880x2880",
      "quality": "high",
      "outputFormat": "png",
      "background": "opaque",
      "outputCompression": 100,
      "inputFidelity": "high",
      "imageStyle": "default",
      "moderation": "low",
      "createdAt": 1759990100000,
      "styleTag": "电影海报",
      "batchIndex": 0,
      "seed": 3203,
      "elapsedSec": 10,
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-b.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-b.png"
    },
    {
      "id": "hist-5",
      "prompt": "产品棚拍样张，银色耳机置于磨砂台面，柔光箱高光干净，商业摄影",
      "revisedPrompt": "极简背景、金属反射控制、轻微俯拍、留白构图",
      "mode": "generate",
      "size": "2048x2048",
      "quality": "medium",
      "outputFormat": "png",
      "background": "auto",
      "outputCompression": 100,
      "inputFidelity": "auto",
      "imageStyle": "natural",
      "moderation": "auto",
      "createdAt": 1759986800000,
      "styleTag": "商业摄影",
      "batchIndex": 0,
      "seed": 3204,
      "elapsedSec": 11,
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-c.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-a.png"
    },
    {
      "id": "hist-6",
      "prompt": "复古未来主义列车站台，雨雾、钠灯、金属结构、低角度广角",
      "revisedPrompt": "同一提示词的第二版，强化雨滴高光与面部轮廓",
      "mode": "edit",
      "size": "2048x2048",
      "quality": "medium",
      "outputFormat": "png",
      "background": "auto",
      "outputCompression": 100,
      "inputFidelity": "high",
      "imageStyle": "natural",
      "moderation": "auto",
      "createdAt": 1759983500000,
      "styleTag": "胶片人像",
      "batchIndex": 0,
      "seed": 3205,
      "elapsedSec": 12,
      "thumbPath": "/tmp/home/Pictures/Image Studio/thumb-a.png",
      "savedPath": "/tmp/home/Pictures/Image Studio/demo-b.png"
    }
  ],
  "historyFull": [
    {
      "id": "hist-1",
      "imageB64": ""
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

xwininfo -display :99 -root -tree > /src/.tmp-gio-history-xwininfo.txt || true
window_id=$(awk '/"Image Studio Gio"/ {print $1; exit}' /src/.tmp-gio-history-xwininfo.txt)
if [ -n "$window_id" ]; then
  import -display :99 -window "$window_id" /src/.tmp-gio-history-shot.png
fi

kill "$app_pid" >/dev/null 2>&1 || true
wait "$app_pid" || true
