package cli

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Vector struct {
	e []float64
}

func getRandomVector(n int) Vector {
	v := Vector{e: make([]float64, n)}

	for i := 0; i < n; i++ {
		v.e[i] = rand.Float64()*2 - 1
	}

	return v
}

func (v Vector) toString() string {
	vs := make([]string, len(v.e))
	for i := 0; i < len(v.e); i++ {
		vs[i] = strconv.FormatFloat(v.e[i], 'g', 6, 32)

	}

	s := strings.Join(vs, ",")

	return fmt.Sprintf("[%s]", s)
}

var vectorCmd = &cobra.Command{
	Use:   "vector",
	Short: "Do some vector things",
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := cmd.Flags().GetInt("n")
		dim, _ := cmd.Flags().GetInt("dim")
		start := time.Now()
		for i := 0; i < count; i++ {
			v := getRandomVector(dim)

			v.toString()
		}
		elapsed := time.Since(start)

		fmt.Printf("Generate %d %d-dim vectors in %s\n", count, dim, elapsed)

	},
}
