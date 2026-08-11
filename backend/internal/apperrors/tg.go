package apperrors

import "errors"

var ErrTelegramDeletePermanent = errors.New("telegram: message cannot be deleted")
