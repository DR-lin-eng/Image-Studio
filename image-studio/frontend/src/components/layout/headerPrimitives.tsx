import { usePlatform } from "../../platform/context";

export function HeaderIconBtn({
  children,
  onClick,
  title,
}: {
  children: React.ReactNode;
  onClick: () => void;
  title: string;
}) {
  const { isWindows, usesAndroidUI } = usePlatform();
  return (
    <button
      type="button"
      onClick={onClick}
      title={title}
      className={`platform-icon-btn no-drag relative flex items-center justify-center text-zinc-600 transition-colors hover:bg-black/[0.05] hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-white/[0.06] dark:hover:text-zinc-100 ${
        usesAndroidUI
          ? "h-10 w-10 rounded-full border border-black/[0.06] bg-[var(--surface)] dark:border-white/[0.08] dark:bg-white/[0.05]"
          : isWindows
            ? "h-8 w-8 rounded-[8px]"
            : "h-8 w-8 rounded-full"
      }`}
    >
      {children}
    </button>
  );
}

export function HeaderToggleBtn({
  active,
  children,
  onClick,
  title,
}: {
  active: boolean;
  children: React.ReactNode;
  onClick: () => void;
  title: string;
}) {
  const { isWindows } = usePlatform();
  return (
    <button
      type="button"
      onClick={onClick}
      title={title}
      className={`platform-chip no-drag flex h-7 w-7 items-center justify-center transition-all ${
        active
          ? "active bg-white text-zinc-900 shadow-sm dark:bg-zinc-900 dark:text-zinc-100"
          : "text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
      } ${isWindows ? "rounded-[7px]" : "rounded-full"}`}
    >
      {children}
    </button>
  );
}
