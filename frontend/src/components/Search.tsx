import { useEffect, useState } from "react";
import Result from "./Result";
import { debounce } from "debounce";

export default function Search() {
  const [statement, setStatement] = useState<string>();
  const [source, setSource] = useState<string>();
  const [categories, setCategories] = useState<string>();
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
    if (!whitespace.test(categories || "")) {
      const c = (categories || "").trim();
      q += `@categories:(${c})`;
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
      <p className="my-3 mx-0 text-xs sm:text-sm max-w-fit">
        Type the text you want to search for (e.g. {queryExample("complex")} or{" "}
        {queryExample("polynomial")}), or you can use redisearch's querying
        capabilities. For example, to just search for AIME introductory
        problems, you might do{" "}
        {queryExample("@source:(AIME) @categories:(easy)")}. To search for AMC
        10 Problems with "mean", search {queryExample("@source:(AMC 10) mean")}.
        Or for USAMO or USAJMO Geometry Problems with "prove cyclic" in their
        statement, search{" "}
        {queryExample(
          "@source:(USAMO|USAJMO) @categories:(geo) @statement:(prove cyclic)"
        )}
        . Wildcard searching is also allowed, such as {queryExample("*count*")}.
        You can also mix and match all of the above, such as{" "}
        {queryExample("@source:(*MO) @statement:(equi*) *gle")}
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
            placeholder="Categories or Level (e.g. Geometry, Intermediate, Olympiad...)"
            onChange={(e) => {
              e.preventDefault();
              setCategories(e.target.value);
            }}
            className="w-full rounded-md my-1 inline-block col-span-2 md:col-span-3"
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
        {results.length
          ? results.map((el, i) => <Result key={i} {...el} />)
          : null}
      </div>
    </>
  );
}
