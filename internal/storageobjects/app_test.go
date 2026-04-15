package storageobjects_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/v2/internal/storageobjects"
)

func TestStorageObjects(t *testing.T) {
	var (
		objId   string = uuid.NewString()
		objName string = gofakeit.Animal()
		objType string = gofakeit.AnimalType()
	)

	size := 15
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()

	s, err := storageobjects.New(
		storageobjects.WithChannelSize[Element](5),
		storageobjects.WithTimeTick[Element](1),
		storageobjects.WithTimeToLive[Element](30),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.Start(ctx)

	t.Run("Test 1. Check count elements", func(t *testing.T) {
		// 0 elements
		assert.Equal(t, s.Len(), 0)

		// add any size elements
		for n := range size {
			s.AddObject(
				3, // 3 секундам
				storageobjects.StorageObjectDataSettings[Element]{
					Id:         uuid.NewString(),
					ObjectType: fmt.Sprintf("test_case:%d", n),
					Data: Element{
						UUID: gofakeit.UUID(),
						Type: gofakeit.AnimalType(),
						Name: gofakeit.Animal(),
					}})
		}

		// add 1 element
		s.AddObject(
			3, // 3 секундам
			storageobjects.StorageObjectDataSettings[Element]{
				Id:         objId,
				ObjectType: "test_case",
				Data: Element{
					UUID: gofakeit.UUID(),
					Type: gofakeit.AnimalType(),
					Name: gofakeit.Animal(),
				}})

		// check count elements
		assert.Equal(t, s.Len(), size+1)

		// modified 1 element
		s.AddObject(
			3, // 3 секундам
			storageobjects.StorageObjectDataSettings[Element]{
				Id:         objId,
				ObjectType: "test_case",
				Data: Element{
					UUID: gofakeit.UUID(),
					Type: objType,
					Name: objName,
				}})

		assert.Equal(t, s.Len(), size+1)
	})

	t.Run("Test 2. Object modified", func(t *testing.T) {
		var num int
		for obj := range s.GetObjects() {
			num++

			//проверяем, что объект был модифицирован
			if obj.Id == objId {
				fmt.Printf("objId '%s' == '%s' obj.Id\n", objId, obj.Id)

				assert.Equal(t, obj.ObjectType, "test_case")
				assert.Equal(t, obj.Data.Type, objType)
				assert.Equal(t, obj.Data.Name, objName)
			}

			fmt.Printf("%d. time: '%s', index: '%s', object type: '%s', object: '%+v'\n", num, obj.TimeCreated, obj.Id, obj.ObjectType, obj.Data)

			if num == size+1 {
				break
			}
		}
	})

	t.Run("Test 3. Check storage size", func(t *testing.T) {
		assert.Equal(t, s.Len(), 0)
	})

	ctx, cancel := context.WithCancel(context.Background())
	ctx.Done()
	assert.Nil(t, ctx.Err())

	cancel()
	assert.Error(t, ctx.Err())
}

type Element struct {
	UUID string
	Type string
	Name string
}
