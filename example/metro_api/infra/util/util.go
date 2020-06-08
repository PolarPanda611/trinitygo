package util

// GetCodeUtil get code util
func GetCodeUtil(prefix, format, seq string, length int) string {
	codesequence := format + seq
	code := prefix + codesequence[len(codesequence)-length:]
	return code
}
