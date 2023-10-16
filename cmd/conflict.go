package main

import (
	"encoding/json"
	"log"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

var (
	ChosenPath      string
	ChosenCandidate string
)

func init() {
	ConflictCmd := ConflictCmd()
	ConflictCmd.AddCommand(ShowConflictListCmd())

	ChooseOneCmd := ChooseOneCmd()
	ChooseOneCmd.Flags().StringVarP(&ChosenPath, "path", "p", "", "chosen local path")
	ChooseOneCmd.Flags().StringVarP(&ChosenCandidate, "candidate", "c", "", "chosen candidate uuid")
	if err := ChooseOneCmd.MarkFlagRequired("path"); err != nil {
		log.Println(err)
	}
	if err := ChooseOneCmd.MarkFlagRequired("candidate"); err != nil {
		log.Println(err)
	}

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

// e.g. qic conflict choose --path /home/username/sync/test.txt --candidate 1234567890abc
func ChooseOneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "choose",
		Short: "choose one among options",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			defer restClient.Close()

			chosenFilePath := &types.ChosenFilePathHTTP3{
				FilePath:  ChosenPath,
				Candidate: ChosenCandidate,
			}
			body, err := json.Marshal(chosenFilePath)
			if err != nil {
				log.Println(err)
			}
			bres := restClient.PostRequest("/api/v1/conflict/choose", "application/json", body)
			log.Println(bres.String())

		},
	}
}
