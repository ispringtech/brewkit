# CLI

BrewKit has command-line-interface for running build, manipulate cache and config.
<br/>
You can learn more about each command by running command with `-h` flag

## build

Manipulates builds

| Command          | Description                                                                                 |
|------------------|---------------------------------------------------------------------------------------------|
| <target-name>    | Runs specified target                                                                       |
| definition       | Print full parsed and verified build-definition in JSON to stdout                           |
| definition-debug | Print compiled build definition in raw JSON, useful for debugging complex build definitions |

Examples:

Build concrete targets
```shell
brewkit build generate compile
```

## config

Manipulate host config

| Command | Description                   |
|---------|-------------------------------|
| init    | Create default brewkit config |

## version

Print BrewKit version

```shell
brewkit version
```

## cache

Manipulate buildkit cache

| Command  | Description                             |
|----------|-----------------------------------------|
| clear    | Clear docker builder cache              |
| clear -a | Clear all cache, not just dangling ones |

## fmt

Pretty format jsonnet files with jsonnetfmt

```shell
brewkit fmt brewkit.jsonnet
```