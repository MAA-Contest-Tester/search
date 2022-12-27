import { useEffect, useState } from "react";
import Result from "./Result";
import { debounce } from "debounce";

export default function Search() {
  const [statement, setStatement] = useState<string>();
  const [source, setSource] = useState<string>();
  const [query, setQuery] = useState("");
  const [error, setError] = useState<any | null>(null);
  const [results, setResults] = useState<any[]>([]);
  const debounced_api = debounce(() => {
    fetch(`/search?query=${encodeURI(query)}`)
      .then(async (data) => {
        if (data.status != 200) {
          setResults([]);
          setError(await data.text());
        } else {
          setError(null);
          data.json().then((json) => setResults(json));
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
    setQuery(q);
  }, [statement, source]);

  const queryExample = (q: string) => (
    <span
      className="underline rounded-md hover:text-blue-800 decoration-blue-800 inline-block font-mono"
      onClick={() => setQuery(q)}
    >
      {q}
    </span>
  );
  return (
    <>
      <p className="my-3 mx-0 text-xs max-w-fit">
        Type the text you want to search for (e.g. {queryExample("complex")} or{" "}
        {queryExample("polynomial")}), or you can use redisearch's querying
        capabilities. For example, to just search for AIME problems, you might
        do {queryExample("@source:(AIME)")}. To search for AMC 10 Problems with
        "mean", search {queryExample("@source:(AMC 10) mean")}. Or for USAMO or
        USAJMO Problems with "prove cyclic" in their statement, search{" "}
        {queryExample("@source:(USAMO|USAJMO) @statement:(prove cyclic)")}.
        Wildcard searching is also allowed, such as {queryExample("*count*")}.
        You can also mix and match all of the above, such as{" "}
        {queryExample("@source:(JBMO) @statement:(equi*) *gle")}
      </p>
      <div className="border-gray-200 rounded-lg p-3 my-2 border">
        <h2 className="font-extrabold text-xl">Query Helper</h2>
        <label className="flex justify-between items-center">
          <span className="inline mr-3">Problem Source</span>
          <input
            type="text"
            placeholder="Problem Source"
            onChange={(e) => {
              e.preventDefault();
              setSource(e.target.value);
            }}
            className="w-9/12 rounded-md my-1 inline-block"
          />
        </label>
        <label className="flex justify-between items-center">
          <span className="inline mr-3">Problem Statement</span>
          <input
            type="text"
            placeholder="Problem Statement"
            onChange={(e) => {
              e.preventDefault();
              setStatement(e.target.value);
            }}
            className="w-9/12 rounded-md my-1 inline-block"
          />
        </label>
      </div>
      <input
        type="text"
        value={query}
        placeholder="Raw Query"
        onChange={(e) => {
          e.preventDefault();
          setQuery(e.target.value);
        }}
        className="w-full rounded-md my-1"
      />
      {error ? <p className="text-red-600 my-2">{error}</p> : null}
      <div className="w-full">
        {results.map((el, i) => (
          <Result key={i} {...el} />
        ))}
      </div>
    </>
  );
}
