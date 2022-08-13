package main

import (
	"strings"
	"time"
)

var (
	replaceMap = map[string]string{
		"&nbsp;": " ",
		"\n":     " ",
	}

	splitMap = map[string]int{
		"今日完成：":     1,
		"下一个工作日计划：": 0,
	}
)

func TransMillTimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

func ParseAndConcatSingleDayMail(content string) string {
	result := getAfterSplitContent(content)
	result = getAfterRepalceContent(result)
	temp := strings.Split(result, "- ")
	items := []string{}
	for _, s := range temp {
		trimed := trimAllSpaces(s)
		if trimed != "" {
			items = append(items, trimed)
		}
	}
	return strings.Join(items, "")
}

func getAfterSplitContent(content string) string {
	result := content
	for k, v := range splitMap {
		result = strings.Split(result, k)[v]
	}
	return result
}

func getAfterRepalceContent(content string) string {
	result := content
	for k, v := range replaceMap {
		result = strings.ReplaceAll(result, k, v)
	}
	return result
}

func trimAllSpaces(content string) string {
	temp := strings.TrimSpace(content)
	if temp == content {
		return temp
	}
	return trimAllSpaces(temp)
}
