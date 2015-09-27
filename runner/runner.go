package runner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
)

//Benchmark is a container for a benchmark result
type Benchmark struct {
	Name       string
	Iterations int
	Allocs     int
	OpTime     time.Duration
}

//ResultSet container for benchmark results
type ResultSet struct {
	Package    string
	Benchmarks []Benchmark
	Pass       bool
	Time       time.Time
	CommitHash string
}

var (
	log *logging.Logger
)

func init() {
	log = logging.MustGetLogger("workbench")
}

//Start the test runner, bool channel for triggers from watcher
func Start(do chan bool, results chan ResultSet) {

	for {
		select {
		case <-do:
			stats, err := Run(".")
			if err != nil {
				continue
			}
			results <- stats
		}

	}
}

//Run the benchmark for the given path, parse the result
func Run(path string) (ResultSet, error) {
	log.Debug("Test run triggered")
	cmd := exec.Command("go", "test", fmt.Sprintf("-bench=%s", path), "-benchmem")
	output, err := cmd.Output()
	if err != nil {
		log.Error(fmt.Sprintf("error! : %s", err.Error()))
	}

	dir, _ := os.Getwd()
	goPath := fmt.Sprintf("%s%s", os.Getenv("GOPATH"), "/src/")
	packageName := dir[len(goPath):]

	hashCmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	longHash, err := hashCmd.Output()
	var shortHash string
	if err != nil {
		log.Error(fmt.Sprintf("error! : %s", err.Error()))
		return ResultSet{}, err
	}
	shortHash = string(longHash[:7])

	benches, pass := parse(output)

	var rs = ResultSet{
		Package:    packageName,
		Benchmarks: benches,
		Pass:       pass,
		Time:       time.Now(),
		CommitHash: shortHash,
	}
	return rs, nil
}

func parse(output []byte) ([]Benchmark, bool) {
	var benchmarks []Benchmark
	var lines [][]byte
	var pass bool
	rd := bufio.NewReader(bytes.NewReader(output))
	for i := 0; true; i++ {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(fmt.Sprintf("error! : %s", err.Error()))
			break
		}
		lines = append(lines, line)
	}
	for i, line := range lines {
		var words []string
		switch i {
		case 0:
			lineScanner := bufio.NewScanner(bytes.NewReader(line))
			lineScanner.Split(bufio.ScanWords)
			for lineScanner.Scan() {
				if lineScanner.Text() != "PASS" {
					log.Error("TESTS FAILED")
					pass = false
				} else {
					pass = true
				}
			}
		case len(lines) - 1:
			log.Debug("do lastline stuff")
		default:
			var bm Benchmark
			lineScanner := bufio.NewScanner(bytes.NewReader(line))
			lineScanner.Split(bufio.ScanWords)
			for lineScanner.Scan() {
				words = append(words, lineScanner.Text())
			}

			bm = Benchmark{}

			bm.Name = words[0]

			iters, err := strconv.ParseInt(words[1], 10, 64)
			if err != nil {
				log.Error(err.Error())
			}
			bm.Iterations = int(iters)

			opTime, err := time.ParseDuration(fmt.Sprintf("%s%s", words[2], strings.Split(words[3], "/")[0]))
			if err != nil {
				log.Error(err.Error())
			}
			bm.OpTime = opTime

			allocs, err := strconv.ParseInt(words[4], 10, 64)
			bm.Allocs = int(allocs)
			if err != nil {
				log.Error(err.Error())
			}
			benchmarks = append(benchmarks, bm)
		}
	}
	return benchmarks, pass
}
