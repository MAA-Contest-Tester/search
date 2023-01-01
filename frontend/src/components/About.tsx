export default function About() {
  const supported = [
    ["AJHSME", 3413],
    ["AHSME", 3415],
    ["AMC 8", 3413],
    ["AMC 10", 3414],
    ["AMC 12", 3415],
    ["AIME", 3416],

    ["CHMMC", 2746308],
    ["CMIMC", 253928],
    ["HMMT", 3417],
    ["Nov HMMT", 2881068],
    ["SMT", 3418],
    ["BMT", 2503467],
    ["PUMAC", 3426],
    ["BAMO", 233906],
    ["USAMTS", 3412],

    ["USAJMO", 3420],
    ["USAMO", 3409],
    ["JBMO", 3227],
    ["Balkan MO", 3225],
    ["Sharygin", 3372],
    ["USA TST", 3411],
    ["USA TSTST", 3424],
    ["China TST", 3282],
    ["EGMO", 3246],
    ["IMO", 3222],
    ["ELMO", 3429],
    ["APMO", 3226],
    ["IMO Shortlist", 3223],
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
        Over <strong>14000</strong> Problems. Supported Contests:
      </p>
      <div className="grid grid-cols-3 sm:grid-cols-5 md:grid-cols-7 mx-0 border-gray-200 border p-1 rounded-lg">
        {supported.map((contest, i) => (
          <a
            className="p-[1px] font-bold"
            href={`https://artofproblemsolving.com/community/c${contest[1]}`}
            key={i}
          >
            {contest[0]}
          </a>
        ))}
      </div>
    </div>
  );
}
