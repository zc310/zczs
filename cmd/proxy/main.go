package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/zc310/alice"
	"github.com/zc310/fasthttprouter"
	"github.com/zc310/fs/middleware"
	fslog "github.com/zc310/log"
	"github.com/zc310/zczs/api"
	"gopkg.in/natefinch/lumberjack.v2"
)

var addr string
var tty bool

func main() {
	flag.StringVar(&addr, "addr", ":8080", ":8080")

	flag.Parse()

	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/std.log",
		MaxSize:    30,
		MaxBackups: 0,
		MaxAge:     30,
	})
	var err error
	router := fasthttprouter.New()

	fslog.SetPath("logs/")
	mwlog := &middleware.Logger{Filename: "logs/access.log",
		MaxSize:    30,
		MaxBackups: 7,
		MaxAge:     7,
		Compress:   false}

	cfg := &middleware.Config{fslog.NewWithPrefix("proxy"), router, "/"}
	if err = mwlog.Init(cfg); err != nil {
		log.Fatal(err)
	}

	sf := middleware.Singleflight{}
	if err = sf.Init(cfg); err != nil {
		log.Fatal(err)
	}
	//	sf.Key = "{document_root}"
	sf.Hash = "sha1"

	mwcache := &middleware.Cache{}
	mwcache.Store = map[string]interface{}{"name": "file", "path": "cache.api"}
	mwcache.Key = ""
	mwcache.Hash = ""
	mwcache.Timeout = "300s"

	if err = mwcache.Init(cfg); err != nil {
		log.Fatal(err)
	}
	mwcache.HashFun = func(b []byte) string {
		return api.RemoveIdSpm(string(b))
	}
	gz := &middleware.Compress{}
	re := &middleware.Recover{}
	if err = re.Init(cfg); err != nil {
		log.Fatal(err)
	}

	mw := []alice.Constructor{re.Process, mwlog.Process,
		gz.Process, mwcache.Process,
		sf.Process}

	router.GET("/", api.Index)
	router.GET("/int/qkjinfo", alice.New(mw...).Then(api.ZczsQkjinfo))
	router.GET("/zczs/issue", alice.New(mw...).Then(api.ZczsIssue))
	router.GET("/zczs/zcmatch", alice.New(mw...).Then(api.ZczsZcmatch))
	router.GET("/sfc/extra", alice.New(mw...).Then(api.SfcExtra))
	router.GET("/sfc/his/360dd", alice.New(mw...).Then(api.His360dd))
	router.GET("/sfc/his/360dd/:id", alice.New(mw...).Then(api.His360dd))
	router.GET("/int/getoupei/", alice.New(mw...).Then(api.GetOupei))
	router.GET("/int/hiszhanji", alice.New(mw...).Then(api.NotOk))
	router.GET("/jczqdata/geteurochange/match/:mid/gcid/:gid", alice.New(mw...).Then(api.NotOk))

	router.NotFound = alice.New(mw...).Then(api.NotFound)

	fmt.Println("足彩助手代理服务")

	log.Fatal(fasthttp.ListenAndServe(addr, router.Handler))

}
