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

func parseCodec(query url.Values) string {
	if query.Has("opus") {
		return "opus"
	}

	return query.Get("codec")
}

func parseID(query url.Values) string {
	if !query.Has("id") {
		return ""
	}

	return query.Get("id")
}
