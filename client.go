package main

import (
	"bufio"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net"
	"os"
	"strconv"
)

const (
	MSG_DISCONNECT = "Disconnected from server\n"
)
var messageChan = make(chan string, 1)


func Write(conn net.Conn, text string) {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(text)
	if err != nil {
		fmt.Println("Error writing client string: ", err)
		os.Exit(1)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing client writer: ", err)
		os.Exit(1)
	}
}

func Read(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(MSG_DISCONNECT)
			return MSG_DISCONNECT, err
		}
		messageChan <- str
	}
}

func SetPort() (int, error) {
	fmt.Println("Please input server port")
	var port int
	_, err := fmt.Scanln(&port)
	if err != nil {
		fmt.Println("Error reading port. Assuming port = 3333")
		return 3333, err
	}
	return port, err
}
func main() {
	var outTE,inTE  *walk.TextEdit
	portStr := ":" + strconv.Itoa(3333)
	conn, err := net.Dial("tcp", portStr)
	if err != nil {
		fmt.Println("Error while creating client connection: ", err)
	}
	mv := MainWindow{
		Title:   "Chat",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &outTE,
						ReadOnly: true,
						MaxSize:  Size{600, 350},
					},
				},
			},
			TextEdit{AssignTo: &inTE,
				MinSize: Size{Height: 50, Width: 600},
				MaxSize: Size{Height: 50, Width: 600}},

			PushButton{
				Text: "Send",
				OnClicked: func() {
					inTE.AppendText("\n")
					Write(conn, inTE.Text())
					inTE.SetText("")
					inTE.SetFocus()
					go func() {
						for {
							select{
							case msg:= <-messageChan:
								outTE.Synchronize(func() {
									outTE.AppendText(msg + "\r\n")
								})
							}
						}
					}()
				},
			},
		},
	}
	go Read(conn)
	mv.Run()
}
