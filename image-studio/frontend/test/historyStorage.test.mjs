import assert from "node:assert/strict";
import test from "node:test";
import "fake-indexeddb/auto";
import * as storage from "../src/lib/storage.ts";
import { mergeHistoryItems } from "../src/lib/history.ts";

const legacy = await new Promise((resolve, reject) => {
  const req = indexedDB.open("keyval-store");
  req.onupgradeneeded = () => req.result.createObjectStore("keyval");
  req.onsuccess = () => resolve(req.result);
  req.onerror = () => reject(req.error);
});

function item(id, createdAt) {
  return { id, createdAt, prompt: `prompt ${id}`, mode: "generate", size: "1024x1024", quality: "high" };
}

test.beforeEach(() => storage.clearHistoryStorage());
test.after(() => legacy.close());

test("hundreds of same-day results with duplicate timestamps remain reachable through bounded pages", async () => {
  const now = new Date(2026, 9, 6, 12).getTime();
  const items = Array.from({ length: 650 }, (_, index) => item(`image-${String(index).padStart(4, "0")}`, now - Math.floor(index / 220) * 86400000));
  await storage.persistHistoryItems(items);
  let cursor = null;
  const loaded = [];
  do {
    const page = await storage.loadHistoryPage({ limit: 18, cursor });
    assert.ok(page.items.length <= 18);
    loaded.push(...page.items);
    cursor = page.nextCursor;
  } while (cursor);
  assert.equal(loaded.length, 650);
  assert.deepEqual(new Set(loaded.map((i) => i.id)), new Set(items.map((i) => i.id)));
  assert.equal(loaded.find((i) => i.id === "image-0649").prompt, "prompt image-0649");
});

test("saving a new result after loading one page preserves older metadata and full images", async () => {
  await storage.persistHistoryItems(Array.from({ length: 300 }, (_, i) => item(`old-${i}`, 1000 + i)));
  await storage.persistHistoryFullImage("old-0", "AQID");
  const page = await storage.loadHistoryPage({ limit: 18 });
  const newest = item("new", 2000);
  const visible = mergeHistoryItems([newest, ...page.items, page.items[0]]);
  assert.equal(visible.length, 19);
  await storage.persistHistoryItem(newest);
  assert.equal((await storage.loadAllHistory()).length, 301);
  assert.equal(await storage.loadHistoryFullImage("old-0"), "AQID");
});

test("date cleanup deletes unloaded old records and their full images while retaining unloaded newer results", async () => {
  const items = Array.from({ length: 300 }, (_, i) => item(`image-${i}`, 1000 + i));
  await storage.persistHistoryItems(items);
  await storage.persistHistoryFullImages([{ id: "image-0", imageB64: "AQID" }, { id: "image-100", imageB64: "BAUG" }]);
  assert.equal((await storage.loadHistoryPage({ limit: 18 })).items.length, 18);
  const removed = await storage.removeHistoryOlderThan(1100);
  assert.equal(removed.length, 100);
  const remaining = await storage.loadAllHistory();
  assert.equal(remaining.length, 200);
  assert.ok(remaining.some((i) => i.id === "image-100"));
  assert.equal(await storage.loadHistoryFullImage("image-0"), "");
  assert.equal(await storage.loadHistoryFullImage("image-100"), "BAUG");
});

test("single history deletion also removes legacy copies so migration cannot restore deleted records", async () => {
  const old = item("legacy", 1000);
  await new Promise((resolve, reject) => {
    const tx = legacy.transaction("keyval", "readwrite");
    tx.objectStore("keyval").put(old, "history:legacy");
    tx.objectStore("keyval").put("AQID", "history-full:legacy");
    tx.oncomplete = resolve;
    tx.onerror = () => reject(tx.error);
  });
  assert.equal((await storage.loadAllHistory()).length, 1);
  await storage.removeHistoryItem("legacy");
  assert.deepEqual(await storage.loadAllHistory(), []);
  assert.equal(await storage.loadHistoryFullImage("legacy"), "");
});

test("history merge keeps every batch result beyond the former 120-item limit", () => {
  const all = Array.from({ length: 650 }, (_, i) => item(`image-${i}`, 1000 + i));
  const merged = mergeHistoryItems([...all, all[0], { ...all[1], prompt: "duplicate" }]);
  assert.equal(merged.length, 650);
  assert.equal(merged[0].id, "image-649");
  assert.equal(merged.at(-1).id, "image-0");
  assert.equal(merged.find((i) => i.id === "image-1").prompt, "prompt image-1");
});
