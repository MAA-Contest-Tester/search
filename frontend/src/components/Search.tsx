import { useEffect, useState } from "react";
import Result from "./Result";

export default function Search() {
  const [query, setQuery] = useState("");
  const [error, setError] = useState<any | null>(null);
  const [results, setResults] = useState<any[]>([]);
  const [showTags, setShowTags] = useState<boolean>(false);

  // indicates whether query has changed.
  const [loading, setLoading] = useState<boolean>(false);
  // indicates whether currently paginating or not.
  const [pageLoading, setPageLoading] = useState<boolean>(false);
  // status for whether pagination returns no more results.
  const [nothing, setNothing] = useState<boolean>(false);

  const apicall = () => {
    fetch(`/search?query=${encodeURI(query)}`)
      .then(async (data) => {
        setLoading(false);
        if (data.status != 200) {
          setResults([]);
          setError(await data.text());
        } else {
          setError(null);
          data.json().then((json: any[]) => {
            setResults(json.map(x => x["_formatted"]));
          });
        }
      })
      .catch((_) => {
        setError(error);
      });
  };
  const nextpage = () => {
    setPageLoading(true);
    const offset = results.length
    fetch(`/search?query=${encodeURI(query)}&offset=${offset}`)
      .then(async (data) => {
        setPageLoading(false);
        if (data.status != 200) {
          setResults([]);
          setError(await data.text());
        } else {
          setError(null);
          data.json().then((json: any[]) => {
            if (json.length === 0) {
              setNothing(true)
            }
            setResults(results.concat(json.map(x => x["_formatted"])));
          });
        }
      })
      .catch((_) => {
        setError(error);
      });
  };
  useEffect(() => {
    setNothing(false)
    setLoading(true);
    const debounced_api = setTimeout(() => apicall(), 200);
    return () => clearTimeout(debounced_api);
  }, [query]);

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
      <div className="print:hidden text-sm">
      Try searching for keywords (e.g. {queryExample("moving points")} or
      {" "}{queryExample("inequality")}). You could also search for problems from a
      specific year or contest ({queryExample("2022 ISL G8")}). Some common
      abbreviations will work ({queryExample("fe")}, {queryExample("nt")}, etc)
        <p className="my-3 mx-0 text-xs sm:text-sm max-w-fit">
        </p>
        <div className="flex flex-row flex-wrap justify-between">
          <input
            type="text"
            value={query}
            readOnly={pageLoading}
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
          <div className="items-center flex">
            <span className="mx-1 text-sm font-bold">Show Tags</span>
            <input
              type="checkbox"
              className="rounded-sm"
              alt="Include when printing?"
              checked={showTags}
              onChange={() => setShowTags(!showTags)}
            />
            </div>
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
          ? results.map((el, i) => <Result key={i} {...el} showtags={showTags} />)
          : null}
      </div>
      {!nothing ?
      <div className="text-center">
        <button
          className="my-1 p-2 hover:bg-blue-800 hover:text-white font-bold rounded-md duration-200 w-fit border border-gray-200 text-sm"
          onClick={() => {if (!loading) nextpage()}}
        >
          Next Page
        </button>
        {pageLoading ? <p className="text-black my-2 font-bold">Loading...</p> : null}
      </div>
      : null}
    </>
  );
}
