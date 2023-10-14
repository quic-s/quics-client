package http3

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-s/quics-client/pkg/sync"
	"github.com/quic-s/quics-client/pkg/types"
)

func SetupHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%#v\n", r)
		w.Write([]byte("hello, Quics Client here"))
	})

	mux.HandleFunc("/api/v1/connect/server", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterClientHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/server] ERROR : " + err.Error()))
			}
			err = sync.ClientRegistration(body.ClientPW, body.Host, body.Port)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/server] ERROR : " + err.Error()))
			} else {
				w.Write([]byte("quics-client : [/api/v1/connect/server] Resp : OK"))
			}
		}
	})

	mux.HandleFunc("/api/v1/connect/root/local", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterRootDirHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				log.Println("quics-client : [/api/v1/connect/root/local] ERROR : cannot marshal", err)

			}

			err = sync.RegistRootDir(body.RootDir, body.RootDirPw, "LOCAL")
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/root/local] ERROR : " + err.Error()))
			}
			log.Println("send response")
			w.Write([]byte("quics-client : [/api/v1/connect/root/local] RESP : OK"))

		}
	})

	mux.HandleFunc("/api/v1/connect/root/remote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.RegisterRootDirHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/root/remote] ERROR : " + err.Error()))
			}

			err = sync.RegistRootDir(body.RootDir, body.RootDirPw, "REMOTE")
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/root/remote] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/connect/root/remote] RESP : OK"))
		}
	})

	// mux.HandleFunc("/api/v1/disconnect/root", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "POST":
	// 		body := &types.DisconnectRootDirHTTP3{}
	// 		err := types.UnmarshalJSONFromRequest(r, body)
	// 		if err != nil {
	// 			log.Println("quics-client : cannot unmarshal")
	// 			log.Println(err)
	// 		}

	// 		connection.UnRegisterRootDirRequest(body.RootDir, body.RootDirPw)
	// 		connection.DirWatchStop(body.RootDir)
	// 	}
	// })

	mux.HandleFunc("/api/v1/connect/list/remote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			rootList, err := sync.GetRemoteRootList()
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/connect/list/remote] ERROR : " + err.Error()))
			}
			result := ""
			for i, root := range rootList.RootDirList {
				result += fmt.Sprintf("%d. %s\n", i, root)
			}
			w.Write([]byte("quics-client : [/api/v1/connect/list/remote] RESP : " + result))

		}
	})

	// mux.HandleFunc("/api/v1/disconnect/server", func(w http.ResponseWriter, r *http.Request) {
	// 	body := &types.DisconnectClientHTTP3{}
	// 	err := types.UnmarshalJSONFromRequest(r, body)
	// 	if err != nil {
	// 		log.Println("quics-client : cannot unmarshal")
	// 		log.Println(err)
	// 	}
	// 	connection.DisconnectClientRequest(body.ClientPw)

	// })

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
				w.Write([]byte("quics-client : [/api/v1/status/root] ERROR : " + err.Error()))
			}

			result, err := sync.ShowStatus(showStatus.Filepath)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/status/root] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/status/root] RESP : " + result))
		}
	})

	mux.HandleFunc("/api/v1/conflict/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			//TODO GET Conflict List from server
		}
	})

	mux.HandleFunc("/api/v1/conflict/choose/server", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			chooseOne := &types.ChosenFilePathHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, chooseOne)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/choose/server] ERROR : " + err.Error()))
			}
			// err = sync.ChooseOne(chooseOne.FilePath, "SERVER")
			// if err != nil {
			// 	w.Write([]byte("quics-client : [/api/v1/conflict/choose/server] ERROR : " + err.Error()))
			// }
		}
	})

	return mux
}
