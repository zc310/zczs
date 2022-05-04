package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
	"github.com/yuin/goldmark"
	"github.com/zc310/alice"
	"github.com/zc310/fasthttprouter"
	"github.com/zc310/fs/middleware"
	"github.com/zc310/fs/middleware/gzip"
	"github.com/zc310/fs/middleware/logger"
	"github.com/zc310/fs/middleware/singleflight"
	fslog "github.com/zc310/log"
	"github.com/zc310/zczs/pkg/api"
	"gopkg.in/natefinch/lumberjack.v2"
)

var addr string

//go:embed assets
var assets embed.FS
var version = "1000"

func main() {
	flag.StringVar(&addr, "addr", ":8080", ":8080")

	flag.Parse()

	log.SetOutput(io.MultiWriter(os.Stderr, &lumberjack.Logger{
		Filename:   filepath.Join("tmp", "logs", "std.log"),
		MaxSize:    30,
		MaxBackups: 0,
		MaxAge:     30,
	}))

	var err error

	router := fasthttprouter.New()
	_ = os.MkdirAll(filepath.Join("tmp", "logs"), os.ModePerm)
	_ = os.Mkdir(filepath.Join("tmp", "cache"), os.ModePerm)

	fslog.SetPath(filepath.Join("tmp", "logs"))

	err = SaveDoc()
	if err != nil {
		log.Println(err)
		return
	}

	mwlog := &logger.Logger{Filename: filepath.Join("tmp", "logs", "access.log"),
		MaxSize:    30,
		MaxBackups: 7,
		MaxAge:     7,
		Compress:   false}

	cfg := &middleware.Config{Logger: fslog.NewWithPrefix("proxy"), Router: router, Path: "/"}
	if err = mwlog.Init(cfg); err != nil {
		log.Fatal(err)
	}

	sf := singleflight.Singleflight{}
	if err = sf.Init(cfg); err != nil {
		log.Fatal(err)
	}

	sf.Hash = "sha1"

	gz := &gzip.Compress{}
	re := &middleware.Recover{}
	if err = re.Init(cfg); err != nil {
		log.Fatal(err)
	}

	mw := []alice.Constructor{re.Process, mwlog.Process,
		gz.Process,
		sf.Process}

	router.GET("/", api.Index)
	router.GET("/int/qkjinfo", alice.New(mw...).Then(api.ZczsQkjinfo))
	router.GET("/zczs/issue", alice.New(mw...).Then(api.ZczsIssue))
	router.GET("/zczs/zcmatch", alice.New(mw...).Then(api.ZczsZcmatch))
	router.GET("/zczs/match", alice.New(mw...).Then(api.NoContent))
	router.GET("/sfc/extra", alice.New(mw...).Then(api.NoContent))
	router.GET("/sfc/his/360dd", alice.New(mw...).Then(api.NotOk))
	router.GET("/sfc/his/360dd/:id", alice.New(mw...).Then(api.NotOk))
	router.GET("/int/getoupei/", alice.New(mw...).Then(api.GetOupei))
	router.GET("/int/hiszhanji", alice.New(mw...).Then(api.NotOk))
	router.GET("/jczqdata/geteurochange/match/:mid/gcid/:gid", alice.New(mw...).Then(api.NotOk))
	router.POST("/int/querybalance", api.NoContent)
	router.POST("/api.php", api.NoContent)

	router.NotFound = alice.New(mw...).Then(api.NotFound)

	fmt.Println("足彩助手代理服务")

	log.Fatal(fasthttp.ListenAndServe(addr, router.Handler))

}

func SaveDoc() error {
	var b []byte
	var err error
	fnv := filepath.Join("tmp", ".version")
	b, err = os.ReadFile(fnv)
	if err == nil {
		if string(b) == version {
			return nil
		}
	}
	if b, err = Doc(); err == nil {
		err = os.WriteFile("看我.html", b, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_ = os.MkdirAll(filepath.Join("tmp", "w360"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join("sfc", "历史"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join("jqc", "历史"), os.ModePerm)
	var fi []fs.DirEntry

	b, err = assets.ReadFile("assets/tmp/1.png")
	if err != nil {
		return err
	}
	err = os.WriteFile("tmp/1.png", b, os.ModePerm)
	if err != nil {
		return err
	}

	if fi, err = assets.ReadDir("assets/sfc"); err != nil {
		return err
	}
	for _, f := range fi {
		if !f.IsDir() {
			b, err = assets.ReadFile(fmt.Sprintf("assets/sfc/%s", f.Name()))
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("sfc/%s", f.Name()), b, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	if fi, err = assets.ReadDir("assets/jqc"); err != nil {
		return err
	}
	for _, f := range fi {
		if !f.IsDir() {
			b, err = assets.ReadFile(fmt.Sprintf("assets/jqc/%s", f.Name()))
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("jqc/%s", f.Name()), b, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	if fi, err = assets.ReadDir(fmt.Sprintf("assets/tmp/w360")); err != nil {
		return err
	}
	for _, f := range fi {
		if !f.IsDir() {
			if "LOCK" == f.Name() {
				continue
			}
			b, err = assets.ReadFile(fmt.Sprintf("assets/tmp/w360/%s", f.Name()))
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("tmp/w360/%s", f.Name()), b, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return os.WriteFile(fnv, []byte(version), os.ModePerm)
}
func Doc() ([]byte, error) {
	var buf bytes.Buffer
	b, err := assets.ReadFile(fmt.Sprintf("assets/看我.md"))
	if err != nil {
		return nil, err
	}
	if err := goldmark.Convert(b, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
