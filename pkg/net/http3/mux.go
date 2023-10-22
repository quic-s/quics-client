package http3

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-s/quics-client/pkg/app"
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

	mux.HandleFunc("/api/v1/reboot", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			err := app.Reboot()

			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/shutdown] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/shutdown] Resp : OK"))

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
		case "GET":
			err := sync.Rescan()
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/rescan] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/rescan] RESP : OK"))
		}
	})

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			showStatus := &types.ShowStatusHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, showStatus)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/status] ERROR : " + err.Error()))
			}

			result, err := sync.ShowStatus(showStatus.Filepath)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/status] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/status] RESP : " + result))
		}
	})

	mux.HandleFunc("/api/v1/conflict/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			result, err := sync.PrintCFOptions()
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/list] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/conflict/list] RESP : " + result))

		}
	})

	mux.HandleFunc("/api/v1/conflict/choose", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			chooseOne := &types.ChosenFilePathHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, chooseOne)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/choose] ERROR : " + err.Error()))
			}
			err = sync.ChooseOne(chooseOne.FilePath, chooseOne.Candidate)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/choose] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/conflict/choose] RESP : OK"))
		}
	})

	mux.HandleFunc("/api/v1/conflict/download", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			conflictDownload := &types.ConflictDownloadHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, conflictDownload)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/download] ERROR : " + err.Error()))
			}
			err = sync.ConflictDownload(conflictDownload.FilePath)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/conflict/download] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/conflict/download] RESP : OK"))
		}
	})

	mux.HandleFunc("/api/v1/share/stop", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.StopShareHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/share/stop] ERROR : " + err.Error()))
			}
			err = sync.StopShare(body.Link)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/share/stop] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/share/stop] RESP : OK"))
		}
	})

	mux.HandleFunc("/api/vi/share/file", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.ShareFileHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/share/file] ERROR : " + err.Error()))
			}
			link, err := sync.GetShareLink(body.FilePath, uint64(body.MaxCnt))
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/share/file] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/share/file] RESP : " + link))
		}
	})

	mux.HandleFunc("/api/v1/history/rollback", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.HistoryRollBackHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/rollback] ERROR : " + err.Error()))
			}
			err = sync.RollBack(body.FilePath, body.Version)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/rollback] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/history/rollback] RESP : OK"))
		}
	})

	mux.HandleFunc("/api/v1/history/show", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.HistoryShowHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/show] ERROR : " + err.Error()))
			}
			historyList, err := sync.HistoryShow(body.FilePath, body.CntFromHead)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/show] ERROR : " + err.Error()))
			}
			result := ""
			for i, history := range historyList {
				result += fmt.Sprintf("%d. %s\n", i, history)
			}
			w.Write([]byte("quics-client : [/api/v1/history/show] RESP : " + result))
		}
	})

	mux.HandleFunc("/api/v1/history/download", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body := &types.HistoryDownloadHTTP3{}
			err := types.UnmarshalJSONFromRequest(r, body)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/download] ERROR : " + err.Error()))
			}
			err = sync.HistoryDownload(body.FilePath, body.Version)
			if err != nil {
				w.Write([]byte("quics-client : [/api/v1/history/download] ERROR : " + err.Error()))
			}
			w.Write([]byte("quics-client : [/api/v1/history/download] RESP : OK"))
		}
	})

	return mux
}
