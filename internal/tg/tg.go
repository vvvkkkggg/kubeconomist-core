package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SendMessageReq struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type apiResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Bot struct {
	token  string
	chatID int64
}

func New(token string, chatID int64) *Bot {
	return &Bot{
		token:  token,
		chatID: chatID,
	}
}

func (b *Bot) SendMessage(ctx context.Context, text string) error {
	//reserved := []string{"_", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "/", "|", "{", "}", ".", "!"}

	reserved := []string{}

	for _, s := range reserved {
		text = strings.ReplaceAll(text, s, "\\"+s)
	}

	return b.sendMessage(b.chatID, text)
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	msg := SendMessageReq{
		ChatID: chatID,
		Text:   text,
	}

	//msg.ParseMode = "MarkdownV2"

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to json marshall: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http post error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is not 200: %d, body: %s", resp.StatusCode, body)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w, body: %s", err, body)
	}

	if !apiResp.Ok {
		return fmt.Errorf("telegram api error: %s (%d)", apiResp.Description, apiResp.ErrorCode)
	}

	return nil
}
