package ics

import (
	"fmt"
	"regexp"
	"testing"
)

func TestICSTrigger(t *testing.T) {
	// lastTime := time.Now().Add(-5 * time.Hour)
	// if !triggerTimer(lastTime, 18, 0) {
	// 	t.Error("trigger wrong")
	// }

	// lastTime1 := time.Now()
	// if triggerTimer(lastTime1, 18, 0) {
	// 	t.Error("trigger1 wrong")
	// }

	stmp := regexp.MustCompile(`^<(\d{4}-\d{2}-\d{2})( [A-Za-z\p{Han}]+)?( \d{2}:\d{2})?(?:-(\d{2}:\d{2}))?( \+\d+[dwmy])?>`)
	if m := stmp.FindStringSubmatch("<2022-01-24 周一 08:00 +1d>"); m != nil {
		fmt.Println(m)
	}

	line := "   DEADLINE: <2022-01-30 周日 23:00> SCHEDULED: <2022-01-24 周一 08:00 +1d>  "
	var scheduleRegexp = regexp.MustCompile(`SCHEDULED:\s*(<[^>]+>)`)
	var deadlineRegexp = regexp.MustCompile(`DEADLINE:\s*(<[^>]+>)`)
	if m := scheduleRegexp.FindStringSubmatch(line); m != nil {
		fmt.Println(m)
	} else if m := deadlineRegexp.FindStringSubmatch(line); m != nil {
		fmt.Println(m)
	}
}
