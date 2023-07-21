import React, { NavLink } from "react-router-dom";

export default function Navbar() {
  return (
  <>
    <nav className="flex flex-row justify-left gap-2 my-2 font-bold print:hidden">
      <NavLink to="/" className={""}>Home</NavLink>
      <NavLink to="/handout" className={""}>Handout</NavLink>
    </nav>
    <hr className="my-2 mb-5 print:hidden"/>
  </>
  )
}
