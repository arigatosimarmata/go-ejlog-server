/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/job"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// schedulerCmd represents the scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "A brief description of your command",
	Long:  `Scheduler filecache untuk consume ejlog`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load(".env")
		if err != nil {
			controller.ErrorLogger.Fatal("Error load file env : ", err)
		}
		err = os.MkdirAll("./cache/"+time.Now().Format("20060102"), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		job.TestingCache3()
		fmt.Println("scheduler executed.")
	},
}

func init() {
	rootCmd.AddCommand(schedulerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schedulerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schedulerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
