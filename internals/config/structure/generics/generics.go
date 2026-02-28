package generics

import (
	"errors"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/secured-signal-api/utils/netutils"
)

// TimeDuration is a wrapper for parsing string durations using time.ParseDuration()
type TimeDuration time.Duration

func (timeDuration *TimeDuration) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

    d, err := time.ParseDuration(str)

	if err != nil {
		logger.Fatal("Invalid duration ", str, ": ", err.Error())
		return err
	}

	*timeDuration = TimeDuration(d)

	return nil
}

// IPOrNet is a wrapper for parsing 1.2.3.4 and 1.2.3.4/24 into net.IPNet (IPs are converted into A.B.C.D/32)
type IPOrNet net.IPNet

func (ipNet *IPOrNet) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

	ip, err := netutils.ParseIPorNet(str)

	if err != nil {
		logger.Fatal("Invalid IP ", str, ": ", err.Error())
		return err
	}

	*ipNet = IPOrNet(*ip)

	return nil
}

// URL is a wrapper for parsing string URLs with url.Parse()
type URL url.URL

func (Url *URL) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

	u, err := url.Parse(str)

	if err != nil {
		logger.Fatal("Invalid URL ", str, ": ", err.Error())
		return err
	}

	*Url = URL(*u)

	return nil
}

func (Url URL) String() string {
	return Url.String()
}

// Enum is a wrapper for enum types
type Enum[T interface{ ParseEnum(string) (T, bool) }] struct {
	Value T
}

func (e *Enum[T]) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

	var zero T
	value, found := zero.ParseEnum(str)

	if !found {
		logger.Fatal("Invalid enum: ", str)
		return errors.New("unsupported enum value: " + str)
	}

	e.Value = value

	return nil
}