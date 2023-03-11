package stdio

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"task-tracker-service/internal/domain/service"
)

type dayToDoHandler struct {
	dayToDoSvc *service.DayToDoService
}

type command struct {
	commandName string
	commandDate string
	commandText string
}

func NewDayToDoHandler(service *service.DayToDoService) *dayToDoHandler {
	return &dayToDoHandler{service}
}

func newCommand(commandName string, commandDate string, commandText string) *command {
	return &command{commandName: commandName, commandDate: commandDate, commandText: commandText}
}

func (h *dayToDoHandler) GetCommand() (string, error) {
	var str string
	var scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str = scanner.Text()
	return str, nil
}

func (h *dayToDoHandler) ParseCommand(str string) (bool, error) {
	const commandLength = 3
	strArray := strings.SplitN(str, " ", commandLength)
	for len(strArray) < commandLength {
		strArray = append(strArray, "")
	}
	cmd := newCommand(strArray[0], strArray[1], strArray[2])

	switch cmd.commandName {
	case "Add":
		if isValid, err := h.validDate(cmd.commandDate); isValid {
			err := h.dayToDoSvc.AddTask(cmd.commandDate, cmd.commandText)
			if err != nil {
				return false, fmt.Errorf("task not added, %v", err)
			}
			fmt.Printf("Task %v in %v added.\n", cmd.commandText, cmd.commandDate)
			return false, nil
		} else {
			return false, fmt.Errorf("task not added, date not valid, %v", err)
		}
	case "Del":
		if cmd.commandText != "" {
			err := h.dayToDoSvc.DeleteTask(cmd.commandDate, cmd.commandText)
			if err != nil {
				return false, fmt.Errorf("task not deleted, %v", err)
			}
			fmt.Printf("Task %v in date %v deleted successfully. \n", cmd.commandText, cmd.commandDate)
			return false, nil
		}
		tasks, err := h.dayToDoSvc.DeleteAllTasks(cmd.commandDate)
		if err != nil {
			return false, err
		}
		fmt.Printf("Deleted %v events. \n", tasks)
		return false, nil

	case "Find":
		tasks, err := h.dayToDoSvc.FindAllTaskDate(cmd.commandDate)
		if err != nil {
			return false, err
		}
		for _, task := range tasks {
			fmt.Println(task)
		}
		return false, nil

	case "Print":
		DaysToDoArray, err := h.dayToDoSvc.FindAllTasks()
		if err != nil {
			return false, err
		}
		for _, day := range DaysToDoArray {
			for _, task := range day.Tasks {
				fmt.Println(day.Date, task)
			}
		}
		return false, nil
	case "Quit":
		return true, nil

	}
	return false, nil
}

func (h *dayToDoHandler) validDate(date string) (bool, error) {
	matchString, err := regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`, date)
	if err != nil {
		return false, err
	}
	return matchString, nil
}
