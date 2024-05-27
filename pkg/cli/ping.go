package cli

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping singlestore db",
	Run: func(cmd *cobra.Command, args []string) {
		connInfo := getConnectionInfo(cmd)
		connection := connInfo.getConnectionString()

		concurrency, _ := cmd.Flags().GetInt("concurrency")
		count, _ := cmd.Flags().GetInt("count")

		jobs := make(chan bool, count)
		results := make(chan bool, count)

		for i := 0; i < concurrency; i++ {
			go ping(connection, jobs, results)
		}

		for i := 0; i < count; i++ {
			jobs <- true
		}

		close(jobs)

		success := 0
		failure := 0
		for i := 0; i < count; i++ {
			if <-results {
				success++
			} else {
				failure++
			}
		}

		fmt.Printf("Total number of connections: %d/%d\n", success, failure)
	},
}

func ping(connection string, jobs <-chan bool, results chan<- bool) {

	for range jobs {
		if db, err := sql.Open("mysql", connection); err != nil {
			results <- false
		} else {
			if err := db.Ping(); err != nil {
				results <- false
			} else {
				results <- true
			}

			db.Close()
		}
	}
}
