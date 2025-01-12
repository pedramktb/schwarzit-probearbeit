package ginRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetID(c *gin.Context, key string) uuid.UUID {
	idAny, _ := c.Get(key)
	id, ok := idAny.(uuid.UUID)
	if !ok {
		return uuid.Nil
	} else {
		return id
	}
}
