run:
	@go run cmd/main.go ;

build_linux:
	@go build -o dist/linux/keizai cmd/main.go ;

build_osx:
	@go build -o dist/osx/keizai cmd/main.go ;

build_win:
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o dist/windows/Keizai.exe -ldflags "-s -w -H=windowsgui" cmd/main.go ;
