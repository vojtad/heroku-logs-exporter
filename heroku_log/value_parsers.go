package herokuLog

import (
	"strconv"
	"strings"
)

func ParseSimpleNumber(value string) float64 {
	number, _ := strconv.ParseFloat(value, 64)
	return number
}

func ParseNumberWithSuffix(value string, suffix string) float64 {
	number, _ := strconv.ParseFloat(strings.Replace(value, suffix, "", 1), 64)
	return number
}

func ParseNumberWithPagesSuffix(value string) float64 {
	return ParseNumberWithSuffix(value, "pages")
}

func ParseMillis(value string) float64 {
	return ParseNumberWithSuffix(value, "ms") / 1000.0
}

func ParseSize(value string) float64 {
	multiplier := 1
	raw_value := value

	if strings.HasSuffix(value, "GB") {
		multiplier = 1024 * 1024 * 1024
		raw_value = strings.Replace(value, "MB", "", 1)
	} else if strings.HasSuffix(value, "MB") {
		multiplier = 1024 * 1024
		raw_value = strings.Replace(value, "MB", "", 1)
	} else if strings.HasSuffix(value, "kB") {
		multiplier = 1024
		raw_value = strings.Replace(value, "kB", "", 1)
	} else if strings.HasSuffix(value, "bytes") {
		multiplier = 1
		raw_value = strings.Replace(value, "bytes", "", 1)
	}

	raw_number_value, _ := strconv.ParseFloat(raw_value, 64)
	bytes := raw_number_value * float64(multiplier)
	return bytes
}
