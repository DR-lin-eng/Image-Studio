export type GenerateOptionsLike = {
  apiKey: string;
  mode: string;
  prompt: string;
  size: string;
  quality: string;
  outputFormat: string;
  imagePaths: string[];
  imagePath: string;
  maskB64: string;
  seed: number;
  negativePrompt: string;
  baseURL: string;
  textModelID: string;
  imageModelID: string;
  apiMode: string;
  requestPolicy: string;
  noPromptRevision: boolean;
  concurrencyLimit?: number;
  requestedJobId?: string;
};

export type PromptOptimizeOptionsLike = {
  apiKey: string;
  prompt: string;
  mode: string;
  baseURL: string;
  textModelID: string;
  imagePaths: string[];
  imagePath: string;
};

export type JobStartedLike = { jobId: string };
export type ImportedImageLike = { path: string; imageB64: string };
export type ImageTransformResultLike = { path: string; acceleration?: string };
export type SelectFileResponseLike = { path: string; size: number; imageB64?: string };

export type HostKind = "wails-desktop" | "android-shell" | "browser";

export type HostCapabilities = {
  localGeneration: boolean;
  promptOptimization: boolean;
  nativeFileDialogs: boolean;
  nativeImageTransforms: boolean;
  nativeHistoryFileIO: boolean;
  nativeOutputDirectoryPicker: boolean;
  secureCredentialStore: boolean;
};

export type KernelRuntimeMode = "auto" | "local" | "remote";
