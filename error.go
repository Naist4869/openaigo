package openaigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type APIErrorType string

const (
	ErrorInsufficientQuota APIErrorType = "insufficient_quota"
	ErrorInvalidRequest    APIErrorType = "invalid_request_error"
)

type APIError struct {
	Message string       `json:"message"`
	Type    APIErrorType `json:"type"`
	Param   interface{}  `json:"param"` // TODO: typing
	Code    interface{}  `json:"code"`  // TODO: typing

	Status     string
	StatusCode int
}

func (err APIError) Error() string {
	return fmt.Sprintf("openai API error: %v: %v (param: %v, code: %v)", err.Type, err.Message, err.Param, err.Code)
}

func (client *Client) apiError(res *http.Response) error {
	errbody := struct {
		Error APIError `json:"error"`
	}{APIError{Status: res.Status, StatusCode: res.StatusCode}}
	if err := json.NewDecoder(res.Body).Decode(&errbody); err != nil {
		return fmt.Errorf("failed to decode error body: %v", err)
	}
	//
	if errbody.Error.StatusCode == http.StatusTooManyRequests {
		rateLimit := parseRateLimit(res)
		return errors.Join(errbody.Error, rateLimit)
	}

	return errbody.Error
}

type RateLimit struct {
	RemainingRequests string
	RemainingTokens   string
	ResetRequests     time.Duration
	ResetTokens       time.Duration
}

func (er RateLimit) Error() string {
	return fmt.Sprintf("openai rate limit exceeded, RemainingRequests: %s, RemainingTokens: %s, ResetRequests: %s, ResetTokens: %s", er.RemainingRequests, er.RemainingTokens, er.ResetRequests.String(), er.ResetTokens.String())
}
func parseRateLimit(resp *http.Response) RateLimit {
	resetRequest := resp.Header.Get("X-Ratelimit-Reset-Requests")
	parsedResetRequest, err := time.ParseDuration(resetRequest)
	if err != nil {
		return RateLimit{}
	}

	resetTokens := resp.Header.Get("X-Ratelimit-Reset-Tokens")
	parsedResetTokens, err := time.ParseDuration(resetTokens)
	if err != nil {
		return RateLimit{}
	}
	return RateLimit{
		RemainingRequests: resp.Header.Get("X-Ratelimit-Remaining-Requests"),
		RemainingTokens:   resp.Header.Get("X-Ratelimit-Remaining-Tokens"),
		ResetRequests:     parsedResetRequest,
		ResetTokens:       parsedResetTokens,
	}
}
