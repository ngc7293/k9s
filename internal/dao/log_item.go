// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"bytes"
	"time"
)

// LogChan represents a channel for logs.
type LogChan chan *LogItem

var ItemEOF = new(LogItem)

// LogItem represents a container log line.
type LogItem struct {
	Pod, Container  string
	SingleContainer bool
	Bytes           []byte
	IsError         bool
}

// The time.RFC3339Nano preset removes trailing zeroes, which breaks the alignment
// of the log lines on screen.
const RFC3339NanoTrailingZeroes = "2006-01-02T15:04:05.000000000Z07:00"

func ConvertTimeZone(ts string, tz *time.Location) string {
	if tz == nil {
		return ts
	}

	dt, err := time.Parse(time.RFC3339Nano, string(ts[:]))

	if err != nil {
		return ts
	}

	return dt.In(tz).Format(RFC3339NanoTrailingZeroes)
}

// NewLogItem returns a new item.
func NewLogItem(bb []byte) *LogItem {
	return &LogItem{
		Bytes: bb,
	}
}

// NewLogItemFromString returns a new item.
func NewLogItemFromString(s string) *LogItem {
	return &LogItem{
		Bytes: []byte(s),
	}
}

// ID returns pod and or container based id.
func (l *LogItem) ID() string {
	if l.Pod != "" {
		return l.Pod
	}
	return l.Container
}

// GetTimestamp fetch log lime timestamp
func (l *LogItem) GetTimestamp() string {
	index := bytes.Index(l.Bytes, []byte{' '})
	if index < 0 {
		return ""
	}
	return string(l.Bytes[:index])
}

// Info returns pod and container information.
func (l *LogItem) Info() string {
	return l.Pod + "::" + l.Container
}

// IsEmpty checks if the entry is empty.
func (l *LogItem) IsEmpty() bool {
	return len(l.Bytes) == 0
}

// Size returns the size of the item.
func (l *LogItem) Size() int {
	return 100 + len(l.Bytes) + len(l.Pod) + len(l.Container)
}

// Render returns a log line as string.
func (l *LogItem) Render(paint string, showTime bool, tz *time.Location, bb *bytes.Buffer) {
	index := bytes.Index(l.Bytes, []byte{' '})
	if showTime && index > 0 {
		bb.WriteString("[gray::b]")
		bb.WriteString(ConvertTimeZone(string(l.Bytes[:index]), tz))
		bb.WriteString(" ")
		if l := 30 - len(l.Bytes[:index]); l > 0 {
			bb.Write(bytes.Repeat([]byte{' '}, l))
		}
		bb.WriteString("[-::-]")
	}

	if l.Pod != "" {
		bb.WriteString("[" + paint + "::]" + l.Pod)
	}

	if !l.SingleContainer && l.Container != "" {
		if len(l.Pod) > 0 {
			bb.WriteString(" ")
		}
		bb.WriteString("[" + paint + "::b]" + l.Container + "[-::-] ")
	} else if len(l.Pod) > 0 {
		bb.WriteString("[-::] ")
	}

	if index > 0 {
		bb.Write(l.Bytes[index+1:])
	} else {
		bb.Write(l.Bytes)
	}
}
