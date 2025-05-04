package localization

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var defaultLang = language.Russian

//go:embed locales/*.toml
var localeFS embed.FS

type Localize struct {
	bundle *i18n.Bundle
}

func NewLocalize() (*Localize, func(), error) {
	bundle := i18n.NewBundle(defaultLang)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if _, err := bundle.LoadMessageFileFS(localeFS, "locales/ru.toml"); err != nil {
		return nil, func() {}, err
	}

	return &Localize{
		bundle: bundle,
	}, func() {}, nil
}

func (l *Localize) GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(l.bundle, lang, defaultLang.String())
}
