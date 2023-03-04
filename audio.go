package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

const (
	TranscriptionModelWhisper1 = "whisper-1"

	ResponseFormatJSON        = "json"
	ResponseFormatText        = "text"
	ResponseFormatSRT         = "srt"
	ResponseFormatVerboseJSON = "verbose_json"
	ResponseFormatVTT         = "vtt"
)

// Transcription returns a Transcription of an image. Convenience methods exist on
// images already returned from Client calls to easily vary those images.
func (c Client) Transcription(ctx context.Context, v TranscriptionReq) (*TranscriptionRes, error) {
	return c.audio(ctx, v, "transcriptions")
}

// Translation uses the OpenAI Translation endpoints to translate audio to english.
// When using this endpoint, the Language parameter of the Req struct is always ignored.
func (c Client) Translation(ctx context.Context, v TranscriptionReq) (*TranscriptionRes, error) {
	v.Language = nil
	return c.audio(ctx, v, "translations")
}

func (c Client) audio(ctx context.Context, v TranscriptionReq, endpoint string) (*TranscriptionRes, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	image, err := w.CreateFormFile("file", "file.mp3")
	if err != nil {
		return nil, fmt.Errorf("error creating audio multipart writer: %w", err)
	}
	io.Copy(image, bytes.NewReader(v.File))

	n, err := w.CreateFormField("model")
	if err != nil {
		return nil, fmt.Errorf("error creating audio multipart writer model: %w", err)
	}
	fmt.Fprint(n, v.Model)

	if v.Prompt != nil {
		n, err := w.CreateFormField("prompt")
		if err != nil {
			return nil, fmt.Errorf("error creating audio multipart writer prompt: %w", err)
		}
		fmt.Fprint(n, *v.Prompt)
	}
	if v.Temperature != nil {
		n, err := w.CreateFormField("temperature")
		if err != nil {
			return nil, fmt.Errorf("error creating audio multipart writer Temperature: %w", err)
		}
		fmt.Fprint(n, *v.Temperature)
	}
	if v.ResponseFormat != nil {
		n, err := w.CreateFormField("response_format")
		if err != nil {
			return nil, fmt.Errorf("error creating audio multipart writer ResponseFormat: %w", err)
		}
		fmt.Fprint(n, *v.ResponseFormat)
	}
	if v.Language != nil {
		n, err := w.CreateFormField("language")
		if err != nil {
			return nil, fmt.Errorf("error creating audio multipart writer Language: %w", err)
		}
		fmt.Fprint(n, *v.Language)
	}
	w.Close()

	var res TranscriptionRes
	err = c.c.R().
		Post("audio/%s", endpoint).
		SetHeader("Content-Type", "multipart/form-data; boundary="+w.Boundary()).
		WithBody(body).
		Do(ctx).
		JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// TranscriptionReq hold the data needed for image variation.
type TranscriptionReq struct {
	File           []byte
	Model          string
	Prompt         *string
	ResponseFormat *string
	Temperature    *float64
	Language       *string
}

type TranscriptionRes struct {
	Text string
}
