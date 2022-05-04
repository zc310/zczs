package api

import "regexp"

var sbre *regexp.Regexp

func init() {
	sbre = regexp.MustCompile(`[?|&]id_spm=\w{24}`)
}
func RemoveIdSpm(s string) string {
	return sbre.ReplaceAllString(s, "")
}
