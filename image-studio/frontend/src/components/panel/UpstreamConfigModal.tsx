import { useState, useEffect } from "react";
import { Modal } from "../common/Modal";
import { useStudioStore } from "../../state/studioStore";

// UpstreamConfigModal — 上游接入的统一配置入口。
//
// 首次启动若 (apiKey || baseURL) 为空会自动弹出;之后可由「设置 → 修改上游配置」
// 或 ControlPanel 顶部的「🔧 上游配置」按钮手动呼起。
// 「保存」只 commit 当前编辑值;「取消」直接关闭(下次启动若仍不完整会再弹)。
export function UpstreamConfigModal({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const {
    apiMode, responsesConfig, imagesConfig,
    setField, setAPIKey,
    testAPIKey, isTestingKey,
  } = useStudioStore();

  // 本地草稿态:每个形态各持一份,modal 内切 mode 时只是切显示草稿,不污染全局。
  // 「保存」时只把当前 draftApiMode 那份草稿 commit 进 store(同时 store 会把它写到对应槽)。
  const [draftApiMode, setDraftApiMode] = useState<"responses" | "images">(apiMode);
  const [draftResponses, setDraftResponses] = useState(responsesConfig);
  const [draftImages, setDraftImages] = useState(imagesConfig);
  const [showKey, setShowKey] = useState(false);

  useEffect(() => {
    if (open) {
      setDraftApiMode(apiMode);
      setDraftResponses(responsesConfig);
      setDraftImages(imagesConfig);
    }
  }, [open, apiMode, responsesConfig, imagesConfig]);

  // 当前可见草稿 = draftApiMode 对应的那份
  const cur = draftApiMode === "responses" ? draftResponses : draftImages;
  const setCur = (patch: Partial<typeof cur>) => {
    if (draftApiMode === "responses") setDraftResponses({ ...draftResponses, ...patch });
    else setDraftImages({ ...draftImages, ...patch });
  };

  const draftBaseURL = cur.baseURL;
  const draftApiKey = cur.apiKey;
  const draftTextModel = cur.textModelID;
  const draftImageModel = cur.imageModelID;

  const canSave = draftBaseURL.trim() && draftApiKey.trim();

  // commit:把两份草稿都写回 store(每个字段调一次 setField,自动落到对应槽)+ 最后切换 apiMode
  function commit() {
    const writeMode = (m: "responses" | "images", cfg: typeof cur) => {
      // 先临时切到 m,再写它的字段(store 的 setField 是按当前 apiMode 决定写哪个槽)
      setField("apiMode", m);
      setField("baseURL", cfg.baseURL.trim());
      setAPIKey(cfg.apiKey.trim());
      setField("textModelID", cfg.textModelID.trim());
      setField("imageModelID", cfg.imageModelID.trim());
    };
    writeMode("responses", draftResponses);
    writeMode("images", draftImages);
    // 最终激活用户在 modal 里选中的那个形态
    setField("apiMode", draftApiMode);
  }

  function save() {
    commit();
    onClose();
  }

  // 在 modal 内点「测试连接」需要先把草稿提交,否则测试用的是旧值。
  function testWithCurrentDraft() {
    if (!canSave) return;
    commit();
    // setAPIKey 是异步触发 setState,但 testAPIKey 在下一个 tick 读 get() 时能拿到新值。
    setTimeout(() => testAPIKey(), 0);
  }

  return (
    <Modal open={open} onClose={onClose} title="上游配置" width={520}>
      <div className="upstream-form">
        {/* API 形态 */}
        <div className="upstream-row">
          <label className="head">API 形态</label>
          <div className="api-mode-grid">
            <button
              className={`api-mode-btn ${draftApiMode === "responses" ? "active" : ""}`}
              onClick={() => setDraftApiMode("responses")}
              type="button"
            >
              <span className="api-mode-title">Responses API</span>
              <span className="api-mode-sub">SSE 保活(CF 超时推荐)</span>
            </button>
            <button
              className={`api-mode-btn ${draftApiMode === "images" ? "active" : ""}`}
              onClick={() => setDraftApiMode("images")}
              type="button"
            >
              <span className="api-mode-title">Images API</span>
              <span className="api-mode-sub">标准 generations / edits</span>
            </button>
          </div>
          <div className="settings-hint">
            {draftApiMode === "responses" ? (
              <>
                通过 <code>/v1/responses</code> 调用模型内置的 <code>image_generation</code> 工具,
                SSE 流式接收 —— 能防 Cloudflare 524/504 超时截断。<br />
                <strong>需要 key 绑定到「拥有 gpt-5.5 模型的分组」</strong>(余额/套餐),不是 image-2 分组。
              </>
            ) : (
              <>
                通过标准 <code>/v1/images/generations</code>(文生图)+ <code>/v1/images/edits</code>
                (图生图,multipart 上传)。一次性 JSON 响应,无 SSE 保活,长推理上 CF 524 风险更高,
                但兼容性最广。<br />
                <strong>可使用标准的 image-2 / image API 分组</strong>(不需要 gpt-5.5 权限)。
              </>
            )}
          </div>
        </div>

        <div className="settings-hint" style={{ background: "var(--accent-soft)", color: "var(--accent)", borderStyle: "solid" }}>
          下方编辑的是 <strong>{draftApiMode === "responses" ? "Responses API" : "Images API"}</strong> 的配置 ——
          两种形态各存一份,切换形态时另一份不动。
        </div>

        {/* BASE_URL */}
        <div className="upstream-row">
          <label className="head">上游 BASE_URL <span className="req">*</span></label>
          <input
            className="input"
            type="text"
            value={draftBaseURL}
            placeholder="https://your-relay.example.com"
            onChange={(e) => setCur({ baseURL: e.target.value })}
            spellCheck={false}
            autoFocus={!draftBaseURL}
          />
        </div>

        {/* API Key */}
        <div className="upstream-row">
          <label className="head">API Key <span className="req">*</span></label>
          <div className="key-input-wrap">
            <input
              className="input"
              type={showKey ? "text" : "password"}
              value={draftApiKey}
              placeholder="sk-..."
              onChange={(e) => setCur({ apiKey: e.target.value })}
              spellCheck={false}
              autoComplete="off"
            />
            <button
              type="button"
              className="key-toggle-btn"
              onClick={() => setShowKey((v) => !v)}
              title={showKey ? "隐藏" : "显示"}
            >
              {showKey ? "🙈" : "👁"}
            </button>
          </div>
        </div>

        {/* 文本模型 ID — 只对 Responses API 有意义 */}
        {draftApiMode === "responses" && (
          <div className="upstream-row">
            <label className="head">文本模型 ID</label>
            <input
              className="input"
              type="text"
              value={draftTextModel}
              placeholder="留空=默认 gpt-5.5"
              onChange={(e) => setCur({ textModelID: e.target.value })}
              spellCheck={false}
            />
          </div>
        )}

        {/* 图像模型 ID */}
        <div className="upstream-row">
          <label className="head">图像模型 ID</label>
          <input
            className="input"
            type="text"
            value={draftImageModel}
            placeholder={
              draftApiMode === "responses"
                ? "留空=默认 gpt-image-2(由 image_generation 工具触发)"
                : "留空=默认 gpt-image-2(直接传给 Images API)"
            }
            onChange={(e) => setCur({ imageModelID: e.target.value })}
            spellCheck={false}
          />
        </div>

        {/* 测试连接 */}
        <div className="upstream-row">
          <button
            className="btn secondary"
            type="button"
            onClick={testWithCurrentDraft}
            disabled={!canSave || isTestingKey}
            style={{ width: "100%" }}
          >
            {isTestingKey ? "测试中..." : "🔌 测试连接(会先保存草稿)"}
          </button>
        </div>

        {/* 操作 */}
        <div className="upstream-actions">
          <button className="btn secondary" type="button" onClick={onClose}>
            稍后再配
          </button>
          <button
            className="btn"
            type="button"
            onClick={save}
            disabled={!canSave}
          >
            保存
          </button>
        </div>
        {!canSave && (
          <div className="settings-hint" style={{ marginTop: 6 }}>
            BASE_URL 和 API Key 至少要填一次才能开始生成。
          </div>
        )}
      </div>
    </Modal>
  );
}
