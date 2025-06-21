package agent

import (
	"fmt"
	"log"

	"github.com/hpcloud/tail"
)

func Watch(path string) {
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Error with tail file: %s", err)
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
