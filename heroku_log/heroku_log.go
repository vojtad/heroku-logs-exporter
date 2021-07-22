package herokuLog

import (
	"strings"
)

type HerokuLog struct {
	AppName string

	Unk1   string
	Unk2   string
	Time   string
	Host   string
	Source string
	Dyno   string

	Line string

	lineValuesParsed bool
	lineValues       map[string]string
}

func ParseHerokuLog(appName string, line string) *HerokuLog {
	parts := strings.SplitN(line, " - ", 2)

	headerParts := strings.Split(parts[0], " ")

	return &HerokuLog{
		appName,
		headerParts[0],
		headerParts[1],
		headerParts[2],
		headerParts[3],
		headerParts[4],
		headerParts[5],
		parts[1],
		false,
		nil,
	}
}

func (l *HerokuLog) parseLineValues() {
	if l.lineValuesParsed {
		return
	}

	tokens := strings.Split(l.Line, " ")
	l.lineValues = make(map[string]string)
	for _, pair := range tokens {
		if strings.Contains(pair, "=") {
			parts := strings.Split(pair, "=")
			l.lineValues[parts[0]] = parts[1]
		}
	}

	l.lineValuesParsed = true
}

func (l *HerokuLog) Value(key string) (string, bool) {
	l.parseLineValues()

	value, ok := l.lineValues[key]
	return value, ok
}

func (l *HerokuLog) ValueOrUnknown(key string) string {
	if value, ok := l.Value(key); ok {
		return value
	}

	return "UNKNOWN"
}
