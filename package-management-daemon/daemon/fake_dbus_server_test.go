package daemon

import (
	"fmt"
	"github.com/godbus/dbus"
)

// FakeDbusServer is a fake implementation of the DbusWrapper interface,
// for use within tests.
type FakeDbusServer struct {
	connectCalled      bool
	namesCalled        bool
	requestNameCalled  bool
	getNameOwnerCalled bool
	exportCalled       bool
	emitCalled         bool

	failConnect      bool
	failNames        bool
	failRequestName  bool
	failGetNameOwner bool
	failExport       bool
	failEmit         bool

	nameAlreadyTaken            bool
	failSpecificExportInterface string
	signals                     chan *dbus.Signal
}

func (server *FakeDbusServer) InitializeSignals() {
	server.signals = make(chan *dbus.Signal)
}

func (server *FakeDbusServer) Connect() error {
	server.connectCalled = true

	if server.failConnect {
		return fmt.Errorf("Failed at user request")
	}

	return nil
}

func (server *FakeDbusServer) Names() []string {
	server.namesCalled = true

	if server.failNames {
		return nil
	}

	return []string{":1.42"}
}

func (server *FakeDbusServer) RequestName(name string, flags dbus.RequestNameFlags) (dbus.RequestNameReply, error) {
	server.requestNameCalled = true

	if server.failRequestName {
		return 0, fmt.Errorf("Failed at user request")
	}

	if server.nameAlreadyTaken {
		return dbus.RequestNameReplyInQueue, nil
	}

	return dbus.RequestNameReplyPrimaryOwner, nil
}

func (server *FakeDbusServer) GetNameOwner(name string) (string, error) {
	server.getNameOwnerCalled = true

	if server.failGetNameOwner {
		return "", fmt.Errorf("Failed at user request")
	}

	return ":1.42", nil
}

func (server *FakeDbusServer) Export(object interface{}, path dbus.ObjectPath, iface string) error {
	server.exportCalled = true

	if server.failExport {
		if server.failSpecificExportInterface == "" || iface == server.failSpecificExportInterface {
			return fmt.Errorf("Failed at user request")
		}
	}

	return nil
}

func (server *FakeDbusServer) Emit(path dbus.ObjectPath, name string, values ...interface{}) error {
	server.emitCalled = true

	if server.failEmit {
		return fmt.Errorf("Failed at user request")
	}

	if server.signals != nil {
		server.signals <- &dbus.Signal{Path: path, Name: name, Body: values}
	}

	return nil
}
