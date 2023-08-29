package main

import (
	"log"

	"github.com/quic-s/quics-client/pkg/sync"
	"github.com/spf13/cobra"
)

const (
	LocalRootDirCmd      = "local"
	LocalRootDirShortCmd = "l"

	RemoteRootDirCmd      = "remote"
	RemoteRootDirShortCmd = "r"

	DisableRootDirCmd      = "disable"
	DisableRootDirShortCmd = "d"

	DirStatusCmd      = "pick"
	DirStatusShortCmd = "p"

	ReScanCmd      = "rescan"
	ReScanShortCmd = "s"
)

var (
	LocalRootDir   string
	RemoteRootDir  string
	DisableRootDir string

	DirForStatus string
	RescanDir    string
)

func init() {
	SyncCmd := SyncCmd()
	SyncCmd.Flags().StringVarP(&LocalRootDir, LocalRootDirCmd, LocalRootDirShortCmd, "", "decide local root directory")
	SyncCmd.Flags().StringVarP(&RemoteRootDir, RemoteRootDirCmd, RemoteRootDirShortCmd, "", "decide remote root directory")
	SyncCmd.Flags().StringVarP(&DisableRootDir, DisableRootDirCmd, DisableRootDirShortCmd, "", "choose witch root be disable ")
	rootCmd.AddCommand(SyncCmd)

	StatusCmd := StatusCmd()
	SyncCmd.Flags().StringVarP(&DirForStatus, DirStatusCmd, DirStatusShortCmd, "", "decide local root directory")
	rootCmd.AddCommand(StatusCmd)

	RescanCmd := RescanCmd()
	RescanCmd.Flags().StringVarP(&RescanDir, ReScanCmd, ReScanShortCmd, "", "decide local root directory")
	rootCmd.AddCommand(RescanCmd)

}

func RescanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rescan",
		Short: "rescan sync dir",
		Run: func(cmd *cobra.Command, args []string) {
			if RescanDir != "" {
				sync.RescanCertainDir(RescanDir)
			} else {
				sync.RescanAllDir()
			}
		},
	}
}

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show status of sync dir",
		Run: func(cmd *cobra.Command, args []string) {
			if DirForStatus != "" {
				log.Println(sync.ShowStatus(DirForStatus))
			} else {
				log.Println(sync.ShowAllStatus())
			}
		},
	}
}

func SyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync local directory(path/nn) to remote directory or vice versa, or disable sync",
		Run: func(cmd *cobra.Command, args []string) {
			if LocalRootDir != "" {
				sync.MakeLocalSync(LocalRootDir)
			}
			if RemoteRootDir != "" {
				sync.MakeRemoteSync(RemoteRootDir)
			}
			if DisableRootDir != "" {
				sync.MakeDisableSync(DisableRootDir)
			}
		},
	}
}
