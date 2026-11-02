package main

import "github.com/shivamkkas/rdpms25-go-rpc-service/pkg"

var (
	version   string = "unknown"
	buildTime string = "unknown"
)

func main() {
	pkg.Start(version, buildTime)
}
