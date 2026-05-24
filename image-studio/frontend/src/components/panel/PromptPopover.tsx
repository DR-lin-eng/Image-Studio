import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { X } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";
import { isWindows, usesAppleUI } from "../../lib/platform";

const PROMPT_TEMPLATES: { label: string; text: string }[] = [
  { label: "写实摄影", text: "photorealistic, professional photography, 35mm, natural lighting, sharp focus, high detail" },
  { label: "电影感", text: "cinematic, dramatic lighting, shallow depth of field, film grain, anamorphic, 2.39:1" },
  { label: "二次元", text: "anime style, vibrant colors, cel shading, detailed illustration" },
  { label: "油画", text: "oil painting, thick brush strokes, classical art style, warm tones" },
  { label: "水彩", text: "watercolor painting, soft edges, pastel colors, paper texture" },
  { label: "扁平插画", text: "flat illustration, minimalist, geometric shapes, vector style" },
  { label: "3D 渲染", text: "3D render, octane render, ray tracing, glossy, studio lighting" },
  { label: "像素风", text: "pixel art, 16-bit, retro game style, limited palette" },
];

// PromptPopover 通过 React Portal 挂到 document.body 顶层,position: fixed +
// 高 z-index 逃出 ControlPanel 的 overflow-y-auto 和邻居(Toolbar / Canvas /
// SourceStrip)各自的 backdrop-blur stacking context。anchorRef 是触发按钮,
// 用它的 getBoundingClientRect 计算 popover 的位置;窗口滚动 / 缩放时跟着重新
// 定位。点 popover 外部任意位置自动关闭。
export function PromptPopover({
  anchorRef,
  onClose,
  onPick,
}: {
  anchorRef: React.RefObject<HTMLElement | null>;
  onClose: () => void;
  onPick: (text: string) => void;
}) {
  const history = useStudioStore((s) => s.promptHistory);
  const [tab, setTab] = useState<"templates" | "history">("templates");
  const popRef = useRef<HTMLDivElement | null>(null);
  const [pos, setPos] = useState<{ top: number; left: number; width: number } | null>(null);

  // 根据 anchor 按钮的位置算 popover 的位置:左对齐按钮,顶部贴按钮下沿 + 6px gap。
  // 窗口尺寸变化 / 用户滚动外层 / 切换 workspace 都触发重算。
  useLayoutEffect(() => {
    const compute = () => {
      const el = anchorRef.current;
      if (!el) return;
      const r = el.getBoundingClientRect();
      const w = Math.max(280, Math.min(360, r.width * 5));
      setPos({ top: r.bottom + 6, left: r.left, width: w });
    };
    compute();
    window.addEventListener("resize", compute);
    window.addEventListener("scroll", compute, true);
    return () => {
      window.removeEventListener("resize", compute);
      window.removeEventListener("scroll", compute, true);
    };
  }, [anchorRef]);

  // 点击外部关闭。escape 也关。
  useEffect(() => {
    const onDocPointer = (e: MouseEvent) => {
      if (!popRef.current) return;
      const target = e.target as Node;
      if (popRef.current.contains(target)) return;
      if (anchorRef.current && anchorRef.current.contains(target)) return;
      onClose();
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    // pointerdown 比 click 更早,避免按下后还没 mouseup 弹窗就闪一下
    document.addEventListener("pointerdown", onDocPointer, true);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("pointerdown", onDocPointer, true);
      document.removeEventListener("keydown", onKey);
    };
  }, [anchorRef, onClose]);

  if (!pos) return null;

  return createPortal(
    <div
      ref={popRef}
      style={{ position: "fixed", top: pos.top, left: pos.left, width: pos.width, zIndex: 9050 }}
      className={`flex max-h-[300px] flex-col overflow-hidden border border-black/[0.08] bg-white/95 shadow-[0_24px_60px_rgb(15_23_42_/_0.16)] backdrop-blur-2xl dark:border-white/[0.08] dark:bg-zinc-900/95 ${usesAppleUI ? "liquid-glass-panel" : ""} ${isWindows ? "rounded-[12px]" : "rounded-[18px]"}`}
    >
      <div className="flex items-center border-b border-black/[0.06] dark:border-white/[0.04]">
        <button
          onClick={() => setTab("templates")}
          className={`flex-1 px-3 py-2 text-xs font-medium transition-colors ${
            tab === "templates"
              ? "border-b-2 border-[var(--accent)] text-[var(--accent)]"
              : "text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200"
          }`}
        >
          模板
        </button>
        <button
          onClick={() => setTab("history")}
          className={`flex-1 px-3 py-2 text-xs font-medium transition-colors ${
            tab === "history"
              ? "border-b-2 border-[var(--accent)] text-[var(--accent)]"
              : "text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200"
          }`}
        >
          历史 ({history.length})
        </button>
        <button
          onClick={onClose}
          title="关闭"
          className={`px-2 py-2 text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-200 ${isWindows ? "rounded-[8px]" : ""}`}
        >
          <X className="w-3.5 h-3.5" />
        </button>
      </div>
      <div className="flex-1 overflow-y-auto p-1.5">
        {tab === "templates" && PROMPT_TEMPLATES.map((t) => (
          <button
            key={t.label}
            onClick={() => { onPick(t.text); onClose(); }}
            className={`w-full px-2.5 py-2 text-left transition-colors hover:bg-[var(--accent-soft)] ${isWindows ? "rounded-[10px]" : "rounded-[14px]"}`}
          >
            <div className="text-xs font-medium text-zinc-900 dark:text-zinc-200 mb-0.5">{t.label}</div>
            <div className="text-[10px] text-zinc-500 leading-relaxed truncate">{t.text}</div>
          </button>
        ))}
        {tab === "history" && (
          history.length === 0 ? (
            <div className="text-xs text-zinc-500 py-6 text-center">还没有提交过 prompt</div>
          ) : (
            history.map((p, i) => (
              <button
                key={i}
                onClick={() => { onPick(p); onClose(); }}
                title="点击使用"
                className={`w-full px-2.5 py-2 text-left transition-colors hover:bg-[var(--accent-soft)] ${isWindows ? "rounded-[10px]" : "rounded-[14px]"}`}
              >
                <div className="text-[11px] text-zinc-700 dark:text-zinc-300 leading-relaxed truncate">{p}</div>
              </button>
            ))
          )
        )}
      </div>
    </div>,
    document.body,
  );
}
