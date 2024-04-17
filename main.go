package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrMalformattedLink = errors.New("the url could not be validated against the facetime url format")
)

type status struct {
	Msg string
}
type ftLink string
type linkInfo struct {
	Link ftLink `json:"link"`
}
type linkJoinInfo struct {
	Link   ftLink `json:"link"`
	Joined bool   `json:"joined"`
}
type deleteLinkResponse struct {
	Link    ftLink `json:"link"`
	Deleted bool   `json:"deleted"`
}

func main() {
	router := gin.Default()
	router.GET("/status", getStatus)

	router.GET("/link", routeGetActiveLinks)
	router.POST("/link/new", routeNewLink)
	router.POST("/link/join", routeJoinLink)
	router.DELETE("/link", routeDeleteLink)

	router.Run("localhost:8080")
}

//TODO set timeouts
//TODO make sure we only process certain requests one at a time

func getStatus(c *gin.Context) {
	fmt.Println("Received status check request")
	c.IndentedJSON(http.StatusOK, status{"Server is running!"})
}

func routeGetActiveLinks(c *gin.Context) {
	fmt.Println("Received request for current facetime links")
	allLinks, err := getAllLinks()
	if err != nil {
		fmt.Println("Failed to retrieve current links: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println("retrieved links:" + strings.Join(allLinks, ", "))
	var allFtLinksInfo []linkInfo
	badLinks := false
	for _, l := range allLinks {
		ftL := ftLink(l)
		ftUrl, err := ftL.getUrl()
		if err != nil {
			badLinks = true
			fmt.Println("Got invalid link: '" + l + "'")
			continue
		}
		newFtLinkInfo := linkInfo{Link: ftUrl}
		allFtLinksInfo = append(allFtLinksInfo, newFtLinkInfo)
	}
	if badLinks {
		fmt.Println("Found links but unable to validate formatting of any.")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, allFtLinksInfo)
}

func routeNewLink(c *gin.Context) {
	newLink, err := makeNewLink()
	if err != nil {
		fmt.Println("Failed to make new link:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	newFtLink := ftLink(newLink)
	_, err = newFtLink.getUrl()
	if err != nil {
		fmt.Print("Got an invalid link:'" + newLink + "'\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, linkInfo{Link: newFtLink})
}

func routeJoinLink(c *gin.Context) {
	var requestInfo linkInfo
	err := c.BindJSON(&requestInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	linkToJoin, err := requestInfo.Link.getUrl()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	idToJoin := linkToJoin.getId()

	err = joinCall(string(idToJoin))
	if err != nil {
		fmt.Println("Failed to join link:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, linkJoinInfo{Link: linkToJoin, Joined: true})
}

func routeDeleteLink(c *gin.Context) {
	var requestInfo linkInfo
	err := c.BindJSON(&requestInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	linkToDelete, err := requestInfo.Link.getUrl()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	idToDelete := linkToDelete.getId()

	wasDeleted, err := deleteCall(string(idToDelete))
	var deleteStatus int
	if err != nil {
		fmt.Println("Link", linkToDelete, "was not able to be deleted due to an error:", err)
		deleteStatus = http.StatusInternalServerError
	} else if !wasDeleted {
		fmt.Println("Link", linkToDelete, "was not found for deletion.")
		deleteStatus = http.StatusNotFound
	} else {
		fmt.Println("Link", linkToDelete, "was successfully deleted.")
		deleteStatus = http.StatusOK
	}
	c.JSON(deleteStatus, deleteLinkResponse{Link: linkToDelete, Deleted: wasDeleted})
}

// Returns a well-formed URL for the FT link. Can be used for validation.
func (f *ftLink) getUrl() (ftLink, error) {
	fStr := string(*f)
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
