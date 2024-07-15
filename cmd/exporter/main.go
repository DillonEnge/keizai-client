package main

import (
	"flag"
	"fmt"

	"github.com/walle/targz"
)

func main() {
	v := flag.String("v", "", "version to be used to name the resulting tar.gz")
	a := flag.String("a", "", "arch to be used to name the resulting tar.gz")
	o := flag.String("o", "", "os to be used to name the resulting tar.gz")
	flag.Parse()
	if *v == "" {
		panic("failed to pass version")
	}
	if *a == "" {
		panic("failed to pass arch")
	}
	if *o == "" {
		panic("failed to pass os")
	}

	path := ""

	switch *o {
	case "darwin":
		switch *a {
		case "amd64":
			path = "dist/Keizai.app"
		case "arm64":
			path = "dist/darwin/arm64/Keizai.app"
		}
	case "windows":
		path = "dist/windows/keizai"
	}

	if path == "" {
		panic("unsupported GOOS + GOARCH")
	}

	err := targz.Compress(path, fmt.Sprintf("dist/assets/keizai-%s-%s-%s.tar.gz", *o, *a, *v))
	if err != nil {
		panic(err)
	}
}
