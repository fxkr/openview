# Building

## Docker (recommended)

Compiling from source only requires [Docker](https://www.docker.com/) and internet access.
This will fetch dependencies, build the software and create packages in one step:

```Shell
make
```

## No docker (not recommended)

Install the prerequisites documented in the [Dockerfile](.circleci/images/openview-ci/Dockerfile) and run:

```Shell
make all-direct
```