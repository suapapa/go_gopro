package ble

import (
	"context"
	"time"

	goble "github.com/go-ble/ble"
	"github.com/pkg/errors"
)

type GoPro struct {
	cln goble.Client
}

func ScanGoPro() (*GoPro, error) {
	ctx := goble.WithSigHandler(context.WithTimeout(context.Background(), 5*time.Second))
	filter := func(a goble.Advertisement) bool {
		svcs := a.Services()
		for _, svc := range svcs {
			if svc.Equal(serviceUUIDControlAndQuery) {
				return true
			}
		}
		return false
	}

	cln, err := goble.Connect(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to GoPro")
	}

	return &GoPro{
		cln: cln,
	}, nil
}

func (g *GoPro) Close() (<-chan struct{}, error) {
	err := g.cln.CancelConnection()
	if err != nil {
		return nil, errors.Wrap(err, "failed to cancel connection")
	}

	return g.cln.Disconnected(), nil
}
