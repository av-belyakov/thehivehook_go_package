package cacherunningfunctions

import (
	"fmt"
	"time"
)

// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningFunctions) SetMethod(id string, f func(v int) bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	//!!!!!!!!!
	//тут можно сделать проверку есть ли объект с таким id, выполняется ли функция и т.д
	//!!!!!!!!!!
	fmt.Println("func 'GetMethod', add func for object with id:", id)

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry: time.Now().Add(crm.ttl),
		cacheFunc:  f,
	}

	fmt.Println("+++ CACHE func 'SetMethod', crm.cacheStorage.storages[id] =", crm.cacheStorage.storages[id])

	return id
}

// GetMethod возвращает данные по полученому id
func (crm *CacheRunningFunctions) GetMethod(id string) (func(int) bool, bool) {
	crm.cacheStorage.mutex.RLock()
	defer crm.cacheStorage.mutex.Unlock()

	if storage, ok := crm.cacheStorage.storages[id]; ok {
		return storage.cacheFunc, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningFunctions) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}

// GetNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) GetNumberAttempts(id string) int {
	crm.cacheStorage.mutex.RLock()
	defer crm.cacheStorage.mutex.RUnlock()

	na, ok := crm.getNumberAttempts(id)
	if !ok {
		return 0
	}

	return na
}

// IncreaseNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) IncreaseNumberAttempts(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.increaseNumberAttempts(id)
}

// SetIsCompletedSuccessfully выполняемая функция завершилась успехом
func (crm *CacheRunningFunctions) SetIsCompletedSuccessfully(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.setIsCompletedSuccessfully(id)
}

// SetIsFunctionExecution функция находится в процессе выполнения
func (crm *CacheRunningFunctions) SetIsFunctionExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.setIsFunctionExecution(id)
}

// SetIsFunctionNotExecution функция не выполняется
func (crm *CacheRunningFunctions) SetIsFunctionNotExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.setIsFunctionNotExecution(id)
}

// getNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) getNumberAttempts(id string) (int, bool) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return 0, ok
	}

	return storage.numberAttempts, ok
}

// increaseNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) increaseNumberAttempts(id string) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.numberAttempts++
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionExecution функция находится в процессе выполнения
func (crm *CacheRunningFunctions) setIsFunctionExecution(id string) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = true
	crm.cacheStorage.storages[id] = storage
}

// setIsCompletedSuccessfully выполняемая функция завершилась успехом
func (crm *CacheRunningFunctions) setIsCompletedSuccessfully(id string) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isCompletedSuccessfully = true
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionNotExecution функция не выполняется
func (crm *CacheRunningFunctions) setIsFunctionNotExecution(id string) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = false
	crm.cacheStorage.storages[id] = storage
}
