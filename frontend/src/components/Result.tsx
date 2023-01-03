import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import "./Result.css";
import { useEffect, useRef } from "react";
import { renderconfig, preprocess } from "../katex_constants";

export default function Result(props: {
  statement?: string;
  solution?: string;
  url?: string;
  source?: string;
}) {
  const ref = useRef(null);
  useEffect(() => {
    if (ref.current) {
      renderMathInElement(ref.current, renderconfig);
    }
  });
  const preprocessed = preprocess(props.statement);
  return (
    <div className="my-5 p-3 border-gray-200 border rounded-lg w-full">
      <a href={props.url} target="_blank" className="mx-3 font-bold text-base">
        {props.source?.replace(new RegExp("Problems Problem"), "Problem")}
      </a>
      <div className="flex flex-wrap flex-row justify-between items-center">
        <a
          href={props.solution}
          target="_blank"
          className="mx-3 font-bold text-base"
        >
          See Solution
        </a>
        <button
          onClick={(_) => navigator.clipboard.writeText(preprocessed)}
          className="mx-3 font-bold text-base hover:bg-blue-800 hover:text-white p-[2px] border-gray-200 rounded-lg border duration-200"
        >
          Copy
        </button>
      </div>
      <div
        ref={ref}
        className="whitespace-pre-wrap md:max-w-3xl sm:max-w-xl max-w-lg overflow-y-hidden overflow-x-auto p-1 text-sm select-text"
      >
        {preprocessed}
      </div>
    </div>
  );
}
