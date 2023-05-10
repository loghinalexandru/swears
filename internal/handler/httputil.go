package handler

import (
	"net/url"

	"github.com/loghinalexandru/swears/internal/encoding"
	"github.com/loghinalexandru/swears/internal/service"
)

func parseLanguage(query url.Values) string {
	if !query.Has("lang") {
		return "en"
	}

	return query.Get("lang")
}

func parseEncoder(query url.Values) service.Encoder {
	if !query.Has("encoder") {
		return nil
	}

	switch query.Get("encoder") {
	case "opus":
		return encoding.NewOpus()
	default:
		return nil
	}
}

func parseID(query url.Values) string {
	if !query.Has("id") {
		return ""
	}

	return query.Get("id")
}
