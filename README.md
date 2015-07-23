# httplz
wrap a stdin->stdout/stderr unix cmd to an http POST->200+body/500+body server

## Installation

```
# after setting up a go environment (https://golang.org/doc/install)
go get github.com/riwsky/httplz
```

## Usage

httplz is the dumbest thing that could possibly work as an http server. Please don't use it in any internet-facing server.

For example, an echo server can specified as simply as:
```
httplz cat
```

A 'file server' of sorts (POST a file name to get its contents) would be:
```
httplz xargs cat
```

A pattern emerges. A number doubling server:
```
httplz awk '{print ($0+0)*2}'
```

A csv -> json server:
```
httplz python -c 'import sys, csv, json; print json.dumps(list(csv.DictReader(sys.stdin)))'
```

And, inevitably, a 'metaserver' (POST a unix cmd to start a server wrapping that command):
```
httplz --port 8082 xargs httplz
```

## Inspired by
 - [gen_server, from Erlang/OTP](http://www.erlang.org/doc/man/gen_server.html)
 - [interact, from the Haskell Prelude](http://hackage.haskell.org/package/base-4.8.1.0/docs/Prelude.html#v:interact)
