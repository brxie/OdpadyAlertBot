package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Telegram struct {
	token  string
	client http.Client
}

type chatMessageBody struct {
	ChatID              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

type RequestError struct {
	Msg          string
	ResponseCode int
	ResponseBody []byte
}

func (e *RequestError) Error() string {
	return e.Msg
}

func NewTelegramClient(token string, client *http.Client) (*Telegram, error) {
	if client == nil {
		client = &http.Client{}
	}
	if token == "" {
		return nil, errors.New("Token can not be empty")
	}
	return &Telegram{
		token:  token,
		client: *client,
	}, nil

}

func (t *Telegram) SendMessage(chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	message := chatMessageBody{
		ChatID: chatID,
		Text:   text,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		reqErr := &RequestError{
			ResponseCode: resp.StatusCode,
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			reqErr.Msg = "unexpected status code; failed to read response body"
			return reqErr
		}

		reqErr.Msg = fmt.Sprintf("unexpected status code %d", resp.StatusCode)
		reqErr.ResponseBody = respBytes

		return reqErr
	}

	return nil
}

func (t *Telegram) SendPhoto(chatID, text string, photo []byte) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%s&caption=%s",
		t.token, chatID, url.QueryEscape(text))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("photo", "logo.png")
	if err != nil {
		return err
	}

	_, err = io.WriteString(fw, string(photo))
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		reqErr := &RequestError{
			ResponseCode: resp.StatusCode,
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			reqErr.Msg = "unexpected status code; failed to read response body"
			return reqErr
		}

		reqErr.Msg = fmt.Sprintf("unexpected status code %d", resp.StatusCode)
		reqErr.ResponseBody = respBytes

		return reqErr
	}

	return nil
}
