package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ftLink string

// Returns a well-formed URL for the FT link. Can be used for validation.
func (f *ftLink) getUrl() (ftLink, error) {
	fStr := string(*f)
	if !strings.HasPrefix(fStr, "http://") && !strings.HasPrefix(fStr, "https://") {
		fStr = "https://" + fStr
	}
	fUrl, err := url.Parse(fStr)
	if err != nil {
		fmt.Println("Unable to parse link to a valid URL object")
		return "", ErrMalformattedLink
	}
	if fUrl.Host != "facetime.apple.com" {
		fmt.Println("URL host failed validation")
		return "", ErrMalformattedLink
	}
	if fUrl.Path != "/join" {
		fmt.Println("URL path failed validation")
		return "", ErrMalformattedLink
	}
	fragFormat := regexp.MustCompile(`^v\=1\&p\=(.+)`)
	if !(fragFormat.MatchString(fUrl.Fragment)) {
		fmt.Println("URL fragment failed validation")
		return "", ErrMalformattedLink
	}

	fUrl.Scheme = "https"
	fL := ftLink(fUrl.String())
	return fL, err
}

// Strips the scheme from the URL for lookup in the FT call table
func (f *ftLink) getId() ftLink {
	fUrl := string(*f)
	// TODO would be nice to do this w a better parser
	fId := strings.Replace(fUrl, "https://", "", 1)
	return ftLink(fId)
}
