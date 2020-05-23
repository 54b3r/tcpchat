package pkg

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Server represents the data structure for the chat server
type Server struct {
	Rooms    map[string]*Room
	commands chan Command
	Protocol string
	Port     string
}

func NewServer() *Server {
	port := os.Getenv("CHAT_PORT")
	proto := os.Getenv("CHAT_PROTOCOL")

	if port == "" {
		Logger(true, "[ERROR]: CHAT_PORT environment variable has not been exported, please export to start the server properly", nil)
		os.Exit(1)
	}

	if proto == "" {
		Logger(false, "[INFO]: no CHAT_PROTOCOL supplied, defaulting to tcp", nil)
		proto = "tcp"
	}

	server := &Server{
		Rooms:    make(map[string]*Room),
		commands: make(chan Command),
		Port:     port,
		Protocol: proto,
	}
	return server
}
func Logger(isFatal bool, format string, v interface{}) {
	if isFatal == true {
		if v == nil {
			log.Fatalf(format)
		}
		log.Fatalf(format, v)

	} else {
		if v == nil {
			log.Printf(format)
		}
		log.Printf(format, v)
	}
}

func (s *Server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case cmdNICK:
			s.nick(cmd.client, cmd.args[1])
		case cmdJOIN:
			s.join(cmd.client, cmd.args[1])
		case cmdROOMS:
			s.listRooms(cmd.client)
		case cmdMSG:
			s.msg(cmd.client, cmd.args)
		case cmdQUIT:
			s.quit(cmd.client)
		}
	}
}
func (s *Server) NewClient(conn net.Conn) {
	c := &Client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
	Logger(false, "[INFO]: New client has joined "+c.nick+":%s", conn.RemoteAddr().String())
	c.readInput()
}

// set the nickname of the user
func (s *Server) nick(c *Client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("Dopesauce, I'm going to call you %s", c.nick))
}

// join a chat room, if not exitst, will be created (itteration on this should only allow creation if the user is part of a group???)
// in any case we should add the current user to the list of members in the room
func (s *Server) join(c *Client, roomName string) {
	// check if room exists
	r, ok := s.Rooms[roomName]
	// if doesnt exist, then we will create and assign the rooms map
	if !ok {
		r = &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}
		s.Rooms[roomName] = r
	}
	// add client to memebers map, since we are ussing the remote addr of the host, that will be the defining
	r.members[c.conn.RemoteAddr()] = c

	// quit prev room so that user may join another
	s.quitCurrentRoom((c))

	// on client side lets define the room
	c.room = r

	// broadcast a message to the room notifying of user joinng the room
	r.broadcast(c, fmt.Sprintf("%s:$s has joined the room", c.nick, c.conn.RemoteAddr()))
	Logger(false, "[INFO]: "+c.nick+":%s has joined "+c.room.name, c.conn.RemoteAddr())

	// send message to the user welcoming to room
	c.msg(fmt.Sprintf("welcome to %s", r.name))
}

// list available rooms on the chat server to user so they can select what room to /join
func (s *Server) listRooms(c *Client) {
	var rooms []string
	for name := range s.Rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are: [%s]", strings.Join(rooms, ", ")))
}

// server sends a message to the connected clients current room
func (s *Server) msg(c *Client, args []string) {
	// msg := strings.Join(args[1:len(args)], " ")

	msg := strings.Join(args[1:len(args)], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

// leave the chat server(close connction)
func (s *Server) quit(c *Client) {

	// get out of current room
	s.quitCurrentRoom(c)
	Logger(false, "[INFO]: Client %s has disconnected from the server, closing connection", c.conn.RemoteAddr().String())
	// write goodbye message
	c.msg("I see you want to leave this wonderous Chat server, Have a great one, and well see you back here soon!")
	// close user connection
	c.conn.Close()
}

func (s *Server) quitCurrentRoom(c *Client) {
	if c.room != nil {
		oldRoom := s.Rooms[c.room.name]
		delete(s.Rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
		// create log message of user:client
		Logger(false, "[INFO]: "+c.nick+":%s has left "+c.room.name, c.conn.RemoteAddr().String())
	}

}
