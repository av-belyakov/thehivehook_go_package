package supportingfunctions

import "slices"

type mainList[T comparable] []T

// GetUniq возвращает список уникальных элементов которые не встречаются в mainList
func (ml mainList[T]) GetUniq(l []T) []T {
	newList := []T(nil)

	for _, v := range l {
		if slices.Contains(ml, v) {
			continue
		}

		newList = append(newList, v)
	}

	return newList
}

// JoinTwoSlicesUniqValues выполняет объединение двух срезов исключая дублирование элементов
func JoinTwoSlicesUniqValues[T comparable](listMain, listCompare []T) []T {
	ml := mainList[T](listMain)
	ml = append(ml, ml.GetUniq(listCompare)...)

	return ml
}
