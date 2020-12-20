package mature

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/tenntenn/natureremo"
)

type fakeDeviceService struct {
	Devices []*natureremo.Device
	Err     error
}

func (f *fakeDeviceService) GetAll(ctx context.Context) ([]*natureremo.Device, error) {
	if f.Err != nil {
		return []*natureremo.Device{}, f.Err
	}

	return f.Devices, nil
}

func (f *fakeDeviceService) Update(ctx context.Context, device *natureremo.Device) (*natureremo.Device, error) {
	return nil, nil
}

func (f *fakeDeviceService) Delete(ctx context.Context, device *natureremo.Device) error {
	return nil
}

func (f *fakeDeviceService) UpdateTemperatureOffset(ctx context.Context, device *natureremo.Device) (*natureremo.Device, error) {
	return nil, nil
}

func (f *fakeDeviceService) UpdateHumidityOffset(ctx context.Context, device *natureremo.Device) (*natureremo.Device, error) {
	return nil, nil
}

func TestNewNatureRemo(t *testing.T) {
	testcases := []struct {
		accessToken string
	}{
		{accessToken: ""},
		{accessToken: "dummy"},
	}

	for _, tc := range testcases {
		got := NewNatureRemo(tc.accessToken)
		if got == nil {
			t.Error("got nil")
		}
	}
}

func TestFetchAllNewestSensorValue(t *testing.T) {
	testcases := []struct {
		deviceID string
		devices  []*natureremo.Device
		want     map[natureremo.SensorType]natureremo.SensorValue
	}{
		{
			deviceID: "016e00cc-a76a-43be-bdfc-b305722fe4fd",
			devices: []*natureremo.Device{
				{
					DeviceCore: natureremo.DeviceCore{
						ID:   "016e00cc-a76a-43be-bdfc-b305722fe4fd",
						Name: "Living Room",
					},
					NewestEvents: map[natureremo.SensorType]natureremo.SensorValue{
						natureremo.SensorTypeTemperature: {
							Value:     24.5,
							CreatedAt: time.Date(2020, 12, 11, 10, 9, 8, 7, time.UTC),
						},
						natureremo.SensorTypeHumidity: {
							Value:     75,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 12, 13, time.UTC),
						},
						natureremo.SensortypeIllumination: {
							Value:     30,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 10, 9, time.UTC),
						},
					},
				},
				{
					DeviceCore: natureremo.DeviceCore{
						ID:   "9347d9c2-cde0-43be-b931-40b229a6645d",
						Name: "Bed Room",
					},
					NewestEvents: map[natureremo.SensorType]natureremo.SensorValue{
						natureremo.SensorTypeTemperature: {
							Value:     4.3,
							CreatedAt: time.Date(2019, 12, 11, 10, 9, 8, 7, time.UTC),
						},
						natureremo.SensorTypeHumidity: {
							Value:     35,
							CreatedAt: time.Date(2019, 12, 11, 10, 11, 12, 13, time.UTC),
						},
						natureremo.SensortypeIllumination: {
							Value:     15,
							CreatedAt: time.Date(2019, 12, 11, 10, 11, 10, 9, time.UTC),
						},
					},
				},
			},
			want: map[natureremo.SensorType]natureremo.SensorValue{
				natureremo.SensorTypeHumidity: {
					Value:     75,
					CreatedAt: time.Date(2020, 12, 11, 10, 11, 12, 13, time.UTC),
				},
				natureremo.SensorTypeTemperature: {
					Value:     24.5,
					CreatedAt: time.Date(2020, 12, 11, 10, 9, 8, 7, time.UTC),
				},
				natureremo.SensortypeIllumination: {
					Value:     30,
					CreatedAt: time.Date(2020, 12, 11, 10, 11, 10, 9, time.UTC),
				},
			},
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		client := &natureRemoImpl{
			client: &natureremo.Client{
				DeviceService: &fakeDeviceService{
					Devices: tc.devices,
				},
			},
		}

		got, err := client.FetchAllNewestSensorValue(ctx, tc.deviceID)
		if err != nil {
			t.Fatalf("want no error. got: %s", err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("want: %v. got: %v", tc.want, got)
		}
	}
}

func TestFetchAllNewestSensorValue_error(t *testing.T) {
	testcases := []struct {
		subtitle string
		deviceID string
		devices  []*natureremo.Device
		err      error
		want     string
	}{
		{
			subtitle: "api error",
			deviceID: "9347d9c2-cde0-43be-b931-40b229a6645d",
			err:      fmt.Errorf("unexpected error"),
			want:     "cannot get all devices: unexpected error",
		},
		{
			subtitle: "device not found",
			deviceID: "9347d9c2-cde0-43be-b931-40b229a6645d",
			devices: []*natureremo.Device{
				{
					DeviceCore: natureremo.DeviceCore{
						ID:   "016e00cc-a76a-43be-bdfc-b305722fe4fd",
						Name: "Living Room",
					},
					NewestEvents: map[natureremo.SensorType]natureremo.SensorValue{
						natureremo.SensorTypeTemperature: {
							Value:     24.5,
							CreatedAt: time.Date(2020, 12, 11, 10, 9, 8, 7, time.UTC),
						},
						natureremo.SensorTypeHumidity: {
							Value:     75,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 12, 13, time.UTC),
						},
						natureremo.SensortypeIllumination: {
							Value:     30,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 10, 9, time.UTC),
						},
					},
				},
			},
			want: "device does not found. device id: 9347d9c2-cde0-43be-b931-40b229a6645d",
		},
		{
			subtitle: "temperature not found",
			deviceID: "016e00cc-a76a-43be-bdfc-b305722fe4fd",
			devices: []*natureremo.Device{
				{
					DeviceCore: natureremo.DeviceCore{
						ID:   "016e00cc-a76a-43be-bdfc-b305722fe4fd",
						Name: "Living Room",
					},
					NewestEvents: map[natureremo.SensorType]natureremo.SensorValue{
						natureremo.SensorTypeHumidity: {
							Value:     75,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 12, 13, time.UTC),
						},
						natureremo.SensortypeIllumination: {
							Value:     30,
							CreatedAt: time.Date(2020, 12, 11, 10, 11, 10, 9, time.UTC),
						},
					},
				},
			},
			want: "cannot fetch temperature. device id: 016e00cc-a76a-43be-bdfc-b305722fe4fd",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subtitle, func(t *testing.T) {
			ctx := context.Background()
			client := &natureRemoImpl{
				client: &natureremo.Client{
					DeviceService: &fakeDeviceService{
						Devices: tc.devices,
						Err:     tc.err,
					},
				},
			}
			_, err := client.FetchAllNewestSensorValue(ctx, tc.deviceID)
			if err == nil {
				t.Fatal("want error. got nil")
			}

			if err.Error() != tc.want {
				t.Errorf("want: %q. got: %q", tc.want, err.Error())
			}
		})
	}
}
