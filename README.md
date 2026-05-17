## Run

из корня:

```bash
go run ./cmd -config impulse/config.json -events impulse/events
```

## Test

```bash
go test ./...
```


## Build

```bash
go build -o path-of-exile-cli ./cmd
```

запуск бинаря:

```bash
./path-of-exile-cli -config impulse/config.json -events impulse/events
```