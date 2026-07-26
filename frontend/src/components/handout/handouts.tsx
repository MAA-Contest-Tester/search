import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useReducer,
  useState,
} from "react";
import {
  Handout,
  clearLegacyKeys,
  emptyHandout,
  loadHandouts,
  saveHandouts,
} from "./handoutStorage";
import { handoutReducer } from "./handoutReducer";

export type { Handout };

interface HandoutsContextValue {
  handouts: Handout[];
  activeId: string | null;
  activeHandout: Handout | null;
  setActiveId: (id: string) => void;
  createHandout: (partial?: Partial<Handout>) => string;
  deleteHandout: (id: string) => void;
  updateHandout: (id: string, patch: Partial<Handout>) => void;
  addProblem: (handoutId: string, problemId: string) => void;
  removeProblem: (handoutId: string, problemId: string) => void;
  moveProblem: (handoutId: string, from: number, to: number) => void;
  // Convenience helpers for the search results, which always target the
  // currently-active handout.
  toggleProblemInActive: (problemId: string) => void;
  isInActive: (problemId: string) => boolean;
}

const HandoutsContext = createContext<HandoutsContextValue | null>(null);

export function useHandouts(): HandoutsContextValue {
  const ctx = useContext(HandoutsContext);
  if (!ctx) {
    throw new Error("useHandouts must be used within a HandoutsProvider");
  }
  return ctx;
}

export function HandoutsProvider(props: { children: React.ReactNode }) {
  // Load once; `migrated` tells us whether to retire the legacy keys.
  const [initial] = useState(loadHandouts);
  const [state, dispatch] = useReducer(handoutReducer, {
    handouts: initial.handouts,
    activeId: initial.activeId,
  });
  const { handouts, activeId } = state;

  useEffect(() => {
    if (initial.migrated) clearLegacyKeys();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  useEffect(() => {
    saveHandouts(handouts, activeId);
  }, [handouts, activeId]);

  const activeHandout = useMemo(
    () => handouts.find((h) => h.id === activeId) ?? null,
    [handouts, activeId]
  );

  const value = useMemo<HandoutsContextValue>(
    () => ({
      handouts,
      activeId,
      activeHandout,
      setActiveId: (id) => dispatch({ type: "set-active", id }),
      createHandout: (partial) => {
        const handout = emptyHandout(partial);
        dispatch({ type: "create", handout });
        return handout.id;
      },
      deleteHandout: (id) => dispatch({ type: "delete", id }),
      updateHandout: (id, patch) => dispatch({ type: "update", id, patch }),
      addProblem: (handoutId, problemId) =>
        dispatch({ type: "add-problem", handoutId, problemId }),
      removeProblem: (handoutId, problemId) =>
        dispatch({ type: "remove-problem", handoutId, problemId }),
      moveProblem: (handoutId, from, to) =>
        dispatch({ type: "move-problem", handoutId, from, to }),
      toggleProblemInActive: (problemId) =>
        dispatch({
          type: "toggle-active-problem",
          problemId,
          fallback: emptyHandout({ ids: [problemId] }),
        }),
      isInActive: (problemId) =>
        !!activeHandout && activeHandout.ids.includes(problemId),
    }),
    [handouts, activeId, activeHandout]
  );

  return (
    <HandoutsContext.Provider value={value}>
      {props.children}
    </HandoutsContext.Provider>
  );
}
