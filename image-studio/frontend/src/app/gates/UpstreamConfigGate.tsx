import { lazy, Suspense } from "react";
import { useStudioStore } from "../../state/studioStore";

const UpstreamConfigModal = lazy(() => import("../../components/panel/UpstreamConfigModal").then((module) => ({ default: module.UpstreamConfigModal })));

export function UpstreamConfigGate() {
  const open = useStudioStore((state) => state.upstreamModalOpen);
  const close = useStudioStore((state) => state.closeUpstreamConfig);

  if (!open) return null;

  return (
    <Suspense fallback={null}>
      <UpstreamConfigModal open={open} onClose={close} />
    </Suspense>
  );
}
