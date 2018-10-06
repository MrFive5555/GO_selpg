package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	flag "github.com/spf13/pflag"
)

var (
	arg_start       int
	arg_end         int
	arg_use_F       bool
	arg_lineNumber  int
	arg_help        bool
	arg_fileName    string
	arg_destination string
)

func init() {
	flag.IntVarP(&arg_start, "start", "s", 0, "the start page number(include)")
	flag.IntVarP(&arg_end, "end", "e", 0, "the end page number(include)")
	flag.BoolVarP(&arg_use_F, "use_F", "f", false, "whether use the \\f")
	flag.IntVarP(&arg_lineNumber, "lineNumber", "l", 72, "the line number in one page")
	flag.BoolVarP(&arg_help, "help", "h", false, "show the help")
	flag.StringVarP(&arg_destination, "destination", "d", "", "the destination to print, it should be able to use in UNIX command 'lp -d")

	flag.Usage = func() {
		fmt.Printf("Usage: selpg -s Number -e Number [-f|-l Number] [fileName] \n")
		fmt.Printf("  [fileName]: used to specify the file to print. If Not be given, it will get input from stdin\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	arg_fileName = flag.Arg(0)
}

func main() {
	if arg_help {
		fmt.Printf("selpg version 1.0.0\n")
		flag.Usage()
		os.Exit(0)
	}

	if arg_start < 1 || arg_end < 1 {
		fmt.Printf("invalid -s or -e argument\n")
	}

	if arg_start > arg_end {
		println("Invalid argument: the start page can't bigger than end page\n")
		flag.Usage()
		os.Exit(1)
	}

	// confirm the place of input
	var reader *bufio.Reader
	if arg_fileName != "" {
		filePtr, err := os.OpenFile(arg_fileName, os.O_RDONLY, 0755)
		if err != nil {
			println(err.Error())
			os.Exit(2)
		}
		reader = bufio.NewReader(filePtr)
	} else {
		fmt.Fprintf(os.Stderr, "Please enter the content you want to print:\n")
		reader = bufio.NewReader(os.Stdin)
	}

	// confirm the place of output
	var writer *bufio.Writer
	_closeWriter := func() {

	}
	var wg sync.WaitGroup
	defer wg.Wait()
	if arg_destination != "" {
		r, w := io.Pipe()
		_closeWriter = func() {
			w.Close()
		}
		cmd := exec.Command("lp", "-d"+arg_destination)
		cmd.Stdin = r
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		writer = bufio.NewWriter(w)
		wg.Add(1)
		go func() {
			cmd.Run()
			wg.Done()
		}()
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	printPage := func(page_ctr int, page string) {
		if page_ctr >= arg_start && page_ctr <= arg_end {
			_, err := writer.WriteString(page)
			if err != nil {
				println(err.Error())
				os.Exit(2)
			}
		}
	}
	page_ctr := 1
	if arg_use_F {
		for {
			page, err := reader.ReadString('\f')
			if err == io.EOF || page_ctr > arg_end {
				break
			}
			printPage(page_ctr, page)
			page_ctr += 1
		}
	} else {
		for line_ctr := 1; ; {
			line, err := reader.ReadString('\n')
			if err == io.EOF || page_ctr > arg_end {
				break
			}
			printPage(page_ctr, line)
			if line_ctr >= arg_lineNumber {
				line_ctr = 1
				page_ctr += 1
			} else {
				line_ctr += 1
			}
		}
	}
	_closeWriter()
	if page_ctr < arg_start {
		fmt.Fprintf(os.Stderr, "the page number of the document is less than start[%d], nothing was printed\n", arg_end)
	}
	if page_ctr < arg_end {
		fmt.Fprintf(os.Stderr, "the page number of the document is less than end[%d]\n", arg_end)
	}
}
