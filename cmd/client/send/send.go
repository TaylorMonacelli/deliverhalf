/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	client "github.com/taylormonacelli/deliverhalf/cmd/client"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	meta "github.com/taylormonacelli/deliverhalf/cmd/meta"
	sns "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger.Trace("send called")
		send()
	},
}

func init() {
	client.ClientCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func send() {
	data := meta.Fetch()
	jsBytes, _ := json.MarshalIndent(data, "", "    ")
	jsonStr := string(jsBytes)
	sns.SendJsonStr(jsonStr)
}
