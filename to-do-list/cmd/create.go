package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"to-do-list/internal/config"
	"to-do-list/internal/storage"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates new task",
	Long:  `Creates new task with written description`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskDescription := args[0]

		config, err := config.GetConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		repo := storage.NewJsonTaskRepository(config.DataFilePath)

		if err := repo.AddTask(taskDescription); err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println("Task Created!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
