# Localization

Website content and URLs can change based on the requested locale.

Musings source code is preconfigured for Arabic, English, and Mongolian.
These three demonstrate locales with drastically different layout needs.

## Configuration file

To enable localization, configure the locales in "data/config.json" in the `locales` array.
A single entry with an "und" tag can be used to localize Musings without enabling route localization.
For Arabic and English, your configuration file might look like this:

```json
{
  "locales": [
    {
      "tag": "ar",
      "date_format": "02/01/2006",
      "direction": "rtl",
      "digits": "٠١٢٣٤٥٦٧٨٩",
      "native_name": "العربية",
      "writing_mode": "horizontal-tb"
    },
    {
      "tag":"en",
      "direction": "ltr",
      "native_name": "English",
      "writing_mode": "horizontal-tb"
    }
  ]
}
```

The locale properties and their defaults are:

* `tag` — required; BCP 47 tag
* `date_format` — date format string per Go specifications \[default: "2006-01-02"\]
* `digits` — digits 0 through 9 to use in numeric localization \[default: "0123456789"\]
* `direction` — writing direction, one of: "ltr", "rtl", "auto" \[default: "auto"\]
* `native_name` — locale name in its native language \[default per internal table\]
* `writing_mode` — writing mode, one of "horizontal-tb", "vertical-rl", "vertical-lr" \[default: "horizontal-tb"\]


## Route localization

To enable route localization, populate `locales` in the configuration with anything other than "und".
The "und" locale can be used therein to localize Musings without enabling route localization.

Route localization forces URLs to begin with a valid configured locale.
If an unlocalized URL is requested, or if the URL's locale is not canonical,
the server responds with a redirect to a canonical localized path.
The locale used for the redirect is the detected best match for the website visitor.

Route localization applies to everything in the "content" directory tree except for "content/_shared/".

For example, assuming Arabic (ar) is configured for route localization:

 * `/hello` redirects to `/ar/hello`
 * `/_shared/wave.png` is not redirected.
 * `/_static/css/default.css` is not redirected.

### Organize content directory

With route localization enabled, the content directory must be organized by locale.
For the configuration example above, this is how that looks:

* `content/ar/`
* `content/en/`

Asset files that are shared between locales can be stored in a "_shared" subdirectory.
Shared files are served at "/_shared/", but relative paths can be used in markdown files.

* `content/_shared/logo.png`

In "content/en/example.md", this can be included with `![](/_shared/logo.png)` or `![](../_shared/logo.png)`.


## String localization

Developers who are [customizing](customizing.md) Musings have access to string localization in templates.

All view models include a `Translations` field that exactly mirrors the data in "data/translations.json" for the active locale.

For example, if the data file looks like this:

```json
{
  "en": {
    "Website logo": "Website logo",
    "Select language": "Select language"
  },
  "ko": {
    "Website logo": "웹사이트 로고",
    "Select language": "언어 선택"
  }
}
```

Then `{{.Translations.Get "Website logo"}}` will be either "Website logo" or "웹사이트 로고" depending on the active locale.


## Data Localization

Currently only date localization is supported.

### Date localization

Dates are localized per the locale configuration with `{{$.FormatDate <date>}}` where `<date>` is a `time.Time` variable or value.

For example, with the Arabic configuration earlier, `$.FormatDate` returns `١٦/٠٦/٢٠٢٦` when passed the date `2026-06-16`.


## What locales are supported?

Musings supports over 50 common locales.
(To be exact, it supports 79 locales as of June 2026.)
If a configured locale is not supported, Musings is designed to fail to start and output the error.

Avoid using "zh" directly even though it is supported.
The "zh" locale does not distinguish between Simplified Chinese ("zh-Hans") and Traditional Chinese ("zh-Hant").
