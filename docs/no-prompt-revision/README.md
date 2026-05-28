# 不优化提示词功能说明

「不优化提示词」有用，但不是绝对硬禁止。

它用于 Responses API 模式。默认情况下，文本模型会先理解并改写用户输入的 prompt，再把改写后的内容交给 `image_generation` 工具生成图片。开启该开关后，Image Studio 会在 `/v1/responses` 请求顶层加入一段 `instructions`，要求文本模型把用户 prompt 原样传给 `image_generation`，不要重写、扩写、润色或调整措辞。

## 适合什么时候开启

- 你已经精修过 prompt，希望图像模型尽量逐字执行。
- prompt 里有固定格式、专有名词、镜头参数、构图要求或中英文混排内容。
- 你想减少 Responses API 文本模型二次发挥导致的风格漂移。

## 功能边界

- 这是一个模型指令约束，不是上游 API 提供的强制参数。
- 它能明显降低 prompt 被改写的概率，但不能保证所有上游、所有模型 100% 遵守。
- 只对 Responses API 模式有意义；Images API 模式本来就是直接把 prompt 发给图像接口，所以该开关不可用。

## 实现路径

开启后，请求 payload 会包含类似下面的顶层指令：

```text
Pass the user prompt to image_generation VERBATIM.
DO NOT rewrite, expand, polish, or revise it in any way.
Use the exact text the user gave.
```

桌面 Wails 后端、前端 remote kernel、Android/Web 路径和 Cloudflare Worker 共享同一套语义：只要 `noPromptRevision` 为 `true`，Responses API payload 就会带上这条 `VERBATIM` 指令。

## 如何判断它是否生效

生成完成后打开输出目录的 `log/sse-response-*.txt`，查看上游响应中的 `instructions` 和 `revised_prompt`：

- 未开启时，`revised_prompt` 往往会出现被扩写或润色后的版本。
- 开启后，请求侧会带 `VERBATIM` 指令；如果上游遵守，`revised_prompt` 应尽量贴近原始 prompt。

