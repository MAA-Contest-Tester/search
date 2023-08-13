import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import About from "./components/About";
import {
  HandoutGenerator,
  HandoutProvider,
} from "./components/HandoutGenerator";
import Search from "./components/Search";
import "./index.css";
import React, {
  createBrowserRouter,
  Route,
  RouterProvider,
  Routes,
} from "react-router-dom";
import Navbar from "./components/Navbar";
import Metadata from "./components/Metadata";

const router = createBrowserRouter([{ path: "*", Component: Root }]);

function Root() {
  return (
    <div className="w-full min-h-screen px-2">
      <main className="clamp mx-auto py-0 mb-5">
        <HandoutProvider>
          <Navbar />
          <h1 className="font-extrabold text-3xl sm:text-4xl md:text-5xl print:hidden">
            <span className="text-blue-800">Search.</span>
            <span>MAATester.com</span>
          </h1>
          <Routes>
            <Route
              path="/"
              element={
                <>
                  <About />
                  <Search />
                </>
              }
            />
            <Route path="/handout" element={<HandoutGenerator />} />
            <Route path="/meta" element={<Metadata />} />
          </Routes>
        </HandoutProvider>
      </main>
    </div>
  );
}

function App() {
  return <RouterProvider router={router} />;
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
