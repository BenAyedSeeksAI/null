package zero

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Bool is a nullable bool. False input is considered null.
// JSON marshals to false if null.
// Considered null to SQL unmarshaled from a false value.
type Bool struct {
	sql.NullBool
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

// BoolFrom creates a new Bool that will be null if false.
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

// BoolFromPtr creates a new Bool that be null if b is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// NewBoolFromString creates a new Int from a string
func NewBoolFromString(s string, valid bool) Bool {
	if valid {
		if s == "" {
			return NewBool(false, false)
		}
		if s == "1" || s == "true" {
			return NewBool(true, true)
		} else if s == "0" || s == "false" {
			return NewBool(false, true)
		}
	}
	return NewBool(false, false)
}

// BoolFromPtr creates a new Bool that be null if b is nil.
func BoolFromString(s string) Bool {
	if s == "" {
		return NewBool(false, true)
	}
	if s == "1" || s == "true" || s == "True" || s == "TRUE" {
		return NewBool(true, true)
	} else if s == "0" || s == "false" || s == "False" || s == "FALSE" {
		return NewBool(false, true)
	}
	return NewBool(false, false)

}

// BoolFromPtr creates a new Bool that be null if b is nil.
func BoolFromStringExist(s string, b bool) Bool {
	if s == "" || b == false {
		return NewBool(false, true)
	}
	return NewBool(true, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// "false" will be considered a null Bool.
// It also supports unmarshalling a sql.NullBool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = (err == nil) && b.Bool
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a false or not a bool.
// It will return an error if the input is not a float, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = b.Bool
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a poBooler to this Bool's value, or a nil poBooler if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for null or zero Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsZero() bool {
	return !b.Valid || !b.Bool
}

// OverwriteWithIfValid returns nothing. Used for type conversion from sql.Nullstring to zero
func (s *Bool) OverwriteWithIfValid(st bool, v bool) {
	if v {
		s.Bool = st
		s.Valid = v
	}
}

// Add boolean operators
// AND operation
func (s Bool) AND(other Bool) Bool {
	result := Bool{
		NullBool: sql.NullBool{},
	}
	if s.Valid && other.Valid {
		result.Bool = s.Bool && other.Bool
		result.Valid = true
		return result
	}
	result.Valid = false
	return result
}

// OR operation
func (s Bool) OR(other Bool) Bool {
	result := Bool{
		NullBool: sql.NullBool{},
	}
	if s.Valid && other.Valid {
		result.Bool = s.Bool || other.Bool
		result.Valid = true
		return result
	}
	result.Valid = false
	return result
}

// NON operation
func (s *Bool) NON() {
	if s.Valid {
		s.Bool = !s.Bool
	}
}

// XOR operation
func (s Bool) XOR(other Bool) Bool {
	result := Bool{
		NullBool: sql.NullBool{},
	}
	x := s.Bool
	y := other.Bool
	if s.Valid && other.Valid {
		result.Bool = (x || y) && !(x && y)
		result.Valid = true
		return result
	}
	result.Valid = false
	return result
}
