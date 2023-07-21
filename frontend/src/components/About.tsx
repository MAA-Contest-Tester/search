export default function About() {
  return (
    <div className="my-3 mx-0 text-sm print:hidden">
      <p>
        An instant search engine for math olympiad questions. Problems sourced
        from the AoPS Community.
      </p>
      <p>
        {" "}
        Written by{" "}
        <a href="https://github.com/junikimm717" target="_blank">
          Juni Kim.
        </a>{" "}
        See the{" "}
        <a href="https://github.com/MAA-Contest-Tester/search" target="_blank">
          Source Code
        </a>
        .
      </p>
      <p className="mt-3">
        Over <strong>17000</strong> Problems. <strong>Instant Handouts</strong>{" "}
        with printer friendliness.
      </p>
    </div>
  );
}
