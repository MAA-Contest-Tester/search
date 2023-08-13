import { useContext } from "react";
import React, { NavLink } from "react-router-dom";
import { HandoutIdsContext } from "./HandoutGenerator";

export default function Navbar() {
  const { idText } = useContext(HandoutIdsContext);
  const trimmed = idText.trim()
  const length = trimmed.split(/\s+/).length;
  return (
    <div className="sticky top-0 p-2 bg-white z-50 print:hidden border-b border-x border-gray-400 rounded-b-lg mb-5">
      <nav className="flex flex-row justify-left gap-4 my-2 font-bold print:hidden">
        <NavLink to="/" className={""}>
          Search
        </NavLink>
        <span className="">
          <NavLink to="/handout" className={""}>
            Handout
          </NavLink>
          {trimmed.length != 0 ? (
            <span className="translate-y-0 rounded-sm text-xs text-white bg-red-700 my-auto ml-[2px] px-1">
              {length > 99 ? "99+" : length}
            </span>
          ) : null}
        </span>
        <NavLink to="/meta" className={""}>
          Info
        </NavLink>
      </nav>
    </div>
  );
}
