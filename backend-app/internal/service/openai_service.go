package service

import (
	"context"
	"log"
	"runtime"
	"strings"
	"text/template"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"github.com/philippgille/chromem-go"
)

var db *chromem.DB

func init() {

	var err error
	db, err = initChromemDB()
	if err != nil {
		log.Fatalf("Failed to initialize Chromem DB: %v", err)
	}

}

func initChromemDB() (*chromem.DB, error) {
	db, err := chromem.NewPersistentDB("./db", false)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetEmbeddings(value string) {
	log.Println("Getting embeddings for:", value)
	log.Println("OpenAI API Key:", openaiAPIKey)
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)

	embedding, err := client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Input:          openai.F[openai.EmbeddingNewParamsInputUnion](shared.UnionString(value)),
		Model:          openai.F(openai.EmbeddingModelTextEmbeddingAda002),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})

	if err != nil {
		panic(err.Error())
	}

	log.Println("Embedding:", embedding.JSON)

}

func LoadVectors() {
	log.Println("Loading Vector Embeddings")
	ctx := context.Background()
	c, err := db.GetOrCreateCollection("open-ai", nil, chromem.NewEmbeddingFuncOpenAI(openaiAPIKey, chromem.EmbeddingModelOpenAI2Ada))
	if err != nil {
		panic(err)
	}
	if c.Count() == 0 {
		err := c.AddDocuments(ctx, []chromem.Document{
			{
				ID:      "1",
				Content: "Malleable Iron Range Company was a company that existed from 1896 to 1985 and primarily produced kitchen ranges made of malleable iron but also produced a variety of other related products. The company's primary trademark was 'Monarch' and was colloquially often referred to as the Monarch Company or just Monarch.",
			},
			{
				ID:      "2",
				Content: "The American Motor Car Company was a short-lived company in the automotive industry founded in 1906 lasting until 1913. It was based in Indianapolis Indiana United States. The American Motor Car Company pioneered the underslung design.",
			},
		}, runtime.NumCPU())
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("not reading JSON lines because collection was loaded from persistent storage.")
	}
}

func ChatCompletionWithVector(query string) []SearchResult {
	log.Println("Searching for:", query)
	ctx := context.Background()

	c := db.GetCollection("open-ai", chromem.NewEmbeddingFuncOpenAI(openaiAPIKey, chromem.EmbeddingModelOpenAI2Ada))
	if c == nil {
		panic("collection not found")
	}
	// Query for similariy documents
	docRes, err := c.Query(ctx, query, 2, nil, nil)
	if err != nil {
		panic(err)
	}

	//Build Contexts for prompt preparation
	var contexts []string
	for _, res := range docRes {
		log.Printf("Document ID: %s, Content: %s", res.ID, res.Content)
		contexts = append(contexts, res.Content)
	}

	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)

	// Prepare system prompt
	sb := &strings.Builder{}
	err = systemPromptTpl.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}

	log.Printf("System prompt %s", sb.String())

	// Get completion
	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sb.String()),
			openai.UserMessage(query),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Open AI output %s", completion.Choices[0].Message.Content)
	results := []SearchResult{{Content: completion.Choices[0].Message.Content}}
	return results
}

var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
You are a helpful assistant with access to a knowlege base, tasked with answering questions about the world and its history, people, places and other things.

Answer the question in a very concise manner. Use an unbiased and journalistic tone. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
Answer the question solely based on the provided search results from the knowledge base. If the search results from the knowledge base are not relevant to the question at hand, just say that you don't know. Don't make anything up.

Anything between the following 'context' XML blocks is retrieved from the knowledge base, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}

Don't mention the knowledge base, context or search results in your answer.
`))

type SearchResult struct {
	Content string `json:"content"`
}
