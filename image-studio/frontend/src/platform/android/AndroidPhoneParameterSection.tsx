import type { CSSProperties, Dispatch, SetStateAction } from "react";
import type { QualityValue } from "../../types/domain";
import { QUALITY_TIERS, STYLE_CHIPS } from "../../components/panel/panelOptions";
import {
  ASPECT_PRESETS,
  RESOLUTION_PRESETS,
  sizeCapabilityHint,
  type AspectPreset,
  type ResolutionPreset,
} from "../../components/panel/sizeCapabilities";
import { vibrateForPlatform } from "./bridge";

export function AndroidPhoneParameterSection({
  activeAspect,
  activeAspectLabel,
  activeResolution,
  activeResolutionLabel,
  activeQualityLabel,
  activeStyleLabel,
  availableResolutions,
  apiMode,
  batchCount,
  handleAspectSelect,
  handleResolutionSelect,
  imageModelID,
  parametersOpen,
  quality,
  requestPolicy,
  setField,
  setParametersOpen,
  styleTag,
}: {
  activeAspect: AspectPreset;
  activeAspectLabel: string;
  activeResolution: ResolutionPreset;
  activeResolutionLabel: string;
  activeQualityLabel: string;
  activeStyleLabel: string;
  availableResolutions: ResolutionPreset[];
  apiMode: "responses" | "images";
  batchCount: number;
  handleAspectSelect: (aspect: AspectPreset) => void;
  handleResolutionSelect: (resolution: ResolutionPreset) => void;
  imageModelID: string;
  parametersOpen: boolean;
  quality: string;
  requestPolicy: "openai" | "compat";
  setField: (key: "quality" | "styleTag" | "batchCount", value: any) => void;
  setParametersOpen: Dispatch<SetStateAction<boolean>>;
  styleTag: string;
}) {
  const toggleParameters = () => {
    vibrateForPlatform(8);
    setParametersOpen((current) => !current);
  };
  const resolutionHint = sizeCapabilityHint({ apiMode, requestPolicy, imageModelID });

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
            <span>{activeAspectLabel}</span>
            <span>{activeResolutionLabel}</span>
            <span>{activeQualityLabel}</span>
            <span>{batchCount} 张</span>
          </div>
        </div>
        <span className="android-phone-summary-cta">{parametersOpen ? "收起" : "编辑"}</span>
      </button>

      {parametersOpen ? (
        <div className="mt-3 flex flex-col gap-4">
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
              {ASPECT_PRESETS.map((item) => {
                const active = activeAspect === item.value;
                return (
                  <button
                    key={item.value}
                    type="button"
                    onClick={() => { vibrateForPlatform(5); handleAspectSelect(item.value); }}
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

          <DiscreteSlider
            label="分辨率"
            value={activeResolution}
            options={RESOLUTION_PRESETS.filter((item) => availableResolutions.includes(item.value))}
            onChange={handleResolutionSelect}
            note={resolutionHint}
          />

          <DiscreteSlider
            label="画面质量"
            value={quality as QualityValue}
            options={QUALITY_TIERS}
            onChange={(next) => setField("quality", next)}
          />

          <DiscreteSlider
            label="出图张数"
            value={batchCount}
            options={BATCH_COUNT_OPTIONS}
            onChange={(next) => setField("batchCount", next)}
            valueSuffix="张"
          />
        </div>
      ) : null}
    </section>
  );
}

const BATCH_COUNT_OPTIONS = [1, 2, 4, 6, 8, 9].map((count) => ({
  value: count,
  label: String(count),
}));

type SliderValue = string | number;

function DiscreteSlider<T extends SliderValue>({
  label,
  note,
  onChange,
  options,
  value,
  valueSuffix = "",
}: {
  label: string;
  note?: string;
  onChange: (value: T) => void;
  options: Array<{ value: T; label: string }>;
  value: T;
  valueSuffix?: string;
}) {
  const activeIndex = Math.max(0, options.findIndex((item) => item.value === value));
  const denominator = Math.max(1, options.length - 1);
  const progress = `${(activeIndex / denominator) * 100}%`;
  const activeOption = options[activeIndex] ?? options[0];
  const disabled = options.length < 2;

  const commit = (index: number) => {
    const next = options[index];
    if (!next || next.value === value) return;
    vibrateForPlatform(4);
    onChange(next.value);
  };

  return (
    <div className="android-phone-settings-group android-phone-slider-group">
      <div className="android-phone-slider-head">
        <div className="text-[12px] font-medium text-zinc-600 dark:text-zinc-300">{label}</div>
        <output className="font-mono-token android-phone-slider-value">
          {activeOption?.label}{valueSuffix}
        </output>
      </div>
      <div
        className="android-phone-discrete-slider"
        style={{ "--android-slider-progress": progress } as CSSProperties}
      >
        <input
          type="range"
          min={0}
          max={Math.max(0, options.length - 1)}
          step={1}
          value={activeIndex}
          disabled={disabled}
          aria-label={label}
          aria-valuetext={`${activeOption?.label ?? ""}${valueSuffix}`}
          onChange={(event) => commit(Number(event.currentTarget.value))}
        />
      </div>
      <div
        className="android-phone-slider-ticks"
        style={{ "--android-slider-columns": options.length } as CSSProperties}
      >
        {options.map((item, index) => (
          <button
            key={`${item.value}`}
            type="button"
            aria-pressed={index === activeIndex}
            onClick={() => commit(index)}
            className={index === activeIndex ? "active" : ""}
          >
            {item.label}
          </button>
        ))}
      </div>
      {note ? (
        <p className="text-[11px] leading-5 text-zinc-500 dark:text-zinc-400">{note}</p>
      ) : null}
    </div>
  );
}
