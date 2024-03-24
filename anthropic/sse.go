package anthropic

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type sseEvent struct {
	Name  string
	ID    string
	Data  string
	Retry int
	Error error
}

func decodeSSE(ctx context.Context, r io.Reader) chan sseEvent {
	ch := make(chan sseEvent)

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		var event sseEvent
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				if event.Data != "" {
					select {
					case ch <- event:
					case <-ctx.Done():
						return
					}
					event = sseEvent{}
				}
				continue
			}

			// comment
			if strings.HasPrefix(line, ":") {
				continue
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				event.Error = fmt.Errorf("parse error malformed line: %s", line)
				select {
				case ch <- event:
				case <-ctx.Done():
				}
				return
			}

			field, value := parts[0], parts[1]
			value = strings.TrimSpace(value)

			switch field {
			case "event":
				event.Name = value
			case "id":
				event.ID = value
			case "data":
				event.Data += value + "\n"
			case "retry":
			default:
				event.Error = fmt.Errorf("parse error unexpected field: %s (value: %s)", field, value)
				select {
				case ch <- event:
				case <-ctx.Done():
				}
				return
			}
		}

		if err := scanner.Err(); err != nil {
			event.Error = fmt.Errorf("scan error: %s", err)
			select {
			case ch <- event:
			case <-ctx.Done():
			}
			return
		}
	}()

	return ch
}
