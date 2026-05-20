import { useMemo, useState } from "react";
import { Clipboard, Copy, FileText, Info, ListRestart, RotateCw, Sparkles, Split, Trash2, X } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";
import type { HistoryItem, Mode } from "../../types/domain";
import { ContextMenu, MenuItem } from "../common/ContextMenu";
import { RawResponseModal } from "./RawResponseModal";

type ModeFilter = "all" | Mode;
type DateFilter = "all" | "today" | "week";

function inDateFilter(h: HistoryItem, f: DateFilter): boolean {
  if (f === "all") return true;
  const now = Date.now();
  const t = h.createdAt;
  if (f === "today") {
    const d1 = new Date(now); d1.setHours(0, 0, 0, 0);
    return t >= d1.getTime();
  }
  return now - t < 7 * 24 * 3600 * 1000;
}

export function HistoryRail() {
  const {
    history, currentImage, reuseAsSource, deleteHistoryItem, setField,
    compareB, setCompareB, pushToast, fullscreen,
    applyHistoryParams, regenerateFromHistory,
    openResultDetail,
  } = useStudioStore();

  const [q, setQ] = useState("");
  const [modeF, setModeF] = useState<ModeFilter>("all");
  const [dateF, setDateF] = useState<DateFilter>("all");
  const [menu, setMenu] = useState<{ x: number; y: number; h: HistoryItem } | null>(null);
  const [rawPath, setRawPath] = useState<string | null>(null);

  const filtered = useMemo(() => {
    const needle = q.trim().toLowerCase();
    return history.filter((h) => {
      if (modeF !== "all" && h.mode !== modeF) return false;
      if (!inDateFilter(h, dateF)) return false;
      if (!needle) return true;
      const hay = `${h.prompt ?? ""} ${h.revisedPrompt ?? ""}`.toLowerCase();
      return hay.includes(needle);
    });
  }, [history, q, modeF, dateF]);

  function buildMenu(h: HistoryItem): MenuItem[] {
    return [
      { label: "详情", icon: "ℹ", onClick: () => openResultDetail(h) },
      {
        label: "复制 prompt",
        icon: "📋",
        separatorBefore: true,
        onClick: () => navigator.clipboard.writeText(h.prompt ?? "").then(
          () => pushToast("已复制 prompt", "success"),
          () => pushToast("复制失败", "error"),
        ),
      },
      {
        label: "复制本地路径",
        icon: "📁",
        disabled: !h.savedPath,
        onClick: () => navigator.clipboard.writeText(h.savedPath ?? "").then(
          () => pushToast("已复制路径", "success"),
          () => pushToast("复制失败", "error"),
        ),
      },
      { label: "查看 raw 响应", icon: "📄", disabled: !h.rawPath, onClick: () => setRawPath(h.rawPath ?? null) },
      { separatorBefore: true, label: "应用参数(不生成)", icon: "📥", onClick: () => applyHistoryParams(h) },
      { label: "以此参数重新生成", icon: "↻", onClick: () => regenerateFromHistory(h) },
      { separatorBefore: true, label: "设为源图", icon: "→", onClick: () => reuseAsSource(h), disabled: !h.savedPath },
      { label: "用作对比图 (B)", icon: "⇄", onClick: () => setCompareB(h), disabled: currentImage?.id === h.id },
      { label: "删除", icon: "✕", danger: true, separatorBefore: true, onClick: () => {
        if (window.confirm(`确定删除此历史项?\n\n${h.prompt?.slice(0, 60) || "(无 prompt)"}`)) {
          deleteHistoryItem(h.id);
        }
      } },
    ];
  }

  if (fullscreen) return null;

  return (
    <aside className="w-[280px] shrink-0 overflow-y-auto flex flex-col gap-2 p-3 bg-white/85 dark:bg-zinc-900/40 border-l border-black/[0.08] dark:border-white/[0.06]">
      <div className="flex items-center justify-between">
        <h3 className="text-xs font-semibold text-zinc-700 dark:text-zinc-300 uppercase tracking-wide">
          历史 <span className="font-mono-token text-zinc-500">({filtered.length}{filtered.length !== history.length && `/${history.length}`})</span>
        </h3>
        <button
          onClick={() => setField("currentImage", null)}
          title="清空画板(不删历史)"
          className="text-[11px] text-zinc-500 hover:text-emerald-400 transition-colors"
        >
          清空画板
        </button>
      </div>

      <input
        placeholder="搜索 prompt..."
        value={q}
        onChange={(e) => setQ(e.target.value)}
        className="w-full bg-white dark:bg-zinc-950 ring-1 ring-black/[0.08] dark:ring-white/[0.06] rounded-md px-2.5 py-1.5 text-xs text-zinc-900 dark:text-zinc-200 placeholder:text-zinc-400 dark:placeholder:text-zinc-600 focus-ring"
      />
      <div className="flex gap-1.5">
        <select
          value={modeF}
          onChange={(e) => setModeF(e.target.value as ModeFilter)}
          className="flex-1 bg-white dark:bg-zinc-950 ring-1 ring-black/[0.08] dark:ring-white/[0.06] rounded-md px-2 py-1 text-[11px] text-zinc-700 dark:text-zinc-300 focus-ring"
        >
          <option value="all">全部模式</option>
          <option value="generate">文生图</option>
          <option value="edit">图生图</option>
        </select>
        <select
          value={dateF}
          onChange={(e) => setDateF(e.target.value as DateFilter)}
          className="flex-1 bg-white dark:bg-zinc-950 ring-1 ring-black/[0.08] dark:ring-white/[0.06] rounded-md px-2 py-1 text-[11px] text-zinc-700 dark:text-zinc-300 focus-ring"
        >
          <option value="all">全部日期</option>
          <option value="today">今天</option>
          <option value="week">本周</option>
        </select>
      </div>

      <p className="text-[10px] text-zinc-500 leading-relaxed">
        点击查看 · Shift+点击对比 · 双击设源图 · 右键更多
      </p>

      {compareB && (
        <button
          onClick={() => setCompareB(null)}
          className="inline-flex items-center justify-center gap-1.5 px-2.5 py-1.5 rounded-md text-xs text-emerald-400 ring-1 ring-emerald-500/30 bg-emerald-500/8 hover:bg-emerald-500/15 transition-colors"
        >
          <Split className="w-3 h-3" /> 退出对比
        </button>
      )}

      {filtered.length === 0 ? (
        <div className="text-xs text-zinc-500 py-8 text-center">
          {q || modeF !== "all" || dateF !== "all" ? "没有匹配项" : "还没有结果"}
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-1.5">
          {filtered.map((h) => {
            const isCurrent = currentImage?.id === h.id;
            const isCompare = compareB?.id === h.id;
            return (
              <div
                key={h.id}
                title={h.prompt}
                onClick={(e) => {
                  if (e.shiftKey) {
                    if (isCompare) setCompareB(null);
                    else if (currentImage && currentImage.id !== h.id) setCompareB(h);
                  } else {
                    setField("currentImage", h);
                  }
                }}
                onDoubleClick={() => reuseAsSource(h)}
                onContextMenu={(e) => {
                  e.preventDefault();
                  setMenu({ x: e.clientX, y: e.clientY, h });
                }}
                className={`group relative aspect-square overflow-hidden rounded-md ring-1 cursor-pointer transition-all ${
                  isCurrent
                    ? "ring-2 ring-emerald-500"
                    : isCompare
                      ? "ring-2 ring-blue-400"
                      : "ring-black/[0.08] dark:ring-white/[0.06] hover:ring-emerald-500/40"
                }`}
              >
                <img src={`data:image/png;base64,${h.imageB64}`} alt={h.prompt} className="w-full h-full object-cover" />
                <span className="absolute top-0.5 left-0.5 px-1 text-[9px] bg-zinc-950/70 text-white rounded">
                  {h.mode === "edit" ? "✎" : "✨"}
                </span>
                {isCompare && (
                  <span className="absolute top-0.5 right-0.5 px-1 text-[9px] bg-blue-500 text-white rounded">B</span>
                )}
                <button
                  onClick={(e) => { e.stopPropagation(); deleteHistoryItem(h.id); }}
                  title="删除"
                  className="absolute bottom-0.5 right-0.5 w-5 h-5 flex items-center justify-center rounded bg-zinc-950/70 text-white opacity-0 group-hover:opacity-100 hover:bg-red-500 transition-all"
                >
                  <X className="w-3 h-3" />
                </button>
              </div>
            );
          })}
        </div>
      )}

      {menu && <ContextMenu x={menu.x} y={menu.y} items={buildMenu(menu.h)} onClose={() => setMenu(null)} />}
      {rawPath && <RawResponseModal path={rawPath} onClose={() => setRawPath(null)} />}

      {/* 防止 lucide import 未使用警告 */}
      <Clipboard className="hidden" /><Copy className="hidden" /><FileText className="hidden" />
      <Info className="hidden" /><ListRestart className="hidden" /><RotateCw className="hidden" /><Sparkles className="hidden" />
      <Trash2 className="hidden" />
    </aside>
  );
}
