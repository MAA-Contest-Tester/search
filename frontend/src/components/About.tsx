export default function About() {
  const supported = [
    "AJHSME",
    "AHSME",
    "AMC 8",
    "AMC 10",
    "AMC 12",
    "AIME",

    "CHMMC",
    "CMIMC",
    "HMMT",
    "SMT",
    "BMT",
    "PUMAC",
    "BAMO",
    "USAMTS",

    "USAJMO",
    "USAMO",
    "JBMO",
    "Balkan MO",
    "USA TST",
    "USA TSTST",
    "China TST",
    "EGMO",
    "IMO",
    "ELMO",
    "APMO",
    "IMO Shortlist",
  ];
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
        Over <strong>13000</strong> Problems. Supported Contests:
      </p>
      <div className="grid grid-cols-7 mx-0 border-blue-800 border p-1 rounded-lg">
        {supported.map((contest) => (
          <div className="p-[1px] font-bold">{contest}</div>
        ))}
      </div>
    </div>
  );
}
