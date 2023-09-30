local images = import 'images.libsonnet';

local cache = std.native('cache');
local copy = std.native('copy');
local copyFrom = std.native('copyFrom');

// External cache for go compiler, go mod, golangci-lint
local gocache = [
    cache("go-build", "/app/cache"),
    cache("go-mod", "/go/pkg/mod"),
];

// Sources which will be tracked for changes
local gosources = [
    "go.mod",
    "go.sum",
    "cmd",
    "internal",
];

{
    project():: {
        apiVersion: "brewkit/v1",

        vars: {
            gitcommit: {
                from: images.golang,
                workdir: "/app",
                copy: copy('.git', '.git'),
                command: "git -c log.showsignature=false show -s --format=%H:%ct"
            }
        },

        targets: {
            all: ["build", "check", "modulesvendor"],

            gosources: {
                from: "scratch",
                workdir: "/app",
                copy: [copy(source, source) for source in gosources]
            },

            gobase: {
                from: images.golang,
                workdir: "/app",
                env: {
                    GOCACHE: "/app/cache/go-build",
                },
                copy: copyFrom(
                    'gosources',
                    '/app',
                    '/app'
                ),
            },

            build: {
                from: "gobase",
                cache: gocache,
                workdir: "/app",
                dependsOn: ['modules'],
                command: std.format('
                    go build \\
                    -trimpath -v \\
                    -ldflags "-X main.Commit=${gitcommit} -X main.DockerfileImage=%s" \\
                    -o ./bin/brewkit ./cmd/brewkit
                ', [images.dockerfile]),
                output: {
                    artifact: "/app/bin/brewkit",
                    "local": "./bin"
                }
            },

            modules: {
                from: "gobase",
                cache: gocache,
                workdir: "/app",
                command: "go mod tidy",
                output: {
                    artifact: "/app/go.*",
                    "local": ".",
                },
            },

            // export local copy of dependencies for ide index
            modulesvendor: {
                from: "gobase",
                workdir: "/app",
                cache: gocache,
                dependsOn: ['modules'],
                command: "go mod vendor",
                output: {
                    artifact: "/app/vendor",
                    "local": "vendor",
                },
            },

            check: ["test", "lint"],

            test: {
                from: "gobase",
                workdir: "/app",
                cache: gocache,
                command: "go test ./...",
            },

            lint: {
                from: images.golangcilint,
                workdir: "/app",
                cache: gocache,
                copy: [
                    copyFrom(
                        'gosources',
                        '/app',
                        '/app'
                    ),
                    copy('.golangci.yml', '.golangci.yml'),
                ],
                env: {
                    GOCACHE: "/app/cache/go-build",
                    GOLANGCI_LINT_CACHE: "/app/cache/go-build"
                },
                command: "golangci-lint run"
            },
        }
    }
}