import { ImageIcon, Upload } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";

// EmptyState 中间的提示卡。背景动效不在这里实现 —— 它是 stage-host 自带的
// 棋盘格 + CSS keyframes 在 _canvas.css 里。这里只负责中央内容。
export function EmptyState() {
  const importImageFile = useStudioStore((s) => s.importImageFile);

  function onFilePick(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0];
    if (f) importImageFile(f);
    e.target.value = "";
  }

  return (
    <div className="absolute inset-0 flex items-center justify-center px-8 pointer-events-none">
      <div className="relative z-10 text-center max-w-sm pointer-events-auto">
        <div className="inline-flex w-16 h-16 mb-4 rounded-2xl bg-emerald-500/10 ring-1 ring-emerald-500/25 items-center justify-center backdrop-blur-sm">
          <ImageIcon className="w-7 h-7 text-emerald-400" />
        </div>
        <h2 className="text-base font-semibold text-zinc-900 dark:text-zinc-200 mb-1">还没有图片</h2>
        <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed mb-4">
          在左侧填好 prompt 后点「生成」, 或者拖一张本地图片到这里来编辑
        </p>
        <label className="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm text-zinc-700 dark:text-zinc-300 ring-1 ring-black/[0.08] dark:ring-white/[0.06] bg-white/60 dark:bg-zinc-900/60 backdrop-blur-sm hover:ring-emerald-500/40 hover:text-emerald-400 cursor-pointer transition-colors">
          <Upload className="w-3.5 h-3.5" />
          选择本地图片
          <input type="file" accept="image/png,image/jpeg,image/webp" onChange={onFilePick} className="hidden" />
        </label>
      </div>
    </div>
  );
}
