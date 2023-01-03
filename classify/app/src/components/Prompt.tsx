import "katex/dist/katex.min.css";
import renderMathInElement from "katex/dist/contrib/auto-render";
import { useEffect, useRef, useState } from "react";
import {
  preprocess,
  renderconfig,
} from "workspace-frontend/src/katex_constants";

function StreakComp(props: { streak: number }) {
  const streak = props.streak;
  if (streak >= 20) {
    return (
      <h1 className="text-xl font-extrabold text-orange-700">
        Hottt!!! {streak} in a row!
      </h1>
    );
  }
  if (streak >= 10) {
    return (
      <h1 className="text-xl font-extrabold text-green-800">
        Nice! {streak} in a row!
      </h1>
    );
  }
  if (streak >= 3) {
    return (
      <h1 className="text-xl font-extrabold text-blue-800">
        Keep Going! {streak} in a row!
      </h1>
    );
  } else {
    return null;
  }
}

export default function Prompt(props: { streak: number }) {
  const [answer, setAnswer] = useState<number | null>(null);
  const [source, setSource] = useState<string>("");
  const [statement, setStatement] = useState<string>("");
  const [error, setError] = useState<string>();
  const [submitted, setSubmitted] = useState<boolean>(false);
  const possibilities = [
    "Not a Problem",
    "Algebra",
    "Geometry",
    "Number Theory",
    "Combinatorics",
  ];
  const ref = useRef(null);

  useEffect(() => {
    fetch("/api/choose")
      .then((res) => res.json())
      .then((res) => {
        setSource(res["Source"]);
        setStatement(res["Statement"]);
      })
      .catch((e) => setError(e));
  }, []);

  const submit = async () => {
    fetch("/api/add", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        source: source,
        statement: statement,
        answer: answer,
      }),
    })
      .then((res) => res.json())
      .catch((e) => setError(e));
    setSubmitted(true);
  };

  // TODO: write logic to fetch a random data point from the API.

  useEffect(() => {
    if (ref.current) {
      renderMathInElement(ref.current, renderconfig);
    }
  }, [statement]);

  return (
    <div className="rounded-lg prose px-0">
      <StreakComp streak={props.streak} />
      <hr className="mb-2" />
      <div
        ref={ref}
        className="whitespace-pre-wrap md:max-w-3xl sm:max-w-xl max-w-lg overflow-y-hidden overflow-x-auto p-1 text-sm select-text"
      >
        {preprocess(statement)}
      </div>
      <table className="border border-gray-200 rounded-xl prose table justify-left p-2 gap-2 my-3">
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
                  onChange={(_) => {
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
      <div className="flex flex-row flex-wrap justify-start items-center gap-2">
        <button
          disabled={submitted || answer === null}
          onClick={submit}
          className={
            "p-2 text-white my-3" +
            " " +
            (submitted || answer === null ? "bg-gray-500" : "bg-blue-700")
          }
        >
          Submit
        </button>
        {submitted ? <div className="text-green-700">Submitted!</div> : null}
      </div>
      {error ? <div className="p-1 text-red-700">{error}</div> : null}
    </div>
  );
}
