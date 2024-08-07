# Developing

## Prerequisites

This repo uses [Task](https://taskfile.dev/) to manage running tasks for development, similar to `make`.

## Tasks

The best place to check for available tasks is going to be in the [Taskfile.yml directly](../Taskfile.yml).
I will make an effort to keep this up to date, but it's always a good idea to check the Taskfile as well.

| Invocation       | Description                                                                                                                       |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `task build`     | Builds the stax binary.                                                                                                           |
| `task local-api` | Builds the binary and then starts a local API server for testing.                                                                 |
| `task test`      | Runs the tests.                                                                                                                   |
| `task lint`      | Lints the code.                                                                                                                   |
| `task ent`       | Regenerates the [bones](../internal/bones) package. Usually you don't invoke this directly, `task build` will run it when needed. |
