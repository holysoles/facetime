package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const newFaceTimeLinkScript = "./lib/createLink.applscript"
const getActiveFaceTimeLinksScript = "./lib/getLinks.applescript"
const joinLatestFaceTimeLinkScript = "./lib/joinFirstLinkApproveAll.applescript"
const deleteFaceTimeLinkScript = "./lib/deleteLink.applescript"

func makeNewLink() string {
	newLScript := loadAppleScript(newFaceTimeLinkScript)
	newLink, err := execAppleScript(newLScript)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return newLink
}

func getAllLinks() []string {
	getLScript := loadAppleScript(getActiveFaceTimeLinksScript)
	allLinksRaw, err := execAppleScript(getLScript)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	allLinks := strings.Split(allLinksRaw, ", ")
	return allLinks
}

func joinLink(l string) {
	joinScript := loadAppleScript(joinLatestFaceTimeLinkScript)
	_, err := execAppleScript(joinScript)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

func deleteLink(l string) bool {
	deleteScript := loadAppleScript(deleteFaceTimeLinkScript)
	deleteScript = fmt.Sprintf(deleteScript, l)
	deletedStr, err := execAppleScript(deleteScript)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	deleted, _ := strconv.ParseBool(deletedStr)
	return deleted
}

func loadAppleScript(p string) string {
	fB, err := os.ReadFile(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(fB)
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
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Stdout = &outBuff
	cmd.Stderr = &errorBuff

	cmd.Start()
	io.WriteString(stdin, s)
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		fmt.Println(errorBuff.String())
		os.Exit(1)
	}

	//osaScript output has a tailing newline, making any later parse logic difficult
	var re = regexp.MustCompile(`\n$`)
	parsedOutput := re.ReplaceAllString(outBuff.String(), "")

	return parsedOutput, nil
}
