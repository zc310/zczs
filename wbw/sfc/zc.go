package sfc

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils"
	"github.com/zc310/zczs"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func Sfc(issue string, r9 bool) (*zczs.QkjInfo, error) {
	res, err := http.Get(fmt.Sprintf("https://kaijiang.500.com/shtml/sfc/%s.shtml", issue[2:]))
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	doc, err := goquery.NewDocumentFromReader(utils.GB18030Dec(res.Body))
	if err != nil {
		return nil, err
	}
	var s1 []string
	var s0 string
	var info zczs.QkjInfo
	tb := doc.Find(".kjxq_box02 .kj_tablelist02")
	if tb.Length() != 2 {
		return nil, fmt.Errorf("无效开奖期次:%s", issue)
	}
	tr := tb.Eq(0).Find("tr")

	if tr.Length() != 4 {
		return nil, fmt.Errorf("无效开奖期次:%s", issue)
	}
	tr.Eq(2).Find("td").Each(func(i int, s *goquery.Selection) {
		s0 = strings.TrimSpace(s.Text())
		if s0 == "-" {
			s0 = "*"
		}
		if strings.HasPrefix(s0, "{") {
			err = errors.New("无效开奖号")
			return
		}

		s1 = append(s1, s0)

	})
	if err != nil {
		return nil, err
	}
	if len(s1) == 13 {
		s1 = append(s1, "*")
	}

	info.Code = strings.Join(s1, "")
	info.Issue = issue

	re := regexp.MustCompile(`(\d{4})年(\d+)月(\d+)日`)
	r := re.FindAllStringSubmatch(tr.Eq(0).Text(), -1)
	if len(r) == 0 {
		return nil, fmt.Errorf("error %s", tr.Eq(0).Text())
	}
	info.Date = r[0][1] + "-" + fmDay(r[0][2]) + "-" + fmDay(r[0][3])

	//kj

	tr = tb.Eq(1).Find("tr")
	if tr.Length() != 6 {
		return nil, fmt.Errorf("无效开奖期次:%s", issue)
	}
	var td *goquery.Selection
	for i := 2; i < 5; i++ {
		td = tr.Eq(i).Find("td")
		var lv []string
		lv = make([]string, 3)
		if td.Length() == 3 {

			lv[0] = td.Eq(0).Text()
			lv[2] = td.Eq(1).Text()
			if strings.Count(lv[1], "-") > 0 {
				lv[2] = "0"
			}
			lv[1] = strconv.FormatInt(getPoolMoney(td.Eq(2).Text()), 10)
			info.Level = append(info.Level, lv)
		}

	}
	if r9 {
		info.Level = info.Level[2:]
	} else {
		info.Level = info.Level[:2]
	}

	return &info, nil
}
func getPoolMoney(s string) int64 {
	var d int64
	var err error

	d, err = strconv.ParseInt(strings.Replace(strings.TrimSuffix(s, "元"), ",", "", -1), 10, 0)
	if err != nil {
		return 0
	}
	return d
}
func SfcHisIssue() ([]string, error) {
	res, err := http.Get("https://kaijiang.500.com/sfc.shtml")
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	doc, err := goquery.NewDocumentFromReader(utils.GB18030Dec(res.Body))
	if err != nil {
		return nil, err
	}
	var s1 []string
	var t string
	doc.Find(".iSelectList a").Each(func(i int, s *goquery.Selection) {
		t = s.Text()
		if len(t) == 5 {
			s1 = append(s1, "20"+t)
		}
	})

	return s1, nil
}
func fmDay(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func HandlerIssue(ctx *fasthttp.RequestCtx) {

	// {"c":"2019038","l":["2019041","2019040","2019039","2019038","2019037","2019036","2019035","2019034","2019033","2019032","2019031","2019030","2019029","2019028","2019027","2019026","2019025","2019024","2019023","2019022","2019021","2019020","2019019","2019018","2019017","2019016","2019015","2019014","2019013","2019012","2019011","2019010","2019009","2019008","2019007","2019006","2019005","2019004","2019003","2019002","2019001","2018176","2018175","2018174","2018173","2018172","2018171","2018170","2018169","2018168","2018167","2018166","2018165","2018164","2018163","2018162","2018161","2018160","2018159","2018158","2018157","2018156","2018155","2018154","2018153","2018152","2018151","2018150","2018149","2018148","2018147","2018146","2018145","2018144","2018143","2018142","2018141","2018140","2018139","2018138","2018137","2018136","2018135","2018134","2018133","2018132","2018131","2018130","2018129","2018128","2018127","2018126","2018125","2018124","2018123","2018122","2018121","2018120","2018119","2018118","2018117","2018116","2018115"]}

	var q zczs.Issue
	obj, err := GetZcMatchObject("")
	if err != nil {
		return
	}

	for i, issue := range obj.Expect() {
		if i == 0 {
			q.C = "20" + issue
		}
		q.L = append(q.L, "20"+issue)
	}
	l1, err := SfcHisIssue()
	if err == nil {
		q.L = append(q.L, l1...)
	}
	b, err := ffjson.Marshal(q)
	if err != nil {
		return
	}
	ctx.Write(b)
}

func HandlerMatch(ctx *fasthttp.RequestCtx) {
	var q zczs.ZcMatch
	q.Issue = string(ctx.QueryArgs().Peek("issue"))

	obj, err := GetZcMatchObject(q.Issue[2:])
	if err != nil {
		log.Println(q.Issue, err)
		return
	}
	for _, m := range obj.Data.Match {
		var mc zczs.Match
		mc.HomeTeam = zczs.AddSpace(m.Mdata.Homesxname)
		mc.AwayTeam = zczs.AddSpace(m.Mdata.Awaysxname)
		mc.Issue = q.Issue
		mc.LeagueName = m.Mdata.Simpleleague
		mc.LeagueColor = m.Mdata.Bgcolor
		mc.MatchTime = m.Mdata.Fsendtime + ":00"
		mc.MatchID = m.Mdata.Fid
		q.Match = append(q.Match, &mc)
		if mc.MatchTime > q.Endtime {
			q.Endtime = mc.MatchTime
		}
	}

	b, err := ffjson.Marshal(q)
	if err != nil {
		return
	}
	ctx.Write(b)
}
