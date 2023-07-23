import React, { useEffect, useState } from "react";

interface MetaType {
  contestlist: Map<string, {id: number, name: string}[]>,
  problemcount: number,
  date: Date
}

export default function Metadata() {
  const [meta, setMeta] = useState<MetaType | null>(null);
  useEffect(() => {
    fetch(`/backend/meta`).then(data => data.json()).then(json => {
      const pairs: [string, {id: number, name:string}[]][] = []
      for (const key in json.contestlist) {
        pairs.push([key, json.contestlist[key]])
      }
      setMeta({
        contestlist: new Map(pairs),
        problemcount: json["problemcount"],
        date: new Date(json["date"]),
      })
    })
  }, [])
  return (
    <div className="mt-2">
      <p className="text-lg"><strong>{meta?.problemcount}</strong>{" "}Problems</p>
      <p className="text-sm">Last updated on {meta?.date.toDateString()}</p>
      <hr className="my-3"/>
      <h1 className="text-xl font-bold">Included Contests</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3">
      {meta?.contestlist ? Array.from(meta?.contestlist).map(([key, value]) => (
        <div key={key} className="border-gray-200 rounded-lg border p-3 m-2">
          <h2 className="text-lg">{key}</h2>
          <ul className="list-inside list-disc">
            {value.map((contest, index) => (
              <li className="text-xs" key={index}><a href={`https://artofproblemsolving.com/community/c${contest.id}`} target="_blank">{contest.name}</a></li>
            ))}
          </ul>
        </div>
      )): null}
      </div>
    </div>
  )
}
