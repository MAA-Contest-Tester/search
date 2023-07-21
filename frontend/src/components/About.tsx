export default function About() {
  return (
    <div className="my-3 mx-0 text-sm print:hidden">
      <p>
        An instant search engine for math olympiad questions. 
        Over <strong>17000</strong> Problems.
        Problems sourced from the Art of Problem Solving Community.
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
      </p>
    </div>
  );
}
