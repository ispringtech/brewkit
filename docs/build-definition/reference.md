# BrewKit build-definition reference

BrewKit build-definition is a JSON schema.

Full build-definition described in [brewkit/v1](/data/specification/build-definition/v1.json) JSON schema

## apiVersion

API version of build definition. Declared for backward compatibility

## Vars

Build-time variables that calculates in build time

`Vars` supports following directives:
* [from](#from)
* [platform](#platform)
* [workdir](#workdir)
* [env](#env)
* [cache](#cache)
* [copy](#copy)
* [secrets](#secrets)
* [network](#network)
* [ssh](#ssh)
* [command](#command)

## Target

Executable build targets

`Targets` supports following directives:
* [from](#from)
* [dependsOn](#dependsOn)
* [platform](#platform)
* [workdir](#workdir)
* [env](#env)
* [cache](#cache)
* [copy](#copy)
* [secrets](#secrets)
* [network](#network)
* [ssh](#ssh)
* [command](#command)
* [output](#output)

### Composite targets

You can define target, that compose other targets as dependencies:

```jsonnet
    targets: {
        // when runs build - gobuild and golint will run sequentially    
        build: ['gobuild', 'golint'],
        
        gobuild: {},
        golint: {},
    }
```

## Directives

### From

Defines base target or image for target.

Use image:
```jsonnet
     targets: {
        gobuild: {
            // use as base go image         
            from: "golang:1.21.1"
        },
    }
```

Use target:
```jsonnet
     targets: {
        gocompiler: {
            from: "golang:1.21.1",
            // ..
        }
        
        gobuild: {
            // use as base gocompiler target         
            from: "gocompiler"
        },
    }
```

### DependsOn

Target may require to subsequently run another target before the execution. 
So, define `dependsOn` to arrange targets execution

```jsonnet
    targets: {
        gogenerate: {}
        
        // Before running gobuild brewkit will run gogenerate         
        gobuild: {
            dependsOn: ['gogenerate']
        }
    }
```

### Platform

Define platform for target. Supported platforms that supported by buildkit.

Underlying used buildkit platform directive - [buildkit-docs](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md#from)

```jsonnet
    targets: {
        gobuild: {
            from: "golang:1.21.1",
            platform: "linux/amd64"
        },
    }
```

### Workdir

Working directory for target or var


```jsonnet
    targets: {
        gobuild: {
            from: "golang:1.21.1",
            platform: "linux/amd64"
        },
    }
```

### Env

Describes env for target or var in JSON map format

```jsonnet
    targets: {
        gobuild: {
            from: "golang:1.21.1",
            env: {
                GOCACHE: "/app/cache/go-build",
                APP_ID: "contentservice"
            },
        },
    }
```

### Cache

Describes buildkit [cache](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md#run---mounttypecache).

Cache can be reused between builds.

**Cache directive do not affect on container layer**

You can define cache in one-liner via [cache jsonnet extension](jsonnet-extensions.md#cache)

Define cache
```jsonnet
// import cache for single line cache declaration
local cache = std.native('cache');
//...
    targets: {
        gobuild: {
            // set up cache with id go-build and path /app/cache in container
            cache: cache("go-build", "/app/cache"),
        },
    }
```

You can define multiple caches
```jsonnet
local cache = std.native('cache');
//...
    targets: {
        gobuild: {
            cache: [
                cache("go-build", "/app/cache"),
                cache("go-mod", "/go/pkg/mod"),
            ]
        },
    }
```

### Copy

Copy files from host fs into container fs

You can define cache in one-liner via [copy jsonnet extension](jsonnet-extensions.md#copy)

```jsonnet
// import copy for single line copy declaration
local copy = std.native('copy');
//...
    targets: {
        gobuild: {
            copy: [
                copy('cmd', 'cmd'),
                copy('pkg', 'pkg')
            ],
        },
    }
```

### Secrets

Use file as secret in container without copying it into container.
**So file will not be copied into container layer**

To use secret in brewkit target you should define **secret in brewkit config**. See [config overview](/docs/config/overview.md)

You can define cache in one-liner via [copy jsonnet extension](jsonnet-extensions.md#secret)

```jsonnet
// import copy for single line copy declaration
local secret = std.native('secret');
//...
    targets: {
        gobuild: {
            // secret with id sould be already defined id brewkit config            
            secret: secret("aws", "/root/.aws/credentials"),
        },
    }
```

Define secret in `~/.brewkit/config`
```jsonnet
{
    "secrets": [
        {
            "id": "aws",
            // path may contain env variables            
            "path": "${HOME}/.aws/credentials"
        },
    ]
}
```

### SSH

Defines access to ssh socket from host. BrewKit mounts ssg agent from `$SSH_AUTH_SOCK` into container via buildkit [ssh mount](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md#run---mounttypessh)

Ssh defined for now as empty object for further customization
```jsonnet
    targets: {
        gomod: {
            ssh: {},            
        },
    }
```

### Network

Define network for target. Supported all network which supported by docker

### Command

Command to be run in stage. Command runs **in container shell**, **not in exec**

```jsonnet
    targets: {
        gobuild: {
            command: 'go build -o ./bin/brewkit ./cmd/brewkit',            
        },
    }
```

Vars values can be used in commands

```jsonnet
    vars: {
        gitcommit: {}
    }
    
    targets: {
        gobuild: {
            // pass gitcommit to go ldflags             
            command: 'go build -ldflags "-X main.Commit=${gitcommit}" -o ./bin/brewkit ./cmd/brewkit',            
        },
    }
```

### Output

Output artifacts from targets. Artifacts exported with current user id, **so no root owned artifacts**

```jsonnet
    targets: {
        gobuild: {
            command: 'go build -o ./bin/brewkit ./cmd/brewkit',
            output: {
                // export /app/bin/brewkit from container               
                artifact: "/app/bin/brewkit",
                // export to ./bin folder                
                "local": "./bin"
            }            
        },
    }
```