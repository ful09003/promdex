# What is promdex?

Promdex is a self-hosted utility meant to add context and understanding to Prometheus metrics beyond what is exposed from exporters themselves.

## How does promdex work?
Under the hood, Promdex has a few concepts (`what`s) which are important to understand before explaining `how`:

- Metastores, which are one of a selection of storage providers for Promdex to use. At this time, Promdex only supports SQLite.
- Targets, which represent a Prometheus instance (**not** the scraped targets of a Prometheus instance)

On start, Promdex checks for an existing `metastore` at the user-defined path. If the metastore does not exist, Promdex creates it and uses the Prometheus [metadata API](https://prometheus.io/docs/prometheus/latest/querying/api/#querying-metric-metadata) to load initial metadata about metrics.

From that point forward, Promdex offers HTTP endpoints to create (POST) or retrieve (GET) additional metadata. These endpoints are offered at the following path: `{promdex_url}:{promdex_port}/{job}/{metric_name}`. It is important to note that job and metric_name are required; Promdex strictly disallows bulk operations such as `GET {promdex_url}:{promdex_port}/{job}`

## Development
To run all of Promdex's tests, stand up a Prometheus listening to localhost:9090 and use `go test --tags "all_tests localtests"`. `localtests` is a build constraint which runs tests targeting a 'live' Prometheus. Thus:

- If testing new/existing code that does *not* require a live Prometheus add `//+build all_tests` to the top of the test file
- If testing new/existing code that *does* require a live Prometheus add `//+build all_tests localtests` to the top of the test file