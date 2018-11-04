package main

import (
	"fmt"
	"os"

	"github.com/demomon/go-grpc-server/pkg/cmd/server"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
