package sfc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils"
	"github.com/zc310/zczs"
)

type Match struct {
	Data struct {
		//Active       int64  `json:"active,string"`   // 0
		Addendtime string `json:"addendtime"` // 2019-06-05 20:00:00
		//Allowbuy     int64  `json:"allowbuy,string"` // 0
		Buyendtime   string `json:"buyendtime"`    // 2019-06-07 21:57:00
		Endtime      string `json:"endtime"`       // 2019-06-07 21:57:00
		Expect       int64  `json:"expect,string"` // 19074
		FutureExpect []struct {
			//Active        int64  `json:"active,string"`        // 1
			Addendtime string `json:"addendtime"` // 2019-06-01 21:30:00
			//Allowbuy      int64  `json:"allowbuy,string"`      // 1
			Buyendtime string `json:"buyendtime"` // 2019-06-01 21:30:00
			Endtime    string `json:"endtime"`    // 2019-06-01 21:30:00
			Expect     string `json:"expect"`     // 19073
			Opencode   string `json:"opencode"`   //
			//	Shenheprocess int64  `json:"shenheprocess,string"` // 1
		} `json:"future_expect"`
		Match []struct {
			Cdata struct {
				CountAgainst string `json:"count_against"` // 2|2|0|0
				CountAway    string `json:"count_away"`    // 5胜1平4负
				CountHome    string `json:"count_home"`    // 9胜0平1负
				CountScore   string `json:"count_score"`   // 4|0
				CountTzbl    string `json:"count_tzbl"`    //
				Draw         string `json:"draw"`          // 3.88
				Lost         string `json:"lost"`          // 1.87
				OddsURL      string `json:"odds_url"`      // https://live.m.500.com/detail/football/737845/analysis/zj
				Win          string `json:"win"`           // 3.80
			} `json:"cdata"`
			Mdata struct {
				//Awaystanding  int64  `json:"awaystanding,string"`   // 0
				Awaysxname    string `json:"awaysxname"`            // 韩国女
				Bgcolor       string `json:"bgcolor"`               // #E61A42
				CurrentExpect int64  `json:"current_expect,string"` // 19074
				Expect        []struct {
					//Active        int64  `json:"active,string"`        // 1
					//Addendtime    string `json:"addendtime"`           // 2019-06-01 21:30:00
					//Allowbuy      int64  `json:"allowbuy,string"`      // 1
					Buyendtime string `json:"buyendtime"` // 2019-06-01 21:30:00
					Endtime    string `json:"endtime"`    // 2019-06-01 21:30:00
					Expect     string `json:"expect"`     // 19073
					Opencode   string `json:"opencode"`   //
					//Shenheprocess int64  `json:"shenheprocess,string"` // 1
				} `json:"expect"`
				Expert struct {
				} `json:"expert"`
				Fid         int64  `json:"fid,string"`  // 768563
				Fsendtime   string `json:"fsendtime"`   // 2019-06-07 21:57
				Guestteamid string `json:"guestteamid"` // 858
				//Homestanding int64  `json:"homestanding,string"` // 0
				Homesxname   string `json:"homesxname"`   // 法国女
				Hometeamid   string `json:"hometeamid"`   // 790
				Isvalid      string `json:"isvalid"`      // 1
				Ordernum     string `json:"ordernum"`     // 1
				Resultscore  string `json:"resultscore"`  // 06-08 03:00
				Simpleleague string `json:"simpleleague"` // 女世杯
			} `json:"mdata"`
			Pdata struct {
				Draw string `json:"draw"` // 7.38
				Lost string `json:"lost"` // 19.16
				Win  string `json:"win"`  // 1.12
			} `json:"pdata"`
		} `json:"match"`
		Opencode string `json:"opencode"` //
		//Shenheprocess int64  `json:"shenheprocess,string"` // 1
	} `json:"data"`
	Message string `json:"message"`       // OK
	Status  int64  `json:"status,string"` // 100
}

func (p Match) Expect() []string {
	var a []string
	for _, v := range p.Data.FutureExpect {
		a = append(a, v.Expect)
	}

	return a
}

func (p Match) Odds() string {
	var a []string
	for _, m := range p.Data.Match {
		a = append(a, fmt.Sprintf("%s\t%s\t%s", m.Cdata.Win, m.Cdata.Draw, m.Cdata.Lost))
	}
	return strings.Join(a, "\n")
}
func GetZcMatchObject(issue string) (*Match, error) {
	var m Match
	b, err := zczs.GetByte(fmt.Sprintf("https://evs.500.com/esinfo/lotinfo/lot_info_modify?lotid=1&expect=%s&webviewsource=touch&platform=touch", issue))
	if err != nil {
		return &m, err
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		return &m, err
	}
	if m.Status != 100 {
		return &m, errors.New(m.Message)
	}
	return &m, nil
}

func HandlerOdds(ctx *fasthttp.RequestCtx) {
	var b []byte
	var err error
	issue := string(ctx.QueryArgs().Peek("MinIssue"))
	//zc r9

	obj, err := GetZcMatchObject(issue[2:])
	if err != nil {
		return
	}
	q := zczs.OuPei{}
	for _, m := range obj.Data.Match {
		var oi zczs.OuPeiItem
		oi.LastOdds3 = utils.StrToFloat(m.Cdata.Win)
		oi.LastOdds1 = utils.StrToFloat(m.Cdata.Draw)
		oi.LastOdds0 = utils.StrToFloat(m.Cdata.Lost)
		q[m.Mdata.Fid] = &oi
	}
	b, err = ffjson.Marshal(q)
	if err != nil {
		return
	}

	ctx.Write(b)
}
