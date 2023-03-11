package database

import (
	"errors"
	"fmt"
	"task-tracker-service/internal/domain/entity"
)

type DayToDoStorage struct {
	db []entity.DayToDo
}

func NewDayToDoStorage() *DayToDoStorage {
	return &DayToDoStorage{db: make([]entity.DayToDo, 0, 10)}
}

func (s *DayToDoStorage) Add(date string, task string) error {
	if existDate, index := s.isExistDate(date); existDate {
		if existTask, _ := s.isExistTask(index, task); existTask {
			return fmt.Errorf("task %v alredy exist", task)
		}
		s.db[index].Tasks = append(s.db[index].Tasks, task)
		return nil
	}
	s.db = append(s.db, entity.DayToDo{
		Date:  date,
		Tasks: []string{task},
	})
	return nil

}

func (s *DayToDoStorage) DeleteOne(date string, task string) error {
	if existDate, indexDate := s.isExistDate(date); existDate {
		if existTask, indexTask := s.isExistTask(indexDate, task); existTask {
			s.db[indexDate].Tasks = append(s.db[indexDate].Tasks[:indexTask], s.db[indexDate].Tasks[indexTask+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func (s *DayToDoStorage) DeleteAll(date string) (int, error) {
	var countTasks int
	if existDate, indexDate := s.isExistDate(date); existDate {
		countTasks = s.countTasks(indexDate)
		s.db = append(s.db[:indexDate], s.db[indexDate+1:]...)
		return countTasks, nil
	}
	return 0, errors.New("date not found")
}

func (s *DayToDoStorage) FindAllDate(date string) ([]string, error) {
	if existDate, indexDate := s.isExistDate(date); existDate {
		return s.db[indexDate].Tasks, nil
	}
	return nil, errors.New("date not found")
}

func (s *DayToDoStorage) FindAll() ([]entity.DayToDo, error) {
	return s.db, nil
}

func (s *DayToDoStorage) UpdateOne(date string, oldTask string, newTask string) error {
	if existDate, indexDate := s.isExistDate(date); existDate {
		if existTask, indexTask := s.isExistTask(indexDate, oldTask); existTask {
			s.db[indexDate].Tasks[indexTask] = newTask
			return nil
		}
	}
	return errors.New("event not found")
}

func (s *DayToDoStorage) isExistDate(date string) (bool, int) {
	for index, item := range s.db {
		if item.Date == date {
			return true, index
		}
	}
	return false, 0
}

func (s *DayToDoStorage) isExistTask(dateIndex int, task string) (bool, int) {
	for index, item := range s.db[dateIndex].Tasks {
		if task == item {
			return true, index
		}
	}
	return false, 0
}

func (s *DayToDoStorage) countTasks(dateIndex int) int {
	return len(s.db[dateIndex].Tasks)
}
