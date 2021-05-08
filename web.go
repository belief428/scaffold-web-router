package WebRouter

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Web struct {
	*WebConfig
	httpServer *http.Server
}

type WebConfig struct {
	Port, ReadTimeout, WriteTimeout, IdleTimeout int
}

type WebServer func(config *WebConfig) *Web

func (this *Web) stop() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				err = fmt.Errorf("internal error: %v", err)
			}
		}()
		s := <-c
		defer close(c)
		signal.Stop(c)
		fmt.Printf("Http Server Stop：Closed - %d\n", s)
		os.Exit(0)
	}()
}

func (this *Web) Run(handler http.Handler) {
	this.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", this.Port),
		Handler:        handler,
		ReadTimeout:    time.Duration(this.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(this.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(this.IdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	this.stop()

	fmt.Println("Http Server Start")
	fmt.Printf("Http Server Address - %v\n", fmt.Sprintf("http://127.0.0.1:%d", this.Port))

	_ = this.httpServer.ListenAndServe()
}

func NewWeb() WebServer {
	return func(config *WebConfig) *Web {
		return &Web{WebConfig: config}
	}
}
