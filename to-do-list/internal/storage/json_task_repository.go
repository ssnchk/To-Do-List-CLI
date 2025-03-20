package storage

import (
	"encoding/json"
	"fmt"
	"github.com/juju/fslock"
	"os"
	"time"
	"to-do-list/internal/models"
)

type JsonTaskRepository struct {
	filePath string
}

func (repo JsonTaskRepository) findTaskIndex(tasks []models.Task, taskId int) (int, error) {
	taskIndex := -1
	for idx, task := range tasks {
		if task.Id == taskId {
			taskIndex = idx
			break
		}
	}
	if taskIndex == -1 {
		return 0, fmt.Errorf("Couldn't find task with Id: %d", taskId)
	}

	return taskIndex, nil
}

func (repo *JsonTaskRepository) loadTasks(file *os.File) ([]models.Task, error) {
	tasks := make([]models.Task, 0)

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return tasks, nil
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *JsonTaskRepository) saveTasks(file *os.File, tasks []models.Task) error {
	encoder := json.NewEncoder(file)

	err := encoder.Encode(tasks)
	if err != nil {
		return err
	}

	return nil
}

func (repo *JsonTaskRepository) loadFile() (*os.File, *fslock.Lock, error) {
	lock := fslock.New(repo.filePath + ".lock")

	err := lock.Lock()
	if err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(repo.filePath, os.O_CREATE|os.O_RDWR, 0644)

	return file, lock, nil
}

func (repo *JsonTaskRepository) closeFile(file *os.File, lock *fslock.Lock) error {
	if file != nil {
		if err := file.Close(); err != nil {
			return err
		}
	}
	if lock != nil {
		err := lock.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *JsonTaskRepository) AddTask(description string) (err error) {
	file, lock, err := repo.loadFile()
	if err != nil {
		return err
	}
	defer func() {
		closeErr := repo.closeFile(file, lock)
		if closeErr != nil {
			err = closeErr
		}
	}()

	tasks, err := repo.loadTasks(file)
	if err != nil {
		return err
	}

	newTask := models.Task{
		Id: func() int {
			if len(tasks) == 0 {
				return 0
			}
			return tasks[len(tasks)-1].Id + 1
		}(),
		Description: description,
		CreateTime:  time.Now(),
		IsCompleted: false,
	}
	tasks = append(tasks, newTask)

	if err = repo.saveTasks(file, tasks); err != nil {
		return err
	}

	return nil
}

func (repo *JsonTaskRepository) GetAllTasks() (tasks []models.Task, err error) {
	file, lock, err := repo.loadFile()
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := repo.closeFile(file, lock)
		if closeErr != nil {
			err = closeErr
		}
	}()

	tasks, err = repo.loadTasks(file)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *JsonTaskRepository) FindTask(taskId int) (task *models.Task, err error) {
	file, lock, err := repo.loadFile()
	if err != nil {
		return
	}
	defer func() {
		closeErr := repo.closeFile(file, lock)
		if closeErr != nil {
			err = closeErr
		}
	}()

	tasks, err := repo.loadTasks(file)
	if err != nil {
		return
	}

	taskIndex, err := repo.findTaskIndex(tasks, taskId)
	if err != nil {
		return
	}
	task = &tasks[taskIndex]

	return
}

func (repo *JsonTaskRepository) DeleteTask(taskId int) (err error) {
	file, lock, err := repo.loadFile()
	if err != nil {
		return err
	}
	defer func() {
		closeErr := repo.closeFile(file, lock)
		if closeErr != nil {
			err = closeErr
		}
	}()

	tasks, err := repo.loadTasks(file)
	if err != nil {
		return err
	}

	taskIndex, err := repo.findTaskIndex(tasks, taskId)
	if err != nil {
		return err
	}

	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)

	if err = repo.saveTasks(file, tasks); err != nil {
		return err
	}

	return
}

func (repo *JsonTaskRepository) UpdateTask(task *models.Task) error {
	file, lock, err := repo.loadFile()
	if err != nil {
		return err
	}
	defer func() {
		closeErr := repo.closeFile(file, lock)
		if closeErr != nil {
			err = closeErr
		}
	}()

	tasks, err := repo.loadTasks(file)
	if err != nil {
		return err
	}

	taskIndex, err := repo.findTaskIndex(tasks, task.Id)
	if err != nil {
		return err
	}

	tasks[taskIndex] = *task

	if err = repo.saveTasks(file, tasks); err != nil {
		return err
	}

	return nil
}

func NewJsonTaskRepository(filePath string) *JsonTaskRepository {
	return &JsonTaskRepository{filePath: filePath}
}
