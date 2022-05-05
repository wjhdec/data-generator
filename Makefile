.PHONY: all generator
all: generator

BASE_NAME = data-generator
TARGET_PATH = dist/$(BASE_NAME)

clean:
	@go clean
	@rm -rf dist/

generator:
	@GOARCH=amd64 GOOS=$(GO_GOOS)  CGO_ENABLED=0 go build -tags "$(TAGS)" -o $(TARGET_PATH) main.go
	@upx $(TARGET_PATH)
	@cp configs/config.yaml dist/

windows-bin: GO_GOOS=windows
windows-bin: TARGET_PATH=dist/$(BASE_NAME).exe
windows-bin: generator

linux-bin: GO_GOOS=linux
linux-bin: generator

mac-bin: GO_GOOS=darwin
mac-bin: generator
