import React, { useEffect, useRef, useState } from "react";
import { flushSync } from "react-dom";
import Result from "../Result";
import { Handout, useHandouts } from "./handouts";

// Fetch the full problem records for a handout's ids, preserving order. The
// backend returns nulls for ids it can't resolve, which we keep so the row
// index stays aligned with handout.ids.
function useHandoutProblems(ids: string[]) {
  const key = ids.join(" ");
  const [problems, setProblems] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (ids.length === 0) {
      setProblems([]);
      setError(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    fetch(`/backend/handout`, {
      headers: { "Content-Type": "application/json" },
      method: "POST",
      body: JSON.stringify({ ids }),
    })
      .then(async (data) => {
        if (cancelled) return;
        setLoading(false);
        if (data.status !== 200) {
          setProblems([]);
          setError(await data.text());
        } else {
          setError(null);
          setProblems(await data.json());
        }
      })
      .catch(() => {
        if (!cancelled) setError("Could not load problems.");
      });
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key]);
  return { problems, loading, error };
}

const iconButton =
  "font-bold rounded-md duration-200 border border-gray-200 hover:bg-blue-800 hover:text-white disabled:opacity-30 disabled:hover:bg-transparent disabled:hover:text-current";

// The ordered, editable list of problems inside a handout: drag-to-reorder,
// up/down for touch/accessibility, and remove. Print-hidden — it's authoring UI.
function ProblemList(props: { handout: Handout; problems: any[] }) {
  const { handout, problems } = props;
  const { moveProblem, removeProblem } = useHandouts();
  const dragFrom = useRef<number | null>(null);
  const [dragOver, setDragOver] = useState<number | null>(null);

  if (handout.ids.length === 0) {
    return (
      <p className="text-sm text-gray-500 my-2">
        No problems yet — add some from the{" "}
        <span className="font-bold">Search</span> page.
      </p>
    );
  }

  return (
    <ol className="my-2 flex flex-col gap-1">
      {handout.ids.map((id, i) => {
        const problem = problems[i];
        const label: string =
          problem && problem.source ? problem.source : `Unknown (${id})`;
        return (
          <li
            key={id}
            draggable
            onDragStart={() => (dragFrom.current = i)}
            onDragOver={(e) => {
              e.preventDefault();
              setDragOver(i);
            }}
            onDragLeave={() => setDragOver((d) => (d === i ? null : d))}
            onDrop={(e) => {
              e.preventDefault();
              if (dragFrom.current !== null)
                moveProblem(handout.id, dragFrom.current, i);
              dragFrom.current = null;
              setDragOver(null);
            }}
            onDragEnd={() => {
              dragFrom.current = null;
              setDragOver(null);
            }}
            className={
              "flex flex-row items-center gap-2 rounded-md border p-1 text-sm bg-white " +
              (dragOver === i ? "border-blue-800" : "border-gray-200")
            }
          >
            <span
              className="cursor-grab select-none px-1 text-gray-400"
              aria-hidden
            >
              ⠿
            </span>
            <span className="w-6 text-right font-mono text-gray-500">
              {i + 1}.
            </span>
            <span
              className="flex-1 truncate"
              dangerouslySetInnerHTML={{ __html: label }}
            />
            <button
              type="button"
              title="Move up"
              className={iconButton + " px-2 py-[2px]"}
              disabled={i === 0}
              onClick={() => moveProblem(handout.id, i, i - 1)}
            >
              ↑
            </button>
            <button
              type="button"
              title="Move down"
              className={iconButton + " px-2 py-[2px]"}
              disabled={i === handout.ids.length - 1}
              onClick={() => moveProblem(handout.id, i, i + 1)}
            >
              ↓
            </button>
            <button
              type="button"
              title="Remove from handout"
              className={
                "font-bold rounded-md duration-200 border border-gray-200 hover:bg-red-700 hover:text-white px-2 py-[2px]"
              }
              onClick={() => removeProblem(handout.id, id)}
            >
              ✕
            </button>
          </li>
        );
      })}
    </ol>
  );
}

// The printable body of a single handout: heading + problems. Shown on paper.
function HandoutContent(props: { handout: Handout; problems: any[] }) {
  const { handout, problems } = props;
  return (
    <>
      {handout.title && (
        <h2 className="w-full font-bold text-2xl rounded-sm p-[5px] text-center">
          {handout.title}
        </h2>
      )}
      {handout.author && (
        <p className="w-full text-sm rounded-sm p-[5px] text-center">
          {handout.author}
        </p>
      )}
      {handout.description && (
        <p className="w-full text-sm rounded-sm p-[5px] text-left whitespace-pre-wrap">
          {handout.description}
        </p>
      )}
      {problems.map((el, i) => (
        <Result
          key={handout.ids[i] ?? i}
          data={el}
          showtags={false}
          alias={handout.hideSource ? `Problem ${i + 1}` : undefined}
          handout
        />
      ))}
    </>
  );
}

