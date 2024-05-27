package cli

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "run",
	Short: "run a statement on singlestoer",
	Run: func(cmd *cobra.Command, args []string) {
		connInfo := getConnectionInfo(cmd)
		connection := connInfo.getConnectionString()

		sqlStr, _ := cmd.Flags().GetString("statement")

		if db, err := sql.Open("mysql", connection); err != nil {
			panic(err)
		} else {
			defer db.Close()

			if err := db.Ping(); err != nil {
				panic(err)
			} else {
				if duration, err := executeQuery(db, sqlStr); err != nil {
					panic(err)
				} else {
					fmt.Printf("Execution completed in %s\n", duration)
				}
			}
		}

	},
}

func executeQuery(db *sql.DB, sql string) (time.Duration, error) {
	start := time.Now()
	if stmt, err := db.Prepare(sql); err != nil {
		return time.Since(start), err
	} else {

		defer stmt.Close()

		stmt.QueryRow()

		return time.Since(start), nil

	}
}
