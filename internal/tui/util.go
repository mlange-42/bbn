package tui

import (
	"fmt"
	"strings"
)

func ParseEvidence(evidence []string) (map[string]string, error) {
	ev := map[string]string{}
	for _, entry := range evidence {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("syntax error in evidence")
		}
		ev[parts[0]] = parts[1]
	}
	return ev, nil
}
