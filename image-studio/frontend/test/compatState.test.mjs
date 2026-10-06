import assert from "node:assert/strict";
import test from "node:test";
import "fake-indexeddb/auto";
import { clearHistoryStorage, loadAllHistory, persistHistoryItems } from "../src/lib/storage.ts";

const realWindow = globalThis.window;
const realLocalStorage = globalThis.localStorage;

test.beforeEach(() => clearHistoryStorage());

function installStorage() {
  const store = new Map();
  globalThis.localStorage = {
    getItem(key) {
      return store.has(key) ? store.get(key) : null;
    },
    setItem(key, value) {
      store.set(key, String(value));
    },
    removeItem(key) {
      store.delete(key);
    },
  };
  return store;
}

function installService(methods = {}) {
  globalThis.window = {
    go: {
      backend: {
        Service: methods,
      },
    },
  };
}

test.afterEach(() => {
  globalThis.window = realWindow;
  globalThis.localStorage = realLocalStorage;
});

test("compat export preserves previewPath sourcePaths and parentId in history", async () => {
  installStorage();
  await persistHistoryItems(Array.from({ length: 250 }, (_, i) => ({
    id: `unloaded-${i}`, prompt: `older prompt ${i}`, mode: "generate", size: "1024x1024", quality: "high", createdAt: 100,
  })));
  let savedState = null;
  installService({
    SaveCompatibilityState(state) {
      savedState = state;
    },
  });
  const compatState = await import(`../src/lib/compatState.ts?compat-export=${Date.now()}-${Math.random().toString(36).slice(2)}`);

  await compatState.exportCompatibilityStateNow({
    history: [{
      id: "hist-1",
      prompt: "edit cat",
      mode: "edit",
      size: "1024x1024",
      quality: "high",
      outputFormat: "png",
      createdAt: 123,
      savedPath: "/tmp/result.png",
      previewPath: "/tmp/previews/result.png",
      thumbPath: "/tmp/thumbs/result.png",
      parentId: "/tmp/source-a.png",
      sourcePaths: ["/tmp/source-a.png", "/tmp/source-b.png"],
      imageB64: "AAAA",
    }],
    profiles: [],
    activeProfileId: "",
    proxyMode: "system",
    proxyURL: "",
    theme: "system",
    fontScale: 1,
    outputFormat: "png",
    background: "auto",
    outputCompression: 100,
    inputFidelity: "auto",
    imageStyle: "default",
    moderation: "low",
    userIdentifier: "",
    partialImages: 1,
    protectStreamPreview: true,
    autoRetryEnabled: true,
    autoRetryCount: 5,
    promptTemplates: [],
    promptHistory: [],
    presets: [],
    customAspectRatios: [],
    kernelRuntimeMode: "remote",
    keepLogs: false,
    cleanupPreviewCacheOnExit: false,
    ignoredReleaseTag: "",
    completionSound: { enabled: true, mode: "default", customName: "", customDataURL: "" },
    completionNotification: { enabled: false },
  });

  assert.ok(savedState, "compat state should be exported");
  assert.equal(savedState.history.length, 251, "compatibility export must include unloaded history");
  assert.equal(savedState.history[0].previewPath, "/tmp/previews/result.png");
  assert.equal(savedState.history[0].parentId, "/tmp/source-a.png");
  assert.deepEqual(savedState.history[0].sourcePaths, ["/tmp/source-a.png", "/tmp/source-b.png"]);
});

test("importing a partial snapshot from an older client cannot delete the local archive", async () => {
  installStorage();
  await persistHistoryItems(Array.from({ length: 250 }, (_, i) => ({
    id: `old-${i}`, prompt: `older prompt ${i}`, mode: "generate", size: "1024x1024", quality: "high", createdAt: 100,
  })));
  installService({ LoadCompatibilityState() {
    return {
      updatedAt: Date.now(), settings: {}, profiles: [], activeProfileId: "",
      history: [{ id: "new", prompt: "new prompt", mode: "generate", size: "1024x1024", quality: "high", createdAt: 200 }],
    };
  } });
  const compat = await import(`../src/lib/compatState.ts?partial-import=${Date.now()}`);
  assert.equal(await compat.importCompatibilityStateIfNewer(), true);
  const all = await loadAllHistory();
  assert.equal(all.length, 251);
  assert.ok(all.some((item) => item.id === "old-0"));
});

test("compat fingerprint changes when previewPath or sourcePaths change", async () => {
  installStorage();
  installService();
  const compatState = await import(`../src/lib/compatState.ts?compat-fingerprint=${Date.now()}-${Math.random().toString(36).slice(2)}`);

  const base = {
    history: [{
      id: "hist-1",
      prompt: "edit cat",
      mode: "edit",
      size: "1024x1024",
      quality: "high",
      outputFormat: "png",
      createdAt: 123,
      savedPath: "/tmp/result.png",
      previewPath: "/tmp/previews/result-a.png",
      sourcePaths: ["/tmp/source-a.png"],
    }],
    profiles: [],
    activeProfileId: "",
    proxyMode: "system",
    proxyURL: "",
    theme: "system",
    fontScale: 1,
    outputFormat: "png",
    background: "auto",
    outputCompression: 100,
    inputFidelity: "auto",
    imageStyle: "default",
    moderation: "low",
    userIdentifier: "",
    partialImages: 1,
    protectStreamPreview: true,
    autoRetryEnabled: true,
    autoRetryCount: 5,
    promptTemplates: [],
    promptHistory: [],
    presets: [],
    customAspectRatios: [],
    kernelRuntimeMode: "remote",
    keepLogs: false,
    cleanupPreviewCacheOnExit: false,
    ignoredReleaseTag: "",
    completionSound: { enabled: true, mode: "default", customName: "", customDataURL: "" },
    completionNotification: { enabled: false },
  };

  const fingerprintA = compatState.compatibilityExportFingerprint(base);
  const fingerprintB = compatState.compatibilityExportFingerprint({
    ...base,
    history: [{ ...base.history[0], previewPath: "/tmp/previews/result-b.png" }],
  });
  const fingerprintC = compatState.compatibilityExportFingerprint({
    ...base,
    history: [{ ...base.history[0], sourcePaths: ["/tmp/source-a.png", "/tmp/source-b.png"] }],
  });

  assert.notEqual(fingerprintA, fingerprintB);
  assert.notEqual(fingerprintA, fingerprintC);
});
