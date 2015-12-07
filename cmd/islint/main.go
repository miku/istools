package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/miku/islint"
	"github.com/miku/span/finc"
)

var (
	tests   = islint.DefaultTests
	Verbose = false
	version = "0.1.5"
	start   = time.Now()
)

func worker(queue chan [][]byte, out chan []islint.Issue, wg *sync.WaitGroup) {
	defer wg.Done()
	for batch := range queue {
		for _, b := range batch {
			var is finc.IntermediateSchema
			if err := json.Unmarshal(b, &is); err != nil {
				log.Fatal(err)
			}
			var issues []islint.Issue
			for _, t := range tests {
				err := t.TestRecord(is)
				if err != nil {
					if _, ok := err.(islint.Issue); !ok {
						log.Fatalf("invalid error type: %T", err)
					}
					issues = append(issues, err.(islint.Issue))
				}
			}
			out <- issues
		}
	}
}

type KindFrequencyMap map[islint.Kind]int

func (m KindFrequencyMap) MarshalJSON() ([]byte, error) {
	sm := make(map[string]int)
	for k, v := range m {
		sm[k.String()] = v
	}
	return json.Marshal(sm)
}

// Stats keeps basic stats on issues.
type Stats struct {
	IssueDistribution map[islint.Kind]int `json:"issues"`
	IssuesPerRecord   map[int]int         `json:"frequency"`
}

// MarshalJSON calculates a few extra metrics on the fly.
func (s Stats) MarshalJSON() ([]byte, error) {
	var total, damaged int
	errcount := make(map[string]int)

	for k, v := range s.IssuesPerRecord {
		total += v
		if k > 0 {
			damaged += v
		}
		errcount[strconv.Itoa(k)] = v
	}
	ratio := (100 / float64(total)) * float64(damaged)

	dist := make(map[string]int)
	for k, v := range s.IssueDistribution {
		dist[k.String()] = v
	}

	return json.Marshal(map[string]interface{}{
		"dist":     dist,
		"errcount": errcount,
		"total":    total,
		"damaged":  damaged,
		"ratio":    fmt.Sprintf("%0.3f", ratio),
		"start":    start,
		"elapsed":  time.Since(start).Seconds(),
	})
}

func writer(in chan []islint.Issue, done chan bool) {
	stats := Stats{IssueDistribution: make(map[islint.Kind]int), IssuesPerRecord: make(map[int]int)}
	var i int
	for issues := range in {
		stats.IssuesPerRecord[len(issues)]++
		for _, issue := range issues {
			stats.IssueDistribution[issue.Kind]++
			if Verbose {
				fmt.Println(issue.TSV())
			}
		}
		i++
		if i%1000000 == 0 {
			b, err := json.Marshal(stats)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b))
		}
	}
	b, err := json.Marshal(stats)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	done <- true
}

func main() {
	verbose := flag.Bool("verbose", false, "show every error")
	showVersion := flag.Bool("v", false, "show version and exit")
	listTests := flag.Bool("ls", false, "list tests")
	sample := flag.Float64("sample", 1.0, "ratio of records to test")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	Verbose = *verbose

	if *listTests {
		fmt.Println(`CurrencyInTitle
EndPageBeforeStartPage
EtAlAuthorName
ExcessivePunctuation
InvalidCollection
InvalidEndPage
InvalidStartPage
InvalidURL
KeyTooLong
NAInAuthorName
NoPublisher
PublicationDateTooEarly
PublicationDateTooLate
RepeatedSlash
RepeatedSubtitle
ShortAuthorName
SuspiciousPageCount
WhitespaceAuthor`)
		os.Exit(0)
	}

	var r io.Reader

	if flag.NArg() == 0 {
		r = os.Stdin
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		r = file
	}

	reader := bufio.NewReader(r)

	var i int
	var batch [][]byte
	var size = 40000

	queue := make(chan [][]byte)
	out := make(chan []islint.Issue)
	done := make(chan bool)

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(queue, out, &wg)
	}

	go writer(out, done)

	for {
		b, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if rand.Float64() > *sample {
			continue
		}
		if i == size {
			ba := make([][]byte, len(batch))
			copy(ba, batch)
			queue <- ba
			batch = batch[:0]
			i = 0
		}
		batch = append(batch, b)
		i++
	}

	ba := make([][]byte, len(batch))
	copy(ba, batch)
	queue <- ba
	batch = batch[:0]

	close(queue)
	wg.Wait()
	close(out)
	<-done
}
