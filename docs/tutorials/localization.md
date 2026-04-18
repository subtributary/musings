# Localization

Many files—including content, templates, and assets—can be created with localized variants.
When a visitor accesses the website,
Musings will use the most appropriate localized versions of files to serve them the website.

## Enabling localization

Localization is disabled by default.

Enabling locales will prefix all routes with the active locale.
For example, "/hello" becomes "/en/hello".
Requests to non-localized URLs will be redirected to their localized equivalents, 
so existing links will continue to work.

To enable locales, set them via the `MUSINGS_LOCALES` environment variable or the `--locales` parameter.
The locales should be comma-separated strings in BCP 47 format.
For example:

```shell
MUSINGS_LOCALES=da,en-GB,en
```

Language extensions, quality values, and other attributes are not supported in the configuration.
These may work in the current version but are not officially supported and may break in future releases.

## Locale resolution

The active locale is determined in the following order:

 1. The locale specified in the URL path, if it is enabled
 2. The Accept-Language header sent by the client
 3. The first configured locale 

If the URL path begins with a locale that is invalid or not enabled,
then it is treated as part of the path rather than a locale,
and the visitor is redirected to the appropriate localized URL.

Once determined, this locale is used for all localization behavior,
including file selection and template rendering.

## File localization

Files read by Musings as part of serving the website—such as content, templates, and static assets—can be localized.

To localize a file, put the locale as a second extension in front of the file extension.
The locale must correspond to an enabled locale (or its parent).
Files with locales that are not enabled are ignored during localization.

Musings will use the most appropriate localized version of any file.
For example, a French visitor requesting "hello.md" will get "hello.fr.md".

If a specific localized version of a file does not exist,
then Musings will fall back to its closest parent locale 
before ending at the unlocalized file.
For example, if "hello.md" is requested for the "fr-ca" locale, 
then Musings will look for "hello.fr-ca.md", "hello.fr.md", then "hello.md".

### Examples

For example, suppose that these files exist:

* /content/hello.md
* /content/hello.fr.md
* /content/hello.fr-ca.md

If "en,fr,fr-ca" are the enabled locales and "hello.md" is requested:

* Canadian French visitors see "hello.fr-ca.md".
* Other French visitors see "hello.fr.md".
* English visitors see "hello.md".

If "fr,fr-ca" are the enabled locales and "hello.md" is requested:

* Canadian French and French visitors are unchanged.
* English visitors see "hello.fr.md" because English is not supported.

If "en" is the only enabled locale and "hello.md" is requested:

* All visitors will see "hello.md"

In the last example,
the other two files are treated as independent files 
and will be served at their own routes if explicitly requested.
Remember that locales that are not enabled do not count for file localization.

## String localization

Developers who are [customizing](customization.md) Musings have access to limited string localization.

The view models passed to templates contain preset translations of some useful strings for the matching locale.
These can be accessed via `{.Translations}` which has the following structure:

```go
type Strings struct {
    SelectLanguage string
    WebsiteLogo    string
}
```

It is not possible to extend the translations table.
This weakness is intentional to keep localization simple and focused on file localizations.

## What locales are supported?

Musings supports a predefined set of common locales 
and can map many related locales to them.

It is recommended to avoid using "zh" directly, 
because it does not distinguish between Simplified Chinese ("zh-Hans") and Traditional Chinese ("zh-Hant").
It can still be used for files shared between "zh-Hans" and "zh-Hant" even if it is not enabled.

If Musings is started with an unsupported locale, then it will fail with an error.
