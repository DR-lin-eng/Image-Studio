import { base64ToBlob } from "../lib/images.ts";
import type {
  HistoryItem,
  Mode,
  OutputFormatValue,
  QualityValue,
  SizeValue,
  StreamPreview,
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
};

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
  const patch: WorkspacePatch = {
    streamPreview: {
      jobId,
      imageB64: item.imageB64,
      revisedPrompt: item.revisedPrompt,
      partialImageIndex: payload.partialImageIndex,
      updatedAt: Date.now(),
    },
  };
  return {
    workspaces: patchWorkspaceRuntime(state.workspaces, snapshot.workspaceId, patch),
    ...(state.activeWorkspaceId === snapshot.workspaceId
      ? {
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
  return state.streamPreview?.jobId === jobId ? snapshot.currentImage : state.currentImage;
}

export function streamPreviewItemFromWorkspace(
  workspace: Workspace,
  currentImage: HistoryItem | null,
): HistoryItem | null {
  const preview = workspace.streamPreview;
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
  });
}

export function currentImageIdForWorkspaceSnapshot(
  currentImage: HistoryItem | null,
  streamPreview: StreamPreview | null,
  fallbackCurrentImageId: string | null,
): string | null {
  if (streamPreview?.jobId && currentImage?.id === `preview-${streamPreview.jobId}`) {
    return fallbackCurrentImageId;
  }
  return currentImage?.id ?? null;
}
