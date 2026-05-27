import { Dices, X } from "lucide-react";
import type { OutputFormatValue } from "../../types/domain";
import { OUTPUT_FORMAT_OPTIONS } from "../../types/domain";
import { vibrateForPlatform } from "./bridge";

export function AndroidPhoneAdvancedSection({
  advancedOpen,
  apiMode,
  negativePrompt,
  noPromptRevision,
  outputFormat,
  seed,
  setAdvancedOpen,
  setField,
}: {
  advancedOpen: boolean;
  apiMode: "responses" | "images";
  negativePrompt: string;
  noPromptRevision: boolean;
  outputFormat: OutputFormatValue;
  seed: number;
  setAdvancedOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setField: (key: string, value: any) => void;
}) {
  const toggleAdvanced = () => {
    vibrateForPlatform(8);
    setAdvancedOpen((current) => !current);
  };

  return (
    <section>
      <button
        type="button"
        onClick={toggleAdvanced}
        className="platform-card android-phone-advanced-toggle flex w-full items-center justify-between px-4 py-3 text-left text-[12px] text-zinc-500 transition-colors hover:text-zinc-900 dark:text-zinc-300 dark:hover:text-zinc-100"
      >
        <span className="android-phone-kicker !mb-0">高级参数</span>
        <span className="text-[11px] opacity-70">{advancedOpen ? "收起 ▾" : "展开 ▸"}</span>
      </button>
      {advancedOpen ? (
        <div className="platform-card android-phone-advanced-card mt-2 flex flex-col gap-3 p-4">
          <button
            type="button"
            role="switch"
            aria-checked={noPromptRevision}
            onClick={() => {
              if (apiMode !== "responses") return;
              vibrateForPlatform(5);
              setField("noPromptRevision", !noPromptRevision);
            }}
            className={`android-phone-advanced-switch inline-flex min-h-[40px] items-center justify-between gap-3 rounded-[16px] border px-3 py-2 text-[12px] transition-colors ${
              noPromptRevision
                ? "border-[color:var(--accent)]/20 bg-[var(--accent-soft)] text-[var(--accent)]"
                : "border-black/[0.08] bg-[var(--surface)] text-zinc-600 dark:border-white/[0.08] dark:text-zinc-300"
            } ${apiMode !== "responses" ? "cursor-not-allowed opacity-50" : "cursor-pointer"}`}
            title={apiMode === "responses" ? "逐字把当前提示词发给图像模型" : "Images API 不支持该项"}
          >
            <span className="android-phone-advanced-copy min-w-0">
              <span className="android-phone-advanced-title block font-medium">逐字提示词</span>
              <span className="android-phone-advanced-caption mt-0.5 block text-[11px] opacity-75">关闭模型改写，按 prompt 原样出图。</span>
            </span>
            <span className={`android-phone-switch ${noPromptRevision ? "active" : ""}`}>
              <span className={`android-phone-switch-knob ${noPromptRevision ? "active" : ""}`} />
            </span>
          </button>

          <div className="android-phone-advanced-section">
            <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">负向提示词</div>
            <textarea
              value={negativePrompt}
              placeholder="不希望出现的元素"
              onChange={(e) => setField("negativePrompt", e.target.value)}
              className="focus-ring min-h-[72px] w-full resize-none border border-black/[0.08] bg-[var(--surface)] px-4 py-3 text-[13px] leading-6 text-zinc-900 placeholder:text-zinc-400 dark:border-white/[0.08] dark:text-zinc-100 dark:placeholder:text-zinc-500"
            />
          </div>

          <div className="android-phone-advanced-section">
            <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">输出格式</div>
            <div className="grid grid-cols-3 gap-2">
              {OUTPUT_FORMAT_OPTIONS.map((item) => (
                <button
                  key={item.value}
                  type="button"
                  onClick={() => { vibrateForPlatform(5); setField("outputFormat", item.value as OutputFormatValue); }}
                  className={`android-choice-chip ${outputFormat === item.value ? "active" : ""}`}
                >
                  {item.label}
                </button>
              ))}
            </div>
            <p className="mt-2 text-[11px] leading-5 text-zinc-500 dark:text-zinc-400">
              JPEG / WebP 更省空间，导出时 `jpeg` 会写成 `.jpg`。
            </p>
          </div>

          <div className="android-phone-advanced-section">
            <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">Seed</div>
            <div className="flex gap-2">
              <input
                type="number"
                value={seed || ""}
                placeholder="留空为随机"
                min={0}
                onChange={(e) => setField("seed", Number(e.target.value) || 0)}
                className="focus-ring min-h-[42px] flex-1 border border-black/[0.08] bg-[var(--surface)] px-4 text-[13px] font-mono-token text-zinc-900 placeholder:text-zinc-400 dark:border-white/[0.08] dark:text-zinc-100 dark:placeholder:text-zinc-500"
              />
              <button
                type="button"
                onClick={() => { vibrateForPlatform(5); setField("seed", Math.floor(Math.random() * 2_000_000_000)); }}
                title="随机 seed"
                className="platform-action-btn inline-flex min-h-[42px] items-center justify-center border border-black/[0.08] px-3 text-zinc-600 transition-colors hover:border-[color:var(--accent)]/35 hover:text-[var(--accent)] dark:border-white/[0.08] dark:text-zinc-400"
              >
                <Dices className="h-3.5 w-3.5" />
              </button>
              {seed > 0 ? (
                <button
                  type="button"
                  onClick={() => { vibrateForPlatform(5); setField("seed", 0); }}
                  title="清除"
                  className="platform-action-btn inline-flex min-h-[42px] items-center justify-center border border-black/[0.08] px-3 text-zinc-500 transition-colors hover:border-red-400/40 hover:text-red-400 dark:border-white/[0.08]"
                >
                  <X className="h-3.5 w-3.5" />
                </button>
              ) : null}
            </div>
          </div>
        </div>
      ) : null}
    </section>
  );
}
