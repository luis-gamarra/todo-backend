package todo

import (
	"errors"
	"github.com/rs/xid"
	"sync"
)

var (
	list []Todo
	mtx sync.Mutex
	once sync.Once
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Todo{}
}

type Todo struct {
	ID	string	`json:"id"`
	Message string 	`json:"message"`
	Complete bool 	`json:"complete"`
}

func Get() []Todo {
	return list
}

func Add(message string) string {
	todo := newTodo(message)
	mtx.Lock()
	list = append(list, todo)
	mtx.Unlock()
	return todo.ID
}

func Delete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

func Complete(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	setTodoCompleteByLocation(location)

	return nil
}

func newTodo(message string) Todo {
	return Todo {
		ID: xid.New().String(),
		Message: message,
		Complete: false,
	}
}

func findTodoLocation(id string) (int, error) {
	mtx.Lock()
	defer mtx.Unlock()
	for index, todo := range list {
		if isMatchingID(todo.ID, id) {
			return index, nil
		}
	}
	return 0, errors.New("Could not find todo based on id")
}

func removeElementByLocation(index int) {
	mtx.Lock()
	list = append(list[:index], list[index + 1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingID(idA, idB string) bool {
	return idA == idB
}
