// logfmt2json accepts logfmt on stdin and reformats it as JSON on stdout.
package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/kr/logfmt"
)

type entry map[string]interface{}

func (e *entry) HandleLogfmt(key, value []byte) error {
	if len(value) == 0 {
		(*e)[string(key)] = true
	} else {
		(*e)[string(key)] = string(value)
	}

	return nil
}

func main() {
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
				log.Printf("invalid logfmt (%q) %s", string(line), err)
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
