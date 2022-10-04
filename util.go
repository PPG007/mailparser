package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	TITEL_CELL                 = "A3"
	WORK_HOURS_CELL            = "E3"
	WORK_DAYS_CELL             = "F3"
	CREATED_AT_CELL            = "L3"
	UPDATED_AT_CELL            = "M3"
	GLOBAL_PROJECT_NUMBER_CELL = "C3"

	TEMPLATE_XLSX = "template.xlsx"

	DEFAULT_SHEET          = "数据"
	DEFAULT_PROJECT_NUMBER = "PN20220414122230056"
	DEFAULT_WORK_TIME      = 8

	CONTENT_DATE_CELL           = "V%d"
	CONTENT_WORK_TIME_CELL      = "W%d"
	CONTENT_WORK_DETAIL_CELL    = "X%d"
	CONTENT_PROJECT_NUMBER_CELL = "Y%d"

	CONTENT_START_ROW = 3
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

func initTemplateXLSX() (*excelize.File, error) {
	f, err := excelize.OpenFile(TEMPLATE_XLSX)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()
	now := time.Now().Format(TIME_FORMATTER_Y_M_D_H_M_S)
	err = f.SetCellStr(DEFAULT_SHEET, TITEL_CELL, fmt.Sprintf("Koston Zhuang发起的工作报告.表单：%s", now))
	if err != nil {
		return nil, err
	}
	err = f.SetCellStr(DEFAULT_SHEET, CREATED_AT_CELL, now)
	if err != nil {
		return nil, err
	}
	err = f.SetCellStr(DEFAULT_SHEET, UPDATED_AT_CELL, now)
	if err != nil {
		return nil, err
	}
	err = f.SetCellStr(DEFAULT_SHEET, GLOBAL_PROJECT_NUMBER_CELL, DEFAULT_PROJECT_NUMBER)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GenerateXLSX(results []Result) (string, error) {
	var (
		err error
		f   *excelize.File
	)
	f, err = initTemplateXLSX()
	if err != nil {
		return "", err
	}
	defer f.Close()
	sort.Slice(results, func(i, j int) bool {
		return results[i].SendAt < results[j].SendAt
	})
	currentROW := CONTENT_START_ROW
	length := len(results)
	err = f.SetCellInt(DEFAULT_SHEET, WORK_DAYS_CELL, length)
	if err != nil {
		return "", err
	}
	err = f.SetCellInt(DEFAULT_SHEET, WORK_HOURS_CELL, length*DEFAULT_WORK_TIME)
	if err != nil {
		return "", err
	}
	for i, result := range results {
		err = setSingleROW(f, result, currentROW)
		if err != nil {
			return "", err
		}
		currentROW++
		if i < length-2 {
			err = f.InsertRow(DEFAULT_SHEET, currentROW)
			if err != nil {
				return "", err
			}
		}
	}
	fileName := fmt.Sprintf("%s.xlsx", time.Now().Format(time.RFC3339))
	err = f.SaveAs(fmt.Sprintf("%s/%s", STATIC_DIR, fileName))
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func setSingleROW(f *excelize.File, result Result, currentROW int) error {
	var err error
	err = f.SetCellStr(DEFAULT_SHEET, fmt.Sprintf(CONTENT_DATE_CELL, currentROW), result.SendAt)
	if err != nil {
		return err
	}
	err = f.SetCellInt(DEFAULT_SHEET, fmt.Sprintf(CONTENT_WORK_TIME_CELL, currentROW), DEFAULT_WORK_TIME)
	if err != nil {
		return err
	}
	err = f.SetCellStr(DEFAULT_SHEET, fmt.Sprintf(CONTENT_WORK_DETAIL_CELL, currentROW), result.Content)
	if err != nil {
		return err
	}
	err = f.SetCellStr(DEFAULT_SHEET, fmt.Sprintf(CONTENT_PROJECT_NUMBER_CELL, currentROW), DEFAULT_PROJECT_NUMBER)
	return err
}
