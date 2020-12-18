export GO111MODULE=on

cmd/notify/notify: *.go cmd/notify/*.go clean
	go generate && cd cmd/notify && go build -ldflags="-s -w"

.PHONY: dist
dist: clean
	go generate && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-s -w" -o ./dist/notify ./cmd/notify

.PHONY: install
install: cmd/notify/notify
	install cmd/notify/notify ${GOPATH}/bin

.PHONY: test
test:
	go test -race ./...

.PHONY: clean
clean:
	rm -f cmd/notify/notify
	rm -f dist/notify
