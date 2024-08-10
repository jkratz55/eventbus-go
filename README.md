# EventBus

EventBus (`eventbus-go`) is a small, lightweight, zero-dependency library for sending events/messages between components/objects in an application without them being directly coupled. Components are only coupled to the eventbus rather than to each other. 

EventBus uses the type system for mapping events to handlers. Events are modeled as types. This provides a lot of flexibility for passing arbitrary data on events throughout an application.

Example:

```go
package main

import (
	"fmt"

	"github.com/jkratz55/eventbus-go"
)

type Status string

const (
	StatusUp       Status = "UP"
	StatusDegraded Status = "DEGRADED"
	StatusDown     Status = "DOWN"
)

type ApplicationHealthChangedEvent struct {
	NewStatus string
}

func main() {

	// Create a handler for the event. This just needs to be an implementation
	// of eventbus.Handler interface, which can be a type or a eventbus.HandlerFunc.
	handlerFn := eventbus.HandlerFunc[ApplicationHealthChangedEvent](func(event ApplicationHealthChangedEvent) {
		fmt.Println("Application status changed to", event.NewStatus)
		// todo: do something useful with the event
	})

	// Register the handler for the event type ApplicationHealthChangedEvent
	eventbus.Subscribe[ApplicationHealthChangedEvent](handlerFn)

	// Publish some events
	eventbus.MustPublish(ApplicationHealthChangedEvent{NewStatus: string(StatusDegraded)})
	eventbus.MustPublish(ApplicationHealthChangedEvent{NewStatus: string(StatusUp)})
}
```

The above code will generate the following output:

```text
Application status changed to DEGRADED
Application status changed to UP
```

In some cases an event may not need to carry additional data, in that case we can use an empty struct.

```go
package main

import (
	"fmt"

	"github.com/jkratz55/eventbus-go"
)

type ApplicationConfigChangedEvent struct{}

func main() {

	// Create a handler for the event. This just needs to be an implementation
	// of eventbus.Handler interface, which can be a type or a eventbus.HandlerFunc.
	handlerFn := eventbus.HandlerFunc[ApplicationConfigChangedEvent](func(event ApplicationConfigChangedEvent) {
		fmt.Println("Application config changed")
		// todo: do something useful like reload or update based on new config
	})

	// Register the handler for the event type ApplicationHealthChangedEvent
	eventbus.Subscribe[ApplicationConfigChangedEvent](handlerFn)

	// Publish some events
	eventbus.MustPublish(ApplicationConfigChangedEvent{})
	eventbus.MustPublish(ApplicationConfigChangedEvent{})
}
```

The recommended approach is to model your events as explicit as possible using good names like `CustomerCreated` or `CustomerDeleted`. However, the following is valid code, it's just as clear to understand at glance, and as more conditions are added ensuring everywhere the event is handled is updated can be a challenge.

```go
package main

import (
	"fmt"

	"github.com/jkratz55/eventbus-go"
)

type DatabaseStatusEvent int

const (
	DatabaseUp DatabaseStatusEvent = iota
	DatabaseDown
	DatabaseDegraded
)

func main() {

	// Create a handler for the event. This just needs to be an implementation
	// of eventbus.Handler interface, which can be a type or a eventbus.HandlerFunc.
	handlerFn := eventbus.HandlerFunc[DatabaseStatusEvent](func(event DatabaseStatusEvent) {
		fmt.Println("database status event")
		switch event {
		case DatabaseUp:
			fmt.Println("Database is up")
		case DatabaseDown:
			fmt.Println("Database is down")
		case DatabaseDegraded:
			fmt.Println("Database is degraded")
		}
	})

	// Register the handler for the event type ApplicationHealthChangedEvent
	eventbus.Subscribe[DatabaseStatusEvent](handlerFn)

	// Publish some events
	eventbus.MustPublish(DatabaseDegraded)
	eventbus.MustPublish(DatabaseDown)
	eventbus.MustPublish(DatabaseUp)
}
```