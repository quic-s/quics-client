package http3

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-s/quics-client/pkg/connection"
	"github.com/quic-s/quics-client/pkg/sync"
	"github.com/quic-s/quics-client/pkg/types"
)

func SetupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", r)
		w.Write([]byte("hello, Quics Client here"))
	})

	// TODO : add handler for each api
	mux.HandleFunc("/api/v1/connect/server", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterClientHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				log.Println("quics-client : cannot unmarshal")
				log.Println(err)
			}
			connection.RegisterClient(body.ClientPW, body.Host, body.Port)

			w.Write([]byte("quics-client : [/api/v1/connect/server] Resp : OK"))
		}
	})

	mux.HandleFunc("/api/v1/connect/root/local", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterRootDirHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				log.Println("quics-client : cannot unmarshal")
				log.Println(err)
			}

			err = connection.RegisterLocalRootDirRequest(body.RootDir, body.RootDirPw)
			if err != nil {
				log.Println("quics-client : ", err)
			}
			w.Write([]byte("quics-client : [/api/v1/connect/root/local] Resp : OK"))

		}
	})

	mux.HandleFunc("/api/v1/connect/root/remote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterRootDirHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				log.Println("quics-client : cannot unmarshal")
				log.Println(err)
			}

			connection.RegisterRemoteRootDirRequest(body.RootDir, body.RootDirPw)
		}
	})

	mux.HandleFunc("/api/v1/disconnect/root", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.DisconnectRootDirHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				log.Println("quics-client : cannot unmarshal")
				log.Println(err)
			}

			connection.UnRegisterRootDirRequest(body.RootDir, body.RootDirPw)
			connection.DirWatchStop(body.RootDir)
		}
	})

	mux.HandleFunc("/api/v1/connect/list/remote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			connection.ShowListRemoteRootDirRequest()
		}
	})

	mux.HandleFunc("/api/v1/disconnect/server", func(w http.ResponseWriter, r *http.Request) {
		body := &types.DisconnectClientHTTP3{}
		err := types.UnmarshalJSONFromRequest(r, body)
		if err != nil {
			log.Println("quics-client : cannot unmarshal")
			log.Println(err)
		}
		connection.DisconnectClientRequest(body.ClientPw)

	})

	mux.HandleFunc("/api/v1/rescan", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			sync.Rescan()
		}
	})

	mux.HandleFunc("/api/v1/status/root", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			showStatus := &types.ShowStatusHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, showStatus)
			if err != nil {
				log.Println("quics-client : cannot unmarshal")
				log.Println(err)
			}

			sync.ShowStatus(showStatus.Filepath)
		}
	})

	return mux
}
