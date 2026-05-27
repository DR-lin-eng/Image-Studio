import type { QualityValue } from "../../types/domain";
import { ASPECT_OPTIONS, QUALITY_TIERS, STYLE_CHIPS } from "../../components/panel/panelOptions";
import { vibrateForPlatform } from "./bridge";

export function AndroidPhoneParameterSection({
  activeAspectLabel,
  activeQualityLabel,
  activeStyleLabel,
  batchCount,
  parametersOpen,
  quality,
  setField,
  setParametersOpen,
  size,
  styleTag,
}: {
  activeAspectLabel: string;
  activeQualityLabel: string;
  activeStyleLabel: string;
  batchCount: number;
  parametersOpen: boolean;
  quality: string;
  setField: (key: "quality" | "styleTag" | "size" | "batchCount", value: any) => void;
  setParametersOpen: React.Dispatch<React.SetStateAction<boolean>>;
  size: string;
  styleTag: string;
}) {
  const toggleParameters = () => {
    vibrateForPlatform(8);
    setParametersOpen((current) => !current);
  };

  return (
    <section className="platform-card android-phone-summary-card p-4">
      <button
        type="button"
        onClick={toggleParameters}
        className="android-phone-summary-toggle"
      >
        <div className="min-w-0">
          <div className="android-phone-kicker">创作参数</div>
          <div className="mt-1 text-[16px] font-semibold text-zinc-900 dark:text-zinc-100">
            {styleTag ? activeStyleLabel : "默认风格"}
          </div>
          <div className="android-phone-summary-meta mt-2">
            <span>{activeQualityLabel}</span>
            <span>{activeAspectLabel}</span>
            <span>{batchCount} 张</span>
          </div>
        </div>
        <span className="android-phone-summary-cta">{parametersOpen ? "收起" : "编辑"}</span>
      </button>

      {parametersOpen ? (
        <div className="mt-3 flex flex-col gap-4">
          <div className="android-phone-settings-group">
            <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">质量</div>
            <div className="android-phone-settings-list android-phone-quality-list">
              {QUALITY_TIERS.map((item) => (
                <button
                  key={item.value}
                  type="button"
                  onClick={() => { vibrateForPlatform(5); setField("quality", item.value as QualityValue); }}
                  className={`android-choice-chip android-phone-list-choice ${quality === item.value ? "active" : ""}`}
                >
                  {item.label}
                </button>
              ))}
            </div>
          </div>

          <div className="android-phone-settings-group">
            <div className="mb-2 flex items-center justify-between">
              <div className="text-[12px] font-medium text-zinc-600 dark:text-zinc-300">风格</div>
              {styleTag ? (
                <button type="button" onClick={() => setField("styleTag", "")} className="text-[11px] text-[var(--accent)]">
                  清除
                </button>
              ) : null}
            </div>
            <div className="android-phone-settings-list">
              {STYLE_CHIPS.map((item) => {
                const active = styleTag === item.id;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("styleTag", active ? "" : item.id); }}
                    className={`platform-chip android-phone-list-choice inline-flex min-h-[36px] items-center px-3 text-[12px] ${
                      active
                        ? "bg-[var(--accent-soft)] text-[var(--accent)] ring-1 ring-[color:var(--accent)]/20"
                        : "ring-1 ring-black/[0.08] text-zinc-600 hover:text-zinc-900 dark:ring-white/[0.08] dark:text-zinc-400 dark:hover:text-zinc-200"
                    }`}
                  >
                    {item.label}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="android-phone-settings-group">
            <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">比例</div>
            <div className="android-phone-settings-list android-phone-aspect-list">
              {ASPECT_OPTIONS.map((item) => {
                const active = size === item.value;
                return (
                  <button
                    key={item.value}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("size", item.value); }}
                    title={item.auto ? "让上游决定尺寸与比例" : item.value}
                    className={`android-aspect-card ${active ? "active" : ""}`}
                  >
                    <span
                      className={`block rounded-sm border-2 ${item.auto ? "border-dashed" : ""} ${active ? "border-[var(--accent)]" : "border-zinc-400 dark:border-zinc-600"}`}
                      style={{ width: item.w, height: item.h }}
                    />
                    <span className="mt-1 text-[10px]">{item.label}</span>
                  </button>
                );
              })}
            </div>
          </div>

          <div className="android-phone-settings-group">
            <div className="mb-2 flex items-center justify-between">
              <div className="text-[12px] font-medium text-zinc-600 dark:text-zinc-300">出图张数</div>
              <span className="font-mono-token text-[11px] text-zinc-400">{batchCount}x</span>
            </div>
            <div className="android-phone-settings-list android-phone-batch-list">
              {[1, 2, 4, 6, 8, 9].map((count) => (
                <button
                  key={count}
                  type="button"
                  onClick={() => { vibrateForPlatform(5); setField("batchCount", count); }}
                  className={`android-choice-chip android-phone-list-choice ${batchCount === count ? "active" : ""}`}
                >
                  {count}
                </button>
              ))}
            </div>
          </div>
        </div>
      ) : null}
    </section>
  );
}
