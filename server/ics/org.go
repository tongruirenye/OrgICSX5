package ics

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tongruirenye/OrgICSX5/server/config"
	"github.com/tongruirenye/OrgICSX5/server/org"
	"github.com/tongruirenye/OrgICSX5/server/utils"
)

var timestampRegexp = regexp.MustCompile(`^<(\d{4}-\d{2}-\d{2})( [A-Za-z\p{Han}]+)?( \d{2}:\d{2})?(?:-(\d{2}:\d{2}))?( \+\d+[dwmy])?>`)
var timestampRegexp1 = regexp.MustCompile(`--<(\d{4}-\d{2}-\d{2})( [A-Za-z\p{Han}]+)?( \d{2}:\d{2})?(?:-(\d{2}:\d{2}))?( \+\d+[dwmy])?>`)

type IcsWriter struct {
	ExtendingWriter org.Writer
	strings.Builder
	document *org.Document
	log      *log.Logger
}

func NewIcsWriter() *IcsWriter {
	return &IcsWriter{}
}

func (w *IcsWriter) ClearContent() {
	w.Builder.Reset()
}

func (w *IcsWriter) Dump(calName string) error {
	icsf, err := os.Create(fmt.Sprintf("%s/public/ics/%s.ics", config.AppPath, calName))
	if err != nil {
		return err
	}

	icsf.WriteString("BEGIN:VCALENDAR\n")
	icsf.WriteString("VERSION:2.0\n")
	icsf.WriteString(fmt.Sprintf("X-WR-CALNAME:%s\n", calName))
	icsf.WriteString("PRODID:-////Emacs with Org mode//EN\n")
	icsf.WriteString(fmt.Sprintf("X-WR-CALDESC:%s\n", calName))
	icsf.WriteString("CALSCALE:GREGORIAN\n")
	icsf.WriteString("BEGIN:VTIMEZONE\n")
	icsf.WriteString("TZID:Asia/Shanghai\n")
	icsf.WriteString("TZURL:http://tzurl.org/zoneinfo/Asia/Shanghai\n")
	icsf.WriteString("X-LIC-LOCATION:Asia/Shanghai\n")
	icsf.WriteString("BEGIN:STANDARD\n")
	icsf.WriteString("TZOFFSETFROM:+0800\n")
	icsf.WriteString("TZOFFSETTO:+0800\n")
	icsf.WriteString("TZNAME:CST\n")
	icsf.WriteString("DTSTART:19490101T0000\n")
	icsf.WriteString("RDATE:19490101T000000\n")
	icsf.WriteString("END:STANDARD\n")
	icsf.WriteString("END:VTIMEZONE\n")
	icsf.WriteString(w.String())
	icsf.WriteString("END:VCALENDAR")
	return icsf.Close()
}

func (w *IcsWriter) Before(d *org.Document) {
	w.document = d
	w.log = d.Log
	w.log.Printf("nodes size is:%d\n", len(d.Nodes))
}

