package controllers

import (
	"github.com/Lachstec/digsinet-ng/builder"
	"github.com/Lachstec/digsinet-ng/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SiblingController struct{}

type sibling struct {
	ID       string `json:"id"`
	Builder  string `json:"builder"`
	Topology types.Topology
}

// siblings
var siblings = []sibling{}

func (s SiblingController) GetSiblings(c *gin.Context) {
	// TODO: c.JSON()?
	c.IndentedJSON(http.StatusOK, siblings)
}

func (s SiblingController) CreateSibling(c *gin.Context) {
	var newSibling sibling
	if err := c.BindJSON(&newSibling); err != nil {
		log.Print("Failed to bind JSON: ", err)
		return
	}

	// Add the sibling
	siblings = append(siblings, newSibling)

	c.IndentedJSON(http.StatusCreated, newSibling)
}

func (s SiblingController) GetSiblingByID(c *gin.Context) {
	id := c.Param("id")
	for _, s := range siblings {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func (s SiblingController) StartSiblingByID(c *gin.Context) {
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

func (s SiblingController) StopSiblingByID(c *gin.Context) {
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

func (s SiblingController) DeleteSiblingByID(c *gin.Context) {
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

func (s SiblingController) DeleteSiblings(c *gin.Context) {
	siblings = []sibling{}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "all siblings deleted"})
}
