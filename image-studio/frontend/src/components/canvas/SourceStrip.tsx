import { useState } from "react";
import { Plus, X } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";

export function SourceStrip() {
  const sources = useStudioStore((s) => s.sources);
  const removeSource = useStudioStore((s) => s.removeSource);
  const reorderSources = useStudioStore((s) => s.reorderSources);
  const mode = useStudioStore((s) => s.mode);
  const selectSourceImage = useStudioStore((s) => s.selectSourceImage);

  const [dragFrom, setDragFrom] = useState<number | null>(null);
  const [overIdx, setOverIdx] = useState<number | null>(null);

  if (mode !== "edit") return null;
  if (sources.length === 0) return null;

  return (
    <div className="flex items-center gap-2 px-3 py-2 border-b border-black/[0.08] dark:border-white/[0.06] bg-white/70 dark:bg-zinc-950/40 backdrop-blur-sm overflow-x-auto">
      <span className="text-[11px] text-zinc-500 shrink-0">参考图 {sources.length} 张:</span>
      {sources.map((s, i) => (
        <div
          key={s.path}
          draggable
          onDragStart={() => setDragFrom(i)}
          onDragOver={(e) => { e.preventDefault(); setOverIdx(i); }}
          onDragLeave={() => setOverIdx(null)}
          onDrop={(e) => {
            e.preventDefault();
            if (dragFrom != null && dragFrom !== i) reorderSources(dragFrom, i);
            setDragFrom(null);
            setOverIdx(null);
          }}
          onDragEnd={() => { setDragFrom(null); setOverIdx(null); }}
          title={`${i + 1}. ${s.name}\n${s.path}`}
          className={`relative group w-12 h-12 rounded-md overflow-hidden ring-1 transition-all cursor-grab shrink-0 ${
            overIdx === i
              ? "ring-emerald-500 ring-2 scale-105"
              : "ring-black/[0.08] dark:ring-white/[0.06] hover:ring-emerald-500/40"
          }`}
        >
          <span className="absolute top-0 left-0 z-10 px-1 text-[9px] bg-zinc-950/70 text-white rounded-br">
            {i + 1}
          </span>
          {s.imageB64 ? (
            <img src={`data:image/png;base64,${s.imageB64}`} alt={s.name} className="w-full h-full object-cover" />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-[10px] text-zinc-500 bg-zinc-100 dark:bg-zinc-800">
              {s.name.split(".").slice(-1)[0].toUpperCase()}
            </div>
          )}
          <button
            onClick={(e) => { e.stopPropagation(); removeSource(i); }}
            title="移除"
            className="absolute top-0 right-0 w-4 h-4 flex items-center justify-center bg-zinc-950/70 text-white opacity-0 group-hover:opacity-100 hover:bg-red-500 rounded-bl transition-all"
          >
            <X className="w-2.5 h-2.5" />
          </button>
        </div>
      ))}
      <button
        onClick={selectSourceImage}
        title="添加参考图"
        className="w-12 h-12 rounded-md flex items-center justify-center ring-1 ring-dashed ring-zinc-300 dark:ring-zinc-700 text-zinc-500 hover:ring-emerald-500/40 hover:text-emerald-400 transition-colors shrink-0"
      >
        <Plus className="w-4 h-4" />
      </button>
    </div>
  );
}
