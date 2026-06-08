# Localization

Content can be served per locale.

If localization is enabled, then routes will be prefixed with the active locale.
If no locale is active for a website visitor, 
then the best match is chosen among the locales configured for the website, 
and they are redirected to that locale.

## Enabling localization

Localization is disabled by default.
To enable it, add the locales you want to support in "data/config.json" in the `locales` array.
For English and Korean, your configuration file might look like this:

```json
{
  "locales": ["en", "ko"]
}
```

In this configuration, requests to "example.com" will be redirected to "example.com/en/" or "example.com/ko/".


## Content localization

After enabling localization, the content directory needs to be split by the configured locales.

 * `content/en/`
 * `content/ko/`

Otherwise, no content will be served.


## Template localization

Template localization is not currently supported.
This will change before release or a workaround will be documented.


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

Musings supports many common locales and can map many related locales to them.
If Musings is started with an unsupported locale, then it will quickly fail with an error.

Avoid using "zh" directly even though it is supported.
The "zh" locale does not distinguish between Simplified Chinese ("zh-Hans") and Traditional Chinese ("zh-Hant").
