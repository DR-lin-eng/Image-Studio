import { Suspense, lazy, useDeferredValue, useMemo, useRef, useState } from "react";
import { ChevronDown, ChevronRight, Filter, GalleryVerticalEnd, Split } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";
import type { HistoryItem, Mode } from "../../types/domain";
import { ContextMenu } from "../common/ContextMenu";
import { RawResponseModal } from "./RawResponseModal";
import { usePlatform } from "../../platform/context";
import {
  isHistoryInDateFilter,
  matchesHistorySearch,
  type RelativeHistoryDateFilter,
} from "./historyFilters";
import { HistoryTile } from "./HistoryTile";
import { useHistoryContextMenu } from "./useHistoryContextMenu";

type ModeFilter = "all" | Mode;
type DateFilter = RelativeHistoryDateFilter;

export function HistoryRail() {
  const {
    history, currentImage, reuseAsSource, deleteHistoryItem, setField,
    compareB, setCompareB, pushToast, fullscreen,
    applyHistoryParams, regenerateFromHistory,
    openResultDetail, apiKey, baseURL, apiMode,
    profiles, activeProfileId, setActiveProfile,
    openUpstreamConfig, openHistoryTimeline, testAPIKey, isTestingKey,
    historyRailCollapsed, setHistoryRailCollapsed,
  } = useStudioStore();

  const [q, setQ] = useState("");
  const deferredQ = useDeferredValue(q);
  const [modeF, setModeF] = useState<ModeFilter>("all");
  const [dateF, setDateF] = useState<DateFilter>("all");
  const [filtersOpen, setFiltersOpen] = useState(false);
  const { isAndroidPhone, isAndroidPad, isMac, isWindows, usesAndroidUI, usesAppleUI } = usePlatform();
  // 防快速连点产生竞态:每次点击递增 epoch,后台 materialize 全图 resolve
  // 时跟当前 epoch 比对,过时的就丢弃。之前的写法是先 await 再 setField,
  // 慢的请求会在用户已经点了另一张图之后把画布盖回去。
  const selectEpochRef = useRef(0);

  const filtered = useMemo(() => {
    return history.filter((h) => {
      if (modeF !== "all" && h.mode !== modeF) return false;
      if (!isHistoryInDateFilter(h.createdAt, dateF)) return false;
      return matchesHistorySearch(h, deferredQ);
    });
  }, [history, deferredQ, modeF, dateF]);
  const recentHistory = filtered.slice(0, 6);
  const desktopFilterThreshold = isMac ? 8 : 4;
  const showHistoryFilters = !isMac && (history.length > desktopFilterThreshold || q.trim().length > 0 || modeF !== "all" || dateF !== "all");
  const historyFiltersActive = q.trim().length > 0 || modeF !== "all" || dateF !== "all";
  const showPhoneFilterToggle = isAndroidPhone && (history.length > 4 || historyFiltersActive);
  const showFilterControls = !isAndroidPhone ? showHistoryFilters : (filtersOpen || historyFiltersActive);

  async function selectCurrent(h: HistoryItem) {
    const myEpoch = ++selectEpochRef.current;
    // 1) 立即把(可能只是预览的)项摆上画布 —— 给用户即时反馈,不等磁盘 IO
    setField("currentImage", h);
    // 2) 关键:从历史栏选图 = 显式单图选择,退出批量结果网格 overlay。否则
    //    刚生成完 9 张批量,grid 一直罩在画板上,用户在历史栏怎么点都只是
    //    切 grid 里的高亮项,视觉上像「卡在第一张」。grid 可以从工具栏的
    //    openResultGrid 重新打开。
    if (useStudioStore.getState().resultGridOpen) {
      useStudioStore.getState().closeResultGrid();
    }
    // 3) previewOnly 需要后台从磁盘 / IndexedDB 读全图;读完只在 epoch 没变
    //    时才提交全图替换。epoch 变了说明用户已经点了别的图,这次结果作废。
    if (h.previewOnly) {
      try {
        const full = await useStudioStore.getState().materializeCurrentImage?.(h);
        if (selectEpochRef.current === myEpoch && full) {
          setField("currentImage", full);
        }
      } catch {
        // 读不出来就维持预览,用户可以再点一次
      }
    }
  }

  const {
    buildMenu,
    closeMenu,
    closeRaw,
    menu,
    openMenu,
    rawPath,
  } = useHistoryContextMenu({
    currentImageId: currentImage?.id ?? null,
    compareItemId: compareB?.id ?? null,
    onOpenDetail: openResultDetail,
    onApplyParams: applyHistoryParams,
    onRegenerate: (item) => void regenerateFromHistory(item),
    onReuseAsSource: (item) => void reuseAsSource(item),
    onToggleCompare: (item) => setCompareB(compareB?.id === item.id ? null : item),
    onDelete: (item) => {
      if (window.confirm(`确定删除此历史项?\n\n${item.prompt?.slice(0, 60) || "(无 prompt)"}`)) {
        deleteHistoryItem(item.id);
      }
    },
    pushToast,
  });

  if (fullscreen) return null;

  return (
    <aside className={`history-rail box-border flex w-[332px] shrink-0 flex-col overflow-y-auto border-l border-[var(--border)] bg-[var(--inspector)] px-4 py-4 backdrop-blur-2xl ${usesAppleUI ? "liquid-sidebar" : ""} ${usesAndroidUI && !isAndroidPhone ? "android-surface-pane" : ""} ${isAndroidPad ? "android-pad-history" : ""}`}>
      <div className={`history-rail-stack ${isAndroidPad ? "android-pad-history-stack" : "history-rail-stack-compact"}`}>
      <div className={`platform-card history-rail-summary-card border border-black/[0.05] bg-white/70 shadow-[var(--shadow-card)] dark:border-white/[0.06] dark:bg-white/[0.03] ${isAndroidPhone ? "p-2.5" : "p-3.5"} ${isWindows ? "rounded-[12px]" : "rounded-[20px]"}`}>
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1.5">
            <h3 className="text-[11px] font-semibold uppercase tracking-[0.12em] text-zinc-500 dark:text-zinc-200">
              上游
            </h3>
            <span className={`h-1.5 w-1.5 rounded-full ${apiKey && baseURL ? "bg-[var(--accent)] shadow-[0_0_6px_rgb(0_122_255_/_0.55)]" : "bg-red-500"}`} />
            <span className={`text-[11px] font-medium ${apiKey && baseURL ? "text-[var(--accent)]" : "text-red-400"}`}>
              {apiKey && baseURL ? "已配置" : "未配置"}
            </span>
          </div>
          <span className="text-[11px] text-zinc-500 dark:text-zinc-400">当前连接</span>
        </div>

        {profiles.length > 0 ? (
          <div className="mt-3">
            <select
              value={activeProfileId}
              onChange={(e) => {
                const id = e.target.value;
                if (id === "__manage__") {
                  openUpstreamConfig("app");
                  return;
                }
                if (id) void setActiveProfile(id);
              }}
              className={`focus-ring w-full border border-black/[0.08] bg-[var(--surface)] px-3 py-2.5 text-[12px] font-medium text-zinc-800 dark:border-white/[0.08] dark:text-zinc-100 ${isWindows ? "rounded-[8px]" : "rounded-[16px]"}`}
              title="切换上游配置 / 管理"
            >
              {profiles.map((profile) => (
                <option key={profile.id} value={profile.id}>
                  {profile.name} · {profile.apiMode === "responses" ? "Responses" : "Images"}
                </option>
              ))}
              <option value="__manage__">⚙ 管理配置...</option>
            </select>
          </div>
        ) : (
          <p className="mt-3 text-[12px] leading-relaxed text-zinc-500 dark:text-zinc-300">
            还没有上游配置，先建一条再开始生成。
          </p>
        )}

          <div className={`mt-2 flex ${isAndroidPhone ? "gap-1" : "gap-1.5"} ${isMac ? "items-stretch" : ""}`}>
          <button
            onClick={() => openUpstreamConfig("app")}
            className={`platform-action-btn flex-1 inline-flex min-h-[34px] items-center justify-center gap-1.5 border border-black/[0.08] px-3 text-[12px] font-medium text-zinc-700 transition-colors hover:border-[color:var(--accent)]/35 hover:text-[var(--accent)] dark:border-white/[0.08] dark:text-zinc-300 ${isAndroidPhone ? "py-1.5" : isMac ? "py-2.5" : "py-2"} ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
          >
            上游配置
          </button>
          <button
            onClick={testAPIKey}
            disabled={!apiKey.trim() || !baseURL.trim() || isTestingKey}
            title="验证当前配置是否可连通"
            className={`platform-action-btn inline-flex min-h-[34px] min-w-[84px] items-center justify-center gap-1.5 border border-black/[0.08] px-3 text-[12px] font-medium text-zinc-700 transition-colors hover:border-[color:var(--accent)]/35 hover:text-[var(--accent)] disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/[0.08] dark:text-zinc-300 ${isAndroidPhone ? "py-1.5" : isMac ? "py-2.5" : "py-2"} ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
          >
            {isTestingKey ? "检查中..." : isAndroidPhone ? "连通性" : "测试"}
          </button>
        </div>

        {!isAndroidPhone ? (
          <div className="mt-2 flex items-center justify-between gap-2">
            <p className="min-w-0 text-[11px] leading-relaxed text-zinc-500 dark:text-zinc-300">
              {apiMode === "responses" ? "Responses API" : "Images API"}
            </p>
          </div>
        ) : null}
      </div>

      <div className={`platform-card history-rail-summary-card border border-black/[0.05] bg-white/70 shadow-[var(--shadow-card)] dark:border-white/[0.06] dark:bg-white/[0.03] ${isAndroidPhone ? "p-2.5" : "p-3.5"} ${isWindows ? "rounded-[12px]" : "rounded-[20px]"}`}>
        <div className={`flex items-center justify-between ${isMac ? "gap-2.5" : "gap-2"}`}>
          <h3 className={`${isMac ? "text-[12px]" : "text-[11px]"} font-semibold uppercase tracking-[0.12em] text-zinc-500 dark:text-zinc-200`}>
            历史 <span className="font-mono-token text-zinc-500 dark:text-zinc-400">({filtered.length}{filtered.length !== history.length && `/${history.length}`})</span>
          </h3>
          <div className={`history-rail-header-actions flex items-center ${isMac ? "gap-1.5 flex-wrap justify-end" : "gap-2"} shrink-0`}>
            {showPhoneFilterToggle ? (
              <button
                type="button"
                onClick={() => setFiltersOpen((v) => !v)}
                className={`platform-pill inline-flex min-h-[30px] items-center gap-1.5 px-2.5 py-1 text-[11px] font-medium transition-colors ${
                  filtersOpen || historyFiltersActive
                    ? "bg-[var(--accent-soft)] text-[var(--accent)] ring-1 ring-[color:var(--accent)]/20"
                    : "text-zinc-500 hover:bg-[var(--accent-soft)] hover:text-[var(--accent)] dark:text-zinc-300"
                } ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
              >
                <Filter className="h-3 w-3" /> 筛选
              </button>
            ) : null}
            {!isAndroidPhone && filtered.length > 6 ? (
              <button
                type="button"
                onClick={openHistoryTimeline}
                className={`history-rail-header-btn platform-pill inline-flex min-h-[30px] items-center justify-center gap-1.5 ${isMac ? "min-w-[78px] px-2.5 py-1.5 text-[12px]" : "px-2.5 py-1 text-[11px]"} font-medium text-zinc-500 transition-colors hover:text-[var(--accent)] ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
              >
                <GalleryVerticalEnd className="h-3 w-3" /> 更多
              </button>
            ) : null}
            {!isAndroidPhone ? (
              <button
                type="button"
                onClick={() => setHistoryRailCollapsed(!historyRailCollapsed)}
                className={`history-rail-header-btn platform-pill inline-flex min-h-[30px] items-center justify-center gap-1.5 ${isMac ? "min-w-[78px] px-2.5 py-1.5 text-[12px]" : "px-2.5 py-1 text-[11px]"} font-medium text-zinc-500 transition-colors hover:text-[var(--accent)] ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
              >
                {historyRailCollapsed ? <ChevronRight className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
                {historyRailCollapsed ? "展开" : "折叠"}
              </button>
            ) : null}
          </div>
        </div>

        {showFilterControls && (
          <>
            <input
              placeholder="搜索 prompt..."
              value={q}
              onChange={(e) => setQ(e.target.value)}
              className={`focus-ring ${isAndroidPhone ? "mt-1.5" : "mt-3"} w-full border border-black/[0.08] bg-[var(--surface)] px-3 py-2.5 text-[12px] text-zinc-900 placeholder:text-zinc-400 dark:border-white/[0.08] dark:text-zinc-100 dark:placeholder:text-zinc-500 ${isWindows ? "rounded-[10px]" : "rounded-[16px]"}`}
            />
            <div className={`mt-2 flex ${isAndroidPhone ? "gap-1" : "gap-1.5"}`}>
              <select
                value={modeF}
                onChange={(e) => setModeF(e.target.value as ModeFilter)}
                className={`focus-ring flex-1 border border-black/[0.08] bg-[var(--surface)] px-3 ${isAndroidPhone ? "py-1.5" : "py-2.5"} text-[12px] text-zinc-700 dark:border-white/[0.08] dark:text-zinc-300 ${isWindows ? "rounded-[10px]" : "rounded-[16px]"}`}
              >
                <option value="all">全部模式</option>
                <option value="generate">文生图</option>
                <option value="edit">图生图</option>
              </select>
              <select
                value={dateF}
                onChange={(e) => setDateF(e.target.value as DateFilter)}
                className={`focus-ring flex-1 border border-black/[0.08] bg-[var(--surface)] px-3 ${isAndroidPhone ? "py-1.5" : "py-2.5"} text-[12px] text-zinc-700 dark:border-white/[0.08] dark:text-zinc-300 ${isWindows ? "rounded-[10px]" : "rounded-[16px]"}`}
              >
                <option value="all">全部日期</option>
                <option value="today">今天</option>
                <option value="week">本周</option>
              </select>
            </div>
          </>
        )}

        {!isAndroidPhone && !isMac && (
          <p className="mt-2 text-[11px] leading-relaxed text-zinc-500 dark:text-zinc-300">
            {isAndroidPad
              ? "点缩略图查看，Shift 可对比，双击可设为源图。"
              : "点击查看 · Shift+点击对比 · 双击设源图 · 更多菜单"}
          </p>
        )}

        {isAndroidPad && filtered.length > 0 && (
          <p className="mt-2 text-[11px] leading-relaxed text-zinc-500 dark:text-zinc-400">
            历史单独收纳，回溯参数、继续变体都从这里进入。
          </p>
        )}
      </div>

      {compareB && (
        <button
          onClick={() => setCompareB(null)}
          className={`platform-pill inline-flex items-center justify-center gap-1.5 border border-[color:var(--accent)]/20 bg-[var(--accent-soft)] px-2.5 py-2 text-xs text-[var(--accent)] transition-colors hover:opacity-90 ${isWindows ? "rounded-[8px]" : "rounded-full"}`}
        >
          <Split className="w-3 h-3" /> 退出对比
        </button>
      )}

      {isMac && !historyRailCollapsed && recentHistory.length > 0 ? (
        <p className="text-[11px] leading-relaxed text-zinc-500 dark:text-zinc-400">
          单击查看 · 双击设源图 · Shift+点击对比 · 更多菜单查看完整操作
        </p>
      ) : null}

      {historyRailCollapsed ? null : recentHistory.length === 0 ? (
        <div className={`platform-card border border-black/[0.05] bg-white/70 text-center text-[12px] text-zinc-500 shadow-[var(--shadow-card)] dark:border-white/[0.06] dark:bg-white/[0.03] dark:text-zinc-300 ${isAndroidPhone ? "py-4" : "py-8"} ${isWindows ? "rounded-[12px]" : "rounded-[20px]"}`}>
          {q || modeF !== "all" || dateF !== "all" ? "没有匹配项" : "还没有结果"}
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-2.5">
          {recentHistory.map((h) => (
            <HistoryTile
              key={h.id}
              item={h}
              isCurrent={currentImage?.id === h.id}
              isCompare={compareB?.id === h.id}
              onSelect={selectCurrent}
              onToggleCompare={(next) => setCompareB(next)}
              onReuse={reuseAsSource}
              onDelete={deleteHistoryItem}
              onOpenMenu={(x, y) => openMenu(h, x, y)}
            />
          ))}
        </div>
      )}

      {menu && <ContextMenu x={menu.x} y={menu.y} items={buildMenu(menu.item)} onClose={closeMenu} />}
      {rawPath && <RawResponseModal path={rawPath} onClose={closeRaw} />}
      </div>
    </aside>
  );
}
