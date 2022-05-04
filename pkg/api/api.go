package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/valyala/fasthttp"
	"github.com/zc310/utils"
	"github.com/zc310/utils/floatf"
	"github.com/zc310/zczs/pkg"
	"github.com/zc310/zczs/pkg/api/jczq"
)

const Max360Issue = 2019081
const Max360IssueJqc = 2019099

func Index(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString(ctx.Request.String())
}
func GetByte(ctx *fasthttp.RequestCtx) ([]byte, error) {
	host := ctx.Request.URI().String()
	if strings.Count(host, "360.cn") > 0 {
		ctx.Request.URI().SetScheme("https")
	} else if strings.Count(host, "zc310.tech") == 0 {
		return []byte{}, nil
	}
	return pkg.GetByte(RemoveIdSpm(ctx.Request.URI().String()))
}
func NotFound(ctx *fasthttp.RequestCtx) {
	b, err := GetByte(ctx)
	if err != nil {
		return
	}
	//log.Println(ctx.Request.String(), "\n", string(b))

	_, _ = ctx.Write(b)
}
func NoContent(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetStatusCode(http.StatusNoContent)
}
func NotOk(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("{}")
}
func NotFoundCache(ctx *fasthttp.RequestCtx) {
	key := RemoveIdSpm(ctx.Request.URI().String())
	if b, ok := pkg.Cache.Get(key); ok {
		_, _ = ctx.Write(b)
		return
	}
	b, err := GetByte(ctx)
	if err != nil {
		_, _ = ctx.WriteString("{}")
		return
	}
	_, _ = ctx.Write(b)
	pkg.Cache.Set(key, b)
}
func HisZcIssue(lotid, issue string) bool {
	return LotZcR9(lotid) && utils.StrToInt(issue) < Max360Issue
}
func HisJqcIssue(lotid, issue string) bool {
	return LotJqc(lotid) && utils.StrToInt(issue) < Max360IssueJqc
}
func LotZcR9(lotid string) bool {
	return lotid == "130011" || lotid == "130019"
}
func LotJqc(lotid string) bool {
	return lotid == "130018"
}
func ZczsQkjinfo(ctx *fasthttp.RequestCtx) {
	// GET http://m.cp.360.cn/int/qkjinfo?LotID=130011&Issue=2019049&id_spm=5cac7d89c283680210020e28
	// Host: m.cp.360.cn
	// {"Issue":"2019048","Code":"30131311031331","Date":"2019-04-07","Pool":"10476909","Sales":"23385958","Level1":"0:0","Level2":"187087:24","Level3":"0:0","Level":[["\u4e00\u7b49\u5956","0","0"],["\u4e8c\u7b49\u5956","187087","24"]]}

	issue := string(ctx.QueryArgs().Peek("Issue"))
	lotid := string(ctx.QueryArgs().Peek("LotID"))

	//历史数据时

	if HisZcIssue(lotid, issue) || HisJqcIssue(lotid, issue) {
		NotFoundCache(ctx)
		return
	}

	var b []byte
	var err error
	var q1 pkg.QkjInfo
	if LotZcR9(lotid) {
		b, err = pkg.ReadFile(filepath.Join("sfc", fmt.Sprintf("%s_奖号.toml", issue)))
	} else if LotJqc(lotid) {
		b, err = pkg.ReadFile(filepath.Join("jqc", fmt.Sprintf("%s_奖号.toml", issue)))
	}
	if err != nil {
		log.Println(err)
		return
	}
	err = toml.Unmarshal(b, &q1)
	if err != nil {
		log.Println(err)
		return
	}
	b, err = json.Marshal(q1)
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = ctx.Write(b)
}
func ZczsIssue(ctx *fasthttp.RequestCtx) {
	// /zczs/issue?lotid=130011
	// Host: cp.360.cn
	// {"c":"2019038","l":["2019041","2019040","2019039","2019038","2019037","2019036","2019035","2019034","2019033","2019032","2019031","2019030","2019029","2019028","2019027","2019026","2019025","2019024","2019023","2019022","2019021","2019020","2019019","2019018","2019017","2019016","2019015","2019014","2019013","2019012","2019011","2019010","2019009","2019008","2019007","2019006","2019005","2019004","2019003","2019002","2019001","2018176","2018175","2018174","2018173","2018172","2018171","2018170","2018169","2018168","2018167","2018166","2018165","2018164","2018163","2018162","2018161","2018160","2018159","2018158","2018157","2018156","2018155","2018154","2018153","2018152","2018151","2018150","2018149","2018148","2018147","2018146","2018145","2018144","2018143","2018142","2018141","2018140","2018139","2018138","2018137","2018136","2018135","2018134","2018133","2018132","2018131","2018130","2018129","2018128","2018127","2018126","2018125","2018124","2018123","2018122","2018121","2018120","2018119","2018118","2018117","2018116","2018115"]}
	lotid := string(ctx.QueryArgs().Peek("lotid"))
	var q1 pkg.Issue
	var b []byte
	var err error
	if LotZcR9(lotid) {
		b, err = pkg.ReadFile(filepath.Join("sfc", "期号.toml"))
	}
	if lotid == "130018" {
		b, err = pkg.ReadFile(filepath.Join("jqc", "期号.toml"))
	}
	if lotid == "130041" {
		b, err = pkg.ReadFile(filepath.Join("dc", "期号.toml"))
	}
	if lotid == "130042" {
		jczq.HandlerIssue(ctx)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	err = toml.Unmarshal(b, &q1)
	if err != nil {
		log.Println(err)
		return
	}
	b, err = json.Marshal(q1)
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = ctx.Write(b)
}

func ZczsZcmatch(ctx *fasthttp.RequestCtx) {
	// GET /zczs/zcmatch?lotid=130011&issue=2019038&id_spm=5c92840a32797804bc65593b HTTP/1.1
	// Host: cp.360.cn
	// {"endtime":"2019-03-20 22:45:00","match":[{"Issue":"2019038","ItemID":"1","MatchID":"1685467","MatchState":"0","LeagueID":"1366","LeagueSimpName":"\u53cb\u8c0a\u8d5b","LeagueName":"\u53cb\u8c0a\u8d5b","LeagueColor":"4666bb","HomeTeam":"\u5fb7\u56fd","AwayTeam":"\u585e\u5c14\u7ef4","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-21 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"2","MatchID":"1685468","MatchState":"0","LeagueID":"1366","LeagueSimpName":"\u53cb\u8c0a\u8d5b","LeagueName":"\u53cb\u8c0a\u8d5b","LeagueColor":"4666bb","HomeTeam":"\u5a01\u5c14\u58eb","AwayTeam":"\u7279\u7acb\u5c3c","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-21 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"3","MatchID":"1685471","MatchState":"0","LeagueID":"1366","LeagueSimpName":"\u53cb\u8c0a\u8d5b","LeagueName":"\u53cb\u8c0a\u8d5b","LeagueColor":"4666bb","HomeTeam":"\u79d1\u7d22\u6c83","AwayTeam":"\u4e39\u9ea6","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 02:00:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"4","MatchID":"1685472","MatchState":"0","LeagueID":"1366","LeagueSimpName":"\u53cb\u8c0a\u8d5b","LeagueName":"\u53cb\u8c0a\u8d5b","LeagueColor":"4666bb","HomeTeam":"\u7f8e\u56fd","AwayTeam":"\u5384\u74dc\u591a","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 08:00:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"5","MatchID":"1647096","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u54c8\u8428\u514b","AwayTeam":"\u82cf\u683c\u5170","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-21 23:00:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"6","MatchID":"1647097","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u585e\u6d66\u8def","AwayTeam":"\u5723\u9a6c\u529b","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 01:00:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"7","MatchID":"1646946","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u8377\u5170","AwayTeam":"\u767d\u4fc4","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"8","MatchID":"1646947","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u5317\u7231","AwayTeam":"\u7231\u6c99\u5c3c","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"9","MatchID":"1646987","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u514b\u7f57\u5730","AwayTeam":"\u963f\u585e\u62dc","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"10","MatchID":"1646986","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u65af\u6d1b\u4f10","AwayTeam":"\u5308\u7259\u5229","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"11","MatchID":"1647038","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u5965\u5730\u5229","AwayTeam":"\u6ce2\u5170","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"12","MatchID":"1647036","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u4ee5\u8272\u5217","AwayTeam":"\u65af\u6d1b\u6587","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"13","MatchID":"1647037","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u9a6c\u5176\u987f","AwayTeam":"\u62c9\u8131\u7ef4","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"},{"Issue":"2019038","ItemID":"14","MatchID":"1647098","MatchState":"0","LeagueID":"67","LeagueSimpName":"\u6b27\u6d32\u676f","LeagueName":"\u6b27\u6d32\u676f","LeagueColor":"660000","HomeTeam":"\u6bd4\u5229\u65f6","AwayTeam":"\u4fc4\u7f57\u65af","LotLose":"0","LotEndTime":"2019-03-20 23:00:00","MatchTime":"2019-03-22 03:45:00","VsReverseFlag":"0","DisableFlag":"0"}],"issue":"2019038","state":"-1"}

	lotid := string(ctx.QueryArgs().Peek("lotid"))
	issue := string(ctx.QueryArgs().Peek("issue"))

	if (!LotZcR9(lotid) && !LotJqc(lotid)) || len(issue) == 0 {
		NotFound(ctx)
		return
	}
	var q1 pkg.ZcMatch
	var b []byte
	var err error
	if HisZcIssue(lotid, issue) || HisJqcIssue(lotid, issue) {
		NotFoundCache(ctx)
		return
	}

	if LotZcR9(lotid) {
		b, err = pkg.ReadFile(filepath.Join("sfc", fmt.Sprintf("%s_对阵.toml", issue)))
	}
	if LotJqc(lotid) {
		b, err = pkg.ReadFile(filepath.Join("jqc", fmt.Sprintf("%s_对阵.toml", issue)))
	}
	if err != nil {
		log.Println(err)
		return
	}

	if err = toml.Unmarshal(b, &q1); err != nil {
		return
	}

	q1.Issue = issue
	for i, m := range q1.Match {
		m.HomeTeam = pkg.AddSpace(m.HomeTeam)
		m.AwayTeam = pkg.AddSpace(m.AwayTeam)
		m.Issue = q1.Issue
		m.MatchID = int64(i) + 1
		m.ItemID = int64(i) + 1
	}
	b, err = json.Marshal(q1)
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = ctx.Write(b)

}
func SfcExtra(ctx *fasthttp.RequestCtx) {
	// GET /sfc/extra?key=/sfc/360dd,/sfc/help,/sfc/360dd/2019038&id_spm=5c92840a32797804bc65593c HTTP/1.1
	// Host: cp.360.cn
	//	{
	//		"/sfc/360dd": "",
	//		"/sfc/help": "{\"top\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-3629447-1-1.html\"},\"zsgl2\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-3506968-1-1.html\"},\"pxys\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2964532-1-1.html\"},\"plpx\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2853637-1-1.html\"},\"gsfb\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2996563-1-1.html\"},\"gsdb\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2983003-1-1.html\"},\"2dfb\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2996563-1-1.html\"},\"2ddb\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2983003-1-1.html\"},\"pkyun\":{\"url\":\"http:\\/\\/yunpan.cn\\/Q75UTWktrdu3R\",\"text\":\"\\u70b9\\u51fb\\u4e0b\\u8f7d\\u5168\\u90e8\\u6570\\u636e(\\u622a\\u6b62\\u524d1\\u5c0f\\u65f6\\u540e\\u5f00\\u59cb\\u4e0b\\u8f7d,\\u63d0\\u53d6\\u78010132)\"},\"pjdb2\":{\"url\":\"http:\\/\\/bbs.360safe.com\\/thread-2911191-1-1.html\"}}"
	//	}

	NotFoundCache(ctx)
}

func GetOupei(ctx *fasthttp.RequestCtx) {
	//http://cp.360.cn/int/getoupei/?LotID=130011&MinIssue=2017114&GcID=9999&id_spm=5cbec9cec28368146c4a9bd6
	//{"1399706":{"MatchID":"1399706","FirstOdds3":"2.06","FirstOdds1":"3.47","FirstOdds0":"3.29","LastOdds3":"2.16","Odds3Trend":"1.00","LastOdds1":"3.57","Odds1Trend":"-1.00","LastOdds0":"3.04","Odds0Trend":"-1.00"},"1399707":{"MatchID":"1399707","FirstOdds3":"2.00","FirstOdds1":"3.39","FirstOdds0":"3.52","LastOdds3":"1.89","Odds3Trend":"-1.00","LastOdds1":"3.37","Odds1Trend":"-1.00","LastOdds0":"4.03","Odds0Trend":"1.00"},"1399709":{"MatchID":"1399709","FirstOdds3":"2.39","FirstOdds1":"3.22","FirstOdds0":"2.87","LastOdds3":"2.55","Odds3Trend":"1.00","LastOdds1":"3.10","Odds1Trend":"-1.00","LastOdds0":"2.80","Odds0Trend":"-1.00"},"1399710":{"MatchID":"1399710","FirstOdds3":"2.37","FirstOdds1":"3.33","FirstOdds0":"2.82","LastOdds3":"2.38","Odds3Trend":"-1.00","LastOdds1":"3.35","Odds1Trend":"-1.00","LastOdds0":"2.83","Odds0Trend":"1.00"},"1436070":{"MatchID":"1436070","FirstOdds3":"2.57","FirstOdds1":"3.08","FirstOdds0":"2.77","LastOdds3":"2.29","Odds3Trend":"-1.00","LastOdds1":"3.00","Odds1Trend":"-1.00","LastOdds0":"3.37","Odds0Trend":"1.00"},"1436071":{"MatchID":"1436071","FirstOdds3":"2.93","FirstOdds1":"3.30","FirstOdds0":"2.31","LastOdds3":"2.93","Odds3Trend":"-1.00","LastOdds1":"3.30","Odds1Trend":"-1.00","LastOdds0":"2.37","Odds0Trend":"1.00"},"1436072":{"MatchID":"1436072","FirstOdds3":"1.67","FirstOdds1":"3.49","FirstOdds0":"5.23","LastOdds3":"1.53","Odds3Trend":"-1.00","LastOdds1":"3.74","Odds1Trend":"1.00","LastOdds0":"6.72","Odds0Trend":"1.00"},"1436073":{"MatchID":"1436073","FirstOdds3":"2.51","FirstOdds1":"3.10","FirstOdds0":"2.80","LastOdds3":"2.77","Odds3Trend":"-1.00","LastOdds1":"2.98","Odds1Trend":"-1.00","LastOdds0":"2.70","Odds0Trend":"1.00"},"1436074":{"MatchID":"1436074","FirstOdds3":"3.42","FirstOdds1":"3.31","FirstOdds0":"2.06","LastOdds3":"3.29","Odds3Trend":"-1.00","LastOdds1":"3.29","Odds1Trend":"-1.00","LastOdds0":"2.17","Odds0Trend":"1.00"},"1436075":{"MatchID":"1436075","FirstOdds3":"2.66","FirstOdds1":"3.25","FirstOdds0":"2.54","LastOdds3":"2.40","Odds3Trend":"-1.00","LastOdds1":"3.31","Odds1Trend":"-1.00","LastOdds0":"2.86","Odds0Trend":"1.00"},"1436076":{"MatchID":"1436076","FirstOdds3":"1.35","FirstOdds1":"4.79","FirstOdds0":"7.74","LastOdds3":"1.26","Odds3Trend":"-1.00","LastOdds1":"5.51","Odds1Trend":"1.00","LastOdds0":"10.72","Odds0Trend":"1.00"},"1436077":{"MatchID":"1436077","FirstOdds3":"1.39","FirstOdds1":"4.44","FirstOdds0":"7.51","LastOdds3":"1.40","Odds3Trend":"1.00","LastOdds1":"4.32","Odds1Trend":"-1.00","LastOdds0":"8.13","Odds0Trend":"-1.00"},"1436078":{"MatchID":"1436078","FirstOdds3":"1.34","FirstOdds1":"4.65","FirstOdds0":"8.94","LastOdds3":"1.30","Odds3Trend":"-1.00","LastOdds1":"4.90","Odds1Trend":"1.00","LastOdds0":"10.79","Odds0Trend":"1.00"},"1436079":{"MatchID":"1436079","FirstOdds3":"1.56","FirstOdds1":"3.66","FirstOdds0":"6.19","LastOdds3":"1.61","Odds3Trend":"1.00","LastOdds1":"3.57","Odds1Trend":"-1.00","LastOdds0":"6.13","Odds0Trend":"-1.00"}}
	lotid := string(ctx.QueryArgs().Peek("LotID"))
	issue := string(ctx.QueryArgs().Peek("MinIssue"))
	GcID := string(ctx.QueryArgs().Peek("GcID"))

	if HisZcIssue(lotid, issue) || HisJqcIssue(lotid, issue) {
		NotFoundCache(ctx)
		return
	}

	var b []byte
	var err error
	//zc r9
	if LotZcR9(lotid) && utils.StrToInt(issue) > Max360Issue && GcID == "9999" {
		b, err = pkg.ReadFile(filepath.Join("sfc", fmt.Sprintf("%s_平均赔率.txt", issue)))
	}
	if LotJqc(lotid) && utils.StrToInt(issue) > Max360IssueJqc && GcID == "9999" {
		b, err = pkg.ReadFile(filepath.Join("jqc", fmt.Sprintf("%s_平均赔率.txt", issue)))
	}
	if lotid == "130041" {
		_, _ = ctx.WriteString("{}")
		return
	}
	if err != nil {
		log.Println(err)
	}
	var o1 []float64
	o1, err = floatf.Parse(string(b))
	if err != nil {
		log.Println(err)
		return
	}

	if LotZcR9(lotid) && len(o1) != 42 {
		log.Println(fmt.Errorf("%s 赔率有误", issue))
		return
	}
	if LotJqc(lotid) && len(o1) < 12 {
		log.Println(fmt.Errorf("%s 赔率有误", issue))
		return
	}
	o2 := make([][]float64, len(o1)/3)
	var t int
	for i := 0; i < len(o1)/3; i++ {
		t = i * 3
		o2[i] = make([]float64, 4)
		o2[i][3] = o1[t]
		o2[i][1] = o1[t+1]
		o2[i][0] = o1[t+2]
	}
	q1 := pkg.OuPei{}
	for i, m := range o2 {
		var oi pkg.OuPeiItem
		oi.LastOdds3 = m[3]
		oi.LastOdds1 = m[1]
		oi.LastOdds0 = m[0]
		oi.MatchID = int64(i + 1)

		q1[int64(i+1)] = &oi
	}
	b, err = json.Marshal(q1)
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = ctx.Write(b)

}

func Dcmatch(ctx *fasthttp.RequestCtx) {
	//https://cp.360.cn/zczs/match?lotid=130041&issue=190904&id_spm=60e876e5d0e4900668046284
	//[
	//	{
	//	"Issue": "190904",
	//	"ItemID": "1",
	//	"MatchID": "1674060",
	//	"MatchState": "99",
	//	"LeagueID": "1292",
	//	"LeagueSimpName": "韩挑K",
	//	"LeagueName": "韩挑K",
	//	"LeagueColor": "0099cc",
	//	"HomeTeam": "大田市民",
	//	"AwayTeam": "釜山偶像",
	//	"LotLose": "1",
	//	"LotEndTime": "2019-09-17 17:55:00",
	//	"MatchTime": "2019-09-17 18:00:00",
	//	"VsReverseFlag": "0",
	//	"DisableFlag": "0",
	//	"result": {
	//	"codespf": "3",
	//	"win": "0",
	//	"halfwin": "-1",
	//	"lose": "0",
	//	"halflose": "-1",
	//	"spspf": "2.570000"
	//	}
	//	},
	lotid := string(ctx.QueryArgs().Peek("lotid"))
	issue := utils.StrToInt(string(ctx.QueryArgs().Peek("issue")))
	sp := string(ctx.QueryArgs().Peek("sp")) != ""
	// dc
	if lotid == "130041" && issue <= 190904 {
		NotFoundCache(ctx)
		return
	}
	if sp && lotid == "130041" {

		return
	}

	if lotid == "130041" {

		return
	}
	NotFoundCache(ctx)
}
