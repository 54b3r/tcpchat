package pkg

type commandID int

const (
	cmdNICK commandID = iota
	cmdJOIN
	cmdROOMS
	cmdMSG
	cmdQUIT
)

// Command represents the needed data structure for a server command
type Command struct {
	id     commandID
	client *Client
	args   []string
}
