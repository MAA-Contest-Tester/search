import React, { NavLink } from "react-router-dom";
import { useHandouts } from "./handouts";

export default function Navbar() {
  const { handouts, activeHandout } = useHandouts();
  const count = activeHandout?.ids.length ?? 0;
  return (
    <div className="sticky top-0 p-2 bg-white z-50 print:hidden border-b border-x border-gray-400 rounded-b-lg mb-5">
      <nav className="flex flex-row justify-left gap-4 my-2 font-bold print:hidden">
        <NavLink to="/" className={""}>
          Search
        </NavLink>
        <span className="">
          <NavLink to="/handout" className={""}>
            Handout{handouts.length > 1 ? "s" : ""}
          </NavLink>
          {count > 0 ? (
            <span className="translate-y-0 rounded-sm text-xs text-white bg-red-700 my-auto ml-[2px] px-1">
              {count > 99 ? "99+" : count}
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
