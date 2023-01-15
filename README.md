# streamline [![go reference](https://pkg.go.dev/badge/go.bobheadxi.dev/streamline.svg)](https://pkg.go.dev/go.bobheadxi.dev/streamline) [![Sourcegraph](https://img.shields.io/badge/view%20on-sourcegraph-A112FE?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAEZklEQVRoQ+2aXWgUZxSG3292sxtNN43BhBakFPyhxSujRSxiU1pr7SaGXqgUxOIEW0IFkeYighYUxAuLUlq0lrq2iCDpjWtmFVtoG6QVNOCFVShVLyxIk0DVjZLMxt3xTGTccd2ZOd/8JBHci0CY9zvnPPN+/7sCIXwKavOwAcy2QgngQiIztDSE0OwQlDPYR1ebiaH6J5kZChyfW12gRG4QVgGTBfMchMbFP9Sn5nlZL2D0JjLD6710lc+z0NfqSGTXQRQ4bX07Mq423yoBL3OSyHSvUxirMuaEvgbJWrdcvkHMoJwxYuq4INUhyuWvQa1jvdMGxAvCxJlyEC9XOBCWL04wwRzpbDoDQ7wfZJzIQLi5Eggk6DiRhZgWIAbE3NrM4A3LPT8Q7UgqAqLqTmLSHLGPkyzG/qXEczhd0q6RH+zaSBfaUoc4iQx19pIClIscrTkNZzG6gd7qMY6eC2Hqyo705ZfTf+eqJmhMzcSbYtQpOXc92ZsZjLVAL4YNUQbJ5Ttg4CQrQdGYj44Xr9m1XJCzmZusFDJOWNpHjmh5x624a2ZFtOKDVL+uNo2TuXE3bZQQZUf8gtgqP31uI94Z/rMqix+IGiRfWw3xN9dCgVx+L3WrHm4Dju6PXz/EkjuXJ6R+IGgyOE1TbZqTq9y1eo0EZo7oMo1ktPu3xjHvuiLT5AFNszUyDULtWpzE2/fEsey8O5TbWuGWwxrs5rS7nFNMWJrNh2No74s9Ec4vRNmRRzPXMP19fBMSVsGcOJ98G8N3Wl2gXcbTjbX7vUBxLaeASDQCm5Cu/0E2tvtb0Ea+BowtskFD0wvlc6Rf2M+Jx7dTu7ubFr2dnKDRaMQe2v/tcIrNB7FH0O50AcrBaApmRDVwFO31ql3pD8QW4dP0feNwl/Q+kFEtRyIGyaWXnpy1OO0qNJWHo1y6iCmAGkBb/Ru+HenDWIF2mo4r8G+tRRzoniSn2uqFLxANhe9LKHVyTbz6egk9+x5w5fK6ulSNNMhZ/Feno+GebLZV6isTTa6k5qNl5RnZ5u56Ib6SBvFzaWBBVFZzvnERWlt/Cg4l27XChLCqFyLekjhy6xJyoytgjPf7opIB8QPx7sYFiMXHPGt76m741MhCKMZfng0nBOIjmoJPsLqWHwgFpe6V6qtfcopxveR2Oy+J0ntIN/zCWkf8QNAJ7y6d8Bq4lxLc2/qJl5K7t432XwcqX5CrI34gzATWuYILQtdQPyePDK3iuOekCR3Efjhig1B1Uq5UoXEEoZX7d1q535J5S9VOeFyYyEBku5XTMXXKQTToX5Rg7OI44nbW5oKYeYK4EniMeF0YFNSmb+grhc84LyRCEP1/OurOcipCQbKxDeK2V5FcVyIDMQvsgz5gwFhcWWwKyRlvQ3gv29RwWoDYAbIofNyBxI9eDlQ+n3YgsgCWnr4MStGXQXmv9pF2La/k3OccV54JEBM4yp9EsXa/3LfO0dGPcYq0Y7DfZB8nJzZw2rppHgKgVHs8L5wvRwAAAABJRU5ErkJggg==)](https://sourcegraph.com/github.com/bobheadxi/streamline)

[![pipeline](https://github.com/bobheadxi/streamline/actions/workflows/pipeline.yaml/badge.svg)](https://github.com/bobheadxi/streamline/actions/workflows/pipeline.yaml)
[![codecov](https://codecov.io/gh/bobheadxi/streamline/branch/main/graph/badge.svg?token=f1VZULSJsT)](https://codecov.io/gh/bobheadxi/streamline)
[![Go Report Card](https://goreportcard.com/badge/go.bobheadxi.dev/streamline)](https://goreportcard.com/report/go.bobheadxi.dev/streamline)

Handle your data, line by line.

```sh
go get go.bobheadxi.dev/streamline
```

> **This package is currently experimental**, and the API may change.

## Overview

When working with data streams in Go, you typically get an `io.Reader`, which is great for arbitrary data - but in many cases, especially when scripting, it's common to end up with data and outputs that are structured line-by-line.
Typically, you might set up an `bufio.Reader` or `bufio.Scanner` to read data line by line, but for cases like `exec.Cmd` you will also need boilerplate to configure the command and set up pipes.

`streamline` offers a variety of primitives that aim to make working with data line by line a breeze:

- `streamline.Stream` offers the ability to add hooks that handle an `io.Reader` line-by-line with `(*Stream).Stream(LineHandler[string])` and `(*Stream).StreamBytes(LineHandler[[]byte])`.
- `pipeline.Pipeline` offers a way to build pipelines that transform the data in a `streamline.Stream`, such as cleaning and mapping data.
  - `jq.Pipeline` can be used to map every line to the output of a JQ query, for example.
- `pipe.NewStream` offers a way to create a buffered pipe between a writer and a `Stream`.
  - Package `streamexec` uses this to attach a `Stream` to an `exec.Cmd`.

## Background

Some of the ideas in this package started in [`sourcegraph/run`](https://github.com/sourcegraph/run), where we were trying to build utilities that [made it easier to write bash-esque scripts using Go](https://github.com/sourcegraph/sourcegraph/blob/main/doc/dev/adr/1652433602-use-go-for-scripting.md), namely being able to do things you would often to in scripts such as grepping and iterating over lines.
