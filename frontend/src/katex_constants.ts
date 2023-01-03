const delimiters = [
  { left: "$$", right: "$$", display: true },
  { left: "\\(", right: "\\)", display: false },
  {
    left: "\\begin{equation}",
    right: "\\end{equation}",
    display: true,
  },
  {
    left: "\\begin{equation*}",
    right: "\\end{equation*}",
    display: true,
  },
  { left: "\\begin{align*}", right: "\\end{align*}", display: true },
  { left: "\\begin{align}", right: "\\end{align}", display: true },
  { left: "\\begin{alignat}", right: "\\end{alignat}", display: true },
  { left: "\\begin{gather}", right: "\\end{gather}", display: true },
  { left: "\\begin{CD}", right: "\\end{CD}", display: true },
  { left: "\\[", right: "\\]", display: true },
  { left: "$", right: "$", display: false },
  { left: "\\(", right: "\\)", display: false },
];

const macros = {
  "\\emph": "\\textit",
  "\\textsc": "",
  "\\textdollar": "\\$",
  "\\overarc": "\\overgroup",
  "\\dfrac": "\\frac",
  "\\O": "\\empty",
};

export const preprocess = (s: string | undefined) => {
  if (!s) return "";
  // eqnarray and tabular modified
  const res = s
    .trim()
    .replace(new RegExp(/\\begin{eqnarray\*}/, "g"), "\\begin{align*}")
    .replace(new RegExp(/\\end{eqnarray\*}/, "g"), "\\end{align*}")
    .replace(new RegExp(/\\begin{tabular}(\[.*?\])?/, "g"), "\\begin{array}")
    .replace(new RegExp(/\\end{tabular}(\[.*?\])?/, "g"), "\\end{array}")
    .replace(new RegExp(/\\makebox(\[.*?\])?/, "g"), "\\begin{array}")
    .replace(new RegExp(/\\mbox/, "g"), "\\text");
  return res;
};

export const renderconfig = {
  delimiters,
  macros,
  trust: true,
  throwonerror: false,
  errorColor: "#cc0000",
};
