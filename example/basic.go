package main

import (
	"fmt"
	"os"

	"github.com/kybin/sequence"
)

func main() {
	fi, err := os.Stat("data")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !fi.IsDir() {
		fmt.Fprintln(os.Stderr, "data should be a directory")
		os.Exit(1)
	}

	dir, err := os.Open(fi.Name())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dir.Close()

	filenames, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	man := sequence.NewManager(sequence.DefaultSplitter, sequence.FmtSharp)
	for _, f := range filenames {
		man.Add(f)
	}

	fmt.Println(man)
	// Output:
	// another.####.exr 1-4 7-10
	// img.####.exr 1-3
}
