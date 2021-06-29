package org

import (
	"fmt"
	"regexp"
	"strings"
)

type Schedule struct {
	Time     string
	Start    string
	End      string
	Interval string
}

type Deadline struct {
	Time string
}

type TimeProperty struct {
	Schedule Node
	Deadline Node
	Closed   Node
}

var scheduleRegexp = regexp.MustCompile(`SCHEDULED:\s*<([^>]+)>`)
var deadlineRegexp = regexp.MustCompile(`DEADLINE:\s*<([^>]+)>`)
var timeDurationRegexp = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})( [A-Za-z\p{Han}]+)?( (\d{2}:\d{2})-(\d{2}:\d{2}))?( \+\d+[dwmy])?`)
var timeNotDurationRegexp = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})( [A-Za-z\p{Han}]+)?( \d{2}:\d{2})?( \+\d+[dwmy])?`)

func lexSchedule(line string) (token, bool) {
	if m := scheduleRegexp.FindStringSubmatch(line); m != nil {
		return token{"schedule", len(m[1]), line, m}, true
	} else if m := deadlineRegexp.FindStringSubmatch(line); m != nil {
		return token{"deadline", len(m[1]), "", m}, true
	}
	return nilToken, false
}

func (d *Document) parseSchedule(i int, parentStop stopFn) (int, Node) {
	tp := TimeProperty{}
	dstring := ""
	if d.tokens[i].kind == "schedule" {
		if m := deadlineRegexp.FindStringSubmatch(d.tokens[i].content); m != nil {
			dstring = m[1]
		}
		if n := timeDurationRegexp.FindStringSubmatch(d.tokens[i].matches[1]); n != nil {
			fmt.Println(n)
			tp.Schedule = Schedule{fmt.Sprintf("%s", n[1]), strings.TrimSpace(n[4]), strings.TrimSpace(n[5]), strings.TrimSpace(n[6])}
		}
	}
	if d.tokens[i].kind == "deadline" {
		dstring = d.tokens[i].matches[1]
	}

	if dstring != "" {
		if n := timeNotDurationRegexp.FindStringSubmatch(dstring); n != nil {
			fmt.Println(n)
			tp.Deadline = Deadline{fmt.Sprintf("%s", n[1])}
		}
	}
	return 1, tp
}

func (n Schedule) String() string     { return orgWriter.WriteNodesAsString(n) }
func (n Deadline) String() string     { return orgWriter.WriteNodesAsString(n) }
func (n TimeProperty) String() string { return orgWriter.WriteNodesAsString(n) }
