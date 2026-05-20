import { ReactNode, useEffect } from "react";
import { X } from "lucide-react";

// 居中 modal:点击背景 / Esc 关闭。
export function Modal({
  open, onClose, title, children, width = 480,
}: {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  width?: number;
}) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;
  return (
    <div
      className="fixed inset-0 z-[9100] flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        style={{ width }}
        onClick={(e) => e.stopPropagation()}
        className="max-h-[85vh] flex flex-col rounded-2xl bg-white dark:bg-zinc-900 ring-1 ring-black/[0.08] dark:ring-white/[0.06] shadow-2xl overflow-hidden"
      >
        {title && (
          <div className="flex items-center justify-between px-5 py-3.5 border-b border-black/[0.06] dark:border-white/[0.04]">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100 m-0">{title}</h3>
            <button
              onClick={onClose}
              title="关闭 (Esc)"
              className="p-1 -mr-1 rounded text-zinc-500 hover:bg-black/5 dark:hover:bg-white/5 hover:text-zinc-900 dark:hover:text-zinc-100"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        )}
        <div className="p-5 overflow-y-auto">{children}</div>
      </div>
    </div>
  );
}
