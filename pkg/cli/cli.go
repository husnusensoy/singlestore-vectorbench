package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	tag       string //last git tag
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
	hostname  string // build host
	goV       string //go compiler version used for build
)

var (
	configFile string = "protop.yaml"
	searchPath string = "."
)

var rootCmd = &cobra.Command{
	Use:   "vectbech [sub]",
	Short: "CLI Benchmark for vector queries",
}

func Initialize(_tag, _sha1ver, _buildTime, _hostname, _goV string) {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	tag = _tag
	sha1ver = _sha1ver
	buildTime = _buildTime
	hostname = _hostname
	goV = _goV
}

func Execute() {

	//rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "protop.yaml", "Configuration yaml for scenario")

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(dbCmd)

	rootCmd.AddCommand(vectorCmd)

	vectorCmd.PersistentFlags().Int("dim", 512, "Vector dimension")
	vectorCmd.PersistentFlags().Int("n", 1000, "# of vectors generated")

	//singlestore -h <client-endpoint> -P <port> -u <database-user> -p<database-user-password> -e status
	dbCmd.PersistentFlags().StringSliceP("hostname", "h", []string{"localhost"}, "Singlestore aggregator hostnames")
	dbCmd.PersistentFlags().StringP("port", "P", "3306", "Singlestore port")
	dbCmd.PersistentFlags().StringP("username", "u", "root", "Singlestore user")
	dbCmd.PersistentFlags().StringP("password", "p", "", "Singlestore password")
	dbCmd.PersistentFlags().StringP("database", "D", "singlestore", "Singlestore database name")
	dbCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	dbCmd.AddCommand(pingCmd)

	pingCmd.PersistentFlags().IntP("concurrency", "C", 1, "Number of concurrent sessions")
	pingCmd.PersistentFlags().IntP("count", "c", 1, "Number of ping requests")

	dbCmd.AddCommand(executeCmd)

	executeCmd.PersistentFlags().String("statement", "select 1", "SQL statement")

	dbCmd.AddCommand(benchCmd)

	benchCmd.PersistentFlags().IntP("concurrency", "C", 1, "Number of concurrent sessions")
	benchCmd.PersistentFlags().IntP("count", "c", 1, "Number of ping requests")
	benchCmd.PersistentFlags().String("statement", "select 1", "SQL statement")

	benchCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	dbCmd.AddCommand(vbenchCmd)

	vbenchCmd.PersistentFlags().IntP("concurrency", "C", 1, "Number of concurrent sessions")
	vbenchCmd.PersistentFlags().IntP("count", "c", 1, "Number of ping requests")
	vbenchCmd.PersistentFlags().String("statement", "select 1", "SQL statement")

	vbenchCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	rootCmd.Execute()
}
