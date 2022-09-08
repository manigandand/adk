# api

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/manigandand/adk/api)

Custom `http.Handler` types which handles `errors.AppError` gracefully in oneplace.

```go
type Handler func(w http.ResponseWriter, r *http.Request) *errors.AppError
```
