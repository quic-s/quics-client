package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Register Client
type RegisterClientHTTP3 struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	ClientPW string `json:"clientPW"`
}

type RegisterRootDirHTTP3 struct { //Local, Remote
	RootDir   string `json:"rootDir"` // LocalAbsPath
	RootDirPw string `json:"rootDirPw"`
}

type DisconnectRootDirHTTP3 struct {
	RootDir   string `json:"rootDir"`
	RootDirPw string `json:"rootDirPw"`
}

type DisconnectClientHTTP3 struct {
	ClientPw string `json:"clientPw"`
}

type ShowStatusHTTP3 struct {
	Filepath string `json:"filepath"`
}

type ChosenFilePathHTTP3 struct {
	FilePath string `json:"filepath"`
}

func UnmarshalJSON(data []byte, dstStruct any) error {
	err := json.Unmarshal(data, dstStruct)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalJSONFromRequest(r *http.Request, dstStruct any) error {
	buf := make([]byte, r.ContentLength)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	if n == 0 {
		return fmt.Errorf("empty body")
	}
	err = json.Unmarshal(buf, dstStruct)
	if err != nil {
		return err
	}

	return nil
}
