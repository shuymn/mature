BIN := notify

ifdef update
	u=-u
endif

export GO111MODULE=on

cmd/$(BIN)/$(BIN): *.go cmd/$(BIN)/*.go clean
	cd cmd/$(BIN) && go build -ldflags="-s -w"

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: dist
dist: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-s -w" -o ./dist/$(BIN) ./cmd/$(BIN)

.PHONY: install
install: cmd/$(BIN)/$(BIN)
	install cmd/$(BIN)/$(BIN) ${GOPATH}/bin

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	rm -f cmd/$(BIN)/$(BIN)
	rm -f dist/$(BIN)
