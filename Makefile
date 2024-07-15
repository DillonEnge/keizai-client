build_win:
	@go-winres simply --icon assets/icon.png --file-version git-tag
	@mv rsrc_windows_* cmd/
	@GOOS=windows GOARCH=amd64 go build -o dist/windows/Keizai.exe cmd/main.go
	@rm cmd/rsrc_windows_*

update_grpc:
	@go get github.com/DillonEnge/keizai-grpc@latest
