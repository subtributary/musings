package localization

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

type Translations map[string]string

// Get looks up a translation by key.
// If the translation is not found, the key is returned.
func (t Translations) Get(key string) string {
	if val, ok := t[key]; ok {
		return val
	}
	return key
}

type Store map[language.Tag]Translations

func LoadStore(dataDir string) (Store, error) {
	data, err := os.ReadFile(filepath.Join(dataDir, "translations.json"))
	if err != nil {
		return nil, err
	}

	var store Store
	if err = json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s Store) For(locale language.Tag) Translations {
	return s[locale]
}
