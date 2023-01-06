import { useState } from "react";

export default function About() {
  const supported: [string, [string, number][]][] = [
    [
      "USA",
      [
        ["AJHSME", 3413],
        ["AHSME", 3415],
        ["AMC 8", 3413],
        ["AMC 10", 3414],
        ["AMC 12", 3415],
        ["AIME", 3416],
        ["MPFG", 3427],
        ["MPFG Olympiad", 953466],
        ["USAMTS", 3412],
        ["BAMO", 233906],
      ],
    ],

    [
      "College-Hosted",
      [
        ["CHMMC", 2746308],
        ["CMIMC", 253928],
        ["HMMT", 3417],
        ["Nov HMMT", 2881068],
        ["SMT", 3418],
        ["BMT", 2503467],
        ["PUMAC", 3426],
      ],
    ],

    [
      "National Olympiads",
      [
        ["Canada MO", 3277],
        ["Korea MO", 3383],
        ["KJMO", 603052],
        ["China MO", 3284],
        ["China GMO", 3287],
        ["China Round 2", 3288],
        ["All-Russian Olympiad", 3371],
        ["USAJMO", 3420],
        ["USAMO", 3409],
        ["ELMO", 3429],
        ["Sharygin", 3372],
      ],
    ],
    [
      "IMO Team Selection Tests",
      [
        ["USA TST", 3411],
        ["USA TSTST", 3424],
        ["China TST", 3282],
        ["Korea TST", 3384],
      ],
    ],
    [
      "International Olympiads",
      [
        ["IMO", 3222],
        ["IMO Shortlist", 3223],
        ["APMO", 3226],
        ["RMM", 3238],
        ["Baltic Way", 3231],
        ["Balkan MO", 3225],
        ["JBMO", 3227],
        ["EGMO", 3246],
      ],
    ],
  ];
  const [open, setOpen] = useState(false);
  return (
    <div className="my-3 mx-0 text-sm">
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
        Over <strong>17000</strong> Problems.
      </p>
      <div className="mx-0 border-gray-200 border p-1 rounded-lg transition-all transform duration-1000">
        <button
          onClick={(e) => {
            e.preventDefault();
            setOpen(!open);
          }}
          className="my-2 py-1 w-full text-left font-extrabold flex justify-between text-sm"
        >
          <span>Supported Contests:</span>
          <span className="mx-2">{open ? "-" : "+"}</span>
        </button>
        {open
          ? supported.map((category) => (
              <>
                <h1 className="text-lg font-bold text-center mt-2">
                  {category[0]}
                </h1>
                <hr />
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 mt-2 text-xs">
                  {category[1].map((contest, i) => (
                    <a
                      className="p-[1px] font-bold"
                      href={`https://artofproblemsolving.com/community/c${contest[1]}`}
                      target="_blank"
                      key={i}
                    >
                      {contest[0]}
                    </a>
                  ))}
                </div>
              </>
            ))
          : null}
      </div>
    </div>
  );
}
