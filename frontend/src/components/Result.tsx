import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import { useEffect, useRef } from "react";

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
};

const preprocess = (s: string | undefined) => {
  if (!s) return "";
  // eqnarray and tabular modified
  const res = s
    .replace(new RegExp(/\s+/, "g"), " ")
    .replace(new RegExp(/\\begin{eqnarray\*}/, "g"), "\\begin{align*}")
    .replace(new RegExp(/\\end{eqnarray\*}/, "g"), "\\end{align*}")
    .replace(new RegExp(/\\begin{tabular}(\[.*?\])?/, "g"), "\\begin{array}")
    .replace(new RegExp(/\\end{tabular}(\[.*?\])?/, "g"), "\\end{array}");
  return res;
};

export default function Result(props: {
  statement?: string;
  solution?: string;
  url?: string;
  source?: string;
}) {
  const ref = useRef(null);
  useEffect(() => {
    if (ref.current) {
      renderMathInElement(ref.current, {
        delimiters,
        macros,
        throwonerror: false,
      });
    }
  });
  return (
    <div className="my-5 p-3 border-gray-200 border rounded-lg w-full">
      <div className="flex flex-wrap flex-row justify-between">
        <a
          href={props.url}
          target="_blank"
          className="mx-3 font-bold text-base"
        >
          {props.source?.replace(new RegExp("Problems Problem"), "Problem")}
        </a>
        <a
          href={props.solution}
          target="_blank"
          className="mx-3 font-bold text-base"
        >
          See Solution
        </a>
      </div>
      <div
        ref={ref}
        className="whitespace-pre-wrap md:max-w-3xl sm:max-w-xl max-w-lg overflow-y-hidden overflow-x-auto p-1 text-sm"
      >
        {preprocess(props.statement)}
      </div>
    </div>
  );
}
