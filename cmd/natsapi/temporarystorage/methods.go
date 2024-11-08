package temporarystoarge

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// NewCell создает новую ячейку в хранилище в которую осуществляется запись,
// возвращает id новой ячейки, без корректного id запись в хранилище не возможна
func (ts *TemporaryStorage) NewCell() (id string) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	id = uuid.New().String()
	ts.ttlStorage.storage[id] = repository{
		timeExpiry: time.Now().Add(ts.ttl),
	}

	return id
}

// SetService добавляет наименование сервиса
func (ts *TemporaryStorage) SetService(id, v string) error {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return fmt.Errorf("the parameter 'service' cannot be set, repository with id '%s' was not found", id)
	}

	tmp := ts.ttlStorage.storage[id]
	tmp.service = v
	ts.ttlStorage.storage[id] = tmp

	return nil
}

// GetService возвращает наименование сервиса
func (ts *TemporaryStorage) GetService(id string) (string, error) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return "", fmt.Errorf("repository with id '%s' was not found", id)
	}

	return ts.ttlStorage.storage[id].service, nil
}

// SetCommand добавляет наименование команды
func (ts *TemporaryStorage) SetCommand(id, v string) error {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return fmt.Errorf("the parameter 'command' cannot be set, repository with id '%s' was not found", id)
	}

	tmp := ts.ttlStorage.storage[id]
	tmp.command = v
	ts.ttlStorage.storage[id] = tmp

	return nil
}

// GetCommand возвращает наименование команды
func (ts *TemporaryStorage) GetCommand(id string) (string, error) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return "", fmt.Errorf("repository with id '%s' was not found", id)
	}

	return ts.ttlStorage.storage[id].command, nil
}

// SetRootId добавляет rootId
func (ts *TemporaryStorage) SetRootId(id, v string) error {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return fmt.Errorf("the parameter 'rootId' cannot be set, repository with id '%s' was not found", id)
	}

	tmp := ts.ttlStorage.storage[id]
	tmp.rootId = v
	ts.ttlStorage.storage[id] = tmp

	return nil
}

// GetRootId возвращает rootId
func (ts *TemporaryStorage) GetRootId(id string) (string, error) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return "", fmt.Errorf("repository with id '%s' was not found", id)
	}

	return ts.ttlStorage.storage[id].rootId, nil
}

// SetCaseId добавляет caseId
func (ts *TemporaryStorage) SetCaseId(id, v string) error {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return fmt.Errorf("the parameter 'caseId' cannot be set, repository with id '%s' was not found", id)
	}

	tmp := ts.ttlStorage.storage[id]
	tmp.caseId = v
	ts.ttlStorage.storage[id] = tmp

	return nil
}

// GetCaseId возвращает caseId
func (ts *TemporaryStorage) GetCaseId(id string) (string, error) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return "", fmt.Errorf("repository with id '%s' was not found", id)
	}

	return ts.ttlStorage.storage[id].caseId, nil
}

// SetNsMsg добавляет дескриптор запроса NATS
func (ts *TemporaryStorage) SetNsMsg(id string, v *nats.Msg) error {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return fmt.Errorf("the parameter 'nats message' cannot be set, repository with id '%s' was not found", id)
	}

	tmp := ts.ttlStorage.storage[id]
	tmp.nsMsg = v
	ts.ttlStorage.storage[id] = tmp

	return nil
}

// GetNsMsg возвращает дескриптор запроса NATS
func (ts *TemporaryStorage) GetNsMsg(id string) (*nats.Msg, error) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	if !ts.repositoryIsExist(id) {
		return nil, fmt.Errorf("repository with id '%s' was not found", id)
	}

	return ts.ttlStorage.storage[id].nsMsg, nil
}

// DeleteElement удаляет заданный элемент по его id
func (ts *TemporaryStorage) DeleteElement(id string) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	delete(ts.ttlStorage.storage, id)
}

func (ts *TemporaryStorage) repositoryIsExist(id string) bool {
	if _, ok := ts.ttlStorage.storage[id]; !ok {
		return false
	}

	return true
}
