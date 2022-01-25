package ics

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/tongruirenye/OrgICSX5/server/config"
	"github.com/tongruirenye/OrgICSX5/server/storage"
	"github.com/tongruirenye/go-org/org"
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
	writer *IcsWriter
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
				"TODO":         "TODO | SCHED | DOING | WAITING",
				"EXCLUDE_TAGS": "noexport",
			},
			Log:      l,
			ReadFile: ioutil.ReadFile,
		},
		writer:    new(IcsWriter),
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

func (ics *ICS) clearWriter() {
	if ics.writer == nil {
		return
	}

	ics.logger.Println("start parse")
	ics.writer.ClearContent()
}

func (ics *ICS) Do() {
	// files, _ := storage.AppStorage.ListFileList("org/roam")
	// if files == nil || len(files) == 0 {
	// 	return
	// }

	ics.clearWriter()
	if err := ics.parse(config.AppConfig.Project); err != nil {
		ics.logger.Println(err.Error())
	}
	ics.dump()
	ics.logger.Println("end parse")
}

func (ics *ICS) parse(fileName string) error {
	f, err := storage.AppStorage.ReadFile("org/roam/" + fileName)
	if err != nil {
		return err
	}
	freader := bytes.NewReader(f)
	doc := ics.config.Parse(freader, "")
	doc.Write(ics.writer)
	return nil
}

func (ics *ICS) dump() error {
	if ics.writer == nil {
		return nil
	}

	return ics.writer.Dump()
}
