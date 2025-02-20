package gemini

import (
	"crypto/tls"
	"fmt"
	"strings"
)

var TLSConfig = tls.Config{
	InsecureSkipVerify: true,
}

var PORT int = 1965
var PREFIX = "gemini://"

type GeminiRequest struct {
	Url  string
	Host string
	Port int
}

// gemini://ghost/...:PORT
func ParseGeminiRequest(url string, port int) (GeminiRequest, bool) {
	trimmed_url := strings.TrimPrefix(url, PREFIX)
	host := host_from_string(trimmed_url)
	if trimmed_url == url {
		return GeminiRequest{}, false
	}

	return GeminiRequest{
		Url:  url,
		Host: host,
		Port: port,
	}, true
}

func host_from_string(url string) string {
	return strings.Split(strings.TrimPrefix(url, PREFIX), "/")[0]
}

func read_response(conn *tls.Conn) string {
	var sb strings.Builder
	buf := make([]byte, 1024)

	i, _ := conn.Read(buf)
	for i > 0 {
		sb.Write(buf)
		buf = make([]byte, 1024)
		i, _ = conn.Read(buf)
	}
	return sb.String()
}

func MakeGeminiQuery(req GeminiRequest) string {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port), &TLSConfig)
	if err != nil {
		panic(err)
	}

	url_buf := []byte(fmt.Sprintf("%s\r\n", req.Url))

	_, err = conn.Write(url_buf)
	if err != nil {
		panic(err)
	}
	buf := read_response(conn)

	return buf
}
