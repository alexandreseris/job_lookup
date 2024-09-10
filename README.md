# Job Lookup [WIP]

## Presentation

A simple app to manage your job search notes

## Todo

- Some modification fail on save

- When renaming company_type or equivalent, the modification is not propagated to the relations (like companies), which leads to an inconsistent state. Workaround: when this kind of modification happens, refresh the data from the server

- Not very pretty
  - When clicking on a text area, it should expand over several lines so it can be read easily => idea remove text areas from the table and add a button to pop a modal to view/edit
  - when in view mode, spacing is off

- Some elements have a black background for no reason (seems to be a bug in vuetify)

## Versions

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
