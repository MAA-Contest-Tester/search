import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import "./Result.css";
import { useContext, useEffect, useRef } from "react";
import {HandoutIdsContext} from "./HandoutGenerator"

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
  statement?: string;
  solution?: string;
  url?: string;
  source?: string;
  categories?: string;
  id: string;
  showtags: boolean;
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
  const {idText, setIdText} = useContext(HandoutIdsContext);
  const preprocessed = preprocess(props.statement);
  return (
    <div
      className={
        "my-2 p-3 border-gray-200 border rounded-lg break-before-avoid-page break-inside-avoid-page break-after-avoid-page inline-block w-full"
      }
    >
      <a
        href={props.url}
        target="_blank"
        className="mx-3 font-bold text-base"
        dangerouslySetInnerHTML={{ __html: props.source! }}
      ></a>
      <div className="flex flex-wrap flex-row justify-between items-center print:hidden">
        <a
          href={props.solution}
          target="_blank"
          className="mx-3 font-bold text-base"
        >
          See Discussion
        </a>
        <div className="flex flex-wrap flex-row justify-left items-center print:hidden">
          <button
            onClick={(_) => setIdText(idText + (idText.trim().length != 0 ? "\n" : "") + props.id)}
            className="mx-3 font-bold text-sm hover:bg-blue-800 hover:text-white p-[5px] border-gray-200 rounded-lg border duration-200"
          >
            Add to Handout
          </button>
        </div>
      </div>
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
            <span dangerouslySetInnerHTML={{ __html: props.categories! }} />
          </div>
        </>
      ) : null}
    </div>
  );
}
