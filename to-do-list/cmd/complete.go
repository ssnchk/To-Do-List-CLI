package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"to-do-list/internal/config"
	"to-do-list/internal/storage"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Marks task as completed <taskId>",
	Long:  `Marks task as completed <taskId>`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln(err.Error())
		}

		config, err := config.GetConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		repo := storage.NewJsonTaskRepository(config.DataFilePath)

		task, err := repo.FindTask(taskId)
		if err != nil {
			log.Fatalln(err.Error())
		}

		task.IsCompleted = true

		if err := repo.UpdateTask(task); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
