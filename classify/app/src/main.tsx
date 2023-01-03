import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import Prompt from "./components/Prompt";
import Header from "./components/Header";

function App() {
  return <>Hello</>;
}

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <div className="mx-auto max-w-3xl p-3 mt-5">
      <Header />
      <Prompt />
    </div>
  </React.StrictMode>
);
