package supportingfunctions

type finalyList[T comparable] []T

func (fl finalyList[T]) searchElement(value T) bool {
	for _, v := range fl {
		if value == v {
			return true
		}
	}

	return false
}

// JoinSlises выполняет объединение нескольких срезов исключая дублирование элементов
func JoinSlises[T comparable](lists ...[]T) []T {
	countList := len(lists)
	fl := finalyList[T](nil)

	if countList == 0 {
		return fl
	}

	if countList == 1 {
		fl = append(fl, lists[0]...)

		return fl
	}

	for _, list := range lists {
		for _, v := range list {
			if fl.searchElement(v) {
				continue
			}

			fl = append(fl, v)
		}
	}

	return fl
}
