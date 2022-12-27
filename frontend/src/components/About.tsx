export default function About() {
  return (
    <div className="my-3 mx-0 text-sm">
      <p>
        A fast search engine for browsing math problems to try. All problems
        scraped from the AoPS wiki.
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
        Supports{" "}
        <strong>
          AJHSME, AHSME, AMC 8, AMC 10, AMC 12, AIME, USAJMO, USAMO, JBMO, and
          IMO.
        </strong>
      </p>
      <p>
        Over <strong>6500</strong> Problems.
      </p>
    </div>
  );
}
