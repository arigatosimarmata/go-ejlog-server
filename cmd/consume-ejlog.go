/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// consumeEjlogCmd represents the consumeEjlog command
var consumeEjlogCmd = &cobra.Command{
	Use:   "consume-ejlog",
	Short: "Consume File Ejlog",
	Long:  `Consume File Ejlog already saved in path.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consume-ejlog called")
	},
}

func init() {
	rootCmd.AddCommand(consumeEjlogCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeEjlogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumeEjlogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