// One handout: authoring card (active only, print-hidden) followed by its
// printable content. On screen only the active handout's content previews;
// on paper every handout prints, each starting on a fresh page.
function HandoutBlock(props: {
  handout: Handout;
  index: number;
  isActive: boolean;
  soloPrintId: string | null;
  onPrintSolo: (id: string) => void;
}) {
  const { handout, index, isActive, soloPrintId, onPrintSolo } = props;
  const { problems, loading, error } = useHandoutProblems(handout.ids);
  const { updateHandout, deleteHandout, handouts } = useHandouts();

  // Empty handouts (and, during a solo print, every non-target handout) are
  // kept off paper so they never emit a blank page.
  const printHidden =
    handout.ids.length === 0 ||
    (soloPrintId !== null && soloPrintId !== handout.id);
  const contentClass =
    (isActive ? "block " : "hidden ") +
    (printHidden ? "print:hidden" : "print:block") +
    (index > 0 ? " print:break-before-page" : "");

  return (
    <div>
      {isActive && (
        <div className="my-2 p-3 border-gray-200 border rounded-lg print:hidden">
          <div className="flex flex-row flex-wrap justify-between gap-2">
            <input
              type="text"
              value={handout.title}
              placeholder="Handout Title"
              onChange={(e) =>
                updateHandout(handout.id, { title: e.target.value })
              }
              className="rounded-md block text-sm my-1 flex-1 min-w-[10rem]"
            />
            <input
              type="text"
              value={handout.author}
              placeholder="Handout Author"
              onChange={(e) =>
                updateHandout(handout.id, { author: e.target.value })
              }
              className="rounded-md block text-sm my-1 flex-1 min-w-[10rem]"
            />
          </div>
          <textarea
            rows={2}
            placeholder="Handout Description"
            value={handout.description}
            onChange={(e) =>
              updateHandout(handout.id, { description: e.target.value })
            }
            className="rounded-md my-1 block text-sm w-full"
          />

          <ProblemList handout={handout} problems={problems} />

          <div className="flex flex-row flex-wrap items-center justify-between gap-2 mt-2">
            <label className="items-center flex">
              <span className="m-1 text-sm font-bold">Hide Problem Sources</span>
              <input
                type="checkbox"
                className="rounded-sm"
                checked={handout.hideSource}
                onChange={() =>
                  updateHandout(handout.id, { hideSource: !handout.hideSource })
                }
              />
            </label>
            <div className="flex flex-row flex-wrap items-center gap-2">
              <button
                className={iconButton + " p-2 text-sm"}
                onClick={() => onPrintSolo(handout.id)}
                disabled={handout.ids.length === 0}
              >
                Print this handout
              </button>
              <button
                className={
                  "p-2 text-sm font-bold rounded-md duration-200 border border-gray-200 hover:bg-red-700 hover:text-white disabled:opacity-30"
                }
                onClick={() => {
                  if (
                    handout.ids.length === 0 ||
                    window.confirm(`Delete handout "${handout.title || "Untitled"}"?`)
                  )
                    deleteHandout(handout.id);
                }}
                disabled={handouts.length <= 1 && handout.ids.length === 0}
                title={
                  handouts.length <= 1 && handout.ids.length === 0
                    ? "Nothing to delete"
                    : "Delete this handout"
                }
              >
                Delete
              </button>
            </div>
          </div>

          <div aria-live="polite">
            {loading ? (
              <p className="text-black my-2 font-bold text-sm">Loading…</p>
            ) : null}
            {error ? (
              <p className="text-red-600 my-2 font-bold text-sm">{error}</p>
            ) : null}
          </div>
        </div>
      )}

      <div className={contentClass}>
        <HandoutContent handout={handout} problems={problems} />
      </div>
    </div>
  );
}

export function HandoutGenerator() {
  const { handouts, activeId, setActiveId, createHandout } = useHandouts();
  const [soloPrintId, setSoloPrintId] = useState<string | null>(null);

  const printSolo = (id: string) => {
    // flushSync commits the "only this handout" render before the (synchronous)
    // print dialog reads the DOM, then we restore the full preview.
    flushSync(() => setSoloPrintId(id));
    window.print();
    setSoloPrintId(null);
  };

  return (
    <>
      {/* Toolbar — authoring only, never printed. */}
      <div className="my-2 print:hidden">
        <div className="flex flex-row flex-wrap items-center gap-2">
          <h2 className="font-bold text-lg mr-2">Handouts</h2>
          {handouts.map((h) => {
            const active = h.id === activeId;
            return (
              <button
                key={h.id}
                onClick={() => setActiveId(h.id)}
                title={h.title || "Untitled"}
                className={
                  "flex items-center max-w-[12rem] px-3 py-1 text-sm font-bold rounded-md border duration-200 " +
                  (active
                    ? "bg-blue-800 text-white border-blue-800"
                    : "border-gray-200 hover:bg-blue-800 hover:text-white")
                }
              >
                <span className="truncate">{h.title || "Untitled"}</span>
                <span
                  className={
                    "ml-2 shrink-0 rounded-sm px-1 text-xs " +
                    (active ? "bg-white text-blue-800" : "bg-gray-200 text-black")
                  }
                >
                  {h.ids.length}
                </span>
              </button>
            );
          })}
          <button
            onClick={() => createHandout({ title: "" })}
            className="px-3 py-1 text-sm font-bold rounded-md border border-gray-200 duration-200 hover:bg-blue-800 hover:text-white"
          >
            + New
          </button>
          <button
            onClick={() => {
              flushSync(() => setSoloPrintId(null));
              window.print();
            }}
            className="ml-auto px-3 py-1 text-sm font-bold rounded-md border border-gray-200 duration-200 hover:bg-blue-800 hover:text-white"
          >
            Print all
          </button>
        </div>
      </div>
      <hr className="my-2 print:hidden" />

      {handouts.length === 0 ? (
        <div className="text-sm text-gray-600 my-4 print:hidden">
          <p>
            You don't have any handouts yet. Create one with{" "}
            <span className="font-bold">+ New</span>, then add problems from the
            Search page.
          </p>
        </div>
      ) : (
        handouts.map((h, i) => (
          <HandoutBlock
            key={h.id}
            handout={h}
            index={i}
            isActive={h.id === activeId}
            soloPrintId={soloPrintId}
            onPrintSolo={printSolo}
          />
        ))
      )}
    </>
  );
}
