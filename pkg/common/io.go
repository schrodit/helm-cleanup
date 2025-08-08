package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintReleasesTable(releases []Release) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"NAMESPACE", "NAME"})
	for _, r := range releases {
		t.AppendRow(table.Row{r.Namespace, r.Name})
	}
	t.Render()
}

// Prompt for non-empty user input from STDIN
func InputPrompt(prompt string) string {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		input, err := r.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" || err != nil {
			return input
		}
	}
}
