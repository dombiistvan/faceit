package helpers

import (
	"faceit/common"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	errorLogFile        = "logs/errors.log"
	panicLogFile        = "logs/panics.log"
	notificationLogFile = "logs/notifications.log"
)

var (
	errorStream, notificationStream *os.File = nil, nil
	ErrChan                                  = make(chan error)
	NotifyChan                               = make(chan string)
	sigs                                     = make(chan os.Signal, 1)
)

// set up listeners to channels and make files to stream logs into
func init() {
	var err error
	listenChannels()

	makeDirs(panicLogFile, errorLogFile, notificationLogFile)

	errorStream, err = os.OpenFile(errorLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		ErrChan <- err
		return
	}

	notificationStream, err = os.OpenFile(notificationLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		ErrChan <- err
		return
	}
}

// make directories
func makeDirs(files ...string) {
	for _, file := range files {
		if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
			panic(err)
		}
	}
}

// handle panic
func HandlePanic(errGot any) {
	pFile, err := os.OpenFile(panicLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer pFile.Close()

	if _, err = pFile.WriteString(fmt.Sprintf("%s - panic captured: %+v\r\n", time.Now().Format("2006-01-02 15:04:05"), errGot)); err != nil {
		panic(err)
	}
}

// start listening on channels
func listenChannels() {
	go func() {
		for {
			if err := <-ErrChan; err != nil {
				if _, err = fmt.Fprintf(errorStream, "%s - error captured: %s\r\n", time.Now().Format("2006-01-02 15:04:05"), err.Error()); err != nil {
					panic(err)
				}
			}
		}
	}()

	go func() {
		for {
			if notification := <-NotifyChan; !common.IsEmptyString(notification) {
				if _, err := fmt.Fprintf(notificationStream, "%s - nofication: %s\r\n", time.Now().Format("2006-01-02 15:04:05"), notification); err != nil {
					ErrChan <- err
				}
			}
		}
	}()

	go func() {
		sig := <-sigs
		fmt.Printf("app manually shut down with signal: %s!\r\n", sig.String())
		os.Exit(0)
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
}
