package ics

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tongruirenye/OrgICSX5/server/org"
	"github.com/tongruirenye/OrgICSX5/server/storage"
)

func isSameDay(lhs, rhs time.Time) bool {
	if lhs.Day() != rhs.Day() {
		return false
	}
	if lhs.Month() != rhs.Month() {
		return false
	}
	if lhs.Day() != rhs.Day() {
		return false
	}
	return true
}

func isSmall(dt time.Time, hour, minute int) bool {
	if dt.Hour() < hour {
		return true
	} else if dt.Hour() == hour {
		if dt.Minute() < minute {
			return true
		} else {
			return false
		}
	}
	return false
}

func triggerTimer(lastTime time.Time, hour, minute int) bool {
	now := time.Now()
	if now.Before(lastTime) {
		return false
	}

	if now.Sub(lastTime) >= 24*time.Hour {
		return true
	}

	old_small := isSmall(lastTime, hour, minute)
	new_small := isSmall(now, hour, minute)

	if isSameDay(lastTime, now) {
		if old_small && !new_small {
			return true
		} else {
			return false
		}
	} else {
		if !old_small && !new_small {
			return true
		} else if old_small && new_small {
			return true
		} else {
			return false
		}
	}

	return false
}

type ICS struct {
	config *org.Configuration
	writer map[string]*IcsWriter
	logger *log.Logger

	closeChan chan struct{}
	taskChan  chan struct{}
	cancel    func()
	ctx       context.Context
}

func NewIcs(l *log.Logger) *ICS {
	ics := &ICS{

		config: &org.Configuration{
			AutoLink:            true,
			MaxEmphasisNewLines: 1,
			DefaultSettings: map[string]string{
				"TODO":         "PROJ | TODO | NEXT | INPROGRESS | DONE",
				"EXCLUDE_TAGS": "noexport",
			},
			Log:      l,
			ReadFile: ioutil.ReadFile,
		},
		writer:    make(map[string]*IcsWriter),
		logger:    l,
		closeChan: make(chan struct{}),
		taskChan:  make(chan struct{}, 1),
	}
	ics.ctx, ics.cancel = context.WithCancel(context.Background())
	return ics
}

func (ics *ICS) Run() {
	timer := time.NewTimer(1 * time.Second)
	lastTime := time.Now()
	exist := false
	for !exist {
		select {
		case <-ics.ctx.Done():
			exist = true
		case <-timer.C:
			timer.Reset(1 * time.Second)
			if triggerTimer(lastTime, 3, 0) {
				ics.logger.Println("trigger timer")
				lastTime = time.Now()
				ics.Task()
			}
		case <-ics.taskChan:
			ics.Do()
		}
		if exist {
			break
		}
	}
	ics.closeChan <- struct{}{}
}

func (ics *ICS) Task() {
	select {
	case ics.taskChan <- struct{}{}:
	default:
		break
	}
}

func (ics *ICS) Close() {
	ics.cancel()
	<-ics.closeChan
}

func (ics *ICS) Do() {
	files, _ := storage.AppStorage.ListFileList("org/roam/project")
	if files == nil || len(files) == 0 {
		return
	}

	for _, v := range files {
		ics.logger.Printf("parse file:%s\n", v)
		if err := ics.parse(v); err != nil {
			ics.logger.Println(err.Error())
		}
	}

	ics.dump()
}

func (ics *ICS) DoLocal(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	doc := ics.config.Parse(f, "")
	calName := doc.Get("CATEGORY")
	if calName == "" {
		calName = "其他"
	}
	writer := NewIcsWriter()
	doc.Write(writer)
	return writer.Dump(calName)
}

func (ics *ICS) parse(fileName string) error {
	f, err := storage.AppStorage.ReadFile("org/roam/project/" + fileName)
	if err != nil {
		return err
	}
	freader := bytes.NewReader(f)
	doc := ics.config.Parse(freader, "")
	calName := doc.Get("CATEGORY")
	if calName == "" {
		calName = "其他"
	}
	writer, ok := ics.writer[calName]
	if !ok {
		writer = NewIcsWriter()
		ics.writer[calName] = writer
	}
	doc.Write(writer)
	return nil
}

func (ics *ICS) dump() {
	for k, v := range ics.writer {
		v.Dump(k)
	}
}
