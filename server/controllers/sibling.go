package controllers

import (
	"fmt"

	"github.com/Lachstec/digsinet-ng/builder"
	"github.com/Lachstec/digsinet-ng/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

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
		log.Error().
			Err(err).
			Msg("Failed to bind JSON")
		err := c.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to abort with error")
		}
		return
	}

	// Add the sibling
	siblings = append(siblings, newSibling)

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "sibling created"})
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
					log.Error().
						Err(err).
						Msg("Failed to deploy topology")
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						log.Error().
							Err(err).
							Msg("Failed to abort with error")
					}
					return
				}
				c.IndentedJSON(http.StatusOK, gin.H{"message": "topology deployed"})
				return
			default:
				log.Error().
					Str("builder", s.Builder).
					Msg("Unknown Builder: ")
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Unknown Builder: %s", s.Builder))
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to abort with error")
				}
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
					log.Error().
						Err(err).
						Msg("Failed to destroy topology")
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						log.Error().
							Err(err).
							Msg("Failed to abort with error")
					}
					return
				}
				c.IndentedJSON(http.StatusOK, gin.H{"message": "topology destroyed"})
				return
			default:
				log.Error().
					Str("builder", s.Builder).
					Msg("Unknown Builder: ")
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Unknown Builder: %s", s.Builder))
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to abort with error")
				}
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func (s SiblingController) StartNodeIface(c *gin.Context) {
	id := c.Param("id")
	node := c.Param("node")
	for _, s := range siblings {
		if s.ID == id {
			switch {
			case s.Builder == "clab":
				clab := builder.ClabBuilder{}
				subscriptionID, err := clab.StartNodeIface(s.Topology, node)
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to start node iface")
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						log.Error().
							Err(err).
							Msg("Failed to abort with error")
					}
					return
				} else {
					c.IndentedJSON(http.StatusOK, gin.H{"subscriptionID": subscriptionID})
					return
				}
			default:
				log.Error().
					Str("builder", s.Builder).
					Msg("Unknown Builder: ")
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Unknown Builder: %s", s.Builder))
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to abort with error")
				}
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func (s SiblingController) StopNodeIface(c *gin.Context) {
	id := c.Param("id")
	node := c.Param("node")
	subscriptionID := c.Param("subscriptionID")
	for _, s := range siblings {
		if s.ID == id {
			switch {
			case s.Builder == "clab":
				clab := builder.ClabBuilder{}
				err := clab.StopNodeIface(s.Topology, node, subscriptionID)
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to stop node iface")
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						log.Error().
							Err(err).
							Msg("Failed to abort with error")
					}
					return
				}
				c.IndentedJSON(http.StatusOK, gin.H{"message": "node iface stopped"})
				return
			default:
				log.Error().
					Str("builder", s.Builder).
					Msg("Unknown Builder: ")
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Unknown Builder: %s", s.Builder))
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed to abort with error")
				}
				return
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
			c.IndentedJSON(http.StatusOK, gin.H{"message": "sibling deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sibling not found"})
}

func (s SiblingController) DeleteSiblings(c *gin.Context) {
	siblings = []sibling{}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "all siblings deleted"})
}
