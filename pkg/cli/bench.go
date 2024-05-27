package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/montanaflynn/stats"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "bench a statement on singlestore",
	Run: func(cmd *cobra.Command, args []string) {
		connInfo := getConnectionInfo(cmd)

		sqlStr, _ := cmd.Flags().GetString("statement")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		count, _ := cmd.Flags().GetInt("count")

		jobs := make(chan string, count)
		results := make(chan bool, count)

		db := make([]*sql.DB, 0)

		//fmt.Println("Building up connections")
		bar := progressbar.Default(int64(concurrency), "Building connections")

		for i := 0; i < concurrency; i++ {
			if r, err := connInfo.getDB(); err == nil {
				db = append(db, r)
			} else {
				fmt.Println(err)
			}

			bar.Add(1)
		}
		fmt.Printf("We have %d healty connections\n", len(db))

		start := time.Now()

		for i := 0; i < len(db); i++ {
			go benchQuery(db[i], jobs, results)
		}

		for i := 0; i < count; i++ {
			jobs <- sqlStr
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

		fmt.Printf("Total number of request: %d/%d in %s\n", success, failure, time.Since(start))

	},
}

func benchQuery(db *sql.DB, jobs <-chan string, results chan<- bool) {

	defer db.Close()

	for sql := range jobs {

		if _, err := db.Exec(sql); err == nil {
			results <- true
		} else {
			results <- false
		}

	}

}

type result struct {
	status  bool
	elapsed int64
}

func responseTimeStatistics(elapsed []float64) {
	median, _ := stats.Median(elapsed)
	p95, _ := stats.Percentile(elapsed, 95)
	p99, _ := stats.Percentile(elapsed, 99)
	mean, _ := stats.Mean(elapsed)
	jitter, _ := stats.StandardDeviation(elapsed)

	min, _ := stats.Min(elapsed)
	max, _ := stats.Max(elapsed)

	colorstring.Printf("\nResponse time (ms): \n")
	colorstring.Printf("Mean: %.2f Jitter: %.2f\n", float32(mean/1000), float32(jitter/1000))
	colorstring.Printf("Min: [green]%.2f[reset] p50: [yellow]%.2f[reset] p95: %.2f p99: %.2f Max: [red]%.2f[reset]\n\n", float32(min/1000), float32(median/1000), float32(p95/1000), float32(p99/1000), float32(max/1000))

}

// /
var vbenchCmd = &cobra.Command{
	Use:   "vbench",
	Short: "vector benchmark a statement on singlestore",
	Run: func(cmd *cobra.Command, args []string) {

		connInfo := getConnectionInfo(cmd)

		sqlStr, _ := cmd.Flags().GetString("statement")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		count, _ := cmd.Flags().GetInt("count")

		jobs := make(chan string, count)
		results := make(chan result, count)

		db := make([]*sql.DB, 0)

		fmt.Println()
		bar := progressbar.Default(int64(concurrency), "Building connections")
		for i := 0; i < concurrency; i++ {
			if r, err := connInfo.getDB(); err == nil {
				db = append(db, r)
			} else {
				zap.L().Error("Error in creating a new connection", zap.Error(err))
			}

			bar.Add(1)
		}
		if len(db) == 0 {
			zap.L().Fatal("No healty connections created", zap.Error(errors.New("0 connections in pool")))
		}

		fmt.Printf("We have %d healty connections\n", len(db))

		for i := 0; i < len(db); i++ {
			go benchQueryWithBind(db[i], sqlStr, jobs, results)
		}

		bar2 := progressbar.Default(int64(count), "Preparing bind variables")
		binds := make([]string, 0)

		for i := 0; i < count; i++ {
			binds = append(binds, getRandomVector(512).toString())
			bar2.Add(1)
		}

		start := time.Now()

		for i := 0; i < count; i++ {
			jobs <- binds[i]
		}

		close(jobs)

		success := 0
		failure := 0
		bar3 := progressbar.Default(int64(count), "Number of executions")
		elapsed := make([]float64, count)
		for i := 0; i < count; i++ {
			res := <-results

			if res.status {
				success++
			} else {
				failure++
			}
			elapsed[i] = float64(res.elapsed)
			bar3.Add(1)
		}

		colorstring.Printf("\nTotal number of request: [green]%d[reset]/[red]%d[reset] in %s\n", success, failure, time.Since(start))

		responseTimeStatistics(elapsed)

	},
}

func benchQueryWithBind(db *sql.DB, sqlStr string, binds <-chan string, results chan<- result) {

	defer db.Close()

	if statement, err := db.Prepare(sqlStr); err != nil {
		zap.L().Fatal("Error in creating prepared statement", zap.String("sql", sqlStr), zap.Error(err))
	} else {

		for bind := range binds {

			start := time.Now()

			if _, err := statement.Exec(bind); err == nil {
				results <- result{status: true, elapsed: time.Since(start).Microseconds()}
			} else {
				zap.L().Warn("Error in sql execution", zap.Error(err))
				results <- result{status: false, elapsed: time.Since(start).Microseconds()}
			}

		}
	}

}
