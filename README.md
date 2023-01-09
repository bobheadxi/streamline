# streamline

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
