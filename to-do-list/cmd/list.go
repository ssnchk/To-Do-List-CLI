package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"to-do-list/internal/config"
	"to-do-list/internal/display"
	"to-do-list/internal/models"
	"to-do-list/internal/storage"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of unfinished tasks",
	Long:  `Lists all of unfinished tasks`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		findAllFlag, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatalln(err.Error())
		}

		config, err := config.GetConfig()
		if err != nil {
			log.Fatalln(err.Error())
		}

		repo := storage.NewJsonTaskRepository(config.DataFilePath)

		tasks, err := repo.GetAllTasks()
		if err != nil {
			log.Fatalln(err.Error())
		}

		if findAllFlag {
			tasks = getFinishedTasks(tasks)
		}
		display.DisplayTasks(tasks)
	},
}

func getFinishedTasks(tasks []models.Task) []models.Task {
	finishedTasks := make([]models.Task, 0, len(tasks)/2)
	for _, task := range tasks {
		if !task.IsCompleted {
			finishedTasks = append(finishedTasks, task)
		}
	}

	return finishedTasks
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Lists all tasks")
}
