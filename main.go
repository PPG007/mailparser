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

func init() {
	token = `QXltVG9rZW4tMTEzMjg4MTktY1o5cWtpV3N4c1l0YWYwVFVraXVBVG5uWkt0bzF5ZHpzb0F2UTNjMFlVWjNaOUdPclY`
	cookie = `cna=DHdxG+Z9onwCAd0DBAu7d4LJ; aliyun_choice=CN; t=f01a3004efb070c6e066130d1cabdc10; login_aliyunid_pk=1191551862183498; aliyun_country=CN; aliyun_site=CN; currentRegionId=cn-hangzhou; isg=BHt7D1EWnW-ITKHc1GTIsCyFClbl0I_SoEHGsm04Q3qRzJqu9aOzIocF4myCa-fK; tfstk=c__CBNxXjIvw7b26QHNa4Ho-MGL5ZFu6wD9CORVSoj6qWIfCiiu2h4VJqA8ywC1..; l=eBa1PbOrL7BM36e5BOfZourza7797IRf_uPzaNbMiOCP_D5p5K7dW6xIpLT9CnGVHsKwR3yJlgPkB-YwWzO-nxv9-zq4RqHmndC..; alimail_sso_umid=wV12z62edf1586f5bb9db00076ec277b7; alimail_init_lang=zh_CN; alimail_browser_instance=dC0xMjQ0NzktU1ltbFBi8981; alimail_sid=I1666AA1-M1Z22HIMD1MV6A3WQI1R3-DI7YAR6L-AK91; alimail_sdata0=a24zos5gOAbHitWQr5w%2FAFIJUuxjvn9tro1VNQUp0SDTYxFPuqaJ9Ghd5NAatIf7agYva%2B6TI%2FC0T2sQfy7gr2QBQbpYUBlAzu73GSDSBU%2FnQRZJCk%2FaZvF2ItoVfuR3s865DBAPd9OngeDNbbAirg%3D%3D; core_token="UNIFY:5a97cab2-61a2-48fe-85c3-781bc2de06d0|47f255d8-121b-47d4-8f52-1e2e4bc120a9"; core_heart_beat=1660359348054; alimail_core_session_key=QXltU2Vzc2lvbi00NzI1OS1iYUlKSWpUYURFR0gyWFZRcjA0TVNxNmZIaVpVQjUxRnRNb1VhVG5oZFRsWFdEU3pZSw; alimail_auth_mode=core; alimail_auth_session_key=QXltU2Vzc2lvbi00NzI2MC1OaW5FcEQ3Y3FoSlNrU3VmTEw1YWJxQkhsSDNBUkxxTEgwaHIxS0tZS3FBZmNJT0QwdQ; alimail_session_version_key=0.1.35; _csrf_token_=QXltVG9rZW4tMTEzMjg4MTktY1o5cWtpV3N4c1l0YWYwVFVraXVBVG5uWkt0bzF5ZHpzb0F2UTNjMFlVWjNaOUdPclY; at="koston.zhuang@maiscrm.comXhKEk1660359349571"; alimail_session_template_key=v4; udtoken="koston.zhuang@maiscrm.com:9e2eb7db073bbd12be13c3cf079e8f06:6460581660359354603424"`
	length = 5
}

type AliMailResponse struct {
	AliMailData []AliMailData `json:"dataList"`
}

type AliMailData struct {
	Subject   string `json:"subject"`
	TimeStamp int64  `json:"timestamp"`
	Body      string `json:"encSummary"`
}

type ResultResponse struct {
	SendAt  string `json:"sendAt"`
	Content string `json:"content"`
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
	resps := []ResultResponse{}
	for _, data := range datas {
		resps = append(resps, ResultResponse{
			SendAt:  TransMillTimestampToTime(data.TimeStamp).Format(TIME_FORMATTER),
			Content: ParseAndConcatSingleDayMail(data.Body),
		})
	}
	ctx.JSON(http.StatusOK, resps)
}

func main() {
	e := gin.Default()
	e.GET("/mails", GetMails)
	e.POST("/init", Init)
	e.Run(":8080")
}
