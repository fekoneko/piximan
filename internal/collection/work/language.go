package work

type Language uint8

const (
	LanguageJapanese Language = iota
	LanguageEnglish
	LanguageKorean
	LanguageChineseChina
	LanguageChineseTaiwan
	LanguageIndonesian
	LanguageDanish
	LanguageGerman
	LanguageSpanish
	LanguageSpanishLatinAmerica
	LanguageFilipino
	LanguageFrench
	LanguageCroatian
	LanguageItalian
	LanguageDutch
	LanguagePolish
	LanguagePortugueseBrazil
	LanguagePortuguesePortugal
	LanguageVietnamese
	LanguageTurkish
	LanguageRussian
	LanguageArabic
	LanguageThai
	LanguageOther

	LanguageJapaneseString            = "ja"
	LanguageEnglishString             = "en"
	LanguageKoreanString              = "ko"
	LanguageChineseChinaString        = "zh-cn"
	LanguageChineseTaiwanString       = "zh-tw"
	LanguageIndonesianString          = "id"
	LanguageDanishString              = "da"
	LanguageGermanString              = "de"
	LanguageSpanishString             = "es"
	LanguageSpanishLatinAmericaString = "es-419"
	LanguageFilipinoString            = "tl"
	LanguageFrenchString              = "fr"
	LanguageCroatianString            = "hr"
	LanguageItalianString             = "it"
	LanguageDutchString               = "nl"
	LanguagePolishString              = "pl"
	LanguagePortugueseBrazilString    = "pt-br"
	LanguagePortuguesePortugalString  = "pt-pt"
	LanguageVietnameseString          = "vi"
	LanguageTurkishString             = "tr"
	LanguageRussianString             = "ru"
	LanguageArabicString              = "ar"
	LanguageThaiString                = "th"
	LanguageOtherString               = "other"
)

// Unlike novel, artwork can have only Japanese or English language.
func ValidArtworkLanguageString(l string) bool {
	return l == LanguageJapaneseString || l == LanguageEnglishString
}

func LanguageFromString(l string) Language {
	switch l {
	case LanguageJapaneseString:
		return LanguageJapanese
	case LanguageEnglishString:
		return LanguageEnglish
	case LanguageKoreanString:
		return LanguageKorean
	case LanguageChineseChinaString:
		return LanguageChineseChina
	case LanguageChineseTaiwanString:
		return LanguageChineseTaiwan
	case LanguageIndonesianString:
		return LanguageIndonesian
	case LanguageDanishString:
		return LanguageDanish
	case LanguageGermanString:
		return LanguageGerman
	case LanguageSpanishString:
		return LanguageSpanish
	case LanguageSpanishLatinAmericaString:
		return LanguageSpanishLatinAmerica
	case LanguageFilipinoString:
		return LanguageFilipino
	case LanguageFrenchString:
		return LanguageFrench
	case LanguageCroatianString:
		return LanguageCroatian
	case LanguageItalianString:
		return LanguageItalian
	case LanguageDutchString:
		return LanguageDutch
	case LanguagePolishString:
		return LanguagePolish
	case LanguagePortugueseBrazilString:
		return LanguagePortugueseBrazil
	case LanguagePortuguesePortugalString:
		return LanguagePortuguesePortugal
	case LanguageVietnameseString:
		return LanguageVietnamese
	case LanguageTurkishString:
		return LanguageTurkish
	case LanguageRussianString:
		return LanguageRussian
	case LanguageArabicString:
		return LanguageArabic
	case LanguageThaiString:
		return LanguageThai
	default:
		return LanguageOther
	}
}

func (l Language) String() string {
	switch l {
	case LanguageJapanese:
		return LanguageJapaneseString
	case LanguageEnglish:
		return LanguageEnglishString
	case LanguageKorean:
		return LanguageKoreanString
	case LanguageChineseChina:
		return LanguageChineseChinaString
	case LanguageChineseTaiwan:
		return LanguageChineseTaiwanString
	case LanguageIndonesian:
		return LanguageIndonesianString
	case LanguageDanish:
		return LanguageDanishString
	case LanguageGerman:
		return LanguageGermanString
	case LanguageSpanish:
		return LanguageSpanishString
	case LanguageSpanishLatinAmerica:
		return LanguageSpanishLatinAmericaString
	case LanguageFilipino:
		return LanguageFilipinoString
	case LanguageFrench:
		return LanguageFrenchString
	case LanguageCroatian:
		return LanguageCroatianString
	case LanguageItalian:
		return LanguageItalianString
	case LanguageDutch:
		return LanguageDutchString
	case LanguagePolish:
		return LanguagePolishString
	case LanguagePortugueseBrazil:
		return LanguagePortugueseBrazilString
	case LanguagePortuguesePortugal:
		return LanguagePortuguesePortugalString
	case LanguageVietnamese:
		return LanguageVietnameseString
	case LanguageTurkish:
		return LanguageTurkishString
	case LanguageRussian:
		return LanguageRussianString
	case LanguageArabic:
		return LanguageArabicString
	case LanguageThai:
		return LanguageThaiString
	default:
		return LanguageOtherString
	}
}
