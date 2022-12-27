import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import About from "./components/About";
import Search from "./components/Search";
import "./index.css";

function App() {
  return (
    <div className="max-w-full clamp">
      <main className="m-2 p-3">
        <h1 className="font-extrabold text-2xl sm:text-4xl md:text-5xl">
          <span className="text-blue-800">Search.</span>
          <span>MAATester.com</span>
        </h1>
        <About />
        <Search />
      </main>
    </div>
  );
}

if (import.meta.env.DEV) {
  ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <StrictMode>
      <App />
    </StrictMode>
  );
} else {
  ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <App />
  );
}
