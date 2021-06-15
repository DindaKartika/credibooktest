package handlers

import (
	"math"

	"github.com/jinzhu/gorm"
)

type Paginator struct {
	DB      *gorm.DB
	OrderBy []string
	Page    int
	PerPage int
}

type Data struct {
	TotalRecords int         `json:"total_records"`
	Data         interface{} `json:"data"`
	CurrentPage  int         `json:"current_page"`
	TotalPages   int         `json:"total_pages"`
}

func (p *Paginator) paginate(data interface{}) *Data {
	db := p.DB

	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PerPage <= 0 {
		p.PerPage = 25
	}

	var offset int
	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.PerPage
	}

	cdb := db
	var total int
	// get total record include filters
	cdb.Find(data).Count(&total)

	db.Limit(p.PerPage).Offset(offset).Find(data)

	lastPage := int(math.Ceil(float64(total) / float64(p.PerPage)))

	resp := Data{
		TotalRecords: total,
		CurrentPage:  p.Page,
		TotalPages:   lastPage,
		Data:         &data,
	}

	return &resp
}
