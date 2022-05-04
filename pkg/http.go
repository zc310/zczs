package pkg

import (
	"io/ioutil"
	"log"
	"net/http"
)

const UserAgent = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.27 Safari/537.36"

func init() {
	http.DefaultTransport = &UserAgentTransport{http.DefaultTransport}
}

type UserAgentTransport struct {
	rt http.RoundTripper
}

func (uat UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", UserAgent)

	return uat.rt.RoundTrip(r)
}

func GetByte(url string) ([]byte, error) {
	log.Println(url)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	return ioutil.ReadAll(r.Body)
}
