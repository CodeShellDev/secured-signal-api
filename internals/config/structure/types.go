package structure

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

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

type AllowBlockSlice struct{
	Allow	[]string
	Block	[]string
}

func (splitter *AllowBlockSlice) UnmarshalMapstructure(raw any) error {
    slice, ok := raw.([]any)

    if !ok {
		fmt.Println(raw)
		return errors.New("expected []string, got " + reflect.TypeOf(raw).String())
    }

	for _, item := range slice {
        str, ok := item.(string)

        if !ok {
			return errors.New("expected string, got " + reflect.TypeOf(item).String())
        }

		str, block := strings.CutPrefix(str, "!")

		if block {
			splitter.Block = append(splitter.Block, str)
		} else {
			splitter.Allow = append(splitter.Allow, str)
		}
	}

	return nil
}

type FieldPolicies struct{
	Allow	[]FieldPolicy
	Block	[]FieldPolicy
}

func (splitter *FieldPolicies) UnmarshalMapstructure(raw any) error {
    slice, ok := raw.([]any)

    if !ok {
		fmt.Println(raw)
		return errors.New("expected []FieldPolicy, got " + reflect.TypeOf(raw).String())
    }

	for _, item := range slice {
        policy, ok := item.(FieldPolicy)

        if !ok {
			return errors.New("expected string, got " + reflect.TypeOf(item).String())
        }

		switch strings.ToLower(policy.Action) {
		case "block":
			splitter.Block = append(splitter.Block, policy)
		case "allow":
			splitter.Allow = append(splitter.Allow, policy)
		}
	}

	return nil
}