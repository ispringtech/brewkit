# BrewKit

----

Container-native build system. Meaning that BrewKit is a **single tool** that you should use(across Docker) to build project.

BrewKit focuses on repeatable builds, host agnostic and build process containerization

----

## Key features

* [BuildKit](https://github.com/moby/buildkit) as core of build system. There is following features from BuildKit 
  * Distributed cache - inherited from BuildKit
  * Automatic garbage collection - inherited from BuildKit
* Aggressive-caching
* Mounted secrets
* Host configuration agnostic
* Output artifacts to Host filesystem
* JSON based build-definition format  
* [JSONNET](https://jsonnet.org/) configuration language - BrewKit uses jsonnet build to compile `brewkit.jsonnet` into JSON build-definition

## Naming

`BrewKit` - common style
<br/>
`brewkit` - go-style

## Start with BrewKit

Install BrewKit via go >= 1.20 
```shell
go install github.com/ispringtech/brewkit
```

Create `brewkit.jsonnet`
```shell
touch brewkit.jsonnet
```

Describe simple target
```jsonnet
local app = "service";

local copy = std.native('copy');

{
    apiVersion: "brewkit/v1",
    targets: {
        all: ['gobuild'],
        
        gobuild: {
            from: "golang:1.20",
            workdir: "/app",
            copy: [
                copy('cmd', 'cmd'),
                copy('pkg', 'pkg'),
            ],
            command: std.format("go build -o ./bin/%s ./cmd/%s", [app])
        }
    }
}
```

Run build
```shell
brewkit build

 => [internal] load build definition from Dockerfile                                                                                                                                                           0.1s
 => => transferring dockerfile: 3.45kB                                                                                                                                                                         0.0s
 => [internal] load .dockerignore                                                                                                                                                                              0.1s
 => => transferring context: 2B
# ...
```

## Build BrewKit

When brewkit installed locally
```shell
brewkit build
```

Build from source:
```shell
go build -o ./bin/brewkit ./cmd/brewkit
```

## Documentation and examples

* [Documentaion entrypoint](docs/readme.md)
* [Examples](docs/examples)