package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/chanxuehong/log"

	"golang.org/x/sync/errgroup"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"kubernetes-ingress-controller/common"
	"kubernetes-ingress-controller/logic/server"
	"kubernetes-ingress-controller/logic/watcher"
)

var (
	Host          string
	Port, TLSPort int
)

func main() {
	_, ctx, _ := log.FromContextOrNew(context.Background(), nil)

	flag.StringVar(&Host, "host", "0.0.0.0", "the host to bind")
	flag.IntVar(&Port, "port", 80, "the insecure http port")
	flag.IntVar(&TLSPort, "tsl-port", 443, "the secure https port")
	flag.Parse()

	config := common.GetFactory().Config()
	fmt.Println(ctx, config)

	// 从集群内的token和ca.crt获取 Config
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.ErrorContext(ctx, "InClusterConfig failed", "error", err.Error())
		os.Exit(1)
	}

	// 从 restConfig 中创建一个新的 Clientset
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.ErrorContext(ctx, "NewForConfig failed", "error", err.Error())
		os.Exit(1)
	}
	s := server.New()
	w := watcher.NewWatcher(client, func(payload *watcher.Payload) {
		s.Update(ctx, payload)
	})

	var eg errgroup.Group
	eg.Go(func() error {
		return s.Run(context.TODO())
	})
	eg.Go(func() error {
		return w.Watcher(context.TODO())
	})
	if err := eg.Wait(); err != nil {
		log.ErrorContext(ctx, "Wait failed", "err", err.Error())
	}

}
