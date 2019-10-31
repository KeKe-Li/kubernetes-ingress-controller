package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"kubernetes-ingress-controller/conf"
	"kubernetes-ingress-controller/common"
	"kubernetes-ingress-controller/logic/watcher"

	"github.com/chanxuehong/log"

	"golang.org/x/sync/errgroup"
)

// A Server serves HTTP pages.
type Server struct {
	cfg          *conf.Config
	routingTable atomic.Value

	ready IEvent
}

// New 创建一个新的服务器
func New() *Server {
	s := &Server{
		ready: NewEvent(),
	}
	s.routingTable.Store(NewRoutingTable(nil))
	return s
}

// Run 启动服务器.
func (s *Server) Run(ctx context.Context) error {
	// 直到第一个 payload 数据后才开始监听
	s.ready.Wait(ctx)

	// 启动 80 和 443 两个端口
	var eg errgroup.Group
	eg.Go(func() error {
		// 当前的 Server 实现了 Handler 接口（ServeHTTP函数)
		srv := http.Server{
			Addr:    fmt.Sprintf("%s:%d", common.GetFactory().Config().Host, common.GetFactory().Config().TLSPort),
			Handler: s,
		}
		srv.TLSConfig = &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return s.routingTable.Load().(*RoutingTable).GetCertificate(hello.ServerName)
			},
		}
		log.InfoContext(ctx, "run", "addr", srv.Addr)

		err := srv.ListenAndServeTLS("", "")
		if err != nil {
			return fmt.Errorf("error serving tls: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		srv := http.Server{
			Addr:    fmt.Sprintf("%s:%d", common.GetFactory().Config().Host, common.GetFactory().Config().TLSPort),
			Handler: s,
		}

		log.InfoContext(ctx, "run", "srv", srv)

		err := srv.ListenAndServe()
		if err != nil {
			return fmt.Errorf("error serving non-tls: %w", err)
		}
		return nil
	})
	return eg.Wait()
}

// ServeHTTP serves an HTTP request.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 获取后端的真实服务地址
	backendURL, err := s.routingTable.Load().(*RoutingTable).GetBackend(r.Host, r.URL.Path)
	if err != nil {
		log.ErrorContext(ctx,"GetBackend failed","error",err.Error())
		return
	}
	// 使用 NewSingleHostReverseProxy 进行代理请求
	p := httputil.NewSingleHostReverseProxy(backendURL)
	p.ServeHTTP(w, r)
}

// Update 更新路由表根据新的 Ingress 规则
func (s *Server) Update(ctx context.Context, payload *watcher.Payload) {
	s.routingTable.Store(NewRoutingTable(payload))
	s.ready.Set(ctx)
}
