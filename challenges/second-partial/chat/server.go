// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan string) // all incoming client messages
	admin         = ""
	clientChannel = make(map[string]chan string)
	clientConn    = make(map[string]net.Conn)
	clientList    = make(map[string]bool)
	clientKicked  = make(map[string]bool)
	clients       = make(map[client]bool) // all connected clients
)

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	servername := "irc-server > "
	username := ""
	usrbyte := make([]byte, 100)
	conn.Read(usrbyte)
	for _, element := range usrbyte {
		if element != 0 {
			username = username + string(element)
		} else {
			break
		}

	}
	if val, contains := clientList[username]; contains && val {
		ch <- servername + "Username taken"
		ch <- ""
		conn.Close()
	} else {
		ch <- servername + "Welcome to the Simple IRC Server"
		clientChannel[username] = ch
		clientConn[username] = conn
		clientList[username] = true
		clientKicked[username] = false

		ch <- servername + "Your user [" + username + "] is successfully logged"
		messages <- servername + "New connected user [" + username + "]"
		entering <- ch
		fmt.Printf(servername + "New connected user [" + username + "]\n")
		if len(clientList) == 1 {
			ch <- servername + "Congrats, you were the first user."
			admin = username
			ch <- servername + "You're the new IRC Server ADMIN"
			fmt.Printf(servername + "[" + username + "] was promoted as the channel ADMIN\n")
		}
		input := bufio.NewScanner(conn)
		for input.Scan() {
			if input.Text() == "/users" {
				for element := range clientList {
					ch <- servername + element
				}
			} else if input.Text() == "/time" {
				local, err := time.LoadLocation("Local")
				if err != nil {
					log.Fatal("No timezone available")
				}
				if local.String() == "Local" {
					ch <- servername + "Local Time: America/Mexico_City " + time.Now().Format("15:04")
				} else {
					ch <- servername + "Local Time: " + local.String() + " " + time.Now().Format("15:04")
				}
			}
			params := strings.Split(input.Text(), " ")
			if len(params) == 0 {
				messages <- username + ": " + input.Text()
			} else if params[0] == "/msg" {
				if len(params) > 1 {
					if val, contains := clientList[params[1]]; val && contains {
						msg := ""
						for i := 2; i < len(params); i++ {
							msg = msg + params[i] + " "
						}
						clientChannel[params[1]] <- "Message from [" + username + "] : " + msg
					} else {
						ch <- servername + "User doesn't exist"
					}
				} else {
					ch <- servername + "Incomplete parameters"
				}

			} else if params[0] == "/user" {
				if len(params) > 1 {
					if val, contains := clientList[params[1]]; val && contains {
						usrName := params[1]
						usrIP := clientConn[params[1]].RemoteAddr().String()
						ch <- servername + "username: " + usrName + ", IP:" + usrIP
					} else {
						ch <- servername + "User doesn't exist"
					}
				} else {
					ch <- servername + "Incomplete parameters"
				}
			} else if params[0] == "/kick" {
				if len(params) > 1 {
					if username == admin {
						if val, contains := clientList[params[1]]; val && contains {

							clientChannel[params[1]] <- servername + "You're kicked from this channel\n" + "Bad language is not allowed on this channel"
							clientConn[params[1]].Close()
							messages <- servername + "[" + params[1] + "] was kicked from channel for bad language policy violation"
							fmt.Printf(servername + "[" + params[1] + "] was kicked\n")
							clientKicked[username] = true
						} else {
							ch <- servername + "User doesn't exist"
						}
					} else {
						ch <- servername + "Commnad only allowed to ADMIN"
					}
				} else {
					ch <- servername + "Incomplete parameters"
				}
			} else {
				ch <- servername + "Invalid sub-command"
			}
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	if !clientKicked[username] {
		messages <- username + " has left"
	}
	conn.Close()
	clientList[username] = false
	if admin == username {
		for i := range clientList {
			if clientList[i] == true {
				admin = i
				clientChannel[i] <- servername + "You're the new IRC Server Admin"
			}
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	if len(os.Args) == 5 {
		host := os.Args[2]
		port := os.Args[4]
		listener, err := net.Listen("tcp", host+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("irc-server > Simple IRC Server started at " + host +
			":" + port + "irc-server > Ready for receiving new clients")

		go broadcaster()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			go handleConn(conn)
		}
	} else {
		log.Fatal("Wrong parameter input")
	}
}

//!-main
