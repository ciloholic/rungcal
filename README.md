# build

```
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/rungcal cmd/rungcal/main.go
$ GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/rungcal cmd/rungcal/main.go
```

# execution

```
$ go run ./cmd/rungcal/main.go insert --project=*** --target-date="yyyy-mm-dd" --dry-run --verbose
```

# golangci-lint

```
$ docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v
```
