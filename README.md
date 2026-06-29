# token-swap-workbench

## Run

```bash
mise run test          # run tests with race detector
mise run lint          # run golangci-lint
mise run run           # run the service locally
mise run build         # build all packages
```



## Project layout

```
cmd/token-swap-workbench/    entry point
internal/                  application logic
pkg/                       shared / exported helpers
```

## License

MIT — see [LICENSE](LICENSE)
