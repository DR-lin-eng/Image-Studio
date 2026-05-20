import { AlertTriangle, CheckCircle2, Info, X, XCircle } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";
import type { Toast } from "../../types/domain";

export function ToastContainer() {
  const toasts = useStudioStore((s) => s.toasts);
  const dismiss = useStudioStore((s) => s.dismissToast);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed top-14 right-4 z-[9050] flex flex-col gap-2 max-w-sm">
      {toasts.map((t) => (
        <ToastItem key={t.id} t={t} onClose={() => dismiss(t.id)} />
      ))}
    </div>
  );
}

function toneClasses(kind: Toast["kind"]): string {
  // 浅色模式:用更深的文字色 + 浅且不透明的卡片背景,保证可读
  // 深色模式:沿用原来的高透明 tint + 浅色文字
  switch (kind) {
    case "success":
      return "bg-emerald-50 text-emerald-700 ring-emerald-500/40 dark:bg-emerald-500/12 dark:text-emerald-300 dark:ring-emerald-500/30";
    case "error":
      return "bg-red-50 text-red-700 ring-red-500/40 dark:bg-red-500/12 dark:text-red-300 dark:ring-red-500/30";
    case "warn":
      return "bg-amber-50 text-amber-800 ring-amber-500/40 dark:bg-amber-500/12 dark:text-amber-300 dark:ring-amber-500/30";
    default:
      return "bg-blue-50 text-blue-700 ring-blue-500/40 dark:bg-blue-500/12 dark:text-blue-300 dark:ring-blue-500/30";
  }
}

function ToneIcon({ kind }: { kind: Toast["kind"] }) {
  const c = "w-4 h-4 shrink-0";
  switch (kind) {
    case "success": return <CheckCircle2 className={c} />;
    case "error":   return <XCircle className={c} />;
    case "warn":    return <AlertTriangle className={c} />;
    default:        return <Info className={c} />;
  }
}

function ToastItem({ t, onClose }: { t: Toast; onClose: () => void }) {
  return (
    <div
      className={`flex items-center gap-2 px-3 py-2 rounded-xl ring-1 backdrop-blur shadow-2xl animate-[toast-in_180ms_ease-out] ${toneClasses(t.kind)}`}
      style={{ animation: "toast-in 180ms ease-out" }}
    >
      <ToneIcon kind={t.kind} />
      <span
        onClick={onClose}
        className="flex-1 text-xs leading-snug break-words cursor-pointer"
      >
        {t.text}
      </span>
      {t.action && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            t.action!.onClick();
            onClose();
          }}
          className="px-2 py-1 text-[11px] font-medium rounded-md bg-current/15 hover:bg-current/25 transition-colors whitespace-nowrap"
        >
          {t.action.label}
        </button>
      )}
      <button
        type="button"
        onClick={onClose}
        className="opacity-60 hover:opacity-100"
      >
        <X className="w-3 h-3" />
      </button>
      <style>{`@keyframes toast-in { from { opacity: 0; transform: translateY(-8px); } to { opacity: 1; transform: translateY(0); } }`}</style>
    </div>
  );
}
