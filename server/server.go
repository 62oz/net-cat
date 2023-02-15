package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	openConnections        = make(map[net.Conn]bool)
	UserConn               = make(map[net.Conn]string)
	newConnection          = make(chan net.Conn)
	deadConnection         = make(chan net.Conn)
	chatlog         string = "\n"
	username        string = ""
	server_ON              = true
)

func Server(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Cannot listen to port.")
		log.Fatal(err)
	}
	fmt.Println("Listening to port :" + port)
	defer ln.Close()

	go func() {
		for {

			// Accepting connections
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Cannot accept connection.")
				log.Fatal(err)
			}

			welcomeMsg := "Welcome to TCP-Chat!\n       _nnnn_\n      dGGGGMMb\n     @p~qp~~qMb\n     M|@||@) M|\n     @,----.JM|\n    JS^\\__/  qKL\n   dZP        qKRb\n  dZP          qKKb\n  fZP            SMMb\n  HZM            MMMM\n  FqM            MMMM\n__| \".        |\\dS\"qML\n|    `.       | `' \\Zq\n_)     \\.___.,|     .'\n\\____  )MMMMMP|   .'\n `-'       `--'\n"

			conn.Write([]byte(welcomeMsg))

			if len(UserConn) == 10 {
				conn.Write([]byte("\nServer is full </3\nCome back later :)\n."))
				conn.Close()
			}

			username = AskNick(conn)

			openConnections[conn] = true
			newConnection <- conn
			UserConn[conn] = username
			conn.Write([]byte(chatlog + "\n"))

			chatlog += "\n" + UserConn[conn] + " has joined the chat!"
			for user := range UserConn {
				if user != conn {
					user.Write([]byte("\n" + username + " has joined the chat!"))
					dt := time.Now()
					stamp := fmt.Sprintf("\n[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[user])
					user.Write([]byte(stamp))
				}
			}

			dt := time.Now()
			stamp := fmt.Sprintf("[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[conn])
			conn.Write([]byte(stamp))
		}

		//}
	}()

	for {

		select {
		case conn := <-newConnection:

			//Invoke broadcast message(broadcasts to other conns)
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			//remove/delete the conn
			chatlog += "\n" + UserConn[conn] + " has left the chat!"
			exportChatlog(chatlog)
			for user := range openConnections {
				if user == conn {
					break
				}
			}
			for user := range UserConn {
				if user != conn {
					user.Write([]byte("\n" + UserConn[conn] + " has left the chat!"))
					dt := time.Now()
					stamp := fmt.Sprintf("\n[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[user])
					user.Write([]byte(stamp))
				}
			}
			delete(openConnections, conn)
		}
	}
}

//Ask user to enter a valid username
func AskNick(conn net.Conn) string {
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))

		reader := bufio.NewReader(conn)
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unexpected input reading error.")
			conn.Close()
			return ""
		}

		username = strings.Trim(username, " \r\n")
		for user := range UserConn {
			if UserConn[user] == username {
				conn.Write([]byte("[USERNAME TAKEN]\n"))
			}
			continue
		}
		if username != "" && !containsSpace(username) && len(username) <= 20 && isAlphaNumeric(username) {
			return username
		} else if len(username) > 20 {
			conn.Write([]byte("[MAX 20 CHARS]\n"))
		} else if !isAlphaNumeric(username) {
			conn.Write([]byte("[ALPHANUMERIC CHARACTERS ONLY]\n"))
		} else {
			conn.Write([]byte("[EMPTY USERNAME]\n"))
		}
	}
}

//Broadcast message to other users
func broadcastMessage(conn net.Conn) {
	for {

		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSuffix(message, "\n")
		if message == "\\nick" {
			newUsername := AskNick(conn)
			for user := range UserConn {
				if user == conn {
					UserConn[conn] = newUsername
				}
			}
			for user := range UserConn {
				chatlog += "\n" + UserConn[conn] + " is now " + newUsername
				if user != conn {
					user.Write([]byte("\n" + username + " is now " + newUsername))
					dt := time.Now()
					stamp := fmt.Sprintf("\n[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[user])
					user.Write([]byte(stamp))
				}
			}
			dt := time.Now()
			stamp := fmt.Sprintf("[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[conn])
			conn.Write([]byte(stamp))
		} else {
			for user := range UserConn {
				if user == conn {
					dt := time.Now()
					stamp := fmt.Sprintf("[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[conn])
					conn.Write([]byte(stamp))
					message = "\n" + stamp + message

					chatlog += message

				}
			}

			// exporting chatlog
			exportChatlog(chatlog)

			//loop through all connections and make sure
			//the msg is broadcast to the other users except
			//the conn that sent the message
			for user := range UserConn {
				if user != conn {
					user.Write([]byte(message))
					dt := time.Now()
					stamp := fmt.Sprintf("\n[%s][%s]: ", dt.Format("01-02-2006 15:04:05"), UserConn[user])
					user.Write([]byte(stamp))
				}
			}
		}
	}

	deadConnection <- conn

}

func exportChatlog(chatlog string) {

	f, err := os.Create("chat.log")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(chatlog); err != nil {
		log.Println(err)
	}
}

func isAlphaNumeric(username string) bool {
	r := []rune(username)
	for _, e := range r {
		if e < '0' || (e > '9' && e < 'A') || (e > 'Z' && e < 'a') || e > 'z' {
			return false
		}
	}
	return true
}

func containsSpace(username string) bool {
	r := []rune(username)
	for _, e := range r {
		if e == ' ' {
			return true
		}
	}
	return false
}
