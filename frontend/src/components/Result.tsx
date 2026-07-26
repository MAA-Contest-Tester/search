import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import "./Result.css";
import { useEffect, useRef } from "react";
import { useHandouts } from "./handouts";
import React, {NavLink} from "react-router-dom";

const delimiters = [
  { left: "$$", right: "$$", display: true },
  { left: "\\(", right: "\\)", display: false },
  {
    left: "\\begin{equation}",
    right: "\\end{equation}",
    display: true,
  },
  {
    left: "\\begin{equation*}",
    right: "\\end{equation*}",
    display: true,
  },
  { left: "\\begin{align*}", right: "\\end{align*}", display: true },
  { left: "\\begin{align}", right: "\\end{align}", display: true },
  { left: "\\begin{alignat}", right: "\\end{alignat}", display: true },
  { left: "\\begin{gather}", right: "\\end{gather}", display: true },
  { left: "\\begin{CD}", right: "\\end{CD}", display: true },
  { left: "\\[", right: "\\]", display: true },
  { left: "$", right: "$", display: false },
  { left: "\\(", right: "\\)", display: false },
];

const macros = {
  "\\emph": "\\textit",
  "\\textsc": "",
  "\\textdollar": "\\$",
  "\\overarc": "\\overgroup",
  "\\dfrac": "\\frac",
  "\\O": "\\empty",
  "\\hdots": "\\ldots",
};

const preprocess = (s: string | undefined) => {
  if (!s) return "";
  // eqnarray and tabular modified
  const res = s
    .trim()
    .replace(new RegExp(/\\begin\{eqnarray\*}/, "g"), "\\begin{align*}")
    .replace(new RegExp(/\\end\{eqnarray\*}/, "g"), "\\end{align*}")
    .replace(new RegExp(/\\begin\{tabular}(\[.*?\])?/, "g"), "\\begin{array}")
    .replace(new RegExp(/\\end\{tabular}(\[.*?\])?/, "g"), "\\end{array}")
    .replace(new RegExp(/\\makebox(\[.*?\])?/, "g"), "\\begin{array}")
    .replace(new RegExp(/\\mbox/, "g"), "\\text")
    .replace(new RegExp(/\\bigg\s*\{\\\}\}/, "g"), "\\bigg \\}")
    .replace(new RegExp(/\\bigg\s*\{\\\{\}/, "g"), "\\bigg \\{");
  return res;
};

export default function Result(props: {
  data: {
    statement?: string;
    solution?: string;
    url?: string;
    source?: string;
    categories?: string;
    id: string;
  } | null;
  showtags: boolean;
  alias?: string;
  handout?: boolean;
}) {
  const ref = useRef(null);
  useEffect(() => {
    if (ref.current) {
      renderMathInElement(ref.current, {
        delimiters,
        macros,
        trust: true,
        throwonerror: false,
        errorColor: "#cc0000",
      });
    }
  });
  const { toggleProblemInActive, isInActive, activeHandout } = useHandouts();
  const added = props.data ? isInActive(props.data.id) : false;
  const preprocessed = preprocess(props.data?.statement);
  return (
    <div
      className={
        "my-2 p-3 border-gray-200 border rounded-lg break-before-avoid-page break-inside-avoid-page break-after-avoid-page inline-block w-full"
      }
    >
      {props.data !== null ? (
        <>
        {
          props.alias ?
          <div
            className={"mx-3 font-bold text-base"}
          >
            {props.alias}
          </div>
          : <a
              href={props.data.url}
              target="_blank"
              className={"mx-3 font-bold text-base"}
              dangerouslySetInnerHTML={{ __html: props.alias || props.data.source! }}
            ></a>
        }
          <div className="flex flex-wrap flex-row justify-between items-center print:hidden">
            {
              !props.alias ?
              <a
                href={props.data.solution}
                target="_blank"
                className="mx-3 font-bold text-base"
              >
                See Discussion
              </a>
              : null
            }
            {props.handout ? null : (
              <div className="flex flex-wrap flex-row justify-left items-center group relative print:hidden">
                {"ontouchstart" in window ||
                navigator.maxTouchPoints > 0 ? null : (
                  <span
                    className={
                      "absolute left-0 top-full mt-2 -translate-x-[0.6rem] max-w-[12rem] w-max origin-top scale-0 group-hover:scale-100 transition duration-200 ease-in-out rounded-lg z-50 p-1 font-bold text-sm text-center" +
                      " " +
                      (added
                        ? "bg-green-800 text-white"
                        : "bg-white text-black border border-gray-200")
                    }
                  >
                    {added ? "Remove from " : "Add to "}
                    <NavLink
                      to="/handout"
                      className={added ? "text-white hover:text-white decoration-white" : ""}
                    >
                      {activeHandout?.title || "Handout"}
                    </NavLink>
                  </span>
                )}
                <button
                  title={
                    added
                      ? `Remove from "${activeHandout?.title || "Handout"}"`
                      : `Add to "${activeHandout?.title || "Handout"}"`
                  }
                  onClick={(_) => {
                    if (props.data) toggleProblemInActive(props.data.id);
                  }}
                  className={
                    "mx-3 font-bold text-sm p-[5px] border-gray-200 rounded-lg border duration-200" +
                    " " +
                    (added ? "bg-green-800 text-white" : "")
                  }
                >
                  {added ? (
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      height="1.2em"
                      viewBox="0 0 448 512"
                    >
                      <path
                        d="M438.6 105.4c12.5 12.5 12.5 32.8 0 45.3l-256 256c-12.5 12.5-32.8 12.5-45.3 0l-128-128c-12.5-12.5-12.5-32.8 0-45.3s32.8-12.5 45.3 0L160 338.7 393.4 105.4c12.5-12.5 32.8-12.5 45.3 0z"
                        fill="#ffffff"
                      />
                    </svg>
                  ) : (
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      height="1.2em"
                      viewBox="0 0 448 512"
                    >
                      <path d="M256 80c0-17.7-14.3-32-32-32s-32 14.3-32 32V224H48c-17.7 0-32 14.3-32 32s14.3 32 32 32H192V432c0 17.7 14.3 32 32 32s32-14.3 32-32V288H400c17.7 0 32-14.3 32-32s-14.3-32-32-32H256V80z" />
                    </svg>
                  )}
                </button>
              </div>
            )}
          </div>
        </>
      ) : (
        <h2 className="mx-3 font-bold text-base">404 Not Found</h2>
      )}
      {props.data !== null ? (
        <>
          <div
            ref={ref}
            className="whitespace-pre-wrap w-full overflow-y-hidden overflow-x-auto p-1 text-sm select-text"
          >
            {preprocessed}
          </div>
          {props.showtags ? (
            <>
              <hr className="print:hidden" />
              <div className="whitespace-pre-wrap w-full overflow-y-hidden overflow-x-auto p-1 text-sm select-text print:hidden">
                <strong>Tags: </strong>
                <span
                  dangerouslySetInnerHTML={{ __html: props.data.categories! }}
                />
              </div>
            </>
          ) : null}
        </>
      ) : null}
    </div>
  );
}
