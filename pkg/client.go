package pkg

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client represents the needed data structure for a chat server client
type Client struct {
	conn     net.Conn
	nick     string
	room     *Room
	commands chan<- Command
}

// blocking function used in goroutine to read the input of user, and if command matches well send a message to the channe, if not they will be sent to the user
// the commands will be received by the run function, based on the command id it will execute the command
func (c *Client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- Command{
				id:     cmdNICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- Command{
				id:     cmdJOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- Command{
				id:     cmdROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- Command{
				id:     cmdMSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- Command{
				id:     cmdQUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
			Logger(false, "[ERROR]: Bad command issued %s ", cmd)
		}
	}

}

// Write message (error) to tcp connection
func (c *Client) err(err error) {
	c.conn.Write([]byte("[ERROR]: " + err.Error() + "\n"))
}

// Write a message to the tcp connection
func (c *Client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
