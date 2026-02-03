package domain

import (
	"time"
)

type Building struct {
	ID          int64  // уникальный идентификатор
	Name        string // название или адрес
	Description string
	Floors      int       // количество этажей
	Geometry    Geometry  // форма здания (полигон)
	CreatedAt   time.Time // когда добавлено
}
