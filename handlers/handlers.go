package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/luis-gamarra/todo-backend/todo"
	"io"
	"io/ioutil"
	"net/http"
)

func GetTodoListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, todo.Get())
}

func AddTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)

	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	c.JSON(statusCode, gin.H{"id": todo.Add(todoItem.Message)})
}

func DeleteTodoHandler(c *gin.Context) {
	todoID := c.Param("id")

	if err := todo.Delete(todoID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func CompleteTodoHandler(c *gin.Context) {
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}

	if todo.Complete(todoItem.ID) != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (todo.Todo, int, error){
	body, err := ioutil.ReadAll(httpBody)

	if err != nil {
		return todo.Todo{}, http.StatusInternalServerError, err
	}

	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (todo.Todo, int, error) {
	var todoItem todo.Todo
	err := json.Unmarshal(jsonBody, &todoItem)
	if err != nil {
		return todo.Todo{}, http.StatusBadRequest, err
	}

	return todoItem, http.StatusOK, nil
}
