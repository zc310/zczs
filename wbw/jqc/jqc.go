package jqc

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils"
	"github.com/zc310/zczs"
	"github.com/zc310/zczs/wbw"
	"log"
	"strings"

	"github.com/pkg/errors"

	"strconv"
)

func Match(issue string) (*zczs.ZcMatch, error) {

	doc, err := zczs.NewDoc(fmt.Sprintf("https://trade.500.com/jqc/?expect=%s", issue[2:]))
	if err != nil {
		return nil, err
	}

	var info zczs.ZcMatch
	var t string
	var ok bool
	if t, ok = doc.Find(".zcfilter-qih .chked").Attr("data-expect"); !ok {
		return nil, wbw.ErrorOldIssue
	}
	info.Issue = "20" + t
	if err != nil {
		return nil, errors.Errorf("无效期次信息:%s", t)
	}
	info.Endtime, err = wbw.GetEndTime(doc.Find(".zcfilter-endtime").Text())
	info.Issue = "20" + t
	if err != nil {
		return nil, errors.New("未发现有效截止日期")
	}
	var td, span *goquery.Selection
	doc.Find("#vsTable tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var item zczs.Match

		td = s.Find("td")
		if td.Length() != 10 {
			err = errors.New("对阵表格数据不符")
			return false
		}
		item.LeagueName = td.Eq(1).Text()
		item.LeagueSimpName = item.LeagueName
		item.MatchTime, err = wbw.GetEndTime(td.Eq(2).Text())
		if err != nil {
			return false
		}

		span = td.Eq(3).Find("span")
		if span.Length() != 2 {
			err = errors.New("对阵球队数量不符")
			return false
		}
		item.HomeTeam = span.Eq(0).Find("a").Text()
		item.AwayTeam = span.Eq(1).Find("a").Text()

		if t, ok = td.Eq(7).Find("a").Eq(0).Attr("href"); !ok {
			err = errors.New("无效比赛编号")
			return false
		}
		item.MatchID, err = strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(t, "http://odds.500.com/fenxi/shuju-"), ".shtml"), 10, 0)
		if err != nil {
			err = errors.Wrapf(err, "无效比赛编号")
			return false
		}

		item.HomeTeam = zczs.AddSpace(item.HomeTeam)
		item.AwayTeam = zczs.AddSpace(item.AwayTeam)

		info.Match = append(info.Match, &item)

		return true
	})

	return &info, err
}
func HisIssue() ([]string, error) {
	doc, err := zczs.NewDoc("https://kaijiang.500.com/jq4.shtml")
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

	return s1[:100], nil
}
func GetExpect() ([]string, error) {
	doc, err := zczs.NewDoc("https://trade.500.com/jqc/")
	if err != nil {
		return nil, err
	}
	var s1 []string
	var t string
	doc.Find(".zcfilter-qih li").Each(func(i int, s *goquery.Selection) {
		t = s.AttrOr("data-expect", "")
		if len(t) == 5 {
			s1 = append(s1, "20"+t)
		}
	})

	return s1, nil
}
func HandlerIssue(ctx *fasthttp.RequestCtx) {

	// {"c":"2019038","l":["2019041","2019040","2019039","2019038","2019037","2019036","2019035","2019034","2019033","2019032","2019031","2019030","2019029","2019028","2019027","2019026","2019025","2019024","2019023","2019022","2019021","2019020","2019019","2019018","2019017","2019016","2019015","2019014","2019013","2019012","2019011","2019010","2019009","2019008","2019007","2019006","2019005","2019004","2019003","2019002","2019001","2018176","2018175","2018174","2018173","2018172","2018171","2018170","2018169","2018168","2018167","2018166","2018165","2018164","2018163","2018162","2018161","2018160","2018159","2018158","2018157","2018156","2018155","2018154","2018153","2018152","2018151","2018150","2018149","2018148","2018147","2018146","2018145","2018144","2018143","2018142","2018141","2018140","2018139","2018138","2018137","2018136","2018135","2018134","2018133","2018132","2018131","2018130","2018129","2018128","2018127","2018126","2018125","2018124","2018123","2018122","2018121","2018120","2018119","2018118","2018117","2018116","2018115"]}

	var q zczs.Issue
	obj, err := GetExpect()
	if err != nil {
		return
	}

	for i, issue := range obj {
		if i == 0 {
			q.C = issue
		}
		q.L = append(q.L, issue)
	}
	l1, err := HisIssue()
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

	issue := string(ctx.QueryArgs().Peek("issue"))

	q, err := Match(issue)
	if err != nil {
		return
	}

	b, err := ffjson.Marshal(q)
	if err != nil {
		return
	}
	ctx.Write(b)
}

type Row struct {
	Fixtureid int64  `xml:"fixtureid,attr"`
	Pl        string `xml:"pl,attr"`
}
type Result struct {
	XMLName xml.Name `xml:"xml"`

	Row []Row `xml:"row"`
}

func HandlerOdds(ctx *fasthttp.RequestCtx) {
	var b []byte
	var err error
	issue := string(ctx.QueryArgs().Peek("MinIssue"))
	data, err := zczs.GetByte(fmt.Sprintf("https://www.500.com/static/public/jq4/daigou/xml/%s.xml", issue[2:]))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	v := Result{}
	err = xml.Unmarshal([]byte(data), &v)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	q := zczs.OuPei{}
	var t []string
	for _, row := range v.Row {

		//2.36&nbsp;3.35&nbsp;3.09
		t = strings.Split(row.Pl, "&nbsp;")
		if len(t) != 3 {
			continue
		}
		var one zczs.OuPeiItem

		one.MatchID = row.Fixtureid
		one.LastOdds3 = utils.StrToFloat(t[0])
		one.LastOdds1 = utils.StrToFloat(t[1])
		one.LastOdds0 = utils.StrToFloat(t[2])
		q[one.MatchID] = &one
	}

	b, err = ffjson.Marshal(q)
	if err != nil {
		return
	}

	ctx.Write(b)
}
