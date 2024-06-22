package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
	"github.com/psanford/claude/vertex"
	"golang.org/x/oauth2/google"
)

var defaultSystemPrompt = `The assistant is Claude, created by Anthropic. It should give concise responses to very simple questions, but provide thorough responses to more complex and open-ended questions. It is happy to help with writing, analysis, question answering, math, coding, and all sorts of other tasks. It uses markdown for coding.`

var modelName = flag.String("model", claude.Claude3Haiku, fmt.Sprintf("Model name (%s,%s,%s)", claude.Claude3Haiku, claude.Claude3Sonnet, claude.Claude3Opus))
var streaming = flag.Bool("stream", true, "Use streaming response")
var systemPrompt = flag.String("system-prompt", defaultSystemPrompt, "System prompt to use")
var maxTokens = flag.Int("max-tokens", 256, "Max response tokens")
var projectID = flag.String("project-id", "", "GCP Project ID")
var debug = flag.Bool("debug", false, "show debug info")

func main() {
	ctx := context.Background()
	flag.Parse()

	region := os.Getenv("CLOUD_ML_REGION")
	if region == "" {
		log.Fatalf("No CLOUD_ML_REGION region environment variable set")
	}

	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("failed to find default credentials: %s", err)
	}

	if *projectID == "" {
		if cred.ProjectID != "" {
			*projectID = cred.ProjectID
		} else {
			log.Fatalf("-project-id is required")
		}
	}

	fmt.Fprintln(os.Stderr, "Enter prompt. Press ctrl-d to send to API")

	prompt, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	prompt = bytes.TrimSpace(prompt)
	if len(prompt) == 0 {
		log.Fatal("Empty prompt, aborting")
	}
	client := vertex.NewClient(vertex.WithProjectID(*projectID))

	if *streaming {
		streamingResponse(ctx, client, string(prompt))
	} else {
		completeResponse(ctx, client, string(prompt))
	}
}

func completeResponse(ctx context.Context, client clientiface.Client, prompt string) {
	req := claude.MessageRequest{
		Model:     *modelName,
		System:    *systemPrompt,
		MaxTokens: *maxTokens,
		Stream:    false,
		Messages: []claude.MessageTurn{
			{
				Role: "user",
				Content: []claude.TurnContent{
					claude.TextContent(prompt),
				},
			},
		},
	}

	respMeta, err := client.Message(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	resp := <-respMeta.Responses()

	if resp.Type != "message" {
		log.Fatalf("error: %s", resp.Data)
	}

	msg, ok := resp.Data.(*claude.MessageStart)
	if !ok {
		log.Fatalf("message not of type MessageStart: type %T", resp.Data)
	}

	if *debug {
		log.Printf("response: %+v", msg)
	}

	for _, content := range msg.Content {
		fmt.Print(content.TextContent())
	}
	fmt.Println()
}

func streamingResponse(ctx context.Context, client clientiface.Client, prompt string) {
	req := claude.MessageRequest{
		Model:     *modelName,
		System:    *systemPrompt,
		MaxTokens: *maxTokens,
		Stream:    true,
		Messages: []claude.MessageTurn{
			{
				Role: "user",
				Content: []claude.TurnContent{
					claude.TextContent(prompt),
				},
			},
		},
	}
	respMeta, err := client.Message(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	for resp := range respMeta.Responses() {
		if *debug {
			log.Printf("response: %s", resp.Type)
		}
		fmt.Print(resp.Data.Text())

		switch ev := resp.Data.(type) {
		case *claude.MessageStart:
			if *debug {
				log.Printf("MessageStart: %+v", ev)
			}
		case *claude.MessagePing:
			if *debug {
				log.Printf("MessagePing: %+v", ev)
			}
		case *claude.ContentBlockStart:
			if *debug {
				log.Printf("ContentBlockStart: %+v", ev)
			}
		case *claude.ContentBlockDelta:
			if *debug {
				log.Printf("ContentBlockDelta: %+v", ev)
			}
		case *claude.ContentBlockStop:
			if *debug {
				log.Printf("ContentBlockStop: %+v", ev)
			}
		case *claude.MessageDelta:
			if *debug {
				log.Printf("MessageDelta: %+v", ev)
			}
		case *claude.MessageStop:
			if *debug {
				log.Printf("MessageStop: %+v", ev)
			}
		case *claude.ClaudeError:
			log.Fatalf("Error from API: %s", ev)
		case *claude.ClientError:
			log.Fatalf("Client side error: %s", ev)
		case error:
			log.Fatalf("Generic error: %s", ev)
		default:
			log.Fatalf("Unexpected message type: %T %+v", ev, ev)
		}
	}

	fmt.Println()
}
