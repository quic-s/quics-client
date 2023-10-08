package main

import (
	"encoding/json"
	"log"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

const (
	DirStatusCmd      = "path"
	DirStatusShortCmd = "p"
)

var (
	DirForStatus string
)

func init() {
	SyncCmd := SyncCmd()

	StatusCmd := StatusCmd()

	SyncCmd.Flags().StringVarP(&DirForStatus, DirStatusCmd, DirStatusShortCmd, "", "decide local root directory")
	if err := SyncCmd.MarkFlagRequired(DirStatusCmd); err != nil {
		log.Println(err)
	}
	SyncCmd.AddCommand(StatusCmd)
	RescanCmd := RescanCmd()

	SyncCmd.AddCommand(RescanCmd)
	rootCmd.AddCommand(SyncCmd)

}

// e.g. qic sync resacn
func RescanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rescan",
		Short: "rescan sync dir",
		Run: func(cmd *cobra.Command, args []string) {

			restClient := NewRestClient()
			//TODO GET으로 변환 해야하는지 테스트
			restClient.PostRequest("/api/v1/rescan", "application/json", nil)
		},
	}
}

// e.g. qic sync status -p /home/username/sync/test.txt
func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show status of sync file",
		Run: func(cmd *cobra.Command, args []string) {

			restClient := NewRestClient()
			showStatusHTTP3 := &types.ShowStatusHTTP3{
				Filepath: DirForStatus,
			}
			body, err := json.Marshal(showStatusHTTP3)
			if err != nil {
				log.Println(err)
			}

			restClient.PostRequest("/api/v1/status/root", "application/json", body)

		},
	}
}

func SyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync dir",
	}
}
