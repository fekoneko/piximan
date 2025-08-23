package work

type Language uint8

const (
	LanguageJapanese Language = iota
	LanguageEnglish
	LanguageDefault = LanguageJapanese

	LanguageJapaneseString = "ja"
	LanguageEnglishString  = "en"
	LanguageDefaultString  = LanguageJapaneseString
)

func ValidLanguageString(language string) bool {
	return language == LanguageJapaneseString || language == LanguageEnglishString
}

func LanguageFromString(language string) Language {
	switch language {
	case LanguageJapaneseString:
		return LanguageJapanese
	case LanguageEnglishString:
		return LanguageEnglish
	default:
		return LanguageDefault
	}
}

func (l Language) String() string {
	switch l {
	case LanguageJapanese:
		return LanguageJapaneseString
	case LanguageEnglish:
		return LanguageEnglishString
	default:
		return LanguageDefaultString
	}
}
