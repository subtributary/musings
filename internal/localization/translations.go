package localization

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Strings struct {
	SelectLanguage string
	WebsiteLogo    string
}

func InitTranslations() {
	for locale, values := range translations {
		mustSetString(locale, "SelectLanguage", values.SelectLanguage)
		mustSetString(locale, "WebsiteLogo", values.WebsiteLogo)
	}
}

func LoadFor(locale language.Tag) Strings {
	if locale == language.Und {
		locale = language.English
	}

	printer := message.NewPrinter(locale)
	return Strings{
		SelectLanguage: printer.Sprintf("SelectLanguage"),
		WebsiteLogo:    printer.Sprintf("WebsiteLogo"),
	}
}

func mustSetString(locale language.Tag, key string, value string) {
	if err := message.SetString(locale, key, value); err != nil {
		panic(err)
	}
}

// translations is a table of UI terms in various languages.
// The tone is formal and would suit a news website.
// The alphabets are the standard or most common for each language.
var translations = map[language.Tag]Strings{
	language.Afrikaans:          {WebsiteLogo: "Webwerflogo", SelectLanguage: "Kies taal"},
	language.Albanian:           {WebsiteLogo: "Logoja e sajtit", SelectLanguage: "Zgjidhni gjuhën"},
	language.Amharic:            {WebsiteLogo: "የድረ-ገጹ አርማ", SelectLanguage: "ቋንቋ ይምረጡ"},
	language.Arabic:             {WebsiteLogo: "شعار الموقع", SelectLanguage: "اختر اللغة"},
	language.Armenian:           {WebsiteLogo: "Կայքի լոգո", SelectLanguage: "Ընտրեք լեզուն"},
	language.Azerbaijani:        {WebsiteLogo: "Vebsaytın loqosu", SelectLanguage: "Dili seçin"},
	language.Bengali:            {WebsiteLogo: "ওয়েবসাইটের লোগো", SelectLanguage: "ভাষা নির্বাচন করুন"},
	language.Bulgarian:          {WebsiteLogo: "Лого на уебсайта", SelectLanguage: "Изберете език"},
	language.Burmese:            {WebsiteLogo: "ဝဘ်ဆိုက် လိုဂို", SelectLanguage: "ဘာသာစကားကို ရွေးချယ်ပါ"},
	language.Catalan:            {WebsiteLogo: "Logotip del lloc web", SelectLanguage: "Seleccioneu la llengua"},
	language.Chinese:            {WebsiteLogo: "网站徽标", SelectLanguage: "选择语言"},
	language.Croatian:           {WebsiteLogo: "Logotip web-mjesta", SelectLanguage: "Odaberite jezik"},
	language.Czech:              {WebsiteLogo: "Logo webu", SelectLanguage: "Vyberte jazyk"},
	language.Danish:             {WebsiteLogo: "Webstedets logo", SelectLanguage: "Vælg sprog"},
	language.Dutch:              {WebsiteLogo: "Logo van de website", SelectLanguage: "Selecteer taal"},
	language.English:            {WebsiteLogo: "Website logo", SelectLanguage: "Select language"},
	language.Estonian:           {WebsiteLogo: "Veebisaidi logo", SelectLanguage: "Valige keel"},
	language.Filipino:           {WebsiteLogo: "Logo ng website", SelectLanguage: "Piliin ang wika"},
	language.Finnish:            {WebsiteLogo: "Verkkosivuston logo", SelectLanguage: "Valitse kieli"},
	language.French:             {WebsiteLogo: "Logo du site web", SelectLanguage: "Sélectionnez la langue"},
	language.Georgian:           {WebsiteLogo: "ვებსაიტის ლოგო", SelectLanguage: "აირჩიეთ ენა"},
	language.German:             {WebsiteLogo: "Website-Logo", SelectLanguage: "Sprache auswählen"},
	language.Greek:              {WebsiteLogo: "Λογότυπο ιστοτόπου", SelectLanguage: "Επιλέξτε γλώσσα"},
	language.Gujarati:           {WebsiteLogo: "વેબસાઇટનું લોગો", SelectLanguage: "ભાષા પસંદ કરો"},
	language.Hebrew:             {WebsiteLogo: "לוגו האתר", SelectLanguage: "בחר שפה"},
	language.Hindi:              {WebsiteLogo: "वेबसाइट का लोगो", SelectLanguage: "भाषा चुनें"},
	language.Hungarian:          {WebsiteLogo: "A webhely logója", SelectLanguage: "Válasszon nyelvet"},
	language.Icelandic:          {WebsiteLogo: "Merki vefsvæðis", SelectLanguage: "Veldu tungumál"},
	language.Indonesian:         {WebsiteLogo: "Logo situs web", SelectLanguage: "Pilih bahasa"},
	language.Italian:            {WebsiteLogo: "Logo del sito web", SelectLanguage: "Seleziona la lingua"},
	language.Japanese:           {WebsiteLogo: "ウェブサイトのロゴ", SelectLanguage: "言語を選択"},
	language.Kannada:            {WebsiteLogo: "ಜಾಲತಾಣದ ಲೋಗೊ", SelectLanguage: "ಭಾಷೆಯನ್ನು ಆಯ್ಕೆಮಾಡಿ"},
	language.Kazakh:             {WebsiteLogo: "Веб-сайт логотипі", SelectLanguage: "Тілді таңдаңыз"},
	language.Khmer:              {WebsiteLogo: "រូបសញ្ញាគេហទំព័រ", SelectLanguage: "ជ្រើសរើសភាសា"},
	language.Kirghiz:            {WebsiteLogo: "Веб-сайттын логотиби", SelectLanguage: "Тилди тандаңыз"},
	language.Korean:             {WebsiteLogo: "웹사이트 로고", SelectLanguage: "언어 선택"},
	language.Lao:                {WebsiteLogo: "ໂລໂກ້ເວັບໄຊ", SelectLanguage: "ເລືອກພາສາ"},
	language.Latvian:            {WebsiteLogo: "Vietnes logotips", SelectLanguage: "Izvēlieties valodu"},
	language.Lithuanian:         {WebsiteLogo: "Svetainės logotipas", SelectLanguage: "Pasirinkite kalbą"},
	language.Macedonian:         {WebsiteLogo: "Лого на веб-локацијата", SelectLanguage: "Изберете јазик"},
	language.Malay:              {WebsiteLogo: "Logo laman web", SelectLanguage: "Pilih bahasa"},
	language.Malayalam:          {WebsiteLogo: "വെബ്‌സൈറ്റിന്റെ ലോഗോ", SelectLanguage: "ഭാഷ തിരഞ്ഞെടുക്കുക"},
	language.Marathi:            {WebsiteLogo: "वेबसाइटचा लोगो", SelectLanguage: "भाषा निवडा"},
	language.Mongolian:          {WebsiteLogo: "Вэбсайтын лого", SelectLanguage: "Хэл сонгоно уу"},
	language.Nepali:             {WebsiteLogo: "वेबसाइटको लोगो", SelectLanguage: "भाषा चयन गर्नुहोस्"},
	language.Norwegian:          {WebsiteLogo: "Nettstedlogo", SelectLanguage: "Velg språk"},
	language.Persian:            {WebsiteLogo: "نشان‌واره وب‌سایت", SelectLanguage: "زبان را انتخاب کنید"},
	language.Polish:             {WebsiteLogo: "Logo witryny", SelectLanguage: "Wybierz język"},
	language.Portuguese:         {WebsiteLogo: "Logótipo do sítio Web", SelectLanguage: "Selecione o idioma"},
	language.Punjabi:            {WebsiteLogo: "ਵੈੱਬਸਾਈਟ ਦਾ ਲੋਗੋ", SelectLanguage: "ਭਾਸ਼ਾ ਚੁਣੋ"},
	language.Romanian:           {WebsiteLogo: "Sigla site-ului", SelectLanguage: "Selectați limba"},
	language.Russian:            {WebsiteLogo: "Логотип сайта", SelectLanguage: "Выберите язык"},
	language.Serbian:            {WebsiteLogo: "логотип веб-сајта", SelectLanguage: "Изаберите језик"},
	language.SerbianLatin:       {WebsiteLogo: "logotip veb-sajta", SelectLanguage: "Izaberite jezik"},
	language.SimplifiedChinese:  {WebsiteLogo: "网站徽标", SelectLanguage: "选择语言"},
	language.Sinhala:            {WebsiteLogo: "වෙබ් අඩවියේ ලාංඡනය", SelectLanguage: "භාෂාව තෝරන්න"},
	language.Slovak:             {WebsiteLogo: "Logo webovej lokality", SelectLanguage: "Vyberte jazyk"},
	language.Slovenian:          {WebsiteLogo: "Logotip spletnega mesta", SelectLanguage: "Izberite jezik"},
	language.Spanish:            {WebsiteLogo: "Logotipo del sitio web", SelectLanguage: "Seleccione el idioma"},
	language.Swahili:            {WebsiteLogo: "Nembo ya tovuti", SelectLanguage: "Chagua lugha"},
	language.Swedish:            {WebsiteLogo: "Webbplatslogotyp", SelectLanguage: "Välj språk"},
	language.Tamil:              {WebsiteLogo: "இணையதள சின்னம்", SelectLanguage: "மொழியைத் தேர்ந்தெடுக்கவும்"},
	language.Telugu:             {WebsiteLogo: "వెబ్‌సైట్ లోగో", SelectLanguage: "భాషను ఎంచుకోండి"},
	language.Thai:               {WebsiteLogo: "โลโก้เว็บไซต์", SelectLanguage: "เลือกภาษา"},
	language.TraditionalChinese: {WebsiteLogo: "網站標誌", SelectLanguage: "選擇語言"},
	language.Turkish:            {WebsiteLogo: "Web sitesi logosu", SelectLanguage: "Dil seçin"},
	language.Ukrainian:          {WebsiteLogo: "Логотип сайту", SelectLanguage: "Виберіть мову"},
	language.Urdu:               {WebsiteLogo: "ویب سائٹ کا لوگو", SelectLanguage: "زبان منتخب کریں"},
	language.Uzbek:              {WebsiteLogo: "Veb-sayt logotipi", SelectLanguage: "Tilni tanlang"},
	language.Vietnamese:         {WebsiteLogo: "Biểu trưng trang web", SelectLanguage: "Chọn ngôn ngữ"},
	language.Zulu:               {WebsiteLogo: "Ilogo yewebhusayithi", SelectLanguage: "Khetha ulimi"},
}
