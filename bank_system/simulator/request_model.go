package simulator

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	Message      struct {
		Content string `json:"content"`
	} `json:"message"`
}

type Usage struct {
	PromptTokens    int `json:"prompt_tokens"`
	CompletionTasks int `json:"completion_tasks"`
	TotalTokens     int `json:"total_tokens"`
}

type ResponseBody struct {
	Id       string   `json:"id"`
	Provider string   `json:"provider"`
	Object   string   `json:"object"`
	Created  int64    `json:"created"`
	Choices  []Choice `json:"choices"`
	Usage    Usage    `json:"usage"`
}
