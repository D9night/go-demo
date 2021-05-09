package wk3

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	service := http.NewServeMux()
	service.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	serverOut := make(chan struct{})

	// 停止单个goruntime
	service.HandleFunc("/out", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})

	server := http.Server{
		Handler: service,
		Addr:    ":8080",
	}

	// 退出后 context 不再阻塞
	group.Go(func() error {
		return server.ListenAndServe()
	})

	// 退出后调用server.shutdown context 不再阻塞
	group.Go(func() error {

		select {
		case <-ctx.Done():
			log.Println("group fatal...")
		case <-serverOut:
			log.Println("server fatal...")
		}
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		return server.Shutdown(timeoutCtx)
	})

	// 接收退出信号 context 不再阻塞
	group.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return fmt.Errorf("get os signal: %v", sig)
		}
	})

	// 阻塞所有的通过Go加入的goroutine，然后等待他们一个个执行完成
	// 然后返回第一个出错的goroutine的错误信息
	group.Wait()

}


