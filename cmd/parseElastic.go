/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"ejol/ejlog-server/controller"

	"github.com/spf13/cobra"
)

// parseElasticCmd represents the parseElastic command
var parseElasticCmd = &cobra.Command{
	Use:   "parse-elastic",
	Short: "Elastic Search Parse Data",
	Long:  `Elastic Search Parse Data for all brand.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := controller.HyosungParseProcess()
		if err != nil {
			controller.ErrorLogger.Printf("Error Application %s", err)
		}
		// controller.LoadConfigKeyword()
		// fmt.Println("Welcome")
	},
}

func init() {
	rootCmd.AddCommand(parseElasticCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseElasticCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseElasticCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
