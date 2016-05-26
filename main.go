// logfmt2json accepts logfmt on stdin and reformats it as JSON on stdout.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kr/logfmt"
)

type entry map[string]interface{}

func (e *entry) HandleLogfmt(key, value []byte) error {
	if len(value) == 0 {
		return valueLessKey(key)
	}

	(*e)[string(key)] = string(value)
	return nil
}

// the string is just the key name
type valueLessKey string

func (vlk valueLessKey) Error() string {
	return fmt.Sprintf("value-less key %q", string(vlk))
}

func main() {
	var verbose bool
	verboseFlag(&verbose)

	flag.Parse()

	rdr := bufio.NewReader(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	end := false
	for {
		line, err := rdr.ReadBytes('\n')
		if err == io.EOF {
			end = true
		}

		if len(line) > 1 {
			e := make(entry)
			err := logfmt.Unmarshal(line, &e)
			if err != nil {
				if verbose {
					log.Printf("invalid logfmt (%s): %q", err, string(line[:len(line)-1]))
				}
			} else {
				if err := enc.Encode(e); err != nil {
					log.Printf("error encoding JSON: %s", err)
				}
			}
		}

		if end {
			break
		}

		if err != nil {
			log.Fatalf("error reading stdin: %s", err)
		}
	}
}

func verboseFlag(v *bool) {
	usage := "Complain to stderr on invalid logfmt input"
	flag.BoolVar(v, "v", false, usage)
	flag.BoolVar(v, "verbose", false, usage)
}
