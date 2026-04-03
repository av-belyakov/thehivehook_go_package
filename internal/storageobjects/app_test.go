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
	const objectType = "test object"

	size := 15
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		stop()
	}()

	s, err := storageobjects.New(
		storageobjects.WithChannelSize(5),
		storageobjects.WithTimeTick(1),
		storageobjects.WithTimeToLive(10),
		storageobjects.WithTimeDelayToSend(3),
	)
	if err != nil {
		t.Fatal(err)
	}
	s.Start(ctx)

	assert.Equal(t, s.Len(), 0)

	for range size {
		id := gofakeit.UUID()
		s.AddObject(id, objectType, fmt.Appendf(nil, `{
			"ID":   %s,
			"Type": %s,
			"Name": %s,
		}`, id, gofakeit.AnimalType(), gofakeit.Animal()))
	}

	objId := gofakeit.UUID()
	s.AddObject(objId, objectType, fmt.Appendf(nil, `{
			"ID":   %s,
			"Type": %s,
			"Name": %s,
		}`, objId, gofakeit.AnimalType(), gofakeit.Animal()))

	assert.Equal(t, s.Len(), size+1)

	updatedObject := fmt.Appendf(nil, `{
			"ID":   %s,
			"Type": %s,
			"Name": %s,
		}`, objId, gofakeit.AnimalType(), gofakeit.Animal())

	//модифицируем объект
	s.AddObject(objId, objectType, updatedObject)

	assert.Equal(t, s.Len(), size+1)

	var num int
	for obj := range s.GetObjects() {
		num++

		//проверяем, что объект был модифицирован
		if objId == obj.Id {
			assert.Equal(t, obj.Data, updatedObject)
		}

		fmt.Printf("%d. time: '%s', index: '%s', object type: '%s', object: '%s'\n", num, obj.TimeCreated, obj.Id, obj.ObjectType, string(obj.Data))

		if num == size+1-4 {
			break
		}
	}

	assert.Equal(t, s.Len(), 0)
}

type Element struct {
	ID   string
	Type string
	Name string
}
