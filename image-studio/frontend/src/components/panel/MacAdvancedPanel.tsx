import type { BackgroundValue, ImageStyleValue, InputFidelityValue, ModerationValue, OutputFormatValue } from "../../types/domain";
import {
  AdvancedBackgroundField,
  AdvancedCard,
  AdvancedImageStyleField,
  AdvancedInputFidelityField,
  AdvancedModerationField,
  AdvancedNegativePromptField,
  AdvancedPartialImagesField,
  AdvancedOutputCompressionField,
  AdvancedOutputFormatField,
  AdvancedSeedField,
  AdvancedUserIdentifierField,
} from "./AdvancedParameterBlocks";

export function MacAdvancedPanel({
  advancedOpen,
  advancedSummary,
  background,
  imageStyle,
  inputFidelity,
  moderation,
  negativePrompt,
  outputCompression,
  outputFormat,
  partialImages,
  seed,
  userIdentifier,
  setAdvancedOpen,
  setField,
  Seg,
  SegItem,
}: {
  advancedOpen: boolean;
  advancedSummary: string;
  background: BackgroundValue;
  imageStyle: ImageStyleValue;
  inputFidelity: InputFidelityValue;
  moderation: ModerationValue;
  negativePrompt: string;
  outputCompression: number;
  outputFormat: OutputFormatValue;
  partialImages: number;
  seed: number;
  userIdentifier: string;
  setAdvancedOpen: React.Dispatch<React.SetStateAction<boolean>>;
  setField: (key: string, value: any) => void;
  Seg: (props: { children: React.ReactNode }) => React.ReactNode;
  SegItem: (props: { active: boolean; onClick: () => void; children: React.ReactNode }) => React.ReactNode;
}) {
  return (
    <section className="platform-card rounded-[22px] border border-black/[0.05] bg-white/70 p-4.5 shadow-[var(--shadow-card)] dark:border-white/[0.06] dark:bg-white/[0.03]">
      <button
        onClick={() => setAdvancedOpen((v) => !v)}
        type="button"
        className="flex w-full min-w-0 items-center justify-between text-left"
      >
        <div className="min-w-0">
          <div className="text-[11px] uppercase tracking-[0.12em] text-zinc-400 dark:text-zinc-500">高级参数</div>
          <div className="mt-1.5 min-w-0 truncate text-[13px] font-normal leading-6 text-zinc-600 dark:text-zinc-300">
            {advancedSummary}
          </div>
        </div>
        <span className="shrink-0 pl-3 text-[12px] text-zinc-500 dark:text-zinc-400">{advancedOpen ? "收起 ▾" : "展开 ▸"}</span>
      </button>
      {advancedOpen && (
        <div className="mt-4 grid min-w-0 gap-[18px]">
          <div className="grid min-w-0 gap-3">
            <AdvancedCard
              title="负向提示词"
              hint="描述不希望出现的物体、色彩或构图倾向。留空时不做额外限制。"
              variant="mac"
            >
              <AdvancedNegativePromptField
                negativePrompt={negativePrompt}
                onChange={(value) => setField("negativePrompt", value)}
                variant="mac"
              />
            </AdvancedCard>

            <AdvancedCard
              title="输出格式"
              hint="PNG 保留细节最多；JPEG / WebP 更省空间。"
              variant="mac"
            >
              <AdvancedOutputFormatField
                outputFormat={outputFormat}
                onChange={(value) => setField("outputFormat", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="背景"
              hint="仅 GPT 图像模型支持。透明背景需要 PNG/WebP；`gpt-image-2` 当前不支持透明背景。"
              variant="mac"
            >
              <AdvancedBackgroundField
                background={background}
                onChange={(value) => setField("background", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="输出压缩"
              hint="仅 JPEG/WebP 生效，范围 `0-100`，默认 `100`。"
              variant="mac"
            >
              <AdvancedOutputCompressionField
                outputCompression={outputCompression}
                onChange={(value) => setField("outputCompression", value)}
                variant="mac"
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="输入保真"
              hint="用于图生图/参考图流程。`gpt-image-2` 会自动高保真并忽略此项。"
              variant="mac"
            >
              <AdvancedInputFidelityField
                inputFidelity={inputFidelity}
                onChange={(value) => setField("inputFidelity", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="图像风格"
              hint="仅 `dall-e-3` 文生图支持；默认值会省略该字段。"
              variant="mac"
            >
              <AdvancedImageStyleField
                imageStyle={imageStyle}
                onChange={(value) => setField("imageStyle", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="内容审核"
              hint="`low` 更宽松；`auto` 使用官方默认审核强度。仅 GPT 图像模型支持。"
              variant="mac"
            >
              <AdvancedModerationField
                moderation={moderation}
                onChange={(value) => setField("moderation", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="稳定用户标识"
              hint="建议传哈希后的用户名或邮箱。Responses 发 `safety_identifier`，Images API 发 `user`。"
              variant="mac"
            >
              <AdvancedUserIdentifierField
                userIdentifier={userIdentifier}
                onChange={(value) => setField("userIdentifier", value)}
                variant="mac"
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="流式预览帧数"
              hint="官方 `partial_images` 范围 `0-3`。`0` 只返回最终图，`1-3` 会流式回预览帧。"
              variant="mac"
            >
              <AdvancedPartialImagesField
                partialImages={partialImages}
                onChange={(value) => setField("partialImages", value)}
                Seg={Seg}
                SegItem={SegItem}
                noteClassName="text-[11px] leading-6 text-zinc-500 dark:text-zinc-400"
              />
            </AdvancedCard>

            <AdvancedCard
              title="随机种子"
              hint={seed > 0 ? `当前固定为 ${seed}` : "留空即随机，每次生成都会变化。"}
              variant="mac"
            >
              <AdvancedSeedField
                seed={seed}
                onChange={(value) => setField("seed", value)}
                onRandomize={() => setField("seed", Math.floor(Math.random() * 2_000_000_000))}
                onClear={() => setField("seed", 0)}
                variant="mac"
              />
            </AdvancedCard>
          </div>

          <div className="rounded-[18px] border border-black/[0.05] bg-black/[0.025] px-3.5 py-3 text-[11px] leading-[1.65] text-zinc-500 dark:border-white/[0.06] dark:bg-white/[0.025] dark:text-zinc-400">
            `background` / `output_compression` / `input_fidelity` / `style` / `moderation` / `partial_images` / `user`(`safety_identifier`) 都是官方字段；`seed` / `negative prompt` 仍只在兼容中转扩展策略下发送。
          </div>
        </div>
      )}
    </section>
  );
}
