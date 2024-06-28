build_win:
	@go-winres simply --icon assets/youricon.png
	@GOOS=windows GOARCH=amd64 go build -o dist/windows/keizai.exe cmd/main.go

update_grpc:
	@go get github.com/DillonEnge/keizai-grpc@latest
