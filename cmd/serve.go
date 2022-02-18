/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/job"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var server = controller.Server{}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start application ejlog server",
	Long:  `Ejlog-Server adalah aplikasi untuk menerima request data dari ejlog mesin ATM/CRM.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load(".env")
		if err != nil {
			controller.ErrorLogger.Fatal("Error load file env : ", err)
		}
		go job.JobCacheAtmMappings()
		time.Sleep(2 * time.Second)

		server.Run(":3000")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
