package eventbus

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userCreatedEvent struct {
	Name  string
	Email string
}

type userCreatedHandler struct {
	mock.Mock
}

func (h *userCreatedHandler) OnEvent(event userCreatedEvent) {
	h.Called(event)
}

func TestSubscribe(t *testing.T) {
	reset()
	id := Subscribe[userCreatedEvent](new(userCreatedHandler))
	assert.Greater(t, id, uint64(0))
}

func TestUnsubscribe(t *testing.T) {
	reset()
	id := Subscribe[userCreatedEvent](new(userCreatedHandler))
	assert.Greater(t, id, uint64(0))
	assert.True(t, Unsubscribe[userCreatedEvent](id))
}

func TestPublish(t *testing.T) {
	reset()
	h := new(userCreatedHandler)
	h.On("OnEvent", userCreatedEvent{Name: "John Doe", Email: "jdoe@gmail.com"}).Return()

	Subscribe[userCreatedEvent](h)
	err := Publish(userCreatedEvent{Name: "John Doe", Email: "jdoe@gmail.com"})
	assert.NoError(t, err)

	h.AssertNumberOfCalls(t, "OnEvent", 1)
}

func TestPublishAsync(t *testing.T) {
	reset()
	h := new(userCreatedHandler)
	h.On("OnEvent", userCreatedEvent{Name: "John Doe", Email: "jdoe@gmail.com"}).Return()

	Subscribe[userCreatedEvent](h)
	err := PublishAsync(userCreatedEvent{Name: "John Doe", Email: "jdoe@gmail.com"})
	assert.NoError(t, err)

	// Since its invoked async need to wait for it to run
	time.Sleep(1 * time.Second)

	h.AssertNumberOfCalls(t, "OnEvent", 1)
}

func reset() {
	handlers = make(map[reflect.Type][]handlerEntry)
	mu = sync.RWMutex{}
	subscriberId = 0
}
