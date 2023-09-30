# Using BuildKit as a build backend

# Context and Problem Statement

Using a direct container run(via `docker run...` or running docker containers) make brewkit highly dependent on docker daemon.
We want to avoid this dependency, as BrewKit does not need to run containers
It can only be used where there is a docker daemon. Although containers themselves as processes (as for kubernetes, for example) are not needed for Brewkit.

# Decision Drivers

- BrewKit should not depend on the docker daemon
- When using BuildKit, you can get fancy features (caching, rootless artifacts, access to ssh-agent inside the builder) out of the box.

# Considered Options

- Using `docker run ...`
- Using BuildKit

# Decision Outcome

The chosen variant is "Use BuildKit."  Since it is a more modern way to build projects inside the container.

# Pros and Cons of the Options

## Using `docker run ...`.

All project build commands will run inside a container that is started via a direct request to the docker daemon. The container is started and a command is run inside it.

* Good, simple solution
* Bad, dependency on the Docker daemon
* Bad, no caching depending on code changes in the project

## Using BuildKit

Send build request to buildkit to build an image with specific instructions

* Good, no root privileges needed
* Good, already integrated in new versions of Docker to build images.
* Bad, BuildKit not in final version yet (as of 05.2023)
* Bad, BuildKit can't export build cache to another folder (relevant for Golang where package cache is not in project directory)

## More Information

## More information about BuildKit

- [BuildKit github](https://github.com/moby/buildkit)
- [Docker BuildKit documentation](https://docs.docker.com/build/buildkit/)

## Independence from the docker daemon

`BuildKit` is a separate part of docker project used to build images
`Buildkit` itself **does not require root privileges**, and can be shipped completely independent of `Docker`.

Independence from `Docker` and root privileges provides the following benefits:
- Easier and safer _Docker-in-Docker_: you don't need to pass a socket your docker into a container to build images inside the container, you can use a BuildKit image and build images through it
- Running a container with BuildKit does not require `priveledged``: more secure image building in containers