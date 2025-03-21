package parser

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiAPI struct {
	client *genai.Client
	system string
}

func SystemPrompt() string {
	return `You are a resume parser. You will be given a resume in PDF format. Your job is to extract the relevant information from the resume and return it in JSON format.
    Follow this JSON schema:
    {
  "other": "{"Hobbies":"","Languages":""}",
  "first_name": "first_name",
  "last_name": "last_name",
  "email": "email",
  "phone": "phone",
  "social": {
    "link_name1": "link_url1",
    "link_name2": "link_url2",
    "link_name3": "link_url3",
  },
  "summary": "summary in the resume",
  "skills": "comma separated list of skills",
  "work": "[{"id": "generate a random id","company":"company","title":"title","startDate":"start_date","endDate":"end_date","description":"description"}]",
  "education": "[{"id": "generate a random id","degree":"degree","institution":"institution","startDate":"start_date","endDate":"end_date"}]",
  "projects": "[{"id": "generate a random id","name":"name","description":"description"}]",
  "achievements": "[{"id": "generate a random id","name":"name","description":"description"}]",
}`
}

func SystemPromptWithSchema(schema string) string {
	return fmt.Sprintf(`You are a resume parser. You will be given a resume in PDF format. Your job is to extract the relevant information from the resume and return it in JSON format.
    Follow this JSON schema:
    %s
}`, schema)
}

func NewGeminiAPI(apiKey string, systemPrompt string) (*GeminiAPI, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiAPI{client: client, system: systemPrompt}, nil
}

func (g *GeminiAPI) Send(content []*genai.Part) (string, error) {
	result, err := g.client.Models.GenerateContent(context.Background(), "gemini-2.0-flash-exp", []*genai.Content{{Parts: content}}, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: g.system}}},
		ResponseMIMEType:  "application/json",
	})
	if err != nil {
		return "", err
	}

	resp := result.Candidates[0].Content.Parts[0].Text

	return resp, nil
}
