package domain

import (
	"github.com/google/uuid"
)

type Item struct {
	id     uuid.UUID
	typeID uint
}

func NewItem(typeID uint) *Item {
	return &Item{id: uuid.New(), typeID: typeID}
}

func (i *Item) ID() uuid.UUID {
	return i.id
}

func (i *Item) TypeID() uint {
	return i.typeID
}

func (i *Item) SetTypeID(typeID uint) {
	i.typeID = typeID
}

func ItemFromSource(id uuid.UUID, typeID uint) *Item {
	return &Item{
		id:     id,
		typeID: typeID,
	}
}
