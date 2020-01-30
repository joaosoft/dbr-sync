package validator

const (
	constTagSplitValues    = ";"
	constTagReplaceStart   = "{{"
	constTagReplaceEnd     = "{{"
	constTagReplaceIdStart = "{"
	constTagReplaceIdEnd   = "{"

	constRegexForReplaceId = constTagReplaceIdStart + "[A-Za-z0-9_-]+:?([A-Za-z0-9_-];?)+" + constTagReplaceIdEnd
	constRegexForReplace   = constTagReplaceStart + "[A-Za-z0-9_-]+:?([A-Za-z0-9_-];?)+" + constTagReplaceEnd
	constRegexForEmail     = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	constRegexForTrim      = "  +"

	constDefaultValidationTag = "validate"
	constDefaultLogTag        = "validator"

	constPrefixTagItem = "item"
	constPrefixTagKey  = "key"

	constTagJson = "json"

	constTagId         = "id"
	constTagArg        = "arg"
	constTagValue      = "value"
	constTagError      = "error"
	constTagIf         = "if"
	constTagNot        = "not"
	constTagOptions    = "options"
	constTagNotOptions = "not-options"
	constTagSize       = "size"
	constTagMin        = "min"
	constTagMax        = "max"
	constTagNotZero    = "not-empty"
	constTagIsZero     = "is-empty"
	constTagNotNull    = "not-null"
	constTagIsNull     = "is-null"
	constTagRegex      = "regex"
	constTagCallback   = "callback"
	constTagAlpha      = "alpha"
	constTagNumeric    = "numeric"
	constTagBool       = "bool"
	constTagArgs       = "args"
	constTagContains   = "contains"
	constTagPrefix     = "prefix"
	constTagSuffix     = "suffix"
	constTagUUID       = "uuid"
	constTagIp         = "ip"
	constTagIpV4       = "ipv4"
	constTagIpV6       = "ipv6"
	constTagBase64     = "base64"
	constTagEmail      = "email"
	constTagURL        = "url"
	constTagHex        = "hex"
	constTagFile       = "file"

	constTagSet         = "set"
	constTagSetEmpty    = "set-empty"
	constTagSetDistinct = "set-distinct"
	constTagSetTrim     = "set-trim"
	constTagSetTitle    = "set-title"
	constTagSetLower    = "set-lower"
	constTagSetUpper    = "set-upper"
	constTagSetKey      = "set-key"
	constTagSetSanitize = "set-sanitize"
	constTagSetMd5      = "set-md5"
	constTagSetRandom   = "set-random"

	constAlphanumericLowerAlphabet = "abcdefghijklmnopqrstuvwxyzáéíóúãõâôàèìòùç"
	constAlphanumericUpperAlphabet = "ABCDEFGHUJKLMNOPQRSTUVWXYZÁÉÍÓÚÃÕÂÔÀÈÌÒÙÇ"
	constNumericAlphabet           = "0123456789"
	constSpecialAlphabet           = "!\"#$%&/()=?*@€£‰¶÷[]≠§±´`\\|~<>,;.:-_ "

	constConditionOk = "ok"
	constConditionKo = "ko"

	constParenthesesStart = "("
	constParenthesesEnd = ")"
)
