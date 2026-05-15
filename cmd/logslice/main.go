package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/slicer"
)

const timeLayout = "2006-01-02T15:04:05"

func main() {
	var (
		filePath  = flag.String("file", "", "Path to the log file (required)")
		startStr  = flag.String("start", "", "Start time (RFC3339 or 2006-01-02T15:04:05)")
		endStr    = flag.String("end", "", "End time (RFC3339 or 2006-01-02T15:04:05)")
		level     = flag.String("level", "", "Filter by log level (e.g. ERROR, WARN)")
		pattern   = flag.String("pattern", "", "Filter by regex pattern")
		outFile   = flag.String("out", "", "Output file path (default: stdout)")
	)
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "error: --file is required")
		flag.Usage()
		os.Exit(1)
	}

	opts := slicer.Options{
		FilePath: *filePath,
	}

	if *startStr != "" {
		t, err := parseTime(*startStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid --start time: %v\n", err)
			os.Exit(1)
		}
		opts.Start = t
	}

	if *endStr != "" {
		t, err := parseTime(*endStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid --end time: %v\n", err)
			os.Exit(1)
		}
		opts.End = t
	}

	if *level != "" || *pattern != "" {
		f, err := filter.New(filter.Options{
			Level:   *level,
			Pattern: *pattern,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid filter: %v\n", err)
			os.Exit(1)
		}
		opts.Filter = f
	}

	out := os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot open output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if err := slicer.Slice(opts, out); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation(timeLayout, s, time.Local)
}
