package main

import (
	"github.com/Lachstec/digsinet-ng/builder"
	"github.com/Lachstec/digsinet-ng/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type sibling struct {
	ID       string `json:"id"`
	Builder  string `json:"builder"`
	Topology types.Topology
}

// siblings
var siblings = []sibling{}

func getSiblings(c *gin.Context) {
	// TODO: c.JSON()?
	c.IndentedJSON(http.StatusOK, siblings)
}

func postSiblings(c *gin.Context) {
	var newSibling sibling
	if err := c.BindJSON(&newSibling); err != nil {
		log.Print("Failed to bind JSON: ", err)
		return
	}

	// Add the sibling
	siblings = append(siblings, newSibling)

	c.IndentedJSON(http.StatusCreated, newSibling)
}

func getSiblingByID(c *gin.Context) {
	id := c.Param("id")
	for _, s := range siblings {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func postStartSiblingByID(c *gin.Context) {
	id := c.Param("id")
	for _, s := range siblings {
		if s.ID == id {
			switch {
			// TODO: detect supported builders instead of hard-coded cases?
			case s.Builder == "clab":
				clab := builder.ClabBuilder{}
				err := clab.DeployTopology(s.Topology)
				if err != nil {
					log.Printf("Failed to deploy topology: %s", err)
				}
				return
			default:
				log.Printf("Unknown builder: %s", s.Builder)
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func postStopSiblingByID(c *gin.Context) {
	id := c.Param("id")
	for _, s := range siblings {
		if s.ID == id {
			switch {
			case s.Builder == "clab":
				clab := builder.ClabBuilder{}
				err := clab.DestroyTopology(s.Topology)
				if err != nil {
					log.Printf("Failed to destroy topology: %s", err)
				}
				return
			default:
				log.Printf("Unknown builder: %s", s.Builder)
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func deleteSiblingByID(c *gin.Context) {
	id := c.Param("id")
	for i, s := range siblings {
		if s.ID == id {
			//  remove s from siblings
			siblings = append(siblings[:i], siblings[i+1:]...)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func deleteSiblings(c *gin.Context) {
	siblings = []sibling{}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "all siblings deleted"})
}

func main() {
	router := gin.Default()

	router.GET("/siblings", getSiblings)
	router.GET("/siblings/:id", getSiblingByID)

	router.POST("/siblings", postSiblings)

	router.POST("/siblings/:id/start", postStartSiblingByID)
	router.POST("/siblings/:id/stop", postStopSiblingByID)

	router.DELETE("/siblings", deleteSiblings)
	router.DELETE("/siblings/:id", deleteSiblingByID)

	err := router.Run("localhost:8088")
	if err != nil {
		log.Fatal("Unable to start server: ", err)
		return
	}
}
