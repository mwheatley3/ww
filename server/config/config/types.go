package config

import (
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StringSl is a config serializable slice of string values
type StringSl []string

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (s *StringSl) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	val := string(b)

	parts := strings.Split(val, ",")
	*s = make([]string, len(parts))

	for i := range parts {
		(*s)[i] = strings.TrimSpace(parts[i])
	}

	return nil
}

// Int64Sl is a config serializable slice of int64 values
type Int64Sl []int64

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (s *Int64Sl) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	val := string(b)

	parts := strings.Split(val, ",")
	*s = make([]int64, len(parts))

	for i := range parts {
		v, err := strconv.ParseInt(strings.TrimSpace(parts[i]), 10, 64)

		if err != nil {
			return err
		}

		(*s)[i] = v
	}

	return nil
}

// FloatSl is a config serializable slice of float64 values
type FloatSl []float64

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (s *FloatSl) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	val := string(b)

	parts := strings.Split(val, ",")
	*s = make([]float64, len(parts))

	for i := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)

		if err != nil {
			return err
		}

		(*s)[i] = v
	}

	return nil
}

// RegexpSl is a config serializable slice of regexp values
type RegexpSl []*regexp.Regexp

// UnmarshalText implements the encoding.TextUnmarshaler interface
// This will panic if val is not a valid regexp
func (s *RegexpSl) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	val := string(b)

	parts := strings.Split(val, ",")
	*s = make([]*regexp.Regexp, len(parts))

	for i := range parts {
		var err error
		(*s)[i], err = regexp.Compile(strings.TrimSpace(parts[i]))

		if err != nil {
			return err
		}
	}

	return nil
}

// HexBytes is a hex encoded string
type HexBytes []byte

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (h *HexBytes) UnmarshalText(b []byte) error {
	str, err := hex.DecodeString(string(b))

	if err != nil {
		return err
	}

	*h = HexBytes(str)
	return nil
}

// Duration is an alias for time.Duration
type Duration time.Duration

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (d *Duration) UnmarshalText(b []byte) error {
	d2, err := time.ParseDuration(string(b))

	if err != nil {
		return err
	}

	*d = Duration(d2)
	return nil
}
