# Build definition

Build definition describes build-time `Vars` and `Targets` for build

## Overview

Each build-definition has top-level blocks:
* **apiVersion** - describes build-definition apiVersion for backward compatibility
* **vars** - build-time variables that calculates in build time 
* **targets** - executable build targets

Build definition with vars and targets
```jsonnet
local app = "service";

local copy = std.native('copy');

{
    apiVersion: "brewkit/v1",
    
    vars: {
        gitcommit: {
            from: "golang:1.20",
            workdir: "/app",
            copy: copy('.git', '.git'),
            command: "git -c log.showsignature=false show -s --format=%H:%ct"
        }
    },
    
    targets: {
        all: ['gobuild'],
        
        gobuild: {
            from: "golang:1.20",
            workdir: "/app",
            copy: [
                copy('cmd', 'cmd'),
                copy('pkg', 'pkg'),
            ],
            // Use gitcommit as go ldflag            
            command: std.format('go build -ldflags "-X main.Commit=${gitcommit}" -o ./bin/%s ./cmd/%s', [app])
        }
    }
}
```

## Vars

Define a var
```jsonnet
// ...
    vars: {
        gitcommit: {
            from: "golang:1.20",
            workdir: "/app",
            copy: copy('.git', '.git'),
            command: "git -c log.showsignature=false show -s --format=%H:%ct"
        }
    }
// ...
```

Now you can reference it in target. Variable reference format: `${VAR}`
```jsonnet
    targets: {
        gobuild: {
            from: "golang:1.20",
            workdir: "/app",
            copy: [
                copy('cmd', 'cmd'),
                copy('pkg', 'pkg'),
            ],
            // Use ${gitcommit} as go ldflag            
            command: std.format('go build -ldflags "-X main.Commit=${gitcommit}" -o ./bin/%s ./cmd/%s', [app])
        }
    }
```

_Notes about variables_:
* You should clearly separate JSONNET variables, which are calculated when build-definition compiles, and brewkit vars, which calculates in build time. 
* Vars calculate before build starts

### Deference between jsonnet variables

You can define variable in jsonnet and use it in build definition. As in top example with `local app = "service";`

But jsonnet variables can't be changed due to build-time, for example git commit hash or smth else.

So, when you need build-time variables, you can use `vars`. As in top example with `gitcommit` variable.

## API version

Each brewkit build-definition should satisfy build-definition apiVersion.
All `apiVersion` schemas placed - [build-definition](/data/specification/build-definition)

## All target

`All` is special reserved target name which runs when no concrete target name passed.

So when run `brewkit build` brewkit executes `all` target.

## jsonnet

BrewKit use jsonnet as build-definition language. 

As jsonnet is extension of JSON you can write build-definition in JSON and pass to brewkit.

Also, all features of [jsonnet](https://github.com/google/go-jsonnet) are supported

## std.native('copy') - extension functions

JSONNET extension functions can be used to simplify writing build-definition.
<br/>
Extension functions can be accessed via `std.native('<ext>')`

List of jsonnet extension functions - [jsonnet-extensions](jsonnet-extensions.md)

