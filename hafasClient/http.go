package hafasClient

import (
	"net/http"
	"net/url"

	conf "github.com/sa-/schedule/conf"
)

func headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + conf.Conf.VbbAPIKey,
		"Accept":        "application/json",
	}
}

func route(path string) string {
	base, _ := url.Parse(conf.Conf.VbbApiUrl)
	ref, _ := url.Parse(path)
	return base.ResolveReference(ref).String()
}

var client = &http.Client{}
