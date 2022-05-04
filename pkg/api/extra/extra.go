package extra

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

func Handler(ctx *fasthttp.RequestCtx) {
	keys := strings.Split(string(ctx.QueryArgs().Peek("key")), ",")
	var b []byte
	var err error
	if len(keys) == 1 {
		if b, err = os.ReadFile(filepath.Join("tmp", "extra", Key(keys[0]))); err != nil {
			ctx.Response.Header.SetStatusCode(http.StatusNoContent)
			return
		}
		_, _ = ctx.Write(b)
		return
	}
	m1 := make(map[string]string)
	for _, k := range keys {
		if b, err = os.ReadFile(filepath.Join("tmp", "extra", Key(k))); err == nil {
			m1[k] = string(b)
		}
	}
	if b, err = json.Marshal(m1); err != nil {
		ctx.Response.Header.SetStatusCode(http.StatusNoContent)
		return
	}
	_, _ = ctx.Write(b)
}
func Key(s string) string {
	return strings.TrimPrefix(strings.Replace(s, "/", "_", -1), "_")
}
