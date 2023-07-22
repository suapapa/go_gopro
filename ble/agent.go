package ble

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/pkg/errors"
)

var (
	ag *agent.SimpleAgent
)

func getAgent() (*agent.SimpleAgent, error) {
	if ag == nil {
		conn, err := dbus.SystemBus()
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to system bus")
		}

		ag = agent.NewSimpleAgent()
		err = agent.ExposeAgent(conn, ag, agent.CapNoInputNoOutput, true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to expose agent")
		}
	}
	return ag, nil
}
