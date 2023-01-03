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
      <footer className="text-center text-gray-600 text-sm">
        &copy; Juni Kim 2023-{new Date().getFullYear()}
      </footer>
    </div>
  </React.StrictMode>
);
