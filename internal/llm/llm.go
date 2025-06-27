package llm

type Client struct {
	ApiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{ApiKey: apiKey}
}

// TODO: Реализовать функцию для отправки запроса к LLM (например, OpenAI)
