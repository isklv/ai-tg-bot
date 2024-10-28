package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIModel string

const (
	GEMINI_FLASH AIModel = "gemini-1.5-flash"
)

type GeminiService struct {
	key   string
	model *genai.GenerativeModel
}

func NewGeminiService(key string) GeminiService {
	return GeminiService{
		key: key,
	}
}

func (gs *GeminiService) Create(ctx context.Context, m AIModel) error {

	client, err := genai.NewClient(ctx, option.WithAPIKey(gs.key))
	if err != nil {
		return err
	}
	gs.model = client.GenerativeModel(fmt.Sprint(m))

	return nil
}

func (gs *GeminiService) SendText(ctx context.Context, query string) (string, error) {
	if gs.model == nil {
		return "", fmt.Errorf("AIModel is nil")
	}

	r, err := gs.model.GenerateContent(ctx, genai.Text(query))
	if err != nil {
		return "", err
	}

	res := ""

	for _, cand := range r.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				res = res + fmt.Sprintf("%v\n", part)
			}
		}
	}

	return res, nil
}

func (gs *GeminiService) SendAudio(ctx context.Context, query string) (string, error) {
	if gs.model == nil {
		return "", fmt.Errorf("AIModel is nil")
	}

	r, err := gs.model.GenerateContent(ctx, genai.Text(query))
	if err != nil {
		return "", err
	}

	res := ""

	for _, cand := range r.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				res = res + fmt.Sprintf("%v\n", part)
			}
		}
	}

	return res, nil
}
