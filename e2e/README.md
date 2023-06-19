# E2E Tests

This directory contains end-to-end tests for the project. These tests are
intended to be run in a Github Actions workflow if a Pull Request on Github is
created. You can run them locally, but you will need to set up a few things
first. And it's a bit more flaky.

The tests will create a Kubernetes cluster using [k3d](https://k3d.io/). k3d
uses Docker and [k3s](https://k3s.io/) to create a lightweight Kubernetes
cluster for testing.

The current amount of tests is very limited. For Telemetry, it merely tests if
the Telemetry client has successfully sent data to the Telemetry Server, which
stores the data in PostgreSQL. The test verifies that the data is there and can
be retrieved using Telemetry Server. For Telemetry Stats, the test verifies that
the data has been successfully stored in InfluxDB. The test retrieves the data
directly from InfluxDB.

## Configuration

### Telemetry branch to be used for testing

The tests are configured via environment variables. They can be configured on
Github specifically for
[Actions](https://docs.github.com/en/actions/learn-github-actions/variables).
The following variables can be configured:

- `TELEMETRY_BRANCH`: The branch to use for the test of Telemetry (not Telemetry
  Stats). This can be any branch, tag, or commit hash. It defaults to `release`.
  `release` is a special value that will result in the latest released container
  image being downloaded and used for testing. Any other value will be taken as
  a Git reference (branch, tag or reference) and the source code will be cloned,
  built and a container image will be created and uploaded to the Kubernetes
  cluster, so that the images will be available for use.

# Running Tests Locally

## Requirements

- [act](https://github.com/nektos/act)
- [Docker](https://www.docker.com/)

## Configuration

A Personal Access token (PAT) from Github is required to run the tests locally.

## Usage

```console
act --env TELEMETRY_BRANCH=master --secret-file .env
```

Where `master` is the branch you want to test. You can also set it to "release"
to test an already released version, for which only the container image will be
downloaded but not built. In all other cases, the source code will be cloned,
built and a corresponding container image is being created and uploaded to the
Kubernetes cluster.

The `--secret-file` option is required to pass the PAT to act. The file should
contain the following:

```console
GITHUB_TOKEN=your-pat-here
```
