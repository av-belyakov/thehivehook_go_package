package cacherunningfunctions

import (
	"time"
)

// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningFunctions) SetMethod(id string, f func(v int) bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry: time.Now().Add(crm.ttl),
		cacheFunc:  f,
	}

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

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.numberAttempts++
	crm.cacheStorage.storages[id] = storage
}

// SetIsCompletedSuccessfully выполняемая функция завершилась успехом
func (crm *CacheRunningFunctions) SetIsCompletedSuccessfully(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isCompletedSuccessfully = true
	crm.cacheStorage.storages[id] = storage
}

// SetIsFunctionExecution функция находится в процессе выполнения
func (crm *CacheRunningFunctions) SetIsFunctionExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = true
	crm.cacheStorage.storages[id] = storage
}

// SetIsFunctionNotExecution функция не выполняется
func (crm *CacheRunningFunctions) SetIsFunctionNotExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = false
	crm.cacheStorage.storages[id] = storage
}

// getNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) getNumberAttempts(id string) (int, bool) {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return 0, ok
	}

	return storage.numberAttempts, ok
}

/*
// getFunctionCompletionStatus статус завершения функции
func (crm *CacheRunningFunctions) getFunctionCompletionStatus(id string) bool {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return true
	}

	return storage.isCompletedSuccessfully
}

// getFunctionExecutionStatus статус выполнения функции
func (crm *CacheRunningFunctions) getFunctionExecutionStatus(id string) bool {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return false
	}

	return storage.isFunctionExecution
}
*/
