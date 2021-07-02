package util

import (
	"xorm.io/xorm"
)

type MainInput struct {
	PageNum  int64    `json:"page" form:"page" query:"page"`
	PageSize int64    `json:"size" form:"size" query:"size"`
	Ids      []uint64 `json:"ids" form:"ids" query:"ids"`
	// refer field
	Refer map[string]uint64 `json:"refer" form:"refer" query:"refer"`
	// value is true=desc false=asc
	Sort   map[string]bool       `json:"sort" form:"sort" query:"sort"`
	Filter map[string]MainFilter `json:"filter" form:"filter" query:"filter"`
}

type MainFilter struct {
	Value  string `json:"value" form:"value" query:"value"`
	IsLike bool   `json:"like" form:"like" query:"like"`
}

func (input *MainInput) format() {
	if input.PageNum <= 0 {
		input.PageNum = 1
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}
}

func (input *MainInput) ApplyWithoutDisable(sql *xorm.Session) *xorm.Session {
	return input.Apply(sql).Where("is_disable = 0")
}

func (input *MainInput) Apply(sql *xorm.Session) *xorm.Session {
	input.Where(sql)
	input.Order(sql)
	if input.PageSize > 0 {
		sql.Limit(int(input.PageSize), input.Offset())
	}
	return sql
}

func (input *MainInput) Offset() int {
	input.format()
	return int((input.PageNum - 1) * input.PageSize)
}

func (input *MainInput) Order(sql *xorm.Session) *xorm.Session {
	if len(input.Sort) > 0 {
		for key, desc := range input.Sort {
			if "" != key {
				if desc {
					sql.Desc(ToSnakeCase(key))
				} else {
					sql.Asc(ToSnakeCase(key))
				}
			}
		}
	} else {
		sql.Desc("id")
	}
	return sql
}

func (input *MainInput) Where(sql *xorm.Session) *xorm.Session {
	if nil == input {
		return sql
	}
	if len(input.Refer) > 0 {
		for key, val := range input.Refer {
			sql.Where(ToSnakeCase(key)+" = ?", val)
		}
	}
	if len(input.Ids) > 0 {
		sql.In("id", input.Ids)
	} else {
		if len(input.Filter) > 0 {
			for key, val := range input.Filter {
				if "" != key && "" != val.Value {
					if val.IsLike {
						sql.Where(ToSnakeCase(key)+" like ?", "%"+val.Value+"%")
					} else {
						sql.Where(ToSnakeCase(key)+" = ?", val.Value)
					}
				}
			}
		}
	}
	return sql
}
