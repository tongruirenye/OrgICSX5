package org

import (
	"regexp"
)

type Planning struct {
	Schedule Node
	Deadline Node
}

var scheduleRegexp = regexp.MustCompile(`SCHEDULED:\s*(<[^>]+>)`)
var deadlineRegexp = regexp.MustCompile(`DEADLINE:\s*(<[^>]+>)`)

func lexSchedule(line string) (token, bool) {
	if m := scheduleRegexp.FindStringSubmatch(line); m != nil {
		return token{"schedule", len(m[1]), line, m}, true
	} else if m := deadlineRegexp.FindStringSubmatch(line); m != nil {
		return token{"deadline", len(m[1]), "", m}, true
	}
	return nilToken, false
}

func (d *Document) parseSchedule(i int, parentStop stopFn) (int, Node) {
	p := Planning{}
	dstring := ""
	if d.tokens[i].kind == "schedule" {
		if m := deadlineRegexp.FindStringSubmatch(d.tokens[i].content); m != nil {
			dstring = m[1]
		}
		if n := d.parseInline(d.tokens[i].matches[1]); n != nil && len(n) == 1 {
			p.Schedule = n[0]
		} else {
			d.Log.Println(n)
		}
	} else if d.tokens[i].kind == "deadline" {
		dstring = d.tokens[i].matches[1]
	}

	if dstring != "" {
		if n := d.parseInline(dstring); n != nil && len(n) == 1 {
			p.Deadline = n[0]
		}

	}
	return 1, p
}

func (n Planning) String() string { return orgWriter.WriteNodesAsString(n) }
