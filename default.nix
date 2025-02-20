{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  pname = "hyprls";
  version = "0.5.2";
  pwd = ./.;
  src = ./.;

  postBuild = ''
    rm $GOPATH/bin/generate
  '';

  modules = ./gomod2nix.toml;
  checkFlags = ["-skip=TestHighLevelParse"]; # not yet implemented
}
