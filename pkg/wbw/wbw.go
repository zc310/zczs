package wbw

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

func GetEndTime(s string) (string, error) {
	re := regexp.MustCompile(`(\d{2})-(\d{2}) (\d{2}):(\d{2})`)
	r := re.FindAllStringSubmatch(s, -1)
	if len(r) == 0 {
		return "", fmt.Errorf("error %s", s)
	}
	return fmt.Sprintf("%d-%s-%s %s:%s:00", time.Now().Year(), r[0][1], r[0][2], r[0][3], r[0][4]), nil
}

var ErrorOldIssue = errors.New("历史期次")
