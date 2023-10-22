package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func MakeHash(after string, info os.FileInfo) string {
	h := sha512.New()
	h.Write([]byte(after)) // /root/*
	h.Write([]byte(info.ModTime().UTC().String()))
	h.Write([]byte(info.Mode().String()))
	h.Write([]byte(fmt.Sprint(info.Size())))
	return hex.EncodeToString(h.Sum(nil))
}

func CopyFile(src string, dst string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		err := os.MkdirAll(dst, fileInfo.Mode())
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		file, err := os.Open(dst)
		if err != nil {
			return err
		}

		// Set file metadata.
		err = file.Chmod(fileInfo.Mode())
		if err != nil {
			return err
		}
		err = os.Chtimes(dst, time.Now(), fileInfo.ModTime())
		if err != nil {
			log.Println("quics: ", err)
			return err
		}
		return nil
	}

	// copy src file to latest file
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	defer srcFile.Close()

	// Open file with O_TRUNC flag to overwrite the file when the file already exists.
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		// If the file does not exist, create the file.
		if os.IsNotExist(err) {
			dir, _ := filepath.Split(dst)
			if dir != "" {
				err := os.MkdirAll(dir, 0700)
				if err != nil {
					return err
				}
			}
			dstFile, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer dstFile.Close()

	n, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println("quics: ", err)
		return err
	}
	if n != fileInfo.Size() {
		return errors.New("quics: copied file size is not equal to original file size")
	}

	// Set file metadata.
	err = dstFile.Chmod(fileInfo.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(dst, time.Now(), fileInfo.ModTime())
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}
	return nil
}
