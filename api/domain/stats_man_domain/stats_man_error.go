package stats_man_domain

import (
	"encoding/json"
	"net/http"
)

type StatsManErrorInterface interface {
	Status() int
	Message() string
}
type StatsManError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
}

func (w *StatsManError) Status() int {
	return w.Code
}
func (w *StatsManError) Message() string {
	return w.ErrorMessage
}

func NewStatsManError(statusCode int, message string) StatsManErrorInterface {
	return &StatsManError{
		Code:         statusCode,
		ErrorMessage: message,
	}
}
func NewBadRequestError(message string) StatsManErrorInterface {
	return &StatsManError{
		Code:         http.StatusBadRequest,
		ErrorMessage: message,
	}
}

func NewForbiddenError(message string) StatsManErrorInterface {
	return &StatsManError{
		Code:         http.StatusForbidden,
		ErrorMessage: message,
	}
}

func NewApiErrFromBytes(body []byte) (StatsManErrorInterface, error) {
	var result StatsManError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
