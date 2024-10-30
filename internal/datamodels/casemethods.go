package datamodels

import "fmt"

// GetEventId возвращает уникальный id элемента основанный на комбинации некоторых значений EventElement
func (e CaseEventElement) GetEventId() string {
	return fmt.Sprintf("%s:%d:%s", e.ObjectType, e.Object.CreatedAt, e.RootId)
}
