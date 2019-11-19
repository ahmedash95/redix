package server

import (
	"strings"

	"github.com/alash3al/redix/db"
	"github.com/tidwall/redcon"
)

// Options a server related options
type Options struct {
	RESPAddr   string
	DriverName string
	DriverOpts map[string]interface{}
}

// ListenAndServe start listening and serving the incomming requests
func ListenAndServe(opts Options) error {
	return redcon.ListenAndServe(opts.RESPAddr, func(incommingConn redcon.Conn, incommingCommand redcon.Command) {
		if len(incommingCommand.Args) < 1 {
			incommingConn.WriteError(errNoCommand.Error())
			return
		}

		commandName := strings.ToLower(string(incommingCommand.Args[0]))
		handler, ok := Handlers[commandName]

		if !ok {
			incommingConn.WriteError(errUnknownCommand.Error())
			return
		}

		db := &db.DB{}

		c := Context{
			conn:    incommingConn,
			command: incommingCommand,
			db:      db,
		}

		if err := handler.Func(c); err != nil {
			incommingConn.WriteError(err.Error())
			return
		}
	}, nil, nil)
}
