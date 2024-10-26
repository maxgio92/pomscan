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
📦 com.google.guava.guava
artifact-id : guava
group-id : com.google.guava
pom-file : druid-handler/pom.xml
version : ${druid.guava.version}
version-property-name : druid.guava.version
version-property-value : 16.0.1
version-property-declare-path : druid-handler/pom.xml

📦 com.google.guava.guava
artifact-id : guava
group-id : com.google.guava
pom-file : itests/qtest-druid/pom.xml
version : ${druid.guava.version}
version-property-name : druid.guava.version
version-property-value : 16.0.1
version-property-declare-path : druid-handler/pom.xml

📦 com.google.guava.guava
artifact-id : guava
group-id : com.google.guava
pom-file : pom.xml
version : ${guava.version}
version-property-name : guava.version
version-property-value : 22.0
version-property-declare-path : pom.xml

📦 com.google.guava.guava
artifact-id : guava
group-id : com.google.guava
pom-file : standalone-metastore/pom.xml
version : ${guava.version}
version-property-name : guava.version
version-property-value : 22.0
version-property-declare-path : pom.xml

📦 com.google.guava.guava
artifact-id : guava
group-id : com.google.guava
pom-file : storage-api/pom.xml
version : ${guava.version}
scope : test
version-property-name : guava.version
version-property-value : 22.0
version-property-declare-path : pom.xml
```
