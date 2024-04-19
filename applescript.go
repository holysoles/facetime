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

const openFacetimeScript = "./lib/openFacetime.applescript"
const checkFacetimeScript = "./lib/checkFacetime.applescript"
const newLinkScript = "./lib/createLink.applescript"
const getActiveLinksScript = "./lib/getLinks.applescript"
const joinLatestLinkScript = "./lib/joinFirstLink.applescript"
const approveJoinScript = "./lib/approveJoin.applescript"
const joinLatestLinkAndApproveScript = "./lib/joinFirstLinkApproveJoin.applescript"
const deleteLinkScript = "./lib/deleteLink.applescript"

func openFacetime() error {
	oFtScript, err := loadAppleScript(openFacetimeScript)
	if err != nil {
		return err
	}
	_, err = execAppleScript(oFtScript)
	return err
}

func getFacetimeStatus() (bool, error) {
	cFtScript, err := loadAppleScript(checkFacetimeScript)
	if err != nil {
		return false, err
	}
	s, err := execAppleScript(cFtScript)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(s)
	return b, err
}

func makeNewLink() (string, error) {
	newLScript, err := loadAppleScript(newLinkScript)
	if err != nil {
		return "", err
	}
	newLink, err := execAppleScript(newLScript)
	return newLink, err
}

func getAllLinks() ([]string, error) {
	getLScript, err := loadAppleScript(getActiveLinksScript)
	if err != nil {
		return make([]string, 0), nil
	}
	allLinksRaw, err := execAppleScript(getLScript)
	var allLinks []string
	if allLinksRaw != "" {
		allLinks = strings.Split(allLinksRaw, ", ")
	}
	return allLinks, err
}

// Join the latest Facetime Call.
func joinCall() error {
	joinScript, err := loadAppleScript(joinLatestLinkScript)
	if err != nil {
		return err
	}
	_, err = execAppleScript(joinScript)
	if err != nil {
		return err
	}
	return nil
}

// Join the latest Facetime Call and approve all requested entrants. Combining these two actions allows us to leverage the sidebar being automatically in focus.
func joinAndAdmitCall() error {
	joinScript, err := loadAppleScript(joinLatestLinkAndApproveScript)
	if err != nil {
		return err
	}
	_, err = execAppleScript(joinScript)
	if err != nil {
		return err
	}
	return nil
}

func admitActiveCall() error {
	joinScript, err := loadAppleScript(approveJoinScript)
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
	deleteScript, err := loadAppleScript(deleteLinkScript)
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
