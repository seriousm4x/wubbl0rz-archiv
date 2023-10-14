package queries

import (
	"math"
	"os"
	"strconv"

	"github.com/seriousm4x/wubbl0rz-archiv-backend/pkg/logger"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int   `json:"limit,omitempty"`
	Page       int   `json:"page,omitempty"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	// maybe add previous, next, first and last page link
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		limit, err := strconv.Atoi(os.Getenv("API_PAGE_LIMIT"))
		if err != nil {
			logger.Error.Println("[query] Env Error: \"API_PAGE_LIMIT\" doesn't seem to be a valid integer. Defaulting to 50.")
			limit = 50
		}
		p.Limit = limit
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func Paginate(pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	totalPages := int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
