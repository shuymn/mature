export GO111MODULE=on

cmd/mature/mature: *.go cmd/mature/*.go
	cd cmd/mature && go build -ldflags="-s -w"

.PHONY: dist
dist:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0\
		go build -ldflags="-s -w" -o ./dist/mature ./cmd/mature

.PHONY: install
install: cmd/mature/mature
	install cmd/mature/mature ${GOPATH}/bin

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	rm -f cmd/mature/mature
	rm -f dist/mature
