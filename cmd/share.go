package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

var (
	SharePath string
	ShareCnt  string

	ShareStopLink string
)

func init() {
	ShareCmd := ShareCmd()
	FileCmd := FileCmd()
	FileCmd.Flags().StringVarP(&SharePath, "path", "p", "", "path of file")
	if err := FileCmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}
	FileCmd.Flags().StringVarP(&ShareCnt, "cnt", "c", "", "count of share")
	if err := FileCmd.MarkFlagRequired("cnt"); err != nil {
		log.Fatal(err)
	}

	ShareCmd.AddCommand(FileCmd)

	StopCmd := StopCmd()
	StopCmd.Flags().StringVarP(&ShareStopLink, "link", "l", "", "link for stop sharing")
	if err := StopCmd.MarkFlagRequired("link"); err != nil {
		log.Fatal(err)
	}
	ShareCmd.AddCommand(StopCmd)

	rootCmd.AddCommand(ShareCmd)
}

// qic share file —path —cnt
func FileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "file",
		Short: "get link and share file to other client",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			//change ShareCnt type from string to uint
			cnt := uint(0)
			parsedCnt, err := strconv.ParseUint(ShareCnt, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			cnt = uint(parsedCnt)

			chosenFilePath := &types.ShareFileHTTP3{
				FilePath: SharePath,
				MaxCnt:   cnt,
			}
			breq, err := json.Marshal(chosenFilePath)
			if err != nil {
				log.Fatal(err)
			}
			bres := restClient.PostRequest("/api/v1/share/file", "application/json", breq)
			log.Println(bres.String())
		},
	}
}

// qic share stop --link
func StopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "stop share",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			chosenFilePath := &types.StopShareHTTP3{
				Link: ShareStopLink,
			}
			breq, err := json.Marshal(chosenFilePath)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(restClient.PostRequest("/api/v1/share/stop", "application/json", breq))
		},
	}
}

func ShareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "share",
		Short: "share file to other client",
	}
}
