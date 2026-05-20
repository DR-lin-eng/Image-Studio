import { Github, Moon, Plus, Sun } from "lucide-react";
import { useStudioStore } from "../../state/studioStore";
import { OpenExternalURL } from "../../../wailsjs/go/backend/Service";
import { HitokotoStrip } from "./HitokotoStrip";

const REPO_URL = "https://github.com/RoseKhlifa/Image-Studio";

export function AppHeader() {
  const { fullscreen, theme, setTheme, pushToast, workspaces, newWorkspace } = useStudioStore();
  if (fullscreen) return null;

  return (
    <header className="sticky top-0 z-40 h-9 px-3 flex items-center gap-2 bg-white/70 dark:bg-zinc-950/70 backdrop-blur-xl border-b border-black/[0.08] dark:border-white/[0.06]">
      {/* 左:每日一言。文字 "Image Studio" 不放,Windows 标题栏已经有了。 */}
      <HitokotoStrip />

      <div className="flex items-center gap-1 ml-auto">
      <HeaderIconBtn
        onClick={() => newWorkspace()}
        title={workspaces.length > 1 ? `${workspaces.length} 个标签 · 新建` : "新建标签"}
      >
        <Plus className="w-4 h-4" />
        {workspaces.length > 1 && (
          <span className="absolute -top-0.5 -right-0.5 px-1 min-w-[14px] h-[14px] rounded-full bg-emerald-500 text-[9px] text-zinc-950 font-bold flex items-center justify-center">
            {workspaces.length}
          </span>
        )}
      </HeaderIconBtn>
      <HeaderIconBtn
        onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        title={theme === "dark" ? "浅色主题" : "深色主题"}
      >
        {theme === "dark" ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
      </HeaderIconBtn>
      <HeaderIconBtn
        onClick={() => OpenExternalURL(REPO_URL).catch(() => pushToast("无法打开浏览器", "error"))}
        title="GitHub"
      >
        <Github className="w-4 h-4" />
      </HeaderIconBtn>
      </div>
    </header>
  );
}

function HeaderIconBtn({ children, onClick, title }: {
  children: React.ReactNode;
  onClick: () => void;
  title: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      title={title}
      className="relative w-8 h-8 rounded-md flex items-center justify-center text-zinc-600 dark:text-zinc-400 hover:bg-black/5 dark:hover:bg-white/5 hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors"
    >
      {children}
    </button>
  );
}
