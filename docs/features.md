# Features overview

## Aggressive-caching

Aggressive caching lets brewkit achieve faster repeated builds.
When running target BrewKit (or BuildKit under the hood) scans file changes from `copy` and changes
from previous targets to decide if target needs to be executed.

So if there is no dependency changed and command is the same - there is no need to execute target.

Exmaple:

Build-definition for simple go service
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

First launch will pull all images, copy artifacts and runs commands

```shell
brewkit build
```

Output:
```shell
 => [gobuild 4/6] COPY  cmd cmd
 => [gobuild 5/6] COPY  internal internal
 => [gobuild 6/6] RUN   <<EOF go build -o ./bin/service ./cmd/service
```

Second launch will skip target gobuild since there is no changes in dependencies and sources in `cmd` and `pkg` 

```shell
 => CACHED [gobuild 4/6] COPY  cmd cmd
 => CACHED [gobuild 5/6] COPY  internal internal
 => CACHED [gobuild 6/6] RUN   <<EOF go build -o ./bin/service ./cmd/service
```

## Secrets

When executing ssh, aws, kubectl commands you need to pass a secret to get access to resources.
* SSH - agent socket
* aws-cli - `~/.aws` 
* kubectl - `~/.kube`

If you simply copy secret into target it can leak into docker images layer.

So you can mount secret into target without copying into image layer.

Define secret in `~/.brewkit/config`:
```jsonnet
{
    "secrets": [
        {
            // unique id for secret            
            "id": "aws",
            "path": "${HOME}/.aws/credentials"
        },
    ]
}
```

Use secret in target:
```jsonnet
local secret = std.native('secret');
//...
    targets: {
        gobuild: {
            // secret with id sould be already defined id brewkit config            
            secret: secret("aws", "/root/.aws/credentials"),
        },
    }
```
Based on [BuildKit secrets](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md#run---mounttypesecret)