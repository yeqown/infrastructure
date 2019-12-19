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
	DB       *gorm.DB    // DB Data Source
	Limit    uint        // Limit
	Offset   uint        // Offset
	Select   string      // Select
	OrderBy  []string    // Orders
	GroupBy  string      // Group
	Joins    []Join      // Joins
	Wheres   []Where     // Wheres
	ORs      []Where     // ORs
	Not      interface{} // Not
	Preloads []string    // Preload
	ShowSQL  bool        // Debug mode
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

	if len(p.ORs) > 0 {
		for _, w := range p.ORs {
			if w.Value != nil {
				db = db.Or(w.Key, w.Value)
			} else {
				db = db.Or(w.Key)
			}
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
		paginator = new(Paginator)
	)

	if paginator.Error = db.Model(out).
		Count(&paginator.Total).Error; paginator.Error != nil {
		return paginator
	}

	// with preloads
	for _, preloadColumn := range p.Preloads {
		db = db.Preload(preloadColumn)
	}

	paginator.Error = db.Limit(p.Limit).Offset(p.Offset).Find(out).Error
	paginator.Records = out

	return paginator
}
