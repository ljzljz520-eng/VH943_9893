# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	campgear/cmd/server	[no test files]
?   	campgear/internal/staff	[no test files]
ok  	campgear/internal/catalog	0.007s
ok  	campgear/internal/domain	0.003s
ok  	campgear/internal/httpapi	0.007s
ok  	campgear/internal/maintenance	0.007s
--- FAIL: TestRentalRejectsMaintenance (0.00s)
    rental_test.go:66: maintenance equipment must be rejected
FAIL
FAIL	campgear/internal/rental	0.007s
ok  	campgear/internal/reporting	0.004s
ok  	campgear/internal/storage	0.022s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
