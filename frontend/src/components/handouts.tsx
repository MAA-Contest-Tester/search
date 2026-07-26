import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

export interface Handout {
  id: string;
  title: string;
  author: string;
  description: string;
  hideSource: boolean;
  // Ordered list of problem ids that make up the handout.
  ids: string[];
}

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

const STORAGE_KEY = "handouts_v2";
const ACTIVE_KEY = "handouts_active";

function newId(): string {
  const c = globalThis.crypto as (Crypto & { randomUUID?: () => string }) | undefined;
  if (c?.randomUUID) return c.randomUUID();
  // Fallback for older browsers — uniqueness only needs to hold locally.
  return "h-" + Math.random().toString(36).slice(2) + "-" + Date.now().toString(36);
}

export function emptyHandout(partial?: Partial<Handout>): Handout {
  return {
    id: newId(),
    title: "",
    author: "",
    description: "",
    hideSource: false,
    ids: [],
    ...partial,
  };
}

// Load persisted handouts, migrating from the old single-handout localStorage
// keys (handout_ids / handout_title / handout_author / handout_desc) the first
// time this version runs.
function loadInitial(): { handouts: Handout[]; activeId: string } {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw) as Handout[];
      if (Array.isArray(parsed) && parsed.length) {
        const cleaned = parsed
          .filter((h) => h && typeof h.id === "string")
          .map((h) => ({ ...emptyHandout({ id: h.id }), ...h }));
        if (cleaned.length) {
          const stored = localStorage.getItem(ACTIVE_KEY);
          const activeId =
            cleaned.find((h) => h.id === stored)?.id ?? cleaned[0].id;
          return { handouts: cleaned, activeId };
        }
      }
    }
  } catch {
    // fall through to migration / default
  }

  const legacyIds = (localStorage.getItem("handout_ids") || "").trim();
  const migrated = emptyHandout({
    title: localStorage.getItem("handout_title") || "",
    author: localStorage.getItem("handout_author") || "",
    description: localStorage.getItem("handout_desc") || "",
    ids: legacyIds ? legacyIds.split(/\s+/).filter(Boolean) : [],
  });
  return { handouts: [migrated], activeId: migrated.id };
}

export function HandoutsProvider(props: { children: React.ReactNode }) {
  const [initial] = useState(loadInitial);
  const [handouts, setHandouts] = useState<Handout[]>(initial.handouts);
  const [activeId, setActiveIdState] = useState<string | null>(initial.activeId);

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(handouts));
  }, [handouts]);
  useEffect(() => {
    if (activeId) localStorage.setItem(ACTIVE_KEY, activeId);
  }, [activeId]);

  const activeHandout = useMemo(
    () => handouts.find((h) => h.id === activeId) ?? null,
    [handouts, activeId]
  );

  const value = useMemo<HandoutsContextValue>(() => {
    const setActiveId = (id: string) => setActiveIdState(id);

    const createHandout = (partial?: Partial<Handout>) => {
      const handout = emptyHandout(partial);
      setHandouts((prev) => [...prev, handout]);
      setActiveIdState(handout.id);
      return handout.id;
    };

    const deleteHandout = (id: string) => {
      setHandouts((prev) => {
        const next = prev.filter((h) => h.id !== id);
        setActiveIdState((current) =>
          current === id ? next[0]?.id ?? null : current
        );
        return next;
      });
    };

    const updateHandout = (id: string, patch: Partial<Handout>) =>
      setHandouts((prev) =>
        prev.map((h) => (h.id === id ? { ...h, ...patch } : h))
      );

    const addProblem = (handoutId: string, problemId: string) =>
      setHandouts((prev) =>
        prev.map((h) =>
          h.id === handoutId && !h.ids.includes(problemId)
            ? { ...h, ids: [...h.ids, problemId] }
            : h
        )
      );

    const removeProblem = (handoutId: string, problemId: string) =>
      setHandouts((prev) =>
        prev.map((h) =>
          h.id === handoutId
            ? { ...h, ids: h.ids.filter((x) => x !== problemId) }
            : h
        )
      );

    const moveProblem = (handoutId: string, from: number, to: number) =>
      setHandouts((prev) =>
        prev.map((h) => {
          if (h.id !== handoutId) return h;
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
        })
      );

    const toggleProblemInActive = (problemId: string) => {
      // Adding from search when nothing exists yet spins up a first handout.
      if (!activeHandout) {
        createHandout({ ids: [problemId] });
        return;
      }
      if (activeHandout.ids.includes(problemId)) {
        removeProblem(activeHandout.id, problemId);
      } else {
        addProblem(activeHandout.id, problemId);
      }
    };

    const isInActive = (problemId: string) =>
      !!activeHandout && activeHandout.ids.includes(problemId);

    return {
      handouts,
      activeId,
      activeHandout,
      setActiveId,
      createHandout,
      deleteHandout,
      updateHandout,
      addProblem,
      removeProblem,
      moveProblem,
      toggleProblemInActive,
      isInActive,
    };
  }, [handouts, activeId, activeHandout]);

  return (
    <HandoutsContext.Provider value={value}>
      {props.children}
    </HandoutsContext.Provider>
  );
}
