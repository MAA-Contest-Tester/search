import { useEffect, useState } from "react";
import Result from "./Result";
import { debounce } from "debounce";

const mtch = /Problem (.*)/;

function sortResults(data: { source: string }[]) {
  data.sort((a, b) => {
    const x = mtch.exec(a.source);
    const y = mtch.exec(b.source);
    const compare = (x: any, y: any) => {
      if (x > y) {
        return 1;
      } else if (x < y) {
        return -1;
      } else {
        return 0;
      }
    };
    if (x && y) {
      const a_prev = a.source.slice(0, x.index);
      const b_prev = b.source.slice(0, y.index);
      if (compare(b_prev, a_prev)) {
        return compare(b_prev, a_prev);
      }
      if (compare(x[1], y[1])) {
        return compare(x[1], y[1]);
      }
    }
    return 0;
  });
}

export default function Search() {
  const [statement, setStatement] = useState<string>();
  const [source, setSource] = useState<string>();
  const [categories, setCategories] = useState<string>();
  const [query, setQuery] = useState("");
  const [error, setError] = useState<any | null>(null);
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const debounced_api = debounce(() => {
    setLoading(true);
    fetch(`/search?query=${encodeURI(query)}`)
      .then(async (data) => {
        setLoading(false);
        if (data.status != 200) {
          setResults([]);
          setError(await data.text());
        } else {
          setError(null);
          data.json().then((json: any[]) => {
            sortResults(json);
            setResults(json);
          });
        }
      })
      .catch((_) => {
        setError(error);
      });
  }, 100);
  useEffect(() => {
    debounced_api();
  }, [query]);

  useEffect(() => {
    let q = "";
    const whitespace = new RegExp("^s*$");
    if (!whitespace.test(statement || "")) {
      const c = (statement || "").trim();
      q += `@statement:(${c}*)`;
    }
    if (!whitespace.test(source || "")) {
      const c = (source || "").trim();
      q += `@source:(${c})`;
    }
    if (!whitespace.test(categories || "")) {
      const c = (categories || "").trim();
      q += `@categories:(${c}*)`;
    }
    setQuery(q);
  }, [statement, source, categories]);

  const queryExample = (q: string) => (
    <button
      className="underline rounded-md hover:text-blue-800 decoration-blue-800 inline font-mono text-left"
      onClick={() => setQuery(q)}
    >
      {q}
    </button>
  );
  return (
    <>
      <div className="print:hidden">
        <p className="my-3 mx-0 text-xs sm:text-sm max-w-fit">
          Type the text you want to search for (e.g. {queryExample("complex")}{" "}
          or {queryExample("polynomial")}), or you can use redisearch's querying
          capabilities. For example, to just search for USAMO geometry problems,
          type
          {queryExample("@source:(USAMO) @categories:(geometry)")}. To search
          for AMC 10 Problems with "mean", search{" "}
          {queryExample("@source:(AMC 10) mean")}. Or for Olympiad Algebra
          Problems about inequalities, search{" "}
          {queryExample("@source:(*MO) @categories:(algebra inequality)")}.
          Wildcard searching is also allowed, such as {queryExample("*count*")}.
          You can also mix and match all of the above, such as{" "}
          {queryExample(
            "@source:(AIME) @statement:(complex) @categories:(number theory)"
          )}
          , which searches for AIME number theory problems about complex
          numbers.
        </p>
        <div className="border-gray-200 rounded-lg p-3 my-2 border">
          <h2 className="font-extrabold text-xl">Query Helper</h2>
          <label className="grid grid-cols-3 md:grid-cols-4 items-center">
            <span className="inline mr-3 col-span-1">Problem Source</span>
            <input
              type="text"
              placeholder="Source (e.g. Contest Name, Year, Problem Number)"
              onChange={(e) => {
                e.preventDefault();
                setSource(e.target.value);
              }}
              className="w-full rounded-md my-1 inline-block col-span-2 md:col-span-3"
            />
          </label>
          <label className="grid grid-cols-3 md:grid-cols-4 items-center">
            <span className="inline mr-3 col-span-1">Problem Statement</span>
            <input
              type="text"
              placeholder="Text that matches your problem"
              onChange={(e) => {
                e.preventDefault();
                setStatement(e.target.value);
              }}
              className="w-full rounded-md my-1 inline-block col-span-2 md:col-span-3"
            />
          </label>
          <label className="grid grid-cols-3 md:grid-cols-4 items-center">
            <span className="inline mr-3 col-span-1">Categories</span>
            <input
              type="text"
              placeholder="Categories (e.g. Geometry, Complex, Functional...)"
              onChange={(e) => {
                e.preventDefault();
                setCategories(e.target.value);
              }}
              className="w-full rounded-md my-1 inline-block col-span-2 md:col-span-3"
            />
          </label>
        </div>
        <div className="flex flex-row flex-wrap justify-between">
          <input
            type="text"
            value={query}
            placeholder="Raw Query"
            onChange={(e) => {
              e.preventDefault();
              setQuery(e.target.value);
            }}
            className="rounded-md my-1 block w-full"
          />
          <button
            className="my-1 p-2 hover:bg-blue-800 hover:text-white font-bold rounded-md duration-200 w-fit border border-gray-200"
            onClick={() => window.print()}
          >
            Print
          </button>
        </div>

        {loading ? <p className="text-black my-2 font-bold">Loading...</p> : null}
        {error ? <p className="text-red-600 my-2 font-bold">{error}</p> : null}
      </div>
      <div className="w-full break-inside-avoid-page">
        <div className="print:visible invisible text-xs">
          Created with{" "}
          <strong>
            <span className="text-blue-800">Search.</span>
            MAATester.com
          </strong>{" "}
        </div>
        {results.length
          ? results.map((el, i) => <Result key={i} {...el} />)
          : null}
      </div>
    </>
  );
}
