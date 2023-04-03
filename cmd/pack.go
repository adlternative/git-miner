/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/adlternative/git-miner/pkg/pack"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "check pack format",
	Long:  `check git pack file format`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := pack.Verify(args[0]); err != nil {
			log.Printf("verify failed: %v\n", err)
			os.Exit(1)
		}
		log.Printf("%s ok", args[0])
	},
}

func init() {
	rootCmd.AddCommand(packCmd)
}
