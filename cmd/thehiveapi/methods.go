package thehiveapi

// NewSpecialObjectForCache конструктор вспомогательного типа реализующий интерфейс CacheStorageFuncHandler[T any]
func NewSpecialObjectForCache[T any]() *SpecialObjectForCache[T] {
	return &SpecialObjectForCache[T]{}
}

func (o *SpecialObjectForCache[T]) SetID(v string) {
	o.id = v
}

func (o *SpecialObjectForCache[T]) GetID() string {
	return o.id
}

func (o *SpecialObjectForCache[T]) SetObject(v T) {
	o.object = v
}

func (o *SpecialObjectForCache[T]) GetObject() T {
	return o.object
}

func (o *SpecialObjectForCache[T]) SetFunc(f func(int) bool) {
	o.handlerFunc = f
}

func (o *SpecialObjectForCache[T]) GetFunc() func(int) bool {
	return o.handlerFunc
}

// Comparison сравнение содержимого объектов. В данном случае сревнение
// нет, это простая заглушка.
// Для того что бы не досить thehive метод всегда будет возвращать TRUE.
// Соответственно не будет заменять объект в работе.
func (o *SpecialObjectForCache[T]) Comparison(objFromCache T) bool {
	return true
}
