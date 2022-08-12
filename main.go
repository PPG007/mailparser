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
	BASE_URL = `https://qiye.aliyun.com/alimail/ajax/mail/queryMailList.txt?showFrom=0&fragment=1&offset=%d&length=%d&curIncrementId=0&forceReturnData=1&query={"folderIds":["1"]}&_csrf_token_=%s`

	SUBJECT = "日报"
)

var (
	length int
	token  string
	cookie string
)

type AliMailResponse struct {
	AliMailData []AliMailData `json:"dataList"`
}

type AliMailData struct {
	Subject   string `json:"subject"`
	TimeStamp int64  `json:"timestamp"`
	Body      string `json:"encSummary"`
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

func Init(ctx *gin.Context) {
	type Updater struct {
		Length int    `json:"length"`
		Token  string `json:"token"`
		Cookie string `json:"cookie"`
	}
	updater := Updater{}
	err := ctx.ShouldBindBodyWith(&updater, binding.JSON)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	token = updater.Token
	length = updater.Length
	cookie = updater.Cookie
	ctx.String(http.StatusOK, "ok")
}

func GetMails(ctx *gin.Context) {
	datas, err := getDailyMailForSingleWeek()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, datas)
}

func main() {
	e := gin.Default()
	e.GET("/mails", GetMails)
	e.POST("/init", Init)
	e.StaticFile("/index", "index.html")
	e.Run(":8080")
}
