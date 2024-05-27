package cli

import (
	"database/sql"
	"math/rand"

	"github.com/spf13/cobra"
)

type ConnectionInfo struct {
	username string
	password string
	hostname []string
	port     string
	database string
}

func getConnectionInfo(cmd *cobra.Command) ConnectionInfo {

	var c ConnectionInfo

	c.username, _ = cmd.Flags().GetString("username")
	c.password, _ = cmd.Flags().GetString("password")
	c.hostname, _ = cmd.Flags().GetStringSlice("hostname")
	c.port, _ = cmd.Flags().GetString("port")
	c.database, _ = cmd.Flags().GetString("database")

	return c

}

func (c ConnectionInfo) getConnectionString() string {
	i := rand.Intn(len(c.hostname))

	connStr := c.username + ":" + c.password + "@tcp(" + c.hostname[i] + ":" + c.port + ")/" + c.database + "?parseTime=true"
	//fmt.Println(connStr)

	return connStr
}

func (c ConnectionInfo) getDB() (*sql.DB, error) {
	if db, err := sql.Open("mysql", c.getConnectionString()); err != nil {
		return nil, err
	} else {

		if err := db.Ping(); err != nil {
			return nil, err
		} else {
			return db, nil
		}
	}
}
