package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
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

//TODO make routines
//TODO set timeouts
///
//TODO make sure we only process certain requests one at a time
//
// TODO make sure we dont exit 1 on error where possible, and instead abort the request with a 500

func getStatus(c *gin.Context) {
	fmt.Println("Received status check request")
	c.IndentedJSON(http.StatusOK, status{"Server is running!"})
}

func routeGetActiveLinks(c *gin.Context) {
	fmt.Println("Received request for current facetime links")
	allLinks := getAllLinks()
	fmt.Println("retrieved links:" + strings.Join(allLinks, ", "))
	var allFtLinksInfo []linkInfo
	for _, l := range allLinks {
		ftL := ftLink(l)
		if !ftL.isValid() {
			fmt.Println("Got invalid link: '" + l + "'") //TODO warn
			continue
		}
		newFtLinkInfo := linkInfo{Link: ftL.getUrl()}
		allFtLinksInfo = append(allFtLinksInfo, newFtLinkInfo)
	}
	c.JSON(http.StatusOK, allFtLinksInfo)
}

func routeNewLink(c *gin.Context) {
	newLink := makeNewLink()
	newFtLink := ftLink(newLink)
	if !newFtLink.isValid() {
		fmt.Println("Got an invalid link: '" + newLink + "'")
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
	linkToJoin := requestInfo.Link.getId()
	if !linkToJoin.isValid() {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	joinLink(string(linkToJoin))
	c.JSON(http.StatusOK, linkJoinInfo{Link: linkToJoin, Joined: true})
}

func routeDeleteLink(c *gin.Context) {
	var requestInfo linkInfo
	err := c.BindJSON(&requestInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	linkToDelete := requestInfo.Link.getId()
	if !linkToDelete.isValid() {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	wasDeleted := deleteLink(string(linkToDelete))
	var deleteStatus int
	if !wasDeleted {
		fmt.Println("Link", linkToDelete, "was not found for deletion.")
		deleteStatus = http.StatusNotFound
	} else {
		fmt.Println("Link", linkToDelete, "was successfully deleted.")
		deleteStatus = http.StatusOK
	}
	c.JSON(deleteStatus, deleteLinkResponse{Link: linkToDelete, Deleted: wasDeleted})
}

func (f *ftLink) isValid() bool {
	u := (*f).getId()

	linkFormat := regexp.MustCompile(`^facetime\.apple\.com\/join\#v\=1\&p\=(.+)`)
	return linkFormat.MatchString(string(u))
}

// TODO Returns a guaranteed URL for the FT link
func (f *ftLink) getUrl() ftLink {
	fUrl := "https://" + string(*f)
	return ftLink(fUrl)
}

// TODO Strips the scheme from the URL
func (f *ftLink) getId() ftLink {
	fUrl := string(*f)
	// strip protocol. would be nice to do this w a better parser
	fId := strings.Replace(fUrl, "https://", "", 1)
	return ftLink(fId)
}
