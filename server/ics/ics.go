package ics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/tongruirenye/OrgICSX5/server/models"
	"github.com/tongruirenye/OrgICSX5/server/org"
	"github.com/tongruirenye/OrgICSX5/server/storage"
	"github.com/tongruirenye/OrgICSX5/server/utils"
)

var intervalRegexp = regexp.MustCompile(`^\+(\d+)([dwmy])`)

var timeLocation, _ = time.LoadLocation("Local")

type ical struct {
	writer       strings.Builder
	category     string
	currentTitle string
	log          *log.Logger
}

func (i *ical) write() error {
	icsf, err := os.Create(fmt.Sprintf("%s/public/ics/%s.ics", web.AppPath, i.category))
	if err != nil {
		return err
	}
	icsf.WriteString("BEGIN:VCALENDAR\n")
	icsf.WriteString("VERSION:2.0\n")
	icsf.WriteString(fmt.Sprintf("X-WR-CALNAME:%s\n", i.category))
	icsf.WriteString("PRODID:-////Emacs with Org mode//EN\n")
	icsf.WriteString(fmt.Sprintf("X-WR-CALDESC:%s\n", i.category))
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
	icsf.WriteString(i.writer.String())
	icsf.WriteString("END:VCALENDAR")
	icsf.Close()
	return nil
}

func (i *ical) writeEvent(id, summary, description, categories, createTime, startTime, endTime, untilTime, triggerTime string) {
	i.writer.WriteString("BEGIN:VEVENT\n")
	i.writer.WriteString(fmt.Sprintf("UID:%s\n", id))
	i.writer.WriteString(fmt.Sprintf("DTSTAMP:%s\n", createTime))
	i.writer.WriteString(fmt.Sprintf("DTSTART;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", startTime))
	i.writer.WriteString(fmt.Sprintf("DTEND;TZID=Asia/Shanghai;VALUE=DATE-TIME:%s\n", endTime))
	if untilTime != "" {
		i.writer.WriteString(untilTime)
	}
	i.writer.WriteString(fmt.Sprintf("SUMMARY:[%s]%s\n", i.currentTitle, summary))
	if description != "" {
		i.writer.WriteString(fmt.Sprintf("DESCRIPTION:%s\n", description))
	}
	if categories != "" {
		i.writer.WriteString(fmt.Sprintf("CATEGORIES:%s\n", categories))
	}

	if triggerTime != "" {
		i.writer.WriteString("BEGIN:VALARM\n")
		i.writer.WriteString(fmt.Sprintf("TRIGGER:-PT%sM\n", triggerTime))
		i.writer.WriteString("ACTION:DISPLAY\n")
		i.writer.WriteString("DESCRIPTION:Alarm\n")
		i.writer.WriteString("END:VALARM\n")
	}
	i.writer.WriteString("END:VEVENT\n")
}

func (i *ical) parseSection(section *org.Section) {
	if section.Headline != nil {
		i.parseEvent(section.Headline, section)
	}
	if section.Children != nil {
		for _, child := range section.Children {
			i.parseSection(child)
		}
	}
}

func (i *ical) parseEvent(h *org.Headline, s *org.Section) {
	if h.Status == "TODO" || h.Status == "DONE" || h.Status == "INPROGRESS" || h.Status == "NEXT" {
		if h.Time == nil {
			i.log.Println("no time")
			return
		}

		timeProperty, _ := h.Time.(org.TimeProperty)
		if timeProperty.Schedule == nil {
			return
		}

		var (
			err    error
			start  string
			end    string
			create string
			utill  string
		)

		schedule, _ := timeProperty.Schedule.(org.Schedule)
		start, err = getTime(schedule.Time, schedule.Start)
		if err != nil {
			i.log.Println(err)
			return
		}

		end, err = getTime(schedule.Time, schedule.End)
		if err != nil {
			i.log.Println(err)
			return

		}

		if schedule.Interval != "" {
			if deadline, ok := timeProperty.Deadline.(org.Deadline); ok {
				utill, err = getTimeWithInterval(deadline.Time, schedule.Interval)
				if err != nil {
					i.log.Println(err)
					return
				}
			}
		}

		create = time.Now().In(timeLocation).Format("20060102T150400")

		tags := ""
		if h.Tags != nil {
			for i, tag := range h.Tags {
				tags += tag
				if i != len(h.Tags)-1 {
					tags = tags + ","
				}
			}
		}

		hs := ""
		for _, node := range h.Title {
			hs = hs + node.String()
		}

		description, _ := h.Properties.Get("DESCRIPTION")
		if description == "" {
			proj := getProjTodo(s)
			if proj != nil {
				ps := ""
				for _, node := range proj.Headline.Title {
					ps = ps + node.String()
				}
				description = fmt.Sprintf("[%s][%s] %s", ps, hs, tags)
			} else {
				description = fmt.Sprintf("[%s] %s", hs, tags)
			}
		}
		summary, _ := h.Properties.Get("SUMMARY")
		if summary == "" {
			summary = hs
		}
		uid, _ := h.Properties.Get("ID")
		if uid == "" {
			uid = utils.Md5(hs)
		}
		trigger, _ := h.Properties.Get("VTRIGGER")
		if trigger == "" {
			trigger = "10"
		}

		i.writeEvent(uid, summary, description, tags, create, start, end, utill, trigger)
	}
}

