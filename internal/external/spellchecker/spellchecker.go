package spellchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/blankspace9/notes-app/internal/domain/models"
)

type YandexSpellChecker struct {
	url string
}

func New(url string) *YandexSpellChecker {
	return &YandexSpellChecker{
		url: url,
	}
}

func (sc *YandexSpellChecker) CheckSpelling(text string) ([]models.SpellError, error) {
	const op = "external.YandexSpellChecker.CheckSpelling"

	resp, err := http.Get(fmt.Sprintf("%s?text=%s", sc.url, url.QueryEscape(text)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var spellErrors []models.SpellError
	d := json.NewDecoder(resp.Body)

	err = d.Decode(&spellErrors)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return spellErrors, nil
}
