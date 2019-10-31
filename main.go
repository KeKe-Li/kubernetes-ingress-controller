package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/chanxuehong/log"
)

var (
	host          string
	port, tlsPort int
)

func main() {
	_, ctx, _ := log.FromContextOrNew(context.Background(), nil)

	flag.StringVar(&host, "host", "0.0.0.0", "the host to bind")
	flag.IntVar(&port, "port", 80, "the insecure http port")
	flag.IntVar(&tlsPort, "tsl-port", 443, "the secure https port")
	flag.Parse()

	fmt.Println(ctx)

}
