package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	processLock = false // http stdlib uses goroutine to process requests. We don't handle parallel applescript though

	ErrBusy             = errors.New("the server is currently busy processing another request")
	ErrMalformattedLink = errors.New("the url could not be validated against the facetime url format")
)

type status struct {
	ServerStatus   string `json:"server_status"`
	FacetimeStatus bool   `json:"facetime_status"`
	Busy           bool   `json:"busy"`
}
type linkInfo struct {
	Link ftLink `json:"link"`
}
type linkJoinInfo struct {
	Link   ftLink `json:"link"`
	Joined bool   `json:"joined"`
}
type linkAdmitInfo struct {
	Admitted bool `json:"admitted"`
}
type linkLeaveInfo struct {
	Left bool `json:"left"`
}
type linkDeleteInfo struct {
	Link    ftLink `json:"link"`
	Deleted bool   `json:"deleted"`
}

func main() {
	router := gin.Default()
	router.GET("/status", routeGetStatus)

	router.GET("/link", routeGetActiveLinks)
	router.POST("/link/new", routeNewLink)
	router.POST("/link/join", routeJoinLink)
	router.POST("/link/admit", routeAdmitLink)
	router.POST("/link/leave", routeLeaveLink)
	router.DELETE("/link", routeDeleteLink)

	openFacetime()
	router.Run("localhost:8080")
}

func routeGetStatus(c *gin.Context) {
	fmt.Println("Received status check request")
	res := http.StatusOK
	fSt, err := getFacetimeStatus()
	if err != nil {
		fmt.Println(err)
		res = http.StatusInternalServerError
	}
	c.IndentedJSON(res, status{ServerStatus: "running", FacetimeStatus: fSt, Busy: processLock})
}

func routeGetActiveLinks(c *gin.Context) {
	fmt.Println("Received request for current facetime links")
	err := initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	allLinks, err := getAllLinks()
	if err != nil {
		fmt.Println("Failed to retrieve current links:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println("retrieved links:" + strings.Join(allLinks, ", "))
	allFtLinksInfo := []linkInfo{}
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
	err := initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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
	err := initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = joinAndAdmitCall()
	if err != nil {
		fmt.Println("Failed to join link:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//TODO can still return link info
	linkToJoin := ftLink("placeholder")
	c.JSON(http.StatusOK, linkJoinInfo{Link: linkToJoin, Joined: true})
}

func routeAdmitLink(c *gin.Context) {
	err := initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = admitActiveCall()
	if err != nil {
		fmt.Println("Failed to admit participants to link:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, linkAdmitInfo{Admitted: true})
}

func routeLeaveLink(c *gin.Context) {
	err := initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = leaveCall()
	if err != nil {
		fmt.Println("Failed to have host leave call:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, linkLeaveInfo{Left: true})
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

	err = initSession()
	defer closeSession()
	if err == ErrBusy {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	} else if err != nil {
		fmt.Println("Unhandled error in initializing FaceTime:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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
	c.JSON(deleteStatus, linkDeleteInfo{Link: linkToDelete, Deleted: wasDeleted})
}
