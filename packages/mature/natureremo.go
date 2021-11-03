package mature

import (
	"context"

	"github.com/tenntenn/natureremo"
	"golang.org/x/xerrors"
)

var SensorTypeString = map[natureremo.SensorType]string{
	natureremo.SensorTypeTemperature:  "temperature",
	natureremo.SensorTypeHumidity:     "humidity",
	natureremo.SensorTypeIllumination: "illumination",
}

type NatureRemo interface {
	FetchAllNewestSensorValue(ctx context.Context, deviceID string) (map[natureremo.SensorType]natureremo.SensorValue, error)
}

type natureRemoImpl struct {
	client *natureremo.Client
}

func NewNatureRemo(accessToken string) *natureRemoImpl {
	return &natureRemoImpl{client: natureremo.NewClient(accessToken)}
}

func (nr *natureRemoImpl) FetchAllNewestSensorValue(ctx context.Context, deviceID string) (map[natureremo.SensorType]natureremo.SensorValue, error) {
	ds, err := nr.client.DeviceService.GetAll(ctx)
	if err != nil {
		return nil, xerrors.Errorf("cannot get all devices: %w", err)
	}

	var events map[natureremo.SensorType]natureremo.SensorValue
	for _, d := range ds {
		if d.ID != deviceID {
			continue
		}
		events = d.NewestEvents
		break
	}
	if events == nil {
		return nil, xerrors.Errorf("device does not found. device id: %s", deviceID)
	}

	vs := make(map[natureremo.SensorType]natureremo.SensorValue, len(events))
	for t, str := range SensorTypeString {
		v, ok := events[t]
		if !ok {
			return nil, xerrors.Errorf("cannot fetch %s. device id: %s", str, deviceID)
		}
		vs[t] = v
	}

	return vs, nil
}
