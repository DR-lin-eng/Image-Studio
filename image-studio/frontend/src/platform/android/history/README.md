# Android History

This folder owns Android phone-only history UI.

- `AndroidHistoryTile.tsx` is the touch-first result card. It handles long press locally and keeps desktop history tile behavior unchanged.
- `AndroidHistoryActionSheet.tsx` replaces the desktop floating context menu on Android phones with a bottom sheet sized for thumbs.

Keep shared history behavior in `components/history/` when it affects every platform. Phone-only visual layout, touch gestures, and action presentation should stay here and in `_android-history.css`, scoped by `html[data-target-platform="android"]`.
