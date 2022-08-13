package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	BASE_URL       = `https://qiye.aliyun.com/alimail/ajax/mail/queryMailList.txt?showFrom=0&fragment=1&offset=%d&length=%d&curIncrementId=0&forceReturnData=1&query={"folderIds":["1"]}&_csrf_token_=%s`
	TIME_FORMATTER = "2006-01-02 15:04"
	SUBJECT        = "日报"
)

var (
	length int
	token  string
	cookie string
)

type Updater struct {
	Length int    `json:"length" binding:"required"`
	Token  string `json:"token" binding:"required"`
	Cookie string `json:"cookie" binding:"required"`
}

type AliMailResponse struct {
	AliMailData []AliMailData `json:"dataList"`
}

type AliMailData struct {
	Subject   string `json:"subject"`
	TimeStamp int64  `json:"timestamp"`
	Body      string `json:"encSummary"`
}

type Result struct {
	SendAt  string `json:"sendAt"`
	Content string `json:"content"`
}

type ResultResponse struct {
	Items []Result `json:"items"`
	Total int      `json:"total"`
}

func postOnce(offset int) (*AliMailResponse, error) {
	url := fmt.Sprintf(BASE_URL, offset, 1, token)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := AliMailResponse{}
	err = json.Unmarshal(bytes, &result)
	return &result, err
}

func getDailyMailForSingleWeek() ([]AliMailData, error) {
	result := []AliMailData{}
	offset := 0
	for {
		resp, err := postOnce(offset)
		if err != nil {
			return nil, err
		}
		if resp.AliMailData[0].Subject == SUBJECT {
			result = append(result, resp.AliMailData[0])
		}
		offset++
		if len(result) == length {
			break
		}
	}
	return result, nil
}

func SearchMails(ctx *gin.Context) {
	updater := Updater{}
	err := ctx.ShouldBindBodyWith(&updater, binding.JSON)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	token = updater.Token
	length = updater.Length
	cookie = updater.Cookie

	datas, err := getDailyMailForSingleWeek()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	results := []Result{}
	for _, data := range datas {
		results = append(results, Result{
			SendAt:  TransMillTimestampToTime(data.TimeStamp).Format(TIME_FORMATTER),
			Content: ParseAndConcatSingleDayMail(data.Body),
		})
	}
	ctx.JSON(http.StatusOK, ResultResponse{
		Total: len(results),
		Items: results,
	})
}

func main() {
	e := gin.Default()
	e.POST("/mails/search", SearchMails)
	e.Run(":8080")
}
