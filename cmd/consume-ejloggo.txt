/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"ejol/ejlog-server/controller"
	"ejol/ejlog-server/job"
	"ejol/ejlog-server/utils"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// consumeEjlogCmd represents the consumeEjlog command
var consumeEjlogCmd = &cobra.Command{
	Use:   "consume-ejlog",
	Short: "Consume File Ejlog",
	Long:  `Consume File Ejlog already saved in path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consume-ejlog called")
		utils.InitUtils()
		go job.JobCacheAtmMappings()
		time.Sleep(1 * time.Second)
		// err := controller.ConsumeFileEjol()
		err := controller.ConsumeFileEjolSchedule()
		if err != nil {
			// models.ErrorLogger.Printf("Error : %s", err)
			log.Printf("Error : %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumeEjlogCmd)
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Printf("Error load file env : %s", err)
	// }

	// fmt.Println("ENV Loaded.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeEjlogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumeEjlogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
