# Job Lookup [WIP]

## Presentation

A simple app to manage your job search notes

## Todo

- Add stats data:
  - companies:
    - number of applications, last event, next event
  - contact:
    - last event, next event
  - job application:
    - number of event, last event, next event

- Some elements have a black background for no reason (seems to be a bug in vuetify)

## Binaries used and versions

- go version 1.21
- gcc 13.2.0
- npm 10.5.0
- wails 2.9.1 (additional dependencies: <https://wails.io/docs/gettingstarted/installation/>)
- sqlc 1.26.0

## Env vars

```sh
CGO_ENABLED=1 # for build
JOB_LOOKUP_TEST=TRUE # for local testing, replace the db with fixtures
```

## Generate SQL bindings from scripts in queries dir and schema.sql

```sh
sqlc generate
```

## Test

`wails dev`

## Build

`wails build`

## IDE configuration

<https://wails.io/docs/next/guides/vscode/>
<https://wails.io/docs/guides/ides/>
