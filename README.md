[![Get-zap](https://github.com/SiliconLabs/get-zap/actions/workflows/go.yml/badge.svg)](https://github.com/SiliconLabs/get-zap/actions/workflows/go.yml)

# What is this?

The immediate purpose of this program is to retrieve releases of [zap](https://github.com/project-chip/zap) from the release page on Github.

Under the hood, it is a generic retriever of release assets from github from any project.
It also supports Artifactory caching.

So this program can be useful for you in any environment (CI, or personal desktop), where you have a following need:
  - retrieve assets from github project releases
  - retrieve files from Artifactory
  - upload files to Artifactory
  - automate the process of: download a release artifact, if it exists in Artifactory, get it from there, otherwise get it directly from Github releases, along the way caching it to Artifactory for next time

# Who should use this?

There are few use cases:
  - Anyone who wants to get the "correct" version of zap, should be able to use this. Whether it's Matter SDK, or something else.
  - Anyone who wants to download release artifacts from github can use it, but there are also other tools just for that.
  - Anyone who wants to upload and download files from/to Artifactory can use this.
  - Anyone who wants to retrieve a file which may be cached on Artifactory, but is originally an asset from Github.

# Architecture

This is a [Go](https://go.dev/) program. It uses the [JFrog Go client library](https://github.com/jfrog/jfrog-client-go) and the [Github Go client librariy](https://github.com/google/go-github) to perform access to Artifactory and Github. It does not use any other means to talk to Artifactory or Github, so all documentation regarding limitations for those libraries apply.

# Build Instructions

You need go toolchain installed to build it from source code. Many platforms (Linuxes, brew) come with Go toolchains easily installable through your package manager of choice, or you can follow [instructions here](https://go.dev/doc/install).

Once you have go installed:
  1. You can simply run `go run .` to execute the program from source.
  2. Or build it using `go build` and run `get-zap` executable that gets created.
  3. You can run `go install` to build and deploy the executable into your Go bin directory.
  4. If you want to build for a different platform than local, then set the GOOS and GOARCH environment variables as described [here](https://go.dev/doc/install/source#environment) before you run `go build`.

When executing `get-zap` without any arguments, it will by default download the latest stable release of Zap for the local platform.

For all other options, type `get-zap --help` or `go run . --help` if you run from source.

# Setup

You can configure all the arguments to this program either via command line arguments, through environment variables, or via a configuration file `~/.get-zap.json`

Github related environment variables:
  - GET_ZAP_GHTOKEN: Github API access token
  - GET_ZAP_GHOWNER: Owner of the github repo.
  - GET_ZAP_GHREPO: Github repo name.

Artifactory related environment variables in question are:
  - GET_ZAP_RTAPIKEY: Artifactory API key.
  - GET_ZAP_RTPATH: Artifactory path within the repo.
  - GET_ZAP_RTREPO: Artifactory repo..
  - GET_ZAP_RTURL: Artifactory url.
  - GET_ZAP_RTUSER: Artifactory user.

# Examples


1. Default execution without any arguments:
```
[~/git/get-zap (main)]$ ./get-zap
```

2. Download from a different repo than zap (example, latest release of PTI library):
```
[~/git/get-zap (main)]$ ./get-zap --ghOwner SiliconLabs --ghRepo java_packet_trace_library
```

3. List latest zap release:
```
[~/git/get-zap (main)]$ ./get-zap gh list
```

4. List all zap releases:
```
[~/git/get-zap (main)]$ ./get-zap gh list --ghRelease all
```

5. List all releases of a different repo:
```
[~/git/get-zap (main)]$ ./get-zap gh list --ghRelease all --ghOwner SiliconLabs --ghRepo java_packet_trace_library
```

6. Download a specific zap release:
```
[~/git/get-zap (main)]$ ./get-zap --ghRelease v2024.01.05-nightly
```

7. Print help:
```
[~/git/get-zap (main)]$ ./get-zap --help
```