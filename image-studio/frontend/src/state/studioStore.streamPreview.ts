import { base64ToBlob } from "../lib/images.ts";
import type {
  HistoryItem,
  Mode,
  OutputFormatValue,
  QualityValue,
  SizeValue,
  StreamPreview,
  StreamPreviewMap,
  Workspace,
} from "../types/domain.ts";
import type { StudioState } from "./studioStore.types.ts";
import {
  activeRuntimePatch,
  patchWorkspaceRuntime,
  type WorkspacePatch,
} from "./workspaceRuntime.ts";

export type StreamPreviewPayload = {
  imageB64?: string;
  revisedPrompt?: string;
  partialImageIndex?: number;
  mode?: string;
  prompt?: string;
};

export type StreamPreviewSnapshot = {
  workspaceId: string;
  mode: Mode;
  prompt: string;
  size: SizeValue;
  quality: QualityValue;
  outputFormat: OutputFormatValue;
  currentImage: HistoryItem | null;
  batchIndex?: number;
};

export function latestStreamPreview(previews: StreamPreviewMap | null | undefined): StreamPreview | null {
  const list = Object.values(previews ?? {});
  if (list.length === 0) return null;
  return list.reduce((latest, item) => (
    item.updatedAt >= latest.updatedAt ? item : latest
  ));
}

export function streamPreviewItemFromPayload(
  jobId: string,
  payload: StreamPreviewPayload,
  snapshot: StreamPreviewSnapshot,
): HistoryItem | null {
  const imageB64 = typeof payload.imageB64 === "string" ? payload.imageB64.trim() : "";
  if (!imageB64) return null;
  const mode: Mode = payload.mode === "edit" || payload.mode === "generate"
    ? payload.mode
    : snapshot.mode;
  return {
    id: `preview-${jobId}`,
    imageB64,
    imageBlob: base64ToBlob(imageB64),
    prompt: payload.prompt || snapshot.prompt,
    revisedPrompt: payload.revisedPrompt || undefined,
    mode,
    size: snapshot.size,
    quality: snapshot.quality,
    outputFormat: snapshot.outputFormat,
    parentId: mode === "edit" ? snapshot.currentImage?.savedPath : undefined,
    createdAt: Date.now(),
    batchIndex: snapshot.batchIndex,
    previewOnly: true,
  };
}

export function streamPreviewStatePatch(
  state: StudioState,
  jobId: string,
  payload: StreamPreviewPayload,
  snapshot: StreamPreviewSnapshot,
): Partial<StudioState> | null {
  const item = streamPreviewItemFromPayload(jobId, payload, snapshot);
  if (!item) return null;
  const workspace = state.workspaces.find((w) => w.id === snapshot.workspaceId);
  const previousPreviews = state.activeWorkspaceId === snapshot.workspaceId
    ? state.streamPreviews
    : workspace?.streamPreviews ?? {};
  const nextPreview: StreamPreview = {
    jobId,
    imageB64: item.imageB64,
    revisedPrompt: item.revisedPrompt,
    partialImageIndex: payload.partialImageIndex,
    batchIndex: snapshot.batchIndex,
    updatedAt: Date.now(),
  };
  const streamPreviews = { ...previousPreviews, [jobId]: nextPreview };
  const streamPreview = latestStreamPreview(streamPreviews);
  const jobsTotal = state.activeWorkspaceId === snapshot.workspaceId
    ? state.jobsTotal
    : workspace?.jobsTotal ?? 0;
  const useGridPreview = jobsTotal > 1;
  const patch: WorkspacePatch = {
    streamPreview,
    streamPreviews,
  };
  return {
    workspaces: patchWorkspaceRuntime(state.workspaces, snapshot.workspaceId, patch),
    ...(state.activeWorkspaceId === snapshot.workspaceId
      ? useGridPreview
        ? {
            ...activeRuntimePatch(patch),
            compareB: null,
            resultGridOpen: true,
            annotations: [],
            strokes: [],
            maskDataURL: null,
            tool: "pan",
          }
        : {
            ...activeRuntimePatch(patch),
            currentImage: item,
            compareB: null,
            resultGridOpen: false,
            annotations: [],
            strokes: [],
            maskDataURL: null,
            tool: "pan",
          }
      : {}),
  } as Partial<StudioState>;
}

export function restoreCurrentImageAfterPreviewError(
  state: StudioState,
  jobId: string,
  snapshot: StreamPreviewSnapshot,
): HistoryItem | null {
  return state.currentImage?.id === `preview-${jobId}` ? snapshot.currentImage : state.currentImage;
}

export function removeStreamPreview(
  previews: StreamPreviewMap | null | undefined,
  jobId: string,
): { streamPreviews: StreamPreviewMap; streamPreview: StreamPreview | null } {
  const streamPreviews = { ...(previews ?? {}) };
  delete streamPreviews[jobId];
  return {
    streamPreviews,
    streamPreview: latestStreamPreview(streamPreviews),
  };
}

export function streamPreviewItemFromWorkspace(
  workspace: Workspace,
  currentImage: HistoryItem | null,
): HistoryItem | null {
  const preview = workspace.streamPreview ?? latestStreamPreview(workspace.streamPreviews);
  if (!preview) return null;
  return streamPreviewItemFromPayload(preview.jobId, {
    imageB64: preview.imageB64,
    revisedPrompt: preview.revisedPrompt,
    partialImageIndex: preview.partialImageIndex,
    mode: workspace.mode,
    prompt: workspace.prompt,
  }, {
    workspaceId: workspace.id,
    mode: workspace.mode,
    prompt: workspace.prompt,
    size: workspace.size,
    quality: workspace.quality,
    outputFormat: workspace.outputFormat,
    currentImage,
    batchIndex: preview.batchIndex,
  });
}

export function streamPreviewItemsFromPreviews(
  previews: StreamPreviewMap | null | undefined,
  snapshot: StreamPreviewSnapshot,
): HistoryItem[] {
  return Object.values(previews ?? {})
    .sort((a, b) => a.updatedAt - b.updatedAt)
    .map((preview) => streamPreviewItemFromPayload(preview.jobId, {
      imageB64: preview.imageB64,
      revisedPrompt: preview.revisedPrompt,
      partialImageIndex: preview.partialImageIndex,
      mode: snapshot.mode,
      prompt: snapshot.prompt,
    }, { ...snapshot, batchIndex: preview.batchIndex ?? snapshot.batchIndex }))
    .filter((item): item is HistoryItem => !!item);
}

export function currentImageIdForWorkspaceSnapshot(
  currentImage: HistoryItem | null,
  streamPreview: StreamPreview | null,
  streamPreviewsOrFallback: StreamPreviewMap | string | null = {},
  fallbackCurrentImageId?: string | null,
): string | null {
  const streamPreviews = typeof streamPreviewsOrFallback === "string" || streamPreviewsOrFallback === null
    ? {}
    : streamPreviewsOrFallback;
  const fallback = typeof streamPreviewsOrFallback === "string" || streamPreviewsOrFallback === null
    ? streamPreviewsOrFallback
    : fallbackCurrentImageId ?? null;
  if (streamPreview?.jobId && currentImage?.id === `preview-${streamPreview.jobId}`) {
    return fallback;
  }
  if (currentImage?.id && Object.values(streamPreviews).some((preview) => currentImage.id === `preview-${preview.jobId}`)) {
    return fallback;
  }
  return currentImage?.id ?? null;
}
