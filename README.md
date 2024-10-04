## Pomscan

**Scan POM files**

```
Usage:
  pomscan [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dep         Get info about a dependency
  help        Help about any command

Flags:
      --debug                 Sets log level to debug
  -h, --help                  help for pomscan
  -p, --project-path string   Project path (default ".")

Use "pomscan [command] --help" for more information about a command.
```

### `dep` command

The `dep` commands returns info about a dependency, recursively scanning the POM files of a project.

```
Get info about a dependency

Usage:
  pomscan dep [flags]

Flags:
  -a, --artifact-id string   Filter by artifact ID.
  -g, --group-id string      Filter by group ID. It must be combined with artifact ID.
  -h, --help                 help for dep
      --version-only         Print only matches that have the version set. It supports variables.

Global Flags:
      --debug                 Sets log level to debug
  -p, --project-path string   Project path (default ".")
```

**Example**

```
$ pomscan dep -a guava -p . --version-only
📦 com.google.guava.guava found
pom: ./druid-handler/pom.xml
version: 16.0.1

📦 com.google.guava.guava found
pom: ./itests/qtest-druid/pom.xml
version: 16.0.1

📦 com.google.guava.guava found
pom: ./storage-api/pom.xml
version: 22.0
scope: test
```
