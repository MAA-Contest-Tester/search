import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import { useEffect, useRef, useState } from "react";
import {
  preprocess,
  renderconfig,
} from "workspace-frontend/src/katex_constants";

export default function Prompt() {
  const [answer, setAnswer] = useState<number | null>(null);
  const [source, setSource] = useState<string>("");
  const [statement, setStatement] = useState<string>("");
  const possibilities = [
    "Not a Problem",
    "Algebra",
    "Geometry",
    "Number Theory",
    "Combinatorics",
  ];
  const ref = useRef(null);

  // TODO: write logic to fetch a random data point from the API.

  useEffect(() => {
    if (ref.current) {
      renderMathInElement(ref.current, renderconfig);
    }
  }, [statement]);

  return (
    <div className="rounded-lg prose p-2">
      <h1 className="text-2xl font-bold">Source:{source}</h1>
      <div
        ref={ref}
        className="whitespace-pre-wrap md:max-w-3xl sm:max-w-xl max-w-lg overflow-y-hidden overflow-x-auto p-1 text-sm select-text"
      >
        {preprocess(statement)}
      </div>
      <table className="border border-gray-200 rounded-lg prose table justify-left p-1 gap-2">
        <tbody>
          {possibilities.map((v, i) => (
            <tr className="table-row p-1" key={i}>
              <td className="table-cell p-1">
                <input
                  type="radio"
                  className="col-span-1"
                  name="response"
                  value={i}
                  key={i}
                  onChange={(e) => {
                    setAnswer(i);
                  }}
                ></input>
              </td>
              <td className="table-cell p-1">
                <label className="font-bold col-span-9">{v}</label>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
