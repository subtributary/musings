package localization

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type Store interface {
	For(locale Locale) (Translations, error)
}

// NewStore creates a new translations store.
// If `live` is true, the translations are reloaded on every request.
func NewStore(dataDir string, live bool) (Store, error) {
	filename := filepath.Join(dataDir, "translations.json")
	if live {
		return LiveStore{filename: filename}, nil
	}
	return loadCachedStore(filename)
}

// CachedStore keeps and uses an in-memory copy of the translations.
type CachedStore map[string]Translations

func loadCachedStore(filename string) (CachedStore, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var store CachedStore
	if err = json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s CachedStore) For(locale Locale) (Translations, error) {
	return s[locale.Tag], nil
}

// LiveStore reloads the translations every time they are requested.
type LiveStore struct {
	filename string
}

func (s LiveStore) For(locale Locale) (Translations, error) {
	store, err := loadCachedStore(s.filename)
	if err != nil {
		return nil, fmt.Errorf("load store: %w", err)
	}

	return store.For(locale)
}
