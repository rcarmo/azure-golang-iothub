export GOOS?=linux
export GOARCH?=amd64
export BINARY?=client

run: test
	go run main.go iothubHttp.go

test:
	go test -v

build: test
	go build -o $(BINARY)

pack: test
	go build -ldflags="-s -w" -o $(BINARY)
	# this will only work if you've installed UPX (on Ubuntu, apt install upx-ucl)
	-upx $(BINARY)
