package http_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestServerStart(t *testing.T) {

	testServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host != "SomeHost" {
				http.Error(w, fmt.Sprintf("Expected Host header to be '%q', but got '%q'", "SomeHost", r.Host), http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}))

	defer testServer.Close()

	// handler := SetupHandler("www")
	// qconf := quic.Config{}
	// t.Log(handler)

	// server := http3.Server{
	// 	Handler:    handler,
	// 	Addr:       ":6121",
	// 	QuicConfig: &qconf,
	// }

	// var wg sync.WaitGroup
	// wg.Add(1)

	// log.Println(server.ListenAndServe())
	// wg.Wait()

	// client := &http.Client{
	// 	Transport: &http3.RoundTripper{},
	// }

	// //필요시 헤더 추가 가능
	// //req.Header.Add("User-Agent", "Crawler")
	// req, err := http.NewRequest("GET", "http://localhost:6121/", nil)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// // 결과 출력
	// bytes, _ := io.ReadAll(resp.Body)
	// str := string(bytes) //바이트를 문자열로
	// log.Println("str >> ", str)

}

func SetupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("root src >> %#v\n", r)

		w.Write([]byte("<html>hello</html>"))
	})

	return mux
}
