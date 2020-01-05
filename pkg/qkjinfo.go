package pkg

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type QkjInfo struct {
	Issue  string     `json:"Issue"`
	Code   string     `json:"Code"`
	Date   string     `json:"Date"`
	Pool   string     `json:"Pool"`
	Sales  string     `json:"Sales"`
	Level1 string     `json:"Level1"`
	Level2 string     `json:"Level2"`
	Level3 string     `json:"Level3"`
	Level  [][]string `json:"Level,omitempty"`
}

type Issue struct {
	C string   `json:"c"`
	L []string `json:"l"`
}
type Match struct {
	AwayTeam       string `json:"AwayTeam"`             // 布莱顿
	DisableFlag    int64  `json:"DisableFlag,string"`   // 0
	HomeTeam       string `json:"HomeTeam"`             // 托特纳姆热刺
	Issue          string `json:"Issue"`                // 2019056
	ItemID         int64  `json:"ItemID,string"`        // 1
	LeagueColor    string `json:"LeagueColor"`          // FF3333
	LeagueID       int64  `json:"LeagueID,string"`      // 36
	LeagueName     string `json:"LeagueName"`           // 英超
	LeagueSimpName string `json:"LeagueSimpName"`       // 英超
	LotEndTime     string `json:"LotEndTime"`           // 2019-04-23 22:30:00
	LotLose        int64  `json:"LotLose,string"`       // 0
	MatchID        int64  `json:"MatchID,string"`       // 1552470
	MatchState     int64  `json:"MatchState,string"`    // 0
	MatchTime      string `json:"MatchTime"`            // 2019-04-24 02:45:00
	VsReverseFlag  int64  `json:"VsReverseFlag,string"` // 0
}
type ZcMatch struct {
	Endtime string   `json:"endtime"` // 1899-12-29 23:45:00
	Issue   string   `json:"string"`  // 2019056
	Match   []*Match `json:"match"`
	State   int64    `json:"state,string"` // -1
}

type OuPeiItem struct {
	FirstOdds0 float64 `json:"FirstOdds0,string"` // 3.29
	FirstOdds1 float64 `json:"FirstOdds1,string"` // 3.47
	FirstOdds3 float64 `json:"FirstOdds3,string"` // 2.06
	LastOdds0  float64 `json:"LastOdds0,string"`  // 3.04
	LastOdds1  float64 `json:"LastOdds1,string"`  // 3.57
	LastOdds3  float64 `json:"LastOdds3,string"`  // 2.16
	MatchID    int64   `json:"MatchID,string"`    // 1399706
	Odds0Trend float64 `json:"Odds0Trend,string"` // -1.00
	Odds1Trend float64 `json:"Odds1Trend,string"` // -1.00
	Odds3Trend float64 `json:"Odds3Trend,string"` // 1.00
}
type OuPei map[int64]*OuPeiItem

func AddSpace(s string) string {
	//拜　仁
	s = strings.Replace(s, " ", "", -1)
	var a []string
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		a = append(a, fmt.Sprintf("%c", r))
		s = s[size:]
	}
	if len(a) == 2 {
		return fmt.Sprintf("%s　%s", a[0], a[1])
	}
	return strings.Join(a, "")
}
