# Sloth common sli plugins

[![CI](https://github.com/slok/sloth-common-sli-plugins/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/slok/sloth-common-sli-plugins/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/slok/sloth-common-sli-plugins)](https://goreportcard.com/report/github.com/slok/sloth-common-sli-plugins)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/slok/sloth-common-sli-plugins/master/LICENSE)

## Introduction

A collection of common [Sloth] [SLI plugins][sli-plugins] ready to be used on Sloth manifests, or used as examples to create your own.

## Getting started

One of the [Sloth] features are the [SLI plugins][sli-plugins]. These can be loaded dynamically when Sloth is executed and the SLO manifests can reference these plugins and pass some options to set an SLI, instead of writing the SLI query.

To use these plugins you could do:

```bash
# Get the plugins.
git clone https://github.com/slok/sloth-common-sli-plugins.git

# Load the plugins and use them (Sloth can load plugins from multiple dirs).
sloth generate -p ./sloth-common-sli-plugins -i {MY_SLO_MANIFEST}
```

## Plugins

- [Noop](./plugins/noop): A plugin example/placeholder that doesn't do anything.

## Contributing

You can contribute with new plugins the same way the ones in [plugins](./plugins), the process would be this:

- Create a directory group in [plugins](./plugins), only if the group would have more than one SLI plugin (e.g `app-x`, `protocol-y`, `library-z`...).
- Create a directory for each plugin: (e.g `availability`, `latency`...).
- Create a `plugin.go` file with the plugin and a `plugin_test.go` for the unit tests.
- Create an sloth manifest in [`test/integration`](./test/integration) to test that sloth can load and use this plugin correctly.
- Add a `README.md` on the group and/or the plugin dir and reference the plugin in this readme in [plugins](#plugins) section. the readme should have for each plugin at least ([example](./plugins/noop/README.md)):
  - Introduction.
  - Options.
  - Usage examples.

[sloth]: https://github.com/slok/sloth
[sli-plugins]: https://github.com/slok/sloth#sli-plugins
