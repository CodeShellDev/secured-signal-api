package structure

import (
	"errors"
	"net"
	"net/url"
	"reflect"
	"time"

	"github.com/codeshelldev/secured-signal-api/utils/netutils"
)

// TimeDuration is a wrapper struct used to parse string durations using time.ParseDuration()
type TimeDuration struct {
	Duration time.Duration
}

func (timeDuration *TimeDuration) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

    d, err := time.ParseDuration(str)

	if err != nil {
		return err
	}

	timeDuration.Duration = d

	return nil
}

// IPOrNet is a wrapper struct used to parse 1.2.3.4 and 1.2.3.4/24 into net.IPNet (IPs are converted into A.B.C.D/32)
type IPOrNet struct {
	IPNet *net.IPNet
}

func (ipNet *IPOrNet) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

	ip, err := netutils.ParseIPorNet(str)

	if err != nil {
		return err
	}

	ipNet.IPNet = ip

	return nil
}

// URL is a wrapper struct used to parse string URLs with url.Parse()
type URL struct {
	URL *url.URL
}

func (Url *URL) UnmarshalMapstructure(raw any) error {
	str, ok := raw.(string)

	if !ok {
		return errors.New("expected string, got " + reflect.TypeOf(raw).String())
	}

	u, err := url.Parse(str)

	if err != nil {
		return err
	}

	Url.URL = u

	return nil
}

func (Url URL) String() string {
	return Url.URL.String()
}