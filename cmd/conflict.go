package main

import (
	"encoding/json"
	"log"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

var (
	ChosenLocalPath  string
	ChosenRemotePath string
)

func init() {
	ConflictCmd := ConflictCmd()
	ConflictCmd.AddCommand(ShowConflictListCmd())

	ChooseOneCmd := ChooseOneCmd()
	ChooseOneCmd.Flags().StringVarP(&ChosenLocalPath, "local", "l", "", "type local file absolute path")
	ChooseOneCmd.Flags().StringVarP(&ChosenRemotePath, "remote", "r", "", "type remote file absolute path")

	ConflictCmd.AddCommand(ChooseOneCmd)

	rootCmd.AddCommand(ConflictCmd)
}

func ConflictCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "conflict",
		Short: "command about conflict",
	}
}

// e.g. qic conflict list
func ShowConflictListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "show conflict list",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			bres := restClient.GetRequest("/api/v1/conflict/list")
			log.Println(bres.String())
		},
	}
}

// e.g. qic conflict choose --local /home/username/sync/test.txt
// e.g. qic conflict choose --remote /home/username/sync/test.txt
func ChooseOneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "choose",
		Short: "choose one between two options",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()

			if ChosenLocalPath != "" && ChosenRemotePath == "" { // In case Local
				chosenFilePath := &types.ChosenFilePathHTTP3{
					FilePath: ChosenLocalPath,
				}
				body, err := json.Marshal(chosenFilePath)
				if err != nil {
					log.Println(err)
				}
				bres := restClient.PostRequest("/api/v1/conflict/choose/local", "application/json", body)
				log.Println(bres.String())
			}

			if ChosenRemotePath != "" && ChosenLocalPath == "" { // In case Remote
				chosenFilePath := &types.ChosenFilePathHTTP3{
					FilePath: ChosenRemotePath,
				}
				body, err := json.Marshal(chosenFilePath)
				if err != nil {
					log.Println(err)
				}
				bres := restClient.PostRequest("/api/v1/conflict/choose/remote", "application/json", body)
				log.Println(bres.String())
			}
		},
	}
}
