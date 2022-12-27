export default function About() {
  return (
    <div className="my-3 mx-0 text-sm">
      <p>
        An Actually Fast Search Engine for Math Contest Problems, scraped from
        the AoPS Wiki.
      </p>
      <p>
        {" "}
        Written by <a href="https://github.com/junikimm717">Juni Kim.</a> See
        the{" "}
        <a href="https://github.com/MAA-Contest-Tester/search">Source Code</a>.
      </p>
      <p className="mt-3">
        Supports{" "}
        <strong>A(J)HSME, AMC 8/10/12, AIME, USA(J)MO, IMO, and JBMO</strong>.
        Over <strong>6500</strong> Problems.
      </p>
    </div>
  );
}
