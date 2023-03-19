package connection

import (
	"sync"
	"time"
)

type ConnectionID uint64

type ConnectionType int

const (
	HTTP1 ConnectionType = iota
	HTTP2
)

func (ct ConnectionType) String() string {
	switch ct {
	case HTTP1:
		return "HTTP1"

	case HTTP2:
		return "HTTP2"
	}

	return "(Unknown)"
}

type Connection interface {
	ID() ConnectionID
	ConnectionType() ConnectionType
	CreationTime() time.Time
}

type connection struct {
	id             ConnectionID
	connectionType ConnectionType
	creationTime   time.Time
}

func (c connection) ID() ConnectionID {
	return c.id
}

func (c connection) ConnectionType() ConnectionType {
	return c.connectionType
}

func (c connection) CreationTime() time.Time {
	return c.creationTime
}

type ConnectionManager interface {
	AddConnection(connectionType ConnectionType) ConnectionID

	RemoveConnection(id ConnectionID)

	Connections() []Connection
}

type connectionManager struct {
	mutex            sync.RWMutex
	idToConnection   map[ConnectionID]connection
	nextConnectionID ConnectionID
}

func newConnectionManager() connectionManager {
	return connectionManager{
		idToConnection:   make(map[ConnectionID]connection),
		nextConnectionID: 1,
	}
}

func (cm *connectionManager) AddConnection(connectionType ConnectionType) ConnectionID {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	id := cm.nextConnectionID
	cm.nextConnectionID++

	cm.idToConnection[id] = connection{
		id:             id,
		creationTime:   time.Now(),
		connectionType: connectionType,
	}

	return id
}

func (cm *connectionManager) RemoveConnection(id ConnectionID) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.idToConnection, id)
}

func (cm *connectionManager) Connections() []Connection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	connections := make([]Connection, 0, len(cm.idToConnection))

	for _, connection := range cm.idToConnection {
		connections = append(connections, connection)
	}

	return connections
}

var connectionManagerInstance = newConnectionManager()

func ConnectionManagerInstance() ConnectionManager {
	return &connectionManagerInstance
}
