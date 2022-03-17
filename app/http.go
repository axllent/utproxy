package app

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	// LogRequestsTo variable - if set then logfile is written
	LogRequestsTo string

	m sync.Mutex
)

// StartServer starts the http server
func StartServer() {
	listen := viper.GetString("Listen")
	if listen == "" {
		fmt.Println("Listen not set")
		os.Exit(1)
	}

	http.HandleFunc("/", makeGzipHandler(httpResponse))

	cert := viper.GetString("SSLCert")
	key := viper.GetString("SSLKey")

	if cert != "" && key != "" {
		log.Printf("Starting server on https://%s", listen)
		log.Fatal(http.ListenAndServeTLS(listen, cert, key, nil))
	} else {
		log.Printf("Starting server on http://%s", listen)
		log.Fatal(http.ListenAndServe(listen, nil))
	}
}

func httpResponse(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Path[1:]
	if req == "" || req == "favicon.ico" || !viper.Sub("Services").IsSet(req) {
		fourOfour(w)
		writeLog(r, 404)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	err := Check(req)
	if err != nil {
		w.WriteHeader(503)
		fmt.Fprint(w, err)
		writeLog(r, 503)
		return
	}

	fmt.Fprintf(w, "ok")
	writeLog(r, 200)
}

func writeLog(r *http.Request, status int) {

	if LogRequestsTo == "" {
		fmt.Println("no log")
		return
	}

	m.Lock()
	defer m.Unlock()

	t := time.Now()

	ts := t.Format("2/Jan/2006:15:04:05 -0700")

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	l := fmt.Sprintf("%s [%s] \"%s %s\" %d \"%s\"\n", ip, ts, r.Method, r.URL.Path, status, r.UserAgent())

	f, err := os.OpenFile(LogRequestsTo, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(l); err != nil {
		log.Println(err)
	}

}

func fourOfour(w http.ResponseWriter) {
	template := `<!DOCTYPE HTML>
<html lang="en">
<head>
<title>Not found</title>
<meta http-equiv="content-type" content="text/html; charset=UTF-8">
</head>
<body>
Not found
</body>
</html>	
`
	// w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, template)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip wrapper for compression
func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}
