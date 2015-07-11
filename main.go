package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	sync    = flag.Bool("sync", false, "use syncronous cat")
	workers = flag.Int("workers", 4, "workers to spawn")
)

func init() {
	flag.Parse()
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if *sync {
		syncCat(in)
	} else {
		asyncCat(in, *workers)
	}
}

// lineProcessor implements methods for processing a line of text
type lineProcessor struct {
	work chan string
	done chan bool
}

// fakeWork fakes doing work
func fakeWork() {
	for i := 0; i < 2*int(time.Millisecond); i++ {
		i++
	}
}

// Work processes one line of input
func (l *lineProcessor) Work() {
	for w := range l.work {
		fakeWork()
		fmt.Print(w)
	}
	l.done <- true
}

// NewLineProcessor creates a new line processor with given parameters
func NewLineProcessor(in chan string, done chan bool) *lineProcessor {
	return &lineProcessor{in, done}
}

// asyncCat cats stdin asyncronously
func asyncCat(in *bufio.Reader, workerCount int) {
	work := make(chan string, *workers)
	done := make(chan bool)
	for i := 0; i < *workers; i++ {
		worker := NewLineProcessor(work, done)
		go worker.Work()
	}
	for {
		line, err := in.ReadString('\n')
		if line != "" {
			work <- line
		}
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			close(work)
			for i := 0; i < *workers; i++ {
				<-done
			}
			return
		}
	}
}

// syncCat cats stdin syncronously
func syncCat(in *bufio.Reader) {
	for {
		line, err := in.ReadString('\n')
		if line != "" {
			fakeWork()
			fmt.Print(line)
		}
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			return
		}
	}
}
