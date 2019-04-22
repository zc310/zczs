package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/zc310/fasthttprouter"
	"github.com/zc310/zczs"
	"log"
	"regexp"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var addr string

var sbre *regexp.Regexp

func init() {
	sbre = regexp.MustCompile(`[?|&]id_spm=\w{24}`)
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString(ctx.Request.String())
}
func GetByte(ctx *fasthttp.RequestCtx) ([]byte, error) {
	host := string(ctx.Request.URI().String())
	if strings.Count(host, "cp.360.cn") > 0 {
		ctx.Request.URI().SetScheme("https")
	} else if strings.Count(host, "zc310.tech") == 0 {
		return []byte{}, nil
	}
	return zczs.GetByte(removeIdSpm(ctx.Request.URI().String()))
}

func NotFound(ctx *fasthttp.RequestCtx) {
	fmt.Println(ctx.Request.URI().String())
	b, err := GetByte(ctx)
	if err != nil {
		return
	}
	log.Println(ctx.Request.String(), "\n", string(b))

	ctx.Write(b)
}

func removeIdSpm(s string) string {
	return sbre.ReplaceAllString(s, "")
}
func main() {
	flag.StringVar(&addr, "addr", ":8080", ":8080")

	flag.Parse()

	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/std.log",
		MaxSize:    30,
		MaxBackups: 0,
		MaxAge:     30,
	})

	router := fasthttprouter.New()

	router.GET("/", Index)
	router.NotFound = NotFound

	log.Fatal(fasthttp.ListenAndServe(addr, router.Handler))
}
