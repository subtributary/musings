# Localization

Website content can be change based on the requested locale.

Musings source code is preconfigured for Arabic, English, and Mongolian.
These three demonstrate locales with drastically different layout needs.

## How it works

If localization is enabled, the first segment of the path must be the locale.
This applies to both the HTTP request and the file within the content directory.
If the request path is not localized and canonical,
the server responds with a redirect to a canonical localized path.

The locale is decided from these sources in order:

1. The first segment of the path;
2. The `Accept-Language` request header compared against configured locales;
3. The first configured locale

## Enabling localization

### Configure locales:

To enable localization, configure the locales in "data/config.json" in the `locales` array.
For Arabic and English, your configuration file might look like this:

```json
{
  "locales": [
    {
      "tag":"en",
      "native_name": "English",
      "direction": "ltr",
      "writing_mode": "horizontal-tb"
    },
    {
      "tag": "ar",
      "native_name": "العربية",
      "direction": "rtl",
      "writing_mode": "horizontal-tb"
    }
  ]
}
```

In this example configuration, requests to "example.com" will be redirected to "example.com/ar/" or "example.com/en/".

The locale properties and their defaults are:

* `tag` — required; BCP 47 tag
* `native_name` — locale name in its native language \[default per internal table\]
* `direction` — writing direction, one of: "ltr", "rtl", "auto" \[default: "auto"\]
* `writing_mode` — writing mode, one of "horizontal-tb", "vertical-rl", "vertical-lr" \[default: "horizontal-tb"\]

If the configuration is invalid, Musings will fail to start and show a descriptive error.

### Organize content directory:

After enabling localization, the content directory needs to be divided by the configured locales.
Otherwise, no content will be served.
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


## What locales are supported?

Musings supports over 50 common locales.
(To be exact, it supports 79 locales as of June 2026.)
If a configured locale is not supported, Musings is designed to fail to start and output the error.

Avoid using "zh" directly even though it is supported.
The "zh" locale does not distinguish between Simplified Chinese ("zh-Hans") and Traditional Chinese ("zh-Hant").
