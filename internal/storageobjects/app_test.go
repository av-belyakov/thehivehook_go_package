package storageobjects_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/storageobjects"
)

func TestStorageObjects(t *testing.T) {
	size := 15
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()

	s, err := storageobjects.New(
		storageobjects.WithChannelSize[Element](5),
		storageobjects.WithTimeTick[Element](1),
		storageobjects.WithTimeToLive[Element](10),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.Start(ctx)

	assert.Equal(t, s.Len(), 0)

	for range size {
		id := gofakeit.UUID()
		s.AddObject(
			3, // 3 секундам
			storageobjects.StorageObjectDataSettings[Element]{
				Id:         id,
				ObjectType: "test_case",
				Data: Element{
					ID:   id,
					Type: gofakeit.AnimalType(),
					Name: gofakeit.Animal(),
				}})
	}

	objId := gofakeit.UUID()
	s.AddObject(
		3, // 3 секундам
		storageobjects.StorageObjectDataSettings[Element]{
			Id:         objId,
			ObjectType: "test_case",
			Data: Element{
				ID:   objId,
				Type: gofakeit.AnimalType(),
				Name: gofakeit.Animal(),
			}})

	assert.Equal(t, s.Len(), size+1)

	//модифицируем объект
	objType := gofakeit.AnimalType()
	objName := gofakeit.Animal()
	s.AddObject(
		3, // 3 секундам
		storageobjects.StorageObjectDataSettings[Element]{
			Id:         objId,
			ObjectType: "test_case",
			Data: Element{
				ID:   objId,
				Type: objType,
				Name: objName,
			}})

	assert.Equal(t, s.Len(), size+1)

	var num int
	for obj := range s.GetObjects() {
		num++

		//проверяем, что объект был модифицирован
		if objId == obj.Id {
			assert.Equal(t, obj.Data.Type, objType)
			assert.Equal(t, obj.Data.Name, objName)
		}

		fmt.Printf("%d. time: '%s', index: '%s', object: '%+v'\n", num, obj.TimeCreated, obj.Id, obj.Data)

		if num == size+1-4 {
			break
		}
	}

	assert.Equal(t, s.Len(), 0)

	ctx, cancel := context.WithCancel(context.Background())
	ctx.Done()
	assert.Nil(t, ctx.Err())

	cancel()
	assert.Error(t, ctx.Err())
}

type Element struct {
	ID   string
	Type string
	Name string
}
