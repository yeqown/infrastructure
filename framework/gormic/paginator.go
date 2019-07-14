package gormic

import (
	"github.com/jinzhu/gorm"
)

// Where .
// db.Where(key, value)
type Where struct {
	Key   string
	Value interface{}
}

// Join .
// db.Joins("JOIN department on department.id = employee.department_id AND employee.name = ?", "Jenkins")
type Join struct {
	Query string
	Args  []interface{}
}

// PagingParam 分页参数
type PagingParam struct {
	DB      *gorm.DB    // DB Data Source
	Limit   uint        // Limit
	Offset  uint        // Offset
	Select  string      // Select
	OrderBy []string    // Orders
	GroupBy string      // Group
	Joins   []Join      // Joins
	Wheres  []Where     // Wheres
	Not     interface{} // Not
	ShowSQL bool        // Debug mode
	// CountDataSource interface{} //
}

// Paginator 分页返回
type Paginator struct {
	Total   uint        `json:"total"`
	Records interface{} `json:"records"`
	Error   error       `json:"error,omitempty"`
}

// Pagging 分页
func Pagging(p *PagingParam, out interface{}) *Paginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}

	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Select != "" {
		db = db.Select(p.Select)
	}

	if len(p.Joins) > 0 {
		for _, j := range p.Joins {
			db = db.Joins(j.Query, j.Args...)
		}
	}

	if len(p.Wheres) > 0 {
		for _, w := range p.Wheres {
			if w.Value != nil {
				db = db.Where(w.Key, w.Value)
			} else {
				db = db.Where(w.Key)
			}
		}
	}

	if p.Not != nil {
		db = db.Not(p.Not)
	}

	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	if p.GroupBy != "" {
		db = db.Group(p.GroupBy)
	}

	var (
		paginator Paginator
	)

	if err := db.Model(out).
		Count(&paginator.Total).Error; err != nil {
		paginator.Error = err
		return &paginator
	}

	paginator.Error = db.
		Limit(p.Limit).
		Offset(p.Limit).
		Find(paginator.Records).Error

	return &paginator
}
