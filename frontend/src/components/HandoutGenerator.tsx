import React, { createContext, useContext, useEffect, useState } from "react";
import Result from "./Result";

export const HandoutIdsContext = createContext<{
  idText: string;
  setIdText: any;
}>({ idText: "", setIdText: null });

export function HandoutProvider(props: { children: React.ReactNode }) {
  const [idText, setIdText] = useState<string>(
    localStorage.getItem("handout_ids") || ""
  );
  return (
    <HandoutIdsContext.Provider value={{ idText, setIdText }}>
      {props.children}
    </HandoutIdsContext.Provider>
  );
}

export function HandoutGenerator() {
  const [title, setTitle] = useState<string>(
    localStorage.getItem("handout_title") || ""
  );
  const [author, setAuthor] = useState<string>(
    localStorage.getItem("handout_author") || ""
  );
  const [desc, setDesc] = useState<string>(
    localStorage.getItem("handout_desc") || ""
  );
  const { idText, setIdText } = useContext(HandoutIdsContext);
  const [ids, setIds] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [hidesource, setHideSource] = useState<boolean>(false);
  const [error, setError] = useState<any | null>(null);
  const [problems, setProblems] = useState<any[]>([]);
  useEffect(() => {
    localStorage.setItem("handout_ids", idText);
    setIds(idText.trim().split(/\s+/));
  }, [idText]);
  useEffect(() => {
    if (ids.length == 0 || ids[0] === "") {
      setProblems([]);
      return;
    }
    setLoading(true);
    fetch(`/backend/handout`, {
      headers: {
        "Content-Type": "application/json",
      },
      method: "POST",
      body: JSON.stringify({ ids: ids }),
    })
      .then(async (data) => {
        setLoading(false);
        if (data.status != 200) {
          setProblems([]);
          setError(await data.text());
        } else {
          setError(null);
          data.json().then((json: any[]) => {
            setProblems(json);
          });
        }
      })
      .catch((_) => {
        setError(error);
      });
  }, [ids]);
  return (
    <>
      <div className="my-2 p-1 border-gray-200 border rounded-lg break-before-avoid-page break-inside-avoid-page break-after-avoid-page inline-block w-full print:hidden">
        <h2 className="w-full font-bold text-lg rounded-sm duration-200 p-[5px] flex justify-between">
          Handout Generator
        </h2>
        <>
          <div className="flex flex-row flex-wrap justify-between">
            <input
              type="text"
              name="title"
              value={title}
              placeholder="Handout Title"
              onChange={(e) => {
                e.preventDefault();
                localStorage.setItem("handout_title", e.target.value);
                setTitle(e.target.value);
              }}
              className="rounded-md block text-sm my-1"
            />
            <input
              type="text"
              name="author"
              value={author}
              placeholder="Handout Author"
              onChange={(e) => {
                e.preventDefault();
                localStorage.setItem("handout_author", e.target.value);
                setAuthor(e.target.value);
              }}
              className="rounded-md block text-sm my-1"
            />
            <button
              className="my-1 p-2 hover:bg-blue-800 hover:text-white font-bold rounded-md duration-200 w-fit border border-gray-200 text-sm"
              onClick={() => window.print()}
            >
              Print
            </button>
          </div>
          <textarea
            rows={2}
            placeholder={"Handout Description"}
            name="description"
            defaultValue={desc}
            onChange={(e) => {
              localStorage.setItem("handout_desc", e.target.value);
              setDesc(e.target.value);
            }}
            className="rounded-md my-1 block text-sm w-full"
          />
          <textarea
            rows={5}
            placeholder={"Problem IDs go here"}
            value={idText}
            onChange={(e) => {
              setIdText(e.target.value);
            }}
            className="rounded-md my-1 block text-sm w-full"
          />

          <div className="items-center flex">
            <span className="m-1 text-sm font-bold">Hide Problem Sources</span>
            <input
              type="checkbox"
              className="rounded-sm"
              alt="Hide problem sources when printing?"
              checked={hidesource}
              onChange={() => setHideSource(!hidesource)}
            />
            </div>
        </>
      </div>
      <div className="print:hidden">
        {loading ? (
          <p className="text-black my-2 font-bold">Loading...</p>
        ) : null}
        {error ? <p className="text-red-600 my-2 font-bold">{error}</p> : null}
      </div>
      <hr className="my-2 print:hidden" />
      {title && (
        <h2 className="w-full font-bold text-2xl rounded-sm duration-200 p-[5px] text-center">
          {title}
        </h2>
      )}
      {author && (
        <p className="w-full text-sm rounded-sm duration-200 p-[5px] text-center">
          {author}
        </p>
      )}
      {desc && (
        <p className="w-full text-sm rounded-sm duration-200 p-[5px] text-left">
          {desc}
        </p>
      )}
      {problems.length
        ? problems.map((el, i) => (
            <Result key={i} data={el} showtags={false} alias={hidesource ? `Problem ${i+1}` : undefined} handout />
          ))
        : null}
    </>
  );
}
