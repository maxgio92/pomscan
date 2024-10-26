# Pomscan

**Scan POM files** for dependencies.

## Usage

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

## Usage with bumps

### With [`pombump`](https://github.com/chainguard-dev/pombump)

It can be really useful when preparing patches, for example with [`pombump`](https://github.com/chainguard-dev/pombump).

Consider that from the previous example we want to bump `com.google.guava.guava` to the version *24.1.1-jre* because both *22.0* and *16.0.1* both contain CVEs, we now know thanks to `pomscan` that we need to change the following version properties:
* `guava.version` in the root project's `pom.xml`
* `druid.guava-version` in the Druid Handler project's `druid-handler/pom.xml`

So, we can run `pombump` to update the `pom.xml` files accordingly, and we feeed it with the information retrieved from `pomscan`, like below:

```shell
$ pombump --properties="guava.version@24.1.1-jre" pom.xml >pom.bumps.xml
...
2024/10/26 13:04:54 INFO Patching property: guava.version from 22.0 to 24.1.1-jre
$ pombump --properties="druid.guava.version@24.1.1-jre" druid-handler/pom.xml >druid-handler/pom.bumps.xml
...
2024/10/26 13:06:05 INFO Patching property: druid.guava.version from 16.0.1 to 24.1.1-jre
```

We now have the new POMs updated at `pom.bumps.xml` and `druid-handler/pom.bumps.xml`, that will set the new version for all the occurrences of the `com.google.guava.guava` artifact across all the direct dependencies of the Maven project and subprojects.

