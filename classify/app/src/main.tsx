import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import Header from "./components/Header";
import Bank from "./components/Bank";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <div className="mx-auto max-w-3xl p-3 mt-5">
      <Header />
      <Bank />
    </div>
  </React.StrictMode>
);
