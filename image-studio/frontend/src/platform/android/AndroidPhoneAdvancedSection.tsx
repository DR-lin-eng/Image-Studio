import { ChevronDown, ChevronRight, Dices, X } from "lucide-react";
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
    <section className="android-phone-advanced-block">
      <button
        type="button"
        onClick={toggleAdvanced}
        className="platform-card android-phone-advanced-toggle"
      >
        <span className="android-phone-kicker !mb-0">高级参数</span>
        <span className="android-phone-advanced-toggle-state">
          {advancedOpen ? "收起" : "展开"}
          {advancedOpen ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronRight className="h-3.5 w-3.5" />}
        </span>
      </button>
      {advancedOpen ? (
        <div className="platform-card android-phone-advanced-card">
          <button
            type="button"
            role="switch"
            aria-checked={noPromptRevision}
            onClick={() => {
              if (apiMode !== "responses") return;
              vibrateForPlatform(5);
              setField("noPromptRevision", !noPromptRevision);
            }}
            className={`android-phone-advanced-switch ${noPromptRevision ? "active" : ""} ${apiMode !== "responses" ? "disabled" : ""}`}
            title={apiMode === "responses" ? "逐字把当前提示词发给图像模型" : "Images API 不支持该项"}
          >
            <span className="android-phone-advanced-copy">
              <span className="android-phone-advanced-title">逐字提示词</span>
              <span className="android-phone-advanced-caption">按原始 prompt 生成</span>
            </span>
            <span className={`android-phone-switch ${noPromptRevision ? "active" : ""}`}>
              <span className={`android-phone-switch-knob ${noPromptRevision ? "active" : ""}`} />
            </span>
          </button>

          <div className="android-phone-advanced-section">
            <div className="android-phone-advanced-label">负向提示词</div>
            <textarea
              value={negativePrompt}
              placeholder="不希望出现的元素"
              onChange={(e) => setField("negativePrompt", e.target.value)}
              className="focus-ring android-phone-advanced-textarea"
            />
          </div>

          <div className="android-phone-advanced-section">
            <div className="android-phone-advanced-label">输出格式</div>
            <div className="android-phone-format-row">
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
          </div>

          <div className="android-phone-advanced-section">
            <div className="android-phone-advanced-label">Seed</div>
            <div className="android-phone-seed-row">
              <input
                type="number"
                value={seed || ""}
                placeholder="留空为随机"
                min={0}
                onChange={(e) => setField("seed", Number(e.target.value) || 0)}
                className="focus-ring android-phone-seed-input font-mono-token"
              />
              <button
                type="button"
                onClick={() => { vibrateForPlatform(5); setField("seed", Math.floor(Math.random() * 2_000_000_000)); }}
                title="随机 seed"
                className="platform-action-btn android-phone-seed-icon-button"
              >
                <Dices className="h-3.5 w-3.5" />
              </button>
              {seed > 0 ? (
                <button
                  type="button"
                  onClick={() => { vibrateForPlatform(5); setField("seed", 0); }}
                  title="清除"
                  className="platform-action-btn android-phone-seed-icon-button danger"
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
