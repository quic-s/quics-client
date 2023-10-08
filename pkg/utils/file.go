package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
)

func MakeHash(after string, info os.FileInfo) string {
	h := sha512.New()
	h.Write([]byte(after)) // /root/*
	h.Write([]byte(info.ModTime().String()))
	h.Write([]byte(info.Mode().String()))
	h.Write([]byte(fmt.Sprint(info.Size())))
	return hex.EncodeToString(h.Sum(nil))
}
