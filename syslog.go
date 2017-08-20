// Copyright 2017 Szakszon Péter. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syslog generates syslog messages.
package syslog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// The Priority is a combination of the syslog facility and
// severity. For example, USER | NOTICE.
type Priority int

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	EMERG Priority = iota
	ALERT
	CRIT
	ERR
	WARNING
	NOTICE
	INFO
	DEBUG
)

const (
	// Facility.

	// From /usr/include/sys/syslog.h.
	// These are the same up to FTP on Linux, BSD, and OS X.
	KERN Priority = iota << 3
	USER
	MAIL
	DAEMON
	AUTH
	SYSLOG
	LPR
	NEWS
	UUCP
	CRON
	AUTHPRIV
	FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOCAL0
	LOCAL1
	LOCAL2
	LOCAL3
	LOCAL4
	LOCAL5
	LOCAL6
	LOCAL7
)

const version = 1 // defined in RFC 5424.

// Writer generates syslog messages as defined in RFC 5424.
type writer struct {
	out io.Writer
	pri Priority
}

// NewWriter wrappes another io.Writer and returns a new
// io.Writer that generates syslog messages as defined
// in RFC 5424 and writes them to the given io.Writer.
func NewWriter(out io.Writer, pri Priority) io.Writer {
	if pri < 0 || pri > LOCAL7|DEBUG {
		panic("syslog: invalid priority: " + strconv.Itoa(int(pri)))
	}

	return &writer{
		out: out,
		pri: pri,
	}
}

// Write generates and writes a syslog message to the
// underlying io.Writer.
func (w *writer) Write(d []byte) (n int, err error) {
	if len(d) == 0 {
		return 0, nil
	}

	if d[0] != '<' {
		return w.out.Write(w.format(d))
	}

	// don't format a syslog message
	return w.out.Write(d)
}

const rfc3339Milli = "2006-01-02T15:04:05.999-07:00"

func (w *writer) format(d []byte) []byte {
	timestamp := time.Now().Format(rfc3339Milli)
	hostname, _ := os.Hostname()
	appName := os.Args[0]
	procid := os.Getpid()

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "<%d>%d %s %s %s %d - - ",
		w.pri,
		version,
		timestamp,
		hostname,
		appName,
		procid,
	)
	buf.Write(d)

	if d[len(d)-1] != '\n' {
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

type Logger interface {
	Log(sev Priority, msgId string, sd StructuredData, format string, a ...interface{})
}

func NewLogger(w io.Writer) Logger {
	return &logger{w}
}

type logger struct {
	w io.Writer
}

func (l *logger) Log(sev Priority, msgId string, sd StructuredData, format string, a ...interface{}) {

}

type StructuredData map[string]SDElement

func (d StructuredData) Element(id string) SDElement {
	elem, ok := d[id]
	if !ok {
		elem = make(SDElement, 1)
		d[id] = elem
	}
	return elem
}

func (d StructuredData) Ids() []string {
	ids := make([]string, 0, len(d))
	for id := range d {
		if len(d[id]) > 0 {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (d StructuredData) String() string {
	r := strings.NewReplacer(`"`, `\"`, `\`, `\\`, `]`, `\]`)
	buf := &bytes.Buffer{}
	for _, id := range d.Ids() {
		elem := d[id]
		if len(elem) > 0 {
			buf.WriteByte('[')
			buf.WriteString(id)
			for _, name := range elem.Keys() {
				buf.WriteByte(' ')
				fmt.Fprintf(buf, `%s="%s"`, name, r.Replace(elem[name]))
			}
			buf.WriteByte(']')
		}
	}
	return buf.String()
}

type SDElement map[string]string

func (e SDElement) Set(name, value string) SDElement {
	e[name] = value
	return e
}

func (e SDElement) Get(name string) string {
	value, ok := e[name]
	if !ok {
		return ""
	}
	return value
}

func (e SDElement) Keys() []string {
	keys := make([]string, 0, len(e))
	for key := range e {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func Alert(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	l.Log(ALERT, msgId, sd, format, a...)
}

func Error(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	l.Log(ERR, msgId, sd, format, a...)
}

func Info(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	l.Log(INFO, msgId, sd, format, a...)
}

func Debug(l Logger, msgId string, sd StructuredData, format string, a ...interface{}) {
	l.Log(DEBUG, msgId, sd, format, a...)
}
