import { useState } from "react";
import { Plus, X } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";

// Browser-tab style strip. 每个 tab = 独立 workspace,历史栏共享。
// 单 workspace 时不显示。
export function WorkspaceBar() {
  const { workspaces, activeWorkspaceId, newWorkspace, switchWorkspace, closeWorkspace, renameWorkspace, fullscreen } = useStudioStore();
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState("");

  if (fullscreen) return null;
  if (workspaces.length <= 1) return null;

  function startRename(id: string, currentName: string) {
    setEditingId(id);
    setEditingName(currentName);
  }
  function commitRename() {
    if (editingId) {
      renameWorkspace(editingId, editingName.trim() || "未命名");
    }
    setEditingId(null);
  }

  return (
    <div className="flex items-center gap-1 px-3 py-1.5 bg-white/70 dark:bg-zinc-950/70 backdrop-blur-xl border-b border-black/[0.08] dark:border-white/[0.06] overflow-x-auto">
      {workspaces.map((w) => {
        const active = w.id === activeWorkspaceId;
        const isEditing = editingId === w.id;
        return (
          <div
            key={w.id}
            onClick={() => !isEditing && switchWorkspace(w.id)}
            onDoubleClick={() => startRename(w.id, w.name)}
            title="双击重命名"
            className={
              "group flex items-center gap-1.5 px-2.5 h-7 rounded-md text-xs transition-colors cursor-pointer shrink-0 " +
              (active
                ? "bg-emerald-500/12 text-emerald-300 ring-1 ring-emerald-500/30"
                : "text-zinc-600 dark:text-zinc-400 hover:bg-black/5 dark:hover:bg-white/5 hover:text-zinc-900 dark:hover:text-zinc-200")
            }
          >
            {isEditing ? (
              <input
                className="bg-transparent outline-none w-24 text-xs"
                value={editingName}
                autoFocus
                onChange={(e) => setEditingName(e.target.value)}
                onBlur={commitRename}
                onKeyDown={(e) => {
                  if (e.key === "Enter") commitRename();
                  if (e.key === "Escape") setEditingId(null);
                }}
              />
            ) : (
              <span className="truncate max-w-[120px]">{w.name}</span>
            )}
            {!isEditing && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  closeWorkspace(w.id);
                }}
                title="关闭"
                className="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 -mr-1 rounded hover:bg-black/10 dark:hover:bg-white/10"
              >
                <X className="w-3 h-3" />
              </button>
            )}
          </div>
        );
      })}
      <button
        type="button"
        onClick={() => newWorkspace()}
        title="新建标签页"
        className="w-7 h-7 rounded-md flex items-center justify-center text-zinc-500 hover:bg-black/5 dark:hover:bg-white/5 hover:text-zinc-900 dark:hover:text-zinc-200 transition-colors shrink-0"
      >
        <Plus className="w-3.5 h-3.5" />
      </button>
    </div>
  );
}
