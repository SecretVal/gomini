package gemini

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
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

type statuscode int

type StatusCodeRange struct {
	Min statuscode
	Max statuscode
}

var StatusCodes = map[string]StatusCodeRange{
	"InputExpected": StatusCodeRange{Min: 10, Max: 19},
	"Succes":        StatusCodeRange{Min: 20, Max: 29},
	"Redirection":   StatusCodeRange{Min: 30, Max: 39},
	"TempFail":      StatusCodeRange{Min: 40, Max: 49},
	"PermFail":      StatusCodeRange{Min: 50, Max: 59},
	"ClientCer":     StatusCodeRange{Min: 50, Max: 59},
}

func GetStatusCodeRange(code statuscode) StatusCodeRange {
	for key, value := range StatusCodes {
		if code > value.Min && code < value.Max {
			return StatusCodes[key]
		}
	}
	fmt.Println("ERROR: GOT BOGUS AMOGUS MESSAGE FROM SERVER")
	os.Exit(1)
	// Go why?
	return StatusCodeRange{}
}

type GeminiResponse struct {
	StatusCode statuscode
	Body       string
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

func MakeGeminiQuery(req GeminiRequest) GeminiResponse {
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

	code, _ := strconv.Atoi(strings.Split(buf, " ")[0])
	return GeminiResponse{
		StatusCode: statuscode(code),
		Body:       buf,
	}
}
