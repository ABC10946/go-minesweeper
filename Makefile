LDFLAGS := ""

minesweeper: *.go game/*.go cmd/*.go
	CGO_ENABLED=1 go build -trimpath -ldflags "-buildmode=pie" .