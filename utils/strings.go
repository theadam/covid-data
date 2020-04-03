package utils

import "regexp"

func StripString(str string) string {
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    return reg.ReplaceAllString(str, "")
}
