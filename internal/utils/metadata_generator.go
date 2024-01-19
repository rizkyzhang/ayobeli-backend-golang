package utils

import (
	"strings"
	"time"

	"github.com/lucsky/cuid"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

func GenerateMetadata() domain.Metadata {
	now := time.Now().UTC()

	return domain.Metadata{
		UID: func() string {
			return cuid.New()
		},
		Slug: func(str string) string {
			return strings.ToLower(strings.Join(strings.Split(str, " "), "-"))
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
