package work

type Language uint8

const (
	LanguageJapanese Language = iota
	LanguageEnglish
	LanguageChinese
	LanguageKorean
	LanguageDefault = LanguageJapanese

	LanguageJapaneseString = "ja"
	LanguageEnglishString  = "en"
	LanguageChineseString  = "zh"
	LanguageKoreanString   = "ko"
	LanguageDefaultString  = LanguageJapaneseString
)

func ValidLanguageString(language string) bool {
	return language == LanguageJapaneseString || language == LanguageEnglishString ||
		language == LanguageChineseString || language == LanguageKoreanString
}

func LanguageFromString(language string) Language {
	switch language {
	case LanguageJapaneseString:
		return LanguageJapanese
	case LanguageEnglishString:
		return LanguageEnglish
	case LanguageChineseString:
		return LanguageChinese
	case LanguageKoreanString:
		return LanguageKorean
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
	case LanguageChinese:
		return LanguageChineseString
	case LanguageKorean:
		return LanguageKoreanString
	default:
		return LanguageDefaultString
	}
}
