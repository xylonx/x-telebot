package session

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrorNilEncoder   = errors.New("encoder is nil")
	ErrorNilDecoder   = errors.New("decoder is nil")
	ErrorNilSessionID = errors.New("sessionID is nil")
	ErrorExpire       = errors.New("session is expire")
	ErrorNilSession   = errors.New("session is not found")
)

type Option struct {
	Expire time.Duration
	// Encoder - encode input value and persistent
	Encoder Encoder
	// Decodeer - decode persistented data
	Decoder Decoder
}

type OptionFunc func(*Option)

type Encoder interface {
	Encode(interface{}) (string, error)
}
type Decoder interface {
	Decode(string) (interface{}, error)
}

// Session - maintain the session about a user or group
type Session interface {
	// Set - store value of this session
	Set(value interface{}, opt Option) error

	// Get - get persistent value of this session set before
	Get() (value interface{}, err error)

	// Close - close the session
	Close() error
}

type SessionManager interface {
	// New - create a new session
	// it will start a new session
	// if the session is created before, it will overwrite it
	New(sessionID string, opts ...OptionFunc) (Session, error)

	// Load - get the session
	// if the session is created before, it will load the existing session
	// else return error
	Load(sessionID string, opts ...OptionFunc) (Session, error)
}

func defaultOption() *Option {
	return &Option{
		Expire:  -1,
		Encoder: &jsonEncoder{},
		Decoder: &jsonDecoder{},
	}
}

type jsonEncoder struct{}
type jsonDecoder struct{}

var _ Encoder = &jsonEncoder{}
var _ Decoder = &jsonDecoder{}

func (*jsonEncoder) Encode(data interface{}) (string, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (*jsonDecoder) Decode(data string) (interface{}, error) {
	var value interface{}
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return nil, err
	}
	return value, nil
}

func WithExpire(expire time.Duration) OptionFunc {
	return func(o *Option) {
		o.Expire = expire
	}
}
