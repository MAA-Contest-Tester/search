// All persistence for handouts lives here: the on-disk shape, validation of
// whatever we read back, migration from the pre-v2 single-handout keys, and
// guarded access so a hostile/disabled localStorage can never crash the app.
// The React layer (handouts.tsx) and the pure reducer (handoutReducer.ts) never
// touch localStorage directly.

export interface Handout {
  id: string;
  title: string;
  author: string;
  description: string;
  hideSource: boolean;
  // Ordered list of problem ids that make up the handout.
  ids: string[];
}

const STORAGE_KEY = "handouts_v2";
const ACTIVE_KEY = "handouts_active";
const LEGACY_KEYS = [
  "handout_ids",
  "handout_title",
  "handout_author",
  "handout_desc",
] as const;

export function newId(): string {
  const c = globalThis.crypto as
    | (Crypto & { randomUUID?: () => string })
    | undefined;
  // randomUUID only exists in secure contexts; plain-http LAN dev falls back.
  if (c?.randomUUID) return c.randomUUID();
  return (
    "h-" + Math.random().toString(36).slice(2) + "-" + Date.now().toString(36)
  );
}

export function emptyHandout(partial?: Partial<Handout>): Handout {
  const base: Handout = {
    id: newId(),
    title: "",
    author: "",
    description: "",
    hideSource: false,
    ids: [],
  };
  // Only copy defined keys so an explicit `undefined` can't blow away a default.
  if (partial) {
    for (const k of Object.keys(partial) as (keyof Handout)[]) {
      if (partial[k] !== undefined) (base as any)[k] = partial[k];
    }
  }
  return base;
}

// --- guarded localStorage access -------------------------------------------

function safeGet(key: string): string | null {
  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
}

function safeSet(key: string, value: string): void {
  try {
    localStorage.setItem(key, value);
  } catch {
    // storage disabled or over quota — persistence is best-effort
  }
}

function safeRemove(key: string): void {
  try {
    localStorage.removeItem(key);
  } catch {
    // ignore
  }
}

// --- validation ------------------------------------------------------------

function uniqueStrings(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  const seen = new Set<string>();
  const out: string[] = [];
  for (const x of value) {
    if (typeof x === "string" && x !== "" && !seen.has(x)) {
      seen.add(x);
      out.push(x);
    }
  }
  return out;
}

// Coerce one persisted entry into a well-formed Handout, or drop it (null) if
// it isn't even minimally usable. Every field is checked — a corrupt `ids`
// becomes `[]` rather than crashing a later `.map`/`.includes`.
function normalizeHandout(raw: unknown): Handout | null {
  if (!raw || typeof raw !== "object") return null;
  const h = raw as Record<string, unknown>;
  if (typeof h.id !== "string" || h.id === "") return null;
  return {
    id: h.id,
    title: typeof h.title === "string" ? h.title : "",
    author: typeof h.author === "string" ? h.author : "",
    description: typeof h.description === "string" ? h.description : "",
    hideSource: h.hideSource === true,
    ids: uniqueStrings(h.ids),
  };
}

function dedupeById(handouts: Handout[]): Handout[] {
  const seen = new Set<string>();
  return handouts.filter((h) => {
    if (seen.has(h.id)) return false;
    seen.add(h.id);
    return true;
  });
}

// --- load / save -----------------------------------------------------------

export interface LoadResult {
  handouts: Handout[];
  activeId: string | null;
  // True only when the state was reconstructed from the legacy keys, so the
  // caller can clear them once and never migrate again.
  migrated: boolean;
}

function migrateLegacy(): LoadResult {
  const rawIds = (safeGet("handout_ids") || "").trim();
  const title = safeGet("handout_title") || "";
  const author = safeGet("handout_author") || "";
  const description = safeGet("handout_desc") || "";
  if (!rawIds && !title && !author && !description) {
    return { handouts: [], activeId: null, migrated: false };
  }
  const migrated = emptyHandout({
    title,
    author,
    description,
    ids: rawIds ? uniqueStrings(rawIds.split(/\s+/)) : [],
  });
  return { handouts: [migrated], activeId: migrated.id, migrated: true };
}

export function loadHandouts(): LoadResult {
  const raw = safeGet(STORAGE_KEY);
  // Key absent → first run on this version: try the legacy single-handout keys.
  if (raw === null) return migrateLegacy();

  // Key present (even as "[]") is authoritative: an intentionally-empty set
  // must stay empty and must NOT fall back to migration.
  let parsed: unknown = null;
  try {
    parsed = JSON.parse(raw);
  } catch {
    parsed = null;
  }
  if (!Array.isArray(parsed)) return { handouts: [], activeId: null, migrated: false };

  const handouts = dedupeById(
    parsed
      .map(normalizeHandout)
      .filter((h): h is Handout => h !== null)
  );
  const storedActive = safeGet(ACTIVE_KEY);
  const activeId =
    handouts.find((h) => h.id === storedActive)?.id ?? handouts[0]?.id ?? null;
  return { handouts, activeId, migrated: false };
}

export function saveHandouts(handouts: Handout[], activeId: string | null): void {
  safeSet(STORAGE_KEY, JSON.stringify(handouts));
  if (activeId) safeSet(ACTIVE_KEY, activeId);
  else safeRemove(ACTIVE_KEY);
}

export function clearLegacyKeys(): void {
  LEGACY_KEYS.forEach(safeRemove);
}