func (w *IcsWriter) After(d *org.Document) {
	dateStr, ok := w.document.BufferSettings["DATE"]
	if !ok {
		return
	}
	m := timestampRegexp.FindStringSubmatch(dateStr)
	if m == nil {
		return
	}

	t, err := time.Parse("2006-01-02 Mon 15:04", fmt.Sprintf("%s Mon %s", m[1], "00:00"))
	if err != nil {
		return
	}

	var tend *time.Time
	if m1 := timestampRegexp1.FindStringSubmatch(dateStr); m1 != nil {
		if te, e := time.Parse("2006-01-02 Mon 15:04", fmt.Sprintf("%s Mon %s", m1[1], "00:00")); e == nil {
			tend = &te
		}
	}

	start := t
	end := t
	if tend != nil {
		end = *tend
	} else {
		end = t.AddDate(0, 0, 1)
	}

	w.WriteString("BEGIN:VEVENT\n")
	w.WriteString(fmt.Sprintf("UID:%s\n", utils.Md5(time.Now().String())))
	w.WriteString(fmt.Sprintf("DTSTAMP:%s\n", time.Now().Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("DTSTART;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", start.Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("DTEND;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", end.Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("SUMMARY:%s\n", w.document.BufferSettings["TITLE"]))
	w.WriteString("END:VEVENT\n")
}

func (w *IcsWriter) WriterWithExtensions() org.Writer {
	return w
}

func (w *IcsWriter) WriteNodesAsString(nodes ...org.Node) string {
	org.WriteNodes(w, nodes...)
	return ""
}

func (w *IcsWriter) WriteKeyword(k org.Keyword) {

}
func (w *IcsWriter) WriteInclude(i org.Include) {

}
func (w *IcsWriter) WriteComment(c org.Comment) {

}
func (w *IcsWriter) WriteNodeWithMeta(n org.NodeWithMeta) {

}
func (w *IcsWriter) WriteNodeWithName(n org.NodeWithName) {

}
func (w *IcsWriter) WriteHeadline(h org.Headline) {
	if h.Status == "PROJ" || h.Status == "GOAL" {
		org.WriteNodes(w, h.Children...)
		return
	}

	if h.Time == nil {
		org.WriteNodes(w, h.Children...)
		return
	}

	if h.Status != "TODO" && h.Status != "INPROGRESS" {
		org.WriteNodes(w, h.Children...)
		return
	}

	w.log.Println(h.Status, h.Title[0].String())

	isDate := true
	var start, end time.Time
	planning, _ := h.Time.(org.Planning)
	if planning.Schedule != nil && planning.Deadline != nil {
		schedule, _ := planning.Schedule.(org.Timestamp)
		deadline, _ := planning.Deadline.(org.Timestamp)
		start = schedule.Time
		end = deadline.Time
	} else if planning.Schedule != nil {
		schedule, _ := planning.Schedule.(org.Timestamp)
		start = schedule.Time
		end = schedule.Time
		if schedule.EndTime != nil {
			end = *schedule.EndTime
			isDate = false
		} else {
			if schedule.IsDate {
				end = end.AddDate(0, 0, 1)
			} else {
				end = end.Add(1 * time.Hour)
				isDate = false
			}
		}
	} else if planning.Deadline != nil {
		deadline, _ := planning.Deadline.(org.Timestamp)
		start = deadline.Time
		end = deadline.Time
		if deadline.IsDate {
			end = end.AddDate(0, 0, 1)
		} else {
			end = end.Add(1 * time.Hour)
			isDate = false
		}
	} else {
		org.WriteNodes(w, h.Children...)
		return
	}

	description, _ := h.Properties.Get("DESCRIPTION")
	if description == "" {
		description, _ = h.Properties.Get("LOCATION")
	}
	summary, _ := h.Properties.Get("SUMMARY")
	if summary == "" {
		for _, node := range h.Title {
			summary = summary + node.String()
		}
	}
	uid, _ := h.Properties.Get("ID")
	if uid == "" {
		uid = utils.Md5(summary)
	}
	trigger, _ := h.Properties.Get("ALARM")
	if trigger == "" {
		trigger = "10"
	}

	status := "TENTATIVE"
	if h.Status == "TODO" {
		status = "ACCEPTED"
	} else if h.Status == "INPROGRESS" {
		status = "IN-PROCESS"
	} else if h.Status == "NEXT" {
		status = "NEEDS-ACTION"
	} else if h.Status == "DONE" {
		status = "COMPLETED"
	}

	w.WriteString("BEGIN:VEVENT\n")
	w.WriteString(fmt.Sprintf("UID:%s\n", uid))
	w.WriteString(fmt.Sprintf("DTSTAMP:%s\n", time.Now().Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("DTSTART;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", start.Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("DTEND;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", end.Format("20060102T150400")))
	w.WriteString(fmt.Sprintf("SUMMARY:[%s]%s\n", w.document.BufferSettings["TITLE"], summary))
	w.WriteString(fmt.Sprintf("STATUS:%s\n", status))
	if description != "" {
		w.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", description))
	}

	if trigger != "" && !isDate {
		w.WriteString("BEGIN:VALARM\n")
		w.WriteString(fmt.Sprintf("TRIGGER:-PT%sM\n", trigger))
		w.WriteString("ACTION:DISPLAY\n")
		w.WriteString("DESCRIPTION:Alarm\n")
		w.WriteString("END:VALARM\n")
	}
	w.WriteString("END:VEVENT\n")

	org.WriteNodes(w, h.Children...)
}
func (w *IcsWriter) WriteBlock(b org.Block) {

}
func (w *IcsWriter) WriteResult(r org.Result) {

}
func (w *IcsWriter) WriteInlineBlock(i org.InlineBlock) {

}
func (w *IcsWriter) WriteExample(e org.Example) {

}
func (w *IcsWriter) WriteDrawer(d org.Drawer) {

}
func (w *IcsWriter) WritePropertyDrawer(p org.PropertyDrawer) {

}
func (w *IcsWriter) WriteList(l org.List) {

}
func (w *IcsWriter) WriteListItem(l org.ListItem) {

}
func (w *IcsWriter) WriteDescriptiveListItem(d org.DescriptiveListItem) {

}
func (w *IcsWriter) WriteTable(t org.Table) {

}
func (w *IcsWriter) WriteHorizontalRule(h org.HorizontalRule) {

}
func (w *IcsWriter) WriteParagraph(p org.Paragraph) {

}
func (w *IcsWriter) WriteText(t org.Text) {

}
func (w *IcsWriter) WriteEmphasis(e org.Emphasis) {

}
func (w *IcsWriter) WriteLatexFragment(l org.LatexFragment) {

}
func (w *IcsWriter) WriteStatisticToken(s org.StatisticToken) {

}
func (w *IcsWriter) WriteExplicitLineBreak(e org.ExplicitLineBreak) {

}
func (w *IcsWriter) WriteLineBreak(l org.LineBreak) {

}
func (w *IcsWriter) WriteRegularLink(r org.RegularLink) {

}
func (w *IcsWriter) WriteMacro(m org.Macro) {

}
func (w *IcsWriter) WriteTimestamp(t org.Timestamp) {

}
func (w *IcsWriter) WriteFootnoteLink(f org.FootnoteLink) {

}
func (w *IcsWriter) WriteFootnoteDefinition(f org.FootnoteDefinition) {

}

func (w *IcsWriter) WritePlanning(f org.Planning) {

}
