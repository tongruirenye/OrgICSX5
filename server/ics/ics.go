package ics

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/niklasfasching/go-org/org"
	"github.com/tongruirenye/OrgICSX5/server/models"
	"github.com/tongruirenye/OrgICSX5/server/storage"
)

var timestampRegexp = regexp.MustCompile(`^<(\d{4}-\d{2}-\d{2})( [A-Za-z]+)?( \d{2}:\d{2})?( \+\d+[dwmy])?>`)

var timeLocation, _ = time.LoadLocation("GMT")

type icsFile struct {
	config *org.Configuration
	writer strings.Builder
}

func GenIcs(ctx context.Context) {
	subList, _ = models.SubGetList()
	if subList != nil {
		icsf := &icsFile{
			config: &org.Configuration{
				AutoLink:            true,
				MaxEmphasisNewLines: 1,
				DefaultSettings: map[string]string{
					"TODO":         "TODO | NEXT | INPROGRESS | DONE",
					"EXCLUDE_TAGS": "noexport",
				},
				Log:      log.New(os.Stderr, "ics-gen ", 0),
				ReadFile: ioutil.ReadFile,
			},
		}

		for _, sub := range subList {
			f, err := storage.DefaultStorage.ReadFile("org/roam/project/" + sub.Name)
			if err == nil {
				icsf.gen(f)
			}
		}
	}
}

func (f *icsFile) gen(f []byte) {
	freader := bytes.NewReader(f)
	doc := f.config.Parse(freader, "")
	if doc.Error == nil {
		f.writeSection(doc.Outline.Section)
	}
}

func getTime(timeStr string) string {
	if m := timestampRegexp.FindStringSubmatch(timeStr); m != nil {
		ddmmyy, hhmm := m[1], m[3]
		if hhmm == "" {
			hhmm = "00:00"
		}
		t, err := time.Parse("2021-06-27 20:00", fmt.Sprintf("%s %s", ddmmyy, hhmm))

		if err != nil {
			return ""
		}

		return t.In(timeLocation).Format("20060102T000000Z")
	}

	return ""
}

func (f *icsFile) writeSection(section *org.Section) {
	if section.Headline != nil {
		if section.Headline.Status == "TODO" || section.Headline.Status == "DONE" || section.Headline.Status == "INPROGRESS" || section.Headline.Status == "NEXT" {
			start, sok := section.Headline.Properties.Get("VSTART")
			end, eok := section.Headline.Properties.Get("VEND")
			create, cok := section.Headline.Properties.Get("VCREATE")
			description, _ := section.Headline.Properties.Get("DESCRIPTION")
			summary, _ := section.Headline.Properties.Get("SUMMARY")
			uid, _ := section.Headline.Properties.Get("ID")
			trigger, _ := section.Headline.Properties.Get("VTRIGGER")
			tags := ""
			if section.Headline.Tags != nil {
				for i, tag := range section.Headline.Tags {
					tags += tag
					if i != len(section.Headline.Tags)-1 {
						tags = tags + ","
					}
				}
			}
			if sok && eok && cok {
				f.writeEvent(summary, description, tags, getTime(create), getTime(start), getTime(end), trigger, uid)
			}
		}
	}
	if section.Children != nil {
		for _, child := range section.Children {
			f.writeSection(child)
		}
	}
}

func (f *icsFile) writeEvent(summary, description, categories, create, start, end, trigger, uid string) {
	f.writer.WriteString("BEGIN:VEVENT\n")
	f.writer.WriteString(fmt.Sprintf("UID:%s\n", uid))
	f.writer.WriteString(fmt.Sprintf("DTSTAMP:%s\n", create))
	f.writer.WriteString(fmt.Sprintf("DTSTART:%s\n", start))
	f.writer.WriteString(fmt.Sprintf("DTEND:%s\n", end))
	f.writer.WriteString(fmt.Sprintf("SUMMARY:%s\n", summary))
	f.writer.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", description))
	f.writer.WriteString(fmt.Sprintf("CATEGORIES:%s\n", categories))
	if trigger != "" {
		f.writer.WriteString("BEGIN:VALARM\n")
		f.writer.WriteString(fmt.Sprintf("TRIGGER:-PT%sM\n", trigger))
		f.writer.WriteString("ACTION:DISPLAY\n")
		f.writer.WriteString("DESCRIPTION:Reminder\n")
		f.writer.WriteString("END:VALARM\n")
	}
	f.writer.WriteString("END:VEVENT\n")
}
