package rtserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var todoList []string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getCmd(input string) string {
	inputArr := strings.Split(input, " ")
	return inputArr[0]
}

func getMessage(input string) string {
	inputArr := strings.Split(input, " ")
	var result string
	for i := 1; i < len(inputArr); i++ {
		result += inputArr[i]
	}
	return result
}

func updateTodoList(input string) {
	tmpList := todoList
	todoList = []string{}
	for _, val := range tmpList {
		if val == input {
			continue
		}
		todoList = append(todoList, val)
	}
}

func HandleCommunication(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Client connected")

	defer func() {
		fmt.Println("Client disconnected")
		conn.Close()
	}()

	for {
		messageType, p, err := conn.ReadMessage()

		messageToStr := string(p)

		fmt.Println("messageToStr", messageToStr)
		if err != nil {
			return
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}
