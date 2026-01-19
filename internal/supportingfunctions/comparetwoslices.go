package supportingfunctions

import "slices"

type mainList[T comparable] []T

// GetUniq возвращает список уникальных элементов которые не встречаются в mainList
func (ml mainList[T]) GetUniq(l []T) []T {
	newList := []T(nil)

	for _, v := range l {
		var isExist bool
		if slices.Contains(ml, v) {
			isExist = true
		}

		if !isExist {
			newList = append(newList, v)
		}
	}

	return newList
}

// CompareTwoSlices выполняет сравнение двух списков и находит элементы из второго списка
// которые не встречаются первом
func CompareTwoSlices[T comparable](listMain, listCompare []T) []T {
	ml := mainList[T](listMain)

	return ml.GetUniq(listCompare)
}
