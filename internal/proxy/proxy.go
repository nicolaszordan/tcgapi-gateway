package proxy

import (
	"io"
	"net/http"
	"net/url"
)

// ForwardRequest proxies an incoming request to a target service and returns the response.
func ForwardRequest(targetBase string, r *http.Request) (*http.Response, error) {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		return nil, err
	}

	targetURL.Path = r.URL.Path
	targetURL.RawQuery = r.URL.RawQuery

	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header.Clone()

	client := &http.Client{}
	return client.Do(req)
}

// CopyResponse writes the proxied response to the gateway’s response writer.
func CopyResponse(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
