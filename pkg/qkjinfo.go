package pkg

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type QkjInfo struct {
	Issue  string     `json:"Issue" comment:"期号"`
	Code   string     `json:"Code" comment:"奖号"`
	Date   string     `json:"Date" comment:"日期"`
	Pool   string     `json:"Pool"`
	Sales  string     `json:"Sales"`
	Level1 string     `json:"Level1,omitempty"`
	Level2 string     `json:"Level2,omitempty"`
	Level3 string     `json:"Level3,omitempty"`
	Level  [][]string `json:"Level,omitempty"`
}

type Issue struct {
	C string   `json:"c" comment:"当前期"`
	L []string `json:"l" toml:"L,multiline"  comment:"历史期"`
}
type Match struct {
	AwayTeam       string `json:"AwayTeam"  comment:"客队"`                  // 布莱顿
	DisableFlag    int64  `json:"DisableFlag,string,omitempty" toml:"-"`   // 0
	HomeTeam       string `json:"HomeTeam"  comment:"主队"`                  // 托特纳姆热刺
	Issue          string `json:"Issue" toml:"-"`                          // 2019056
	ItemID         int64  `json:"ItemID,string" toml:"-"`                  // 1
	LeagueColor    string `json:"LeagueColor,omitempty" toml:"-"`          // FF3333
	LeagueID       int64  `json:"LeagueID,string,omitempty" toml:"-"`      // 36
	LeagueName     string `json:"LeagueName,omitempty"  toml:"-"`          // 英超
	LeagueSimpName string `json:"LeagueSimpName,omitempty" toml:"-"`       // 英超
	LotEndTime     string `json:"LotEndTime,omitempty" toml:"-"`           // 2019-04-23 22:30:00
	LotLose        int64  `json:"LotLose,string,omitempty" toml:"-"`       // 0
	MatchID        int64  `json:"MatchID,string,omitempty" toml:"-"`       // 1552470
	MatchState     int64  `json:"MatchState,string,omitempty" toml:"-"`    // 0
	MatchTime      string `json:"MatchTime,omitempty" toml:"-"`            // 2019-04-24 02:45:00
	VsReverseFlag  int64  `json:"VsReverseFlag,string,omitempty" toml:"-"` // 0
}
type ZcMatch struct {
	Endtime string   `json:"endtime" comment:"截止时间"` // 1899-12-29 23:45:00
	Issue   string   `json:"issue" toml:"-"`         // 2019056
	Match   []*Match `json:"match"`
	State   int64    `json:"state,string,omitempty" toml:"-"` // -1
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
