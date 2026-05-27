import { ImagePlus, Trash2, Wand2, X } from "lucide-react";
import type { HistoryItem, SourceImage } from "../../types/domain";
import { vibrateForPlatform } from "./bridge";

export function AndroidPhoneSourceSection({
  clearSources,
  currentImage,
  editSourceLabel,
  onSelectSource,
  removeSource,
  sources,
}: {
  clearSources: () => void;
  currentImage: HistoryItem | null;
  editSourceLabel: string;
  onSelectSource: () => void;
  removeSource: (index: number) => void;
  sources: SourceImage[];
}) {
  return (
    <section className="platform-card android-phone-source-card p-4">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="android-phone-kicker">源图片 / 参考图</div>
          <div className="mt-1 text-[15px] font-semibold text-zinc-900 dark:text-zinc-100">{editSourceLabel}</div>
          <p className="mt-1 text-[12px] leading-6 text-zinc-500 dark:text-zinc-300">
            {sources.length > 0
              ? "已添加显式参考图，可继续替换或补充更多图。"
              : currentImage?.savedPath
                ? "当前画板图片会作为隐式源图参与本次编辑。"
                : "先添加一张图，或者从历史里挑一张结果继续编辑。"}
          </p>
        </div>
        <Wand2 className="mt-1 h-4 w-4 shrink-0 text-zinc-400" />
      </div>
      {sources.length > 0 ? (
        <div className="mt-3 flex flex-col gap-2">
          {sources.map((source, index) => (
            <div key={source.path} className="flex items-center gap-2 rounded-[16px] border border-black/[0.06] bg-[var(--surface)] px-3 py-2 dark:border-white/[0.06]">
              <span className="min-w-0 flex-1 truncate text-[12px] text-zinc-700 dark:text-zinc-300" title={source.path}>
                {index + 1}. {source.name}
              </span>
              <button
                type="button"
                onClick={() => { vibrateForPlatform(5); removeSource(index); }}
                title="移除"
                className="rounded-full p-1 text-zinc-400 transition-colors hover:bg-red-500/10 hover:text-red-400"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      ) : null}
      <div className="mt-3 flex gap-2">
        <button
          type="button"
          onClick={onSelectSource}
          className="platform-action-btn inline-flex min-h-[42px] flex-1 items-center justify-center gap-1.5 border border-black/[0.08] px-3 py-2 text-[12px] text-zinc-700 transition-colors hover:border-[color:var(--accent)]/35 hover:text-[var(--accent)] dark:border-white/[0.08] dark:text-zinc-300"
        >
          <ImagePlus className="h-3.5 w-3.5" /> 添加图片
        </button>
        {sources.length > 0 ? (
          <button
            type="button"
            onClick={() => { vibrateForPlatform(5); clearSources(); }}
            className="platform-action-btn inline-flex min-h-[42px] items-center gap-1.5 border border-black/[0.08] px-3 py-2 text-[12px] text-zinc-500 transition-colors hover:border-red-400/40 hover:text-red-400 dark:border-white/[0.08]"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        ) : null}
      </div>
    </section>
  );
}
