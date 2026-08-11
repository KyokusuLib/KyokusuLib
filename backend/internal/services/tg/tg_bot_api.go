package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lanxre/kyokusulib/internal/apperrors"
	"github.com/lanxre/kyokusulib/internal/config"
	"github.com/lanxre/kyokusulib/internal/constants"
	"github.com/lanxre/kyokusulib/internal/models/dto"
)

type TgBotAPI struct {
	cfg     *config.Config
	log     *slog.Logger
	client  *http.Client
	enabled bool
}

func NewTgBotAPI(cfg *config.Config, log *slog.Logger) *TgBotAPI {
	enabled := cfg.TgBotToken != "" && cfg.TgChannelUsername != ""
	if !enabled {
		log.Warn("telegram bot api is disabled: TgBotToken or TgChannelUsername is empty")
	}

	return &TgBotAPI{
		cfg:     cfg,
		log:     log,
		client:  &http.Client{Timeout: constants.TgDeleteTimeout},
		enabled: enabled,
	}
}

func (a *TgBotAPI) DeleteFromTelegram(ctx context.Context, messageID int64) error {
	if !a.enabled {
		return nil
	}

	endpoint := fmt.Sprintf("%s/bot%s/deleteMessage", constants.TgAPIBaseURL, a.cfg.TgBotToken)

	form := url.Values{
		"chat_id":    {"@" + strings.TrimPrefix(a.cfg.TgChannelUsername, "@")},
		"message_id": {strconv.FormatInt(messageID, 10)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build telegram delete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("delete telegram post: %w", err)
	}
	defer resp.Body.Close()

	var body dto.DeleteMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("decode telegram delete response: %w", err)
	}

	if body.Ok {
		return nil
	}

	if body.ErrorCode == 400 && strings.Contains(body.Description, "message to delete not found") {
		a.log.Debug("telegram post already deleted", slog.Int64("message_id", messageID))
		return nil
	}

	if body.ErrorCode == 400 && strings.Contains(body.Description, "message can't be deleted") {
		a.log.Warn(
			"telegram message cannot be deleted",
			slog.Int64("message_id", messageID),
			slog.String("description", body.Description),
		)
		return apperrors.ErrTelegramDeletePermanent
	}

	a.log.Warn(
		"telegram rejected delete request",
		slog.Int64("message_id", messageID),
		slog.Int("status", resp.StatusCode),
		slog.Int("error_code", body.ErrorCode),
		slog.String("description", body.Description),
	)

	return fmt.Errorf("telegram delete failed: %s (code %d)", body.Description, body.ErrorCode)
}
