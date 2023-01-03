import { useState } from "react";
import Prompt from "./Prompt";

export default function Bank() {
  const [count, setCount] = useState(0);
  return (
    <>
      <Prompt key={count} streak={count} />
      <button
        onClick={() => {
          setCount(count + 1);
        }}
        className={"p-2 text-white my-3 bg-blue-700 mx-0"}
      >
        Next
      </button>
    </>
  );
}
