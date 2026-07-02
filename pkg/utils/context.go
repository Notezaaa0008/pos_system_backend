package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetFromCtx(c *gin.Context, name string) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, errors.New("context key name cannot be empty")
	}

	value, exists := c.Get(name)
	if !exists {
		return uuid.Nil, fmt.Errorf("unauthorized: missing %s in context", name)
	}

	parsedUUID, ok := value.(uuid.UUID)
	if ok {
		return parsedUUID, nil
	}

	idStr, ok := value.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("internal: context key %s is not a string", name)
	}

	parsedUUID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid format for %s: %w", name, err)
	}

	return parsedUUID, nil
}