type icsFile struct {
	config *org.Configuration
	doc    *org.Document
	cal    []*ical
}

func (f *icsFile) getCal(name string) *ical {
	if f.cal != nil {
		for _, v := range f.cal {
			if v.category == name {
				return v
			}
		}
	}

	ni := &ical{
		category: name,
	}
	f.cal = append(f.cal, ni)
	return ni
}

func (f *icsFile) writeCal() {
	if f.cal == nil {
		return
	}
	for _, v := range f.cal {
		if err := v.write(); err != nil {
			f.config.Log.Println(err)
		}
	}
}

func GenIcsTask(ctx context.Context) error {
	logger := logs.GetLogger("log")
	logger.Println("start parse ics org file....")

	subList, _ := models.SubGetList()
	if subList == nil || len(subList) == 0 {
		logger.Println("no sub ics org file")
		logger.Println("end parse ics org file....")
		return nil
	}

	icsf := &icsFile{
		config: &org.Configuration{
			AutoLink:            true,
			MaxEmphasisNewLines: 1,
			DefaultSettings: map[string]string{
				"TODO":         "PROJ | TODO | NEXT | INPROGRESS | DONE",
				"EXCLUDE_TAGS": "noexport",
			},
			Log:      logger,
			ReadFile: ioutil.ReadFile,
		},
	}

	for _, sub := range subList {
		logger.Printf("parse file:%s\n", sub.Name)
		f, err := storage.DefaultStorage.ReadFile("org/roam/project/" + sub.Name)
		if err != nil {
			logger.Println(err)
			continue
		}
		icsf.gen(f)
	}

	icsf.writeCal()
	logger.Println("end parse ics org file....")
	return nil
}

func (f *icsFile) gen(b []byte) {
	freader := bytes.NewReader(b)
	f.doc = f.config.Parse(freader, "")
	if f.doc.Error == nil {
		if v, ok := f.doc.BufferSettings["CATEGORY"]; ok {
			cal := f.getCal(v)
			cal.currentTitle = f.doc.BufferSettings["TITLE"]
			cal.log = f.config.Log
			cal.parseSection(f.doc.Outline.Section)
		}
	}
}

func getTime(ddmmyy, hhmm string) (string, error) {
	if hhmm == "" {
		hhmm = "00:00"
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:00", fmt.Sprintf("%s %s:00", ddmmyy, hhmm), timeLocation)

	if err != nil {
		return "", err
	}

	return t.In(timeLocation).Format("20060102T150400"), nil

	return "", errors.New(fmt.Sprintf("invalid time format:%s %s", ddmmyy, hhmm))
}

func getTimeWithInterval(ddmmyy, interval string) (string, error) {
	if m := intervalRegexp.FindStringSubmatch(interval); m != nil {
		num, freq := m[1], m[2]
		if freq == "d" {
			t, err := time.ParseInLocation("2006-01-02 15:04:00", fmt.Sprintf("%s 23:00:00", ddmmyy), timeLocation)
			if err != nil {
				return "", err
			}

			ts := t.In(timeLocation).Format("20060102T150400")

			return fmt.Sprintf("RRULE:FREQ=DAILY;WKST=MO;UNTIL=%s;INTERVAL=%s\n", ts, num), nil
		}
	}
	return "", errors.New(fmt.Sprintf("invalid time format:%s %s", ddmmyy, interval))
}

func getProjTodo(s *org.Section) *org.Section {
	if s.Headline != nil && s.Headline.Status == "PROJ" {
		return s
	}

	if s.Parent != nil {
		return getProjTodo(s.Parent)
	}
	return nil
}
