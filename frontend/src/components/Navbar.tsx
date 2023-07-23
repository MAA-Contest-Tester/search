import React, { NavLink } from "react-router-dom";

export default function Navbar() {
  return (
  <div className="sticky top-0 p-2 bg-white z-50">
    <nav className="flex flex-row justify-left gap-3 my-2 font-bold print:hidden">
      <NavLink to="/" className={""}>Search</NavLink>
      <NavLink to="/handout" className={""}>Handout</NavLink>
      <NavLink to="/meta" className={""}>Info</NavLink>
    </nav>
    <hr className="print:hidden"/>
  </div>
  )
}
