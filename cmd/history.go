package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

var (
	RollBackPath    string
	RollBackVersion string

	ShowPath     string
	ShowFromHead string

	DownloadPath    string
	DownloadVersion string
)

func init() {

	HistoryCmd := HistoryCmd()
	RollBackCmd := RollBackCmd()
	RollBackCmd.Flags().StringVarP(&RollBackPath, "path", "p", "", "path of file")
	if err := RollBackCmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}
	RollBackCmd.Flags().StringVarP(&RollBackVersion, "version", "v", "", "version of file")
	if err := RollBackCmd.MarkFlagRequired("version"); err != nil {
		log.Fatal(err)
	}

	HistoryCmd.AddCommand(RollBackCmd)

	ShowCmd := ShowCmd()
	ShowCmd.Flags().StringVarP(&ShowPath, "path", "p", "", "path of file")
	if err := ShowCmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}
	ShowCmd.Flags().StringVarP(&ShowFromHead, "from-head", "f", "", "show history from head")
	if err := ShowCmd.MarkFlagRequired("from-head"); err != nil {
		log.Fatal(err)
	}

	HistoryCmd.AddCommand(ShowCmd)

	DownloadCmd := DownloadCmd()
	DownloadCmd.Flags().StringVarP(&DownloadPath, "path", "p", "", "path of file")
	if err := DownloadCmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}
	DownloadCmd.Flags().StringVarP(&DownloadVersion, "version", "v", "", "version of file")
	if err := DownloadCmd.MarkFlagRequired("version"); err != nil {
		log.Fatal(err)
	}

	HistoryCmd.AddCommand(DownloadCmd)

	rootCmd.AddCommand(HistoryCmd)

}

// e.g. qic history download —path —version
func DownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "download file from history",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			parsedCnt, err := strconv.ParseUint(DownloadVersion, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			cnt := uint64(parsedCnt)

			downloadHistoryReq := &types.HistoryDownloadHTTP3{
				FilePath: DownloadPath,
				Version:  cnt,
			}
			breq, err := json.Marshal(downloadHistoryReq)
			if err != nil {
				log.Fatal(err)
			}

			// Request to REST Server
			bres, err := restClient.PostRequest("/api/v1/history/download", "application/json", breq)
			if err != nil {
				log.Println(err)
			}
			log.Println(bres.String())

		},
	}
}

// e.g. qic history show --path --from-head
func ShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show history of commands",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			parsedCnt, err := strconv.ParseUint(ShowFromHead, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			cnt := uint64(parsedCnt)

			showHistoryReq := &types.HistoryShowHTTP3{
				FilePath:    ShowPath,
				CntFromHead: cnt,
			}
			breq, err := json.Marshal(showHistoryReq)
			if err != nil {
				log.Fatal(err)
			}

			//Request to REST Server
			bres, err := restClient.PostRequest("/api/v1/history/show", "application/json", breq)
			if err != nil {
				log.Println(err)
			}
			log.Println(bres.String())

		},
	}
}

// e.g. qic history  rollback --path --version
func RollBackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "rollback to previous version",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			parsedCnt, err := strconv.ParseUint(RollBackVersion, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			cnt := uint64(parsedCnt)

			rollbackReq := &types.HistoryRollBackHTTP3{
				FilePath: RollBackPath,
				Version:  cnt,
			}
			breq, err := json.Marshal(rollbackReq)
			if err != nil {
				log.Fatal(err)
			}

			//Request to REST Server
			bres, err := restClient.PostRequest("/api/v1/history/rollback", "application/json", breq)
			if err != nil {
				log.Println(err)
			}
			log.Println(bres.String())
		},
	}
}

func HistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "managing history of commands",
	}
}
