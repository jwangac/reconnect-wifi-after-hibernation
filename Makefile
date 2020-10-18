TARGET := reconnect-wifi-after-hibernation.exe

BUILD := GOPATH=$$(pwd) GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"

all:
	@go fmt
	@$(BUILD) -o $(TARGET)
	@upx -9 $(TARGET) 2>/dev/null; true
