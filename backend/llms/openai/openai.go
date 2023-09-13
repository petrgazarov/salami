package openai

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

var Client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
