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
  "\\ ": " ",
  "\\O": "\\empty",
};

export default function Result(props: {
  statement?: string;
  solution?: string;
  url?: string;
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
    // md:w-[600px] sm:w-[400px] w-[300px]
    <div className="my-5 p-3 border-gray-200 border rounded-lg w-full">
      <div className="flex flex-wrap flex-row">
        <a
          href={props.url}
          target="_blank"
          className="mx-3 font-bold text-base"
        >
          {props.url
            ?.split("index.php/")[1]
            .replace(new RegExp("[#_]", "g"), " ")}
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
        {props.statement}
      </div>
    </div>
  );
}
