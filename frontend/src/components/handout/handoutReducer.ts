// Pure state machine for the handout collection. No React, no localStorage —
// every transition is a pure function of (state, action), so the tricky bits
// (deleting the active handout, toggling a problem when none is active) are
// individually testable and can't fall out of sync with the stored value.
import { Handout } from "./handoutStorage";

export interface HandoutState {
  handouts: Handout[];
  activeId: string | null;
}

export type HandoutAction =
  | { type: "set-active"; id: string }
  // The caller mints the Handout (ids need impure randomness); the reducer only
  // decides where it lands and whether it becomes active.
  | { type: "create"; handout: Handout }
  | { type: "delete"; id: string }
  | { type: "update"; id: string; patch: Partial<Handout> }
  | { type: "add-problem"; handoutId: string; problemId: string }
  | { type: "remove-problem"; handoutId: string; problemId: string }
  | { type: "move-problem"; handoutId: string; from: number; to: number }
  // Toggle a problem in the active handout. `fallback` is a ready-made handout
  // (already containing the problem) used only when nothing is active yet.
  | { type: "toggle-active-problem"; problemId: string; fallback: Handout };

function mapHandout(
  handouts: Handout[],
  id: string,
  fn: (h: Handout) => Handout
): Handout[] {
  return handouts.map((h) => (h.id === id ? fn(h) : h));
}

export function handoutReducer(
  state: HandoutState,
  action: HandoutAction
): HandoutState {
  switch (action.type) {
    case "set-active":
      return state.activeId === action.id
        ? state
        : { ...state, activeId: action.id };

    case "create":
      return {
        handouts: [...state.handouts, action.handout],
        activeId: action.handout.id,
      };

    case "delete": {
      const handouts = state.handouts.filter((h) => h.id !== action.id);
      // Re-point active in the same transition rather than via a nested setState.
      const activeId =
        state.activeId === action.id
          ? handouts[0]?.id ?? null
          : state.activeId;
      return { handouts, activeId };
    }

    case "update":
      return {
        ...state,
        handouts: mapHandout(state.handouts, action.id, (h) => ({
          ...h,
          ...action.patch,
        })),
      };

    case "add-problem":
      return {
        ...state,
        handouts: mapHandout(state.handouts, action.handoutId, (h) =>
          h.ids.includes(action.problemId)
            ? h
            : { ...h, ids: [...h.ids, action.problemId] }
        ),
      };

    case "remove-problem":
      return {
        ...state,
        handouts: mapHandout(state.handouts, action.handoutId, (h) => ({
          ...h,
          ids: h.ids.filter((x) => x !== action.problemId),
        })),
      };

    case "move-problem":
      return {
        ...state,
        handouts: mapHandout(state.handouts, action.handoutId, (h) => {
          const { from, to } = action;
          if (
            from === to ||
            from < 0 ||
            to < 0 ||
            from >= h.ids.length ||
            to >= h.ids.length
          )
            return h;
          const ids = [...h.ids];
          const [moved] = ids.splice(from, 1);
          ids.splice(to, 0, moved);
          return { ...h, ids };
        }),
      };

    case "toggle-active-problem": {
      const active = state.handouts.find((h) => h.id === state.activeId);
      if (!active) {
        // Nothing to add to yet — adopt the caller's pre-built handout.
        return {
          handouts: [...state.handouts, action.fallback],
          activeId: action.fallback.id,
        };
      }
      const has = active.ids.includes(action.problemId);
      return {
        ...state,
        handouts: mapHandout(state.handouts, active.id, (h) => ({
          ...h,
          ids: has
            ? h.ids.filter((x) => x !== action.problemId)
            : [...h.ids, action.problemId],
        })),
      };
    }

    default:
      return state;
  }
}
