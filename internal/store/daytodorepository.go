package store

import (
	"errors"
	"fmt"
	"task-tracker-service/internal/domain/entity"
)

const (
	isTaskExist          = "SELECT exists(SELECT 1 FROM day WHERE date = $1 AND task = $2)"
	isDateExist          = "SELECT exists(SELECT 1 FROM day WHERE date = $1)"
	insertTask           = "INSERT INTO day (date, task) VALUES ($1, $2)"
	deleteTask           = "DELETE FROM day WHERE date = $1 AND task = $2"
	countAllTasksByDate  = "SELECT count(id) FROM day WHERE date = $1"
	deleteAllTasksByDate = "DELETE FROM day WHERE date = $1"
	findAllTasksByDate   = "SELECT task FROM day WHERE date = $1"
	findAllTasks         = "SELECT date, task from day"
	updateTaskByDate     = "UPDATE day SET task = $3 WHERE date = $1 AND task = $2"
)

type DayToDoRepository struct {
	client *ClientDB
}

func NewDayToDoRepository(client *ClientDB) *DayToDoRepository {
	return &DayToDoRepository{client: client}
}

func (r *DayToDoRepository) Add(date string, task string) error {
	var isExist bool
	err := r.client.db.QueryRow(isTaskExist, date, task).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("task %v alredy exist", task)
	}
	err = r.client.db.QueryRow(insertTask, date, task).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *DayToDoRepository) DeleteOne(date string, task string) error {
	var isExist bool
	err := r.client.db.QueryRow(isTaskExist, date, task).Scan(&isExist)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("task not found")
	}
	err = r.client.db.QueryRow(deleteTask, date, task).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *DayToDoRepository) DeleteAll(date string) (int, error) {
	var tasksCount int
	err := r.client.db.QueryRow(countAllTasksByDate, date).Scan(&tasksCount)
	if err != nil {
		return 0, err
	}
	if tasksCount == 0 {
		return 0, errors.New("date not found")
	}
	err = r.client.db.QueryRow(deleteAllTasksByDate, date).Err()
	if err != nil {
		return 0, err
	}
	return tasksCount, nil
}
func (r *DayToDoRepository) FindAllDate(date string) ([]string, error) {
	var isExist bool
	err := r.client.db.QueryRow(isDateExist, date).Scan(&isExist)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.New("date not found")
	}
	rows, err := r.client.db.Query(findAllTasksByDate, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]string, 0, 5)
	for rows.Next() {
		var task string
		rows.Scan(&task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (r *DayToDoRepository) FindAll() ([]entity.DayToDo, error) {
	rows, err := r.client.db.Query(findAllTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//daysToDo := make([]entity.DayToDo, 0, 5)
	var daysToDo []entity.DayToDo
	for rows.Next() {
		var date string
		var task string
		var isExist bool
		rows.Scan(&date, &task)
		for i, day := range daysToDo {
			if day.Date == date {
				daysToDo[i].Tasks = append(daysToDo[i].Tasks, task)
				isExist = true
				break
			}
		}
		if !isExist {
			daysToDo = append(daysToDo, entity.DayToDo{
				Date:  date,
				Tasks: []string{task},
			})
		}
	}
	return daysToDo, nil

}
func (r *DayToDoRepository) UpdateOne(date string, oldTask string, newTask string) error {
	var isExist bool
	err := r.client.db.QueryRow(isTaskExist, date, oldTask).Scan(&isExist)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("task not found")
	}
	err = r.client.db.QueryRow(updateTaskByDate, date, oldTask, newTask).Err()
	if err != nil {
		return err
	}
	return nil
}
