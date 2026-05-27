import { QUALITY_TIERS, STYLE_CHIPS, ASPECT_OPTIONS } from "../../components/panel/panelOptions";
import type { QualityValue } from "../../types/domain";
import { vibrateForPlatform } from "./bridge";

export function AndroidPadParameterSection({
  activeAspectLabel,
  activeQualityLabel,
  activeStyleLabel,
  batchCount,
  batchOpen,
  isMediumPad,
  mode,
  needsUpstreamSetup,
  onOpenUpstream,
  quality,
  qualityOpen,
  setAspectOpen,
  setBatchOpen,
  setField,
  setQualityOpen,
  setStyleOpen,
  size,
  styleOpen,
  styleTag,
  aspectOpen,
}: {
  activeAspectLabel: string;
  activeQualityLabel: string;
  activeStyleLabel: string;
  aspectOpen: boolean;
  batchCount: number;
  batchOpen: boolean;
  isMediumPad: boolean;
  mode: "generate" | "edit";
  needsUpstreamSetup: boolean;
  onOpenUpstream: () => void;
  quality: string;
  qualityOpen: boolean;
  setAspectOpen: (updater: (value: boolean) => boolean) => void;
  setBatchOpen: (updater: (value: boolean) => boolean) => void;
  setField: (key: "quality" | "batchCount" | "styleTag" | "size", value: any) => void;
  setQualityOpen: (updater: (value: boolean) => boolean) => void;
  setStyleOpen: (updater: (value: boolean) => boolean) => void;
  size: string;
  styleOpen: boolean;
  styleTag: string;
}) {
  return (
    <section className="platform-card android-pad-parameter-card p-5">
      <div className="android-pad-parameter-head">
        <div className="android-pad-parameter-copy">
          <div className="android-phone-kicker">创作参数</div>
          <div className="mt-1 text-[16px] font-semibold text-zinc-900 dark:text-zinc-100">
            {activeStyleLabel}
          </div>
          <div className="android-phone-summary-chips mt-2">
            <span>{activeQualityLabel}</span>
            <span>{activeAspectLabel}</span>
            <span>{batchCount} 张</span>
          </div>
        </div>
        {needsUpstreamSetup ? (
          <button
            type="button"
            onClick={onOpenUpstream}
            className="platform-action-btn inline-flex min-h-[42px] items-center gap-1.5 border border-[color:var(--accent)]/20 bg-white/70 px-3 py-2 text-[12px] text-[var(--accent)] dark:bg-white/[0.05]"
          >
            打开设置
          </button>
        ) : null}
      </div>

      {isMediumPad ? (
        <div className="mt-4 flex flex-col gap-3">
          <PaneBlock title="质量">
            <button
              type="button"
              onClick={() => { vibrateForPlatform(8); setQualityOpen((v) => !v); }}
              className="android-pad-medium-toggle"
            >
              <span>{activeQualityLabel}</span>
              <span>{qualityOpen ? "收起 ▾" : "展开 ▸"}</span>
            </button>
            {qualityOpen ? (
              <div className="mt-2 grid grid-cols-4 gap-2">
                {QUALITY_TIERS.map((item) => (
                  <button
                    key={item.value}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("quality", item.value as QualityValue); }}
                    className={`android-choice-chip ${quality === item.value ? "active" : ""}`}
                  >
                    {item.label}
                  </button>
                ))}
              </div>
            ) : null}
          </PaneBlock>

          <PaneBlock title="出图张数">
            <button
              type="button"
              onClick={() => { vibrateForPlatform(8); setBatchOpen((v) => !v); }}
              className="android-pad-medium-toggle"
            >
              <span>{batchCount} 张</span>
              <span>{batchOpen ? "收起 ▾" : "展开 ▸"}</span>
            </button>
            {batchOpen ? (
              <div className="mt-2 grid grid-cols-6 gap-2">
                {[1, 2, 4, 6, 8, 9].map((count) => (
                  <button
                    key={count}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("batchCount", count); }}
                    className={`android-choice-chip ${batchCount === count ? "active" : ""}`}
                  >
                    {count}
                  </button>
                ))}
              </div>
            ) : null}
          </PaneBlock>

          <PaneBlock title="风格">
            <button
              type="button"
              onClick={() => { vibrateForPlatform(8); setStyleOpen((v) => !v); }}
              className="android-pad-medium-toggle"
            >
              <span>{activeStyleLabel}</span>
              <span>{styleOpen ? "收起 ▾" : "展开 ▸"}</span>
            </button>
            {styleOpen ? (
              <div className="mt-2 grid grid-cols-3 gap-2">
                {STYLE_CHIPS.map((item) => {
                  const active = styleTag === item.id;
                  return (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => { vibrateForPlatform(5); setField("styleTag", active ? "" : item.id); }}
                      className={`platform-chip inline-flex min-h-[42px] items-center justify-center px-3 text-[12px] ${
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
            ) : null}
          </PaneBlock>

          <PaneBlock title="比例">
            <button
              type="button"
              onClick={() => { vibrateForPlatform(8); setAspectOpen((v) => !v); }}
              className="android-pad-medium-toggle"
            >
              <span>{activeAspectLabel}</span>
              <span>{aspectOpen ? "收起 ▾" : "展开 ▸"}</span>
            </button>
            {aspectOpen ? (
              <div className="mt-2 grid grid-cols-3 gap-2">
                {ASPECT_OPTIONS.map((item) => {
                  const active = size === item.value;
                  return (
                    <button
                      key={item.value}
                      type="button"
                      onClick={() => { vibrateForPlatform(5); setField("size", item.value); }}
                      className={`android-aspect-card ${active ? "active" : ""}`}
                      title={item.auto ? "让上游决定尺寸与比例" : item.value}
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
            ) : null}
          </PaneBlock>
        </div>
      ) : (
        <>
          <div className="android-pad-parameter-grid mt-4">
            <PaneBlock title="质量">
              <div className="grid grid-cols-2 gap-2">
                {QUALITY_TIERS.map((item) => (
                  <button
                    key={item.value}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("quality", item.value as QualityValue); }}
                    className={`android-choice-chip ${quality === item.value ? "active" : ""}`}
                  >
                    {item.label}
                  </button>
                ))}
              </div>
            </PaneBlock>

            <PaneBlock title="出图张数">
              <div className="grid grid-cols-3 gap-2">
                {[1, 2, 4, 6, 8, 9].map((count) => (
                  <button
                    key={count}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("batchCount", count); }}
                    className={`android-choice-chip ${batchCount === count ? "active" : ""}`}
                  >
                    {count}
                  </button>
                ))}
              </div>
            </PaneBlock>
          </div>

          <PaneBlock title="风格" className="mt-3">
            <div className="flex flex-wrap gap-2">
              {STYLE_CHIPS.map((item) => {
                const active = styleTag === item.id;
                return (
                  <button
                    key={item.id}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("styleTag", active ? "" : item.id); }}
                    className={`platform-chip inline-flex min-h-[38px] items-center px-3 text-[12px] ${
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
          </PaneBlock>

          <PaneBlock title="比例" className="mt-3">
            <div className="grid grid-cols-3 gap-2">
              {ASPECT_OPTIONS.map((item) => {
                const active = size === item.value;
                return (
                  <button
                    key={item.value}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); setField("size", item.value); }}
                    className={`android-aspect-card ${active ? "active" : ""}`}
                    title={item.auto ? "让上游决定尺寸与比例" : item.value}
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
          </PaneBlock>
        </>
      )}
    </section>
  );
}

function PaneBlock({
  title,
  children,
  className = "",
}: {
  title: string;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <section className={className}>
      <div className="mb-2 text-[12px] font-medium text-zinc-600 dark:text-zinc-300">{title}</div>
      {children}
    </section>
  );
}
