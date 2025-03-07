package simulator

import (
	"bank_system/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type LLModel struct {
	Model    string
	messages []Message
	URL      string
	APIKey   string
}

func NewLLModel(messages []Message) *LLModel {
	return &LLModel{
		Model:    viper.GetString("openrouter.model.llama_3-1_8b"),
		messages: messages,
		URL:      viper.GetString("openrouter.url"),
		APIKey:   viper.GetString("openrouter.api_key"),
	}
}

func (l *LLModel) GenerateContent() (string, error) {
	requestBody := RequestBody{
		Model:    l.Model,
		Messages: l.messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", utils.NewBankSystemError(utils.ErrRequest)
	}

	req, err := http.NewRequest("POST", l.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", utils.NewBankSystemError(utils.ErrRequest)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", utils.NewBankSystemError(utils.ErrRequest)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.NewBankSystemError(utils.ErrRequest)
	}

	var response ResponseBody
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", utils.NewBankSystemError(utils.ErrRequest)
	}

	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		return content, nil
	} else {
		return "", utils.NewBankSystemError(utils.ErrGenerateNoContent)
	}
}
