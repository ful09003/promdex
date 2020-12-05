# `What` is promdex?

Promdex is a self-hosted utility meant to add context and understanding to Prometheus metrics beyond what is exposed from exporters themselves.

## `Why` does Promdex need to exist?

Ultimately, metric descriptions in the Prometheus world are hit-and-miss in terms of quality. I've often found myself going down rabbit holes to answer questons like "_how_ is this metric derived? _where_ does it come from? _when_ would I need this?", to the point that it became a good idea to standardize the data model and automate a small bit of the investigation process.

## `How` does promdex work?
Under the hood, Promdex has a few concepts (`what`s) which are important to understand before explaining `how`:

- Metastores, which are one of a selection of storage providers for Promdex to use. At this time, Promdex only supports SQLite.
- Targets, which represent a Prometheus instance (**not** the scraped targets of a Prometheus instance)

On start, Promdex checks for an existing `metastore` at the user-defined path and of the user-defined type. If the metastore does not exist, Promdex creates it and uses the Prometheus [metadata API](https://prometheus.io/docs/prometheus/latest/querying/api/#querying-metric-metadata) to load initial metadata about metrics.

From that point forward, Promdex offers HTTP endpoints to create (POST) or retrieve (GET) additional metadata. These endpoints are offered at the following path: `{promdex_url}:{promdex_port}/{job}/{metric_name}`. It is important to note that job and metric_name are required; Promdex strictly disallows bulk operations such as `GET {promdex_url}:{promdex_port}/{job}`

## Development
To run all of Promdex's tests, stand up a Prometheus listening to localhost:9090 and use `go test --tags "all_tests localtests"`. `localtests` is a build constraint which runs tests targeting a 'live' Prometheus. Thus:

- If testing new/existing code that does *not* require a live Prometheus add `//+build all_tests` to the top of the test file
- If testing new/existing code that *does* require a live Prometheus add `//+build all_tests localtests` to the top of the test file

## Problems that still need solving:
- Tests for at least all happy execution paths
- Clean up data model; we need created/updated times for the metric as an example
- HTTP path handling (add POST, probably DELETE)
- Consider another builtin metastore (maybe object storage?)
- Metadata discovery feels like it should dedupe metrics when their descriptions are the same, e.g.
```
sqlite> select * from promdexMetrics where metric_name="process_open_fds";
metric_name       metric_instance_job  metric_source_desc              
----------------  -------------------  --------------------------------
process_open_fds  node_exporter        Number of open file descriptors.
process_open_fds  prometheus           Number of open file descriptors.
```
- On that topic, Promdex should also be able to build metadata from any number of Prometheus instances