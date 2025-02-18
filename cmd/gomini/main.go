package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
)

var TLSConfig = tls.Config{
	InsecureSkipVerify: true,
}

var PORT int = 1965
var PREFIX = "gemini://"

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

func host_from_string(url string) string {
	return strings.Split(strings.TrimPrefix(url, PREFIX), "/")[0]
}

func make_gemini_query(url string) string {
	host := host_from_string(url)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, PORT), &TLSConfig)
	if err != nil {
		panic(err)
	}

	url_buf := []byte(fmt.Sprintf("%s\r\n", url))

	_, err = conn.Write(url_buf)
	if err != nil {
		panic(err)
	}
	buf := read_response(conn)

	return buf
}

func main() {
	prog := os.Args[0]
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Usage: <%s> <url>\n", prog)
		fmt.Printf("ERROR: Did not specify url\n")
		os.Exit(1)
	}
	url := args[0]

	fmt.Println(make_gemini_query(url))
}
