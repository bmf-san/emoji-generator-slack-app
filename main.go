package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chromedp/chromedp"
)

type Draw struct {
	Text  string
	Color string
}

// TODO: 写経し直してコード理解
func handler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("tpl/index.tpl"))

	tpl.Execute(w, Draw{
		Text:  "Hello World",
		Color: "green",
	})
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	// サーバはブロックするので別の goroutine で実行する
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	// シグナルを待つ
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	go takeScreenShot(sigCh)
	<-sigCh

	// シグナルを受け取ったらShutdown
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Print(err)
	}
}

// スクショ撮影
func takeScreenShot(sigCh chan os.Signal) {
	defer func() {
		sigCh <- syscall.SIGTERM
	}()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// ローカルで起動してるサーバーの特定DOMだけ画像として撮影する
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`http://localhost:9999`, `#target`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("result.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

// 特定のDOMだけ撮影に使う
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}
