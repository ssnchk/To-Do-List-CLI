package display

import (
	"fmt"
	"github.com/mergestat/timediff"
	"os"
	"text/tabwriter"
	"to-do-list/internal/models"
)

func DisplayTasks(tasks []models.Task) {
	if len(tasks) == 0 {
		fmt.Println(`No tasks created. 
Create new tasks using: to-do-list create <description>`)
		return
	}

	writer := tabwriter.NewWriter(os.Stdout,
		0, 8, 1, '\t', 0)

	fmt.Fprintln(writer, "Id\tDescription\tCreated\tStatus")
	for _, task := range tasks {
		fmt.Fprintf(writer, "%d\t%s\t%s\t", task.Id, task.Description, timediff.TimeDiff(task.CreateTime))
		if task.IsCompleted {
			fmt.Fprintln(writer, "Completed")
		} else {
			fmt.Fprintln(writer, "In Progress")
		}
	}
	writer.Flush()
}
