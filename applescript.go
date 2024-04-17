package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const newFaceTimeLinkScript = "./lib/createLink.applescript"
const getActiveFaceTimeLinksScript = "./lib/getLinks.applescript"
const joinLatestFaceTimeLinkScript = "./lib/joinFirstLinkApproveAll.applescript"
const deleteFaceTimeLinkScript = "./lib/deleteLink.applescript"

func makeNewLink() (string, error) {
	newLScript, err := loadAppleScript(newFaceTimeLinkScript)
	if err != nil {
		return "", nil
	}
	newLink, err := execAppleScript(newLScript)
	return newLink, err
}

func getAllLinks() ([]string, error) {
	getLScript, err := loadAppleScript(getActiveFaceTimeLinksScript)
	if err != nil {
		return make([]string, 0), nil
	}
	allLinksRaw, err := execAppleScript(getLScript)
	allLinks := strings.Split(allLinksRaw, ", ")
	return allLinks, err
}

// TODO
func joinCall(l string) error {
	joinScript, err := loadAppleScript(joinLatestFaceTimeLinkScript)
	if err != nil {
		return err
	}
	_, err = execAppleScript(joinScript)
	if err != nil {
		return err
	}
	return nil
}

func deleteCall(id string) (bool, error) {
	deleteScript, err := loadAppleScript(deleteFaceTimeLinkScript)
	if err != nil {
		return false, err
	}
	deleteScript = fmt.Sprintf(deleteScript, id)
	deletedStr, err := execAppleScript(deleteScript)
	if err != nil {
		return false, err
	}
	deleted, _ := strconv.ParseBool(deletedStr)
	return deleted, nil
}

func loadAppleScript(p string) (string, error) {
	fB, err := os.ReadFile(p)
	return string(fB), err
}

func execAppleScript(s string) (string, error) {
	//TODO have a safe FT process init check
	args := "-"
	cmd := exec.Command("osascript", args)

	var stdin io.WriteCloser
	var outBuff, errorBuff bytes.Buffer
	var err error
	stdin, err = cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	cmd.Stdout = &outBuff
	cmd.Stderr = &errorBuff

	cmd.Start()
	io.WriteString(stdin, s)
	stdin.Close()

	err = cmd.Wait()
	// prefer returning errors from script execution
	if errorBuff.Len() != 0 {
		return "", errors.New(errorBuff.String())
	}
	if err != nil {
		return "", err
	}

	//osaScript output has a tailing newline, making any later parse logic difficult
	var re = regexp.MustCompile(`\n$`)
	parsedOutput := re.ReplaceAllString(outBuff.String(), "")

	return parsedOutput, nil
}
