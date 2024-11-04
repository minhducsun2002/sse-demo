# sse-demo
A small demo for Server Sent Events, written with Go 1.23.

Run it using
```shell
$ go run main.go 
```
and then access [localhost:10499](http://localhost:10499).

The Go server reads & caches [`index.html`](./index.html) in the current directory (and will crash if it doesn't find one), so restart it if you tweak the HTML file. 

Licensed under [CC0 1.0](./LICENSE). Feel free to hack this snippet!
