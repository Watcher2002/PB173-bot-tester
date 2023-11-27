package ELI5

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func GetWikiArticleExtract(query string) (string, error) {
	modifiedQuery := strings.ReplaceAll(query, " ", "+")
	resp, err := http.Get("https://en.wikipedia.org/api/rest_v1/page/summary/" + modifiedQuery)

	if err != nil {
		log.Error().Msg("Couldn't connect to wiki.")
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		log.Error().Msg("Couldn't decode message.")
		return "", err
	}

	return data["extract"].(string), nil
}
