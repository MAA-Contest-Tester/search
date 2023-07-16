import { useEffect, useState } from "react";
import Result from "./Result";
import { debounce } from "debounce";


export default function Search() {
  const [query, setQuery] = useState("");
  const [error, setError] = useState<any | null>(null);
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [showTags, setShowTags] = useState<boolean>(false);
  const apicall = () => {
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
            setResults(json.map(x => x["_formatted"]));
          });
        }
      })
      .catch((_) => {
        setError(error);
      });
  };
  useEffect(() => {
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
      <div className="print:hidden">
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
    </>
  );
}
