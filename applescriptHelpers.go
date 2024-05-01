package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func openFacetime() error {
	_, err := execAppleScript(openFacetimeScript)
	return err
}

func getFacetimeStatus() (bool, error) {
	s, err := execAppleScript(checkFacetimeScript)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(s)
	return b, err
}

func makeNewLink() (string, error) {
	newLink, err := execAppleScript(newLinkScript)
	return newLink, err
}

func getAllLinks() ([]string, error) {
	allLinksRaw, err := execAppleScript(getActiveLinksScript)
	var allLinks []string
	if allLinksRaw != "" {
		allLinks = strings.Split(allLinksRaw, ", ")
	}
	return allLinks, err
}

// Join the latest Facetime Call.
func joinCall() error {
	_, err := execAppleScript(joinLatestLinkScript)
	if err != nil {
		return err
	}
	return nil
}

// Join the latest Facetime Call and approve all requested entrants. Combining these two actions allows us to leverage the sidebar being automatically in focus.
func joinAndAdmitCall() error {
	_, err := execAppleScript(joinLatestLinkAndApproveScript)
	if err != nil {
		return err
	}
	return nil
}

func admitActiveCall() error {
	_, err := execAppleScript(approveJoinScript)
	if err != nil {
		return err
	}
	return nil
}

func leaveCall() error {
	_, err := execAppleScript(leaveCallScript)
	if err != nil {
		return err
	}
	return nil
}

func deleteCall(id string) (bool, error) {
	deleteScript := fmt.Sprintf(deleteLinkScript, id)
	deletedStr, err := execAppleScript(deleteScript)
	if err != nil {
		return false, err
	}
	deleted, _ := strconv.ParseBool(deletedStr)
	return deleted, nil
}

func execAppleScript(s string) (string, error) {
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
