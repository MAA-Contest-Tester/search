{
  description = "Search.MAATester.com";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs = { self, nixpkgs, flake-utils }:
    with flake-utils.lib;
    eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "0.1.0";
        name = "Search.MAATester.com";
      in {

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [ go yarn scc redis meilisearch gopls ];
        };

        defaultPackage = pkgs.stdenv.mkDerivation { inherit name version; };

        formatter = nixpkgs.legacyPackages."${system}".nixfmt;

      });
}
