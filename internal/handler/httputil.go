package handler

import (
	"net/url"
)

func parseLanguage(query url.Values) string {
	if !query.Has("lang") {
		return "en"
	}

	return query.Get("lang")
}

func parseEncoder(query url.Values) string {
	if !query.Has("encoder") {
		return "none"
	}

	if query.Has("opus") {
		return "opus"
	}

	return query.Get("encoder")
}

func parseID(query url.Values) string {
	if !query.Has("id") {
		return ""
	}

	return query.Get("id")
}
