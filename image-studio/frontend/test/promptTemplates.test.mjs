import assert from "node:assert/strict";
import test from "node:test";

const promptTemplates = await import("../src/lib/promptTemplates.ts");

test("normalizePromptTemplates keeps only valid templates", () => {
  const normalized = promptTemplates.normalizePromptTemplates([
    { id: "a", label: "模板1", text: "cat", createdAt: 1, updatedAt: 2 },
    { id: "", label: "bad", text: "x" },
    { id: "b", label: " ", text: "x" },
  ]);
  assert.deepEqual(normalized, [
    { id: "a", label: "模板1", text: "cat", createdAt: 1, updatedAt: 2 },
  ]);
});

test("nextDefaultPromptTemplateLabel uses first available numeric slot", () => {
  const label = promptTemplates.nextDefaultPromptTemplateLabel([
    { id: "1", label: "模板1", text: "a", createdAt: 1, updatedAt: 1 },
    { id: "3", label: "模板3", text: "b", createdAt: 1, updatedAt: 1 },
  ]);
  assert.equal(label, "模板2");
});
