import React, { createContext, useContext, useEffect, useState } from "react";

export const HandoutIdsContext = createContext<{
  idText: string;
  setIdText: any;
}>({idText: "", setIdText:null});

export function HandoutProvider(props:{children:React.ReactNode}) {
  const [idText, setIdText] = useState<string>(
    (localStorage.getItem("handout_ids") || "")
  );
  return <HandoutIdsContext.Provider value={{idText, setIdText}}>
    {props.children}
  </HandoutIdsContext.Provider>
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
  const [expanded, setExpanded] = useState<boolean>(false);
  const {idText, setIdText} = useContext(HandoutIdsContext);
  const [ids, setIds] = useState<string[]>([])
  useEffect(() => {
    localStorage.setItem("handout_ids", idText);
    setIds(idText.trim().split(/\s+/));
  }, [idText])
  return (
    <form
      className="my-2 p-2 border-gray-200 border rounded-lg break-before-avoid-page break-inside-avoid-page break-after-avoid-page inline-block w-full"
      method="POST"
      action="/handout"
    >
      <h2
        className="w-full font-bold text-md rounded-sm duration-200 p-[5px] flex justify-between cursor-pointer"
        onClick={() => setExpanded(!expanded)}
      >
        <span>Handout Generator</span> <span>{expanded ? "-" : "+"}</span>
      </h2>
      {expanded ? (
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
              className="rounded-md m-1 block text-sm"
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
              className="rounded-md m-1 block text-sm"
            />
            <button
              className="my-1 p-2 hover:bg-blue-800 hover:text-white font-bold rounded-md duration-200 w-fit border border-gray-200 text-sm"
              type="submit"
              value="submit"
            >
              Generate
            </button>
          </div>
          <textarea
            rows={2}
            placeholder={"Handout Description."}
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
            placeholder={"Problem IDs go here."}
            value={idText}
            onChange={(e) => {
              setIdText(e.target.value);
            }}
            className="rounded-md my-1 block text-sm w-full"
          />
          {ids.map((id) => (
            <input name="id" value={id} type="hidden" key={id} />
          ))}
        </>
      ) : null}
    </form>
  );
}
