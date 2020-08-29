package dbx

import (
	"fmt"
	"github.com/leporo/sqlf"
)

// Пагинатор
type Pagination struct {
	Size int
	From int
}

// Применить параметры пагинации
// 	stmt - постоитель запросов
func (p *Pagination) Apply(stmt *sqlf.Stmt) {
	stmt.Limit(p.Size).Offset(p.From)
}

// Сортировка
type Sort []SortOption

func (p *Sort) Add(field string, direct SortDirect) *Sort {
	*p = append(*p, SortOption{
		Field:  field,
		Direct: direct,
	})

	return p
}

// Применить параметры сортировки
// 	stmt - постоитель запросов
func (p Sort) Apply(stmt *sqlf.Stmt, fieldPrefix string) {
	if len(p) == 0 {
		return
	}

	if fieldPrefix != "" {
		fieldPrefix = fieldPrefix + "."
	}

	sort := make([]string, len(p))
	for i := 0; i < len(p); i++ {
		sort[i] = fmt.Sprintf("%s%s %s", fieldPrefix, p[i].Field, fixSortDirect(p[i].Direct))
	}

	stmt.OrderBy(sort...)
}

func fixSortDirect(direct SortDirect) string {
	if direct == SortDirectAsc {
		return fmt.Sprintf("ASC NULLS FIRST")
	} else if direct == SortDirectDesc {
		return fmt.Sprintf("DESC NULLS LAST")
	}
	return ""
}

type SortOption struct {
	Field  string
	Direct SortDirect
}

type SortDirect string

const (
	SortDirectAsc  SortDirect = "ASC"
	SortDirectDesc SortDirect = "DESC"
)
