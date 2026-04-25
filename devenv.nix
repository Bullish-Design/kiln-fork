{ pkgs, ... }:

{
  languages.go = {
    enable = true;
    package = pkgs.go_1_25;
  };

  languages.javascript = {
    enable = true;
    npm.enable = true;
  };

  packages = with pkgs; [
    gopls
    gofumpt
    golines
    goimports-reviser
    templ
    air

    nodejs_22
    prettier
    tailwindcss_4

    libavif
    libwebp
  ];

  scripts.setup.exec = ''
    npm install
  '';

  scripts.test.exec = ''
    go test ./...
  '';

  scripts.lint.exec = ''
    go vet ./...
  '';

  scripts.fmt.exec = ''
    gofumpt -w .
  '';

  scripts.templ.exec = ''
    templ generate ./internal/templates
  '';

  scripts."tailwind-default".exec = ''
    tailwindcss -i ./assets/default_style_input.css -o ./assets/default_style.css
  '';

  scripts."tailwind-simple".exec = ''
    tailwindcss -i ./assets/simple_style_input.css -o ./assets/simple_style.css
  '';

  scripts.assets.exec = ''
    tailwind-default
    tailwind-simple
  '';

  scripts.build.exec = ''
    go build -o ./kiln ./cmd/kiln
  '';

  scripts.run.exec = ''
    if [ ! -x ./kiln ]; then
      go build -o ./kiln ./cmd/kiln
    fi
    exec ./kiln "$@"
  '';

  scripts.generate.exec = ''
    go run ./cmd/kiln generate "$@"
  '';

  scripts.serve.exec = ''
    go run ./cmd/kiln serve "$@"
  '';

  scripts.dev.exec = ''
    go run ./cmd/kiln dev "$@"
  '';

  scripts.doctor.exec = ''
    go run ./cmd/kiln doctor "$@"
  '';

  scripts."demo-test-setup".exec = ''
    mkdir -p ./demo/.gocache ./demo/.gomodcache ./demo/logs ./demo/out ./demo/public ./demo/tmp
    export GOCACHE="$PWD/demo/.gocache"
    export GOMODCACHE="$PWD/demo/.gomodcache"
    go build -o ./demo/out/kiln ./cmd/kiln
  '';

  scripts."demo-test-run".exec = ''
    demo-test-setup
    export GOCACHE="$PWD/demo/.gocache"
    export GOMODCACHE="$PWD/demo/.gomodcache"

    ./demo/out/kiln doctor --input ./demo/vault > ./demo/logs/demo_doctor.log 2>&1
    ./demo/out/kiln generate --input ./demo/vault --output ./demo/public > ./demo/logs/demo_generate.log 2>&1

    echo "Demo test run complete. See demo/logs/demo_doctor.log and demo/logs/demo_generate.log"
  '';

  scripts.clean.exec = ''
    rm -rf ./kiln ./public ./tmp
    go clean -cache -testcache
  '';

  scripts.air.exec = ''
    air
  '';

  enterShell = ''
    export PATH="$PWD/node_modules/.bin:$PATH"
    echo "Kiln dev shell ready. Run: setup, test, build, dev, serve, generate"
  '';
}
