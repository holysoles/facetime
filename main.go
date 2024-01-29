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

	router.GET("/links", routeGetActiveLinks)
	router.POST("/links/new", routeNewLink)
	router.POST("/links/join", routeJoinLink)
	router.DELETE("/links", routeDeleteLink)

	router.Run("localhost:8080")
}

//TODO make routines
//TODO set timeouts

func getStatus(c *gin.Context) {
	fmt.Println("Received status check request")
	c.IndentedJSON(http.StatusOK, status{"Server is running!"})
}

func routeGetActiveLinks(c *gin.Context) {
	fmt.Println("Received request for current facetime links")
	allLinks := getAllLinks()
	fmt.Println(allLinks)
	var allFtLinksInfo []linkInfo
	for _, l := range allLinks {
		ftL := ftLink(l)
		newFtLinkInfo := linkInfo{Link: ftL.getUrl()}
		allFtLinksInfo = append(allFtLinksInfo, newFtLinkInfo)
	}
	c.JSON(http.StatusOK, allFtLinksInfo)
}

func routeNewLink(c *gin.Context) {
	newLink := makeNewLink()
	newFtLink := ftLink(newLink)

	if !newFtLink.isValid() {
		c.AbortWithStatus(http.StatusInternalServerError)
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
	return
}

func (f *ftLink) isValid() bool {
	linkFormat := regexp.MustCompile(`^facetime\.apple\.com\/join\#v\=1\&p\=(.+)`)
	return linkFormat.MatchString(string(*f))
}

// TODO this should return a URL
func (f *ftLink) getUrl() ftLink {
	fUrl := "https://" + string(*f)
	return ftLink(fUrl)
}

// TODO this should probably take a URL
func (f *ftLink) getId() ftLink {
	fUrl := string(*f)
	// strip protocol. would be nice to do this w a better parser
	fId := strings.Replace(fUrl, "https://", "", 1)
	return ftLink(fId)
}
