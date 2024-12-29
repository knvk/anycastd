package ntpq

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	th "github.com/teran/go-time"
)

type spec struct {
	NtpServer string      `json:"ntpserver"`
	SrcAddr   string      `json:"srcaddr"`
	Tries     uint8       `json:"tries"`
	Offset    th.Duration `json:"offset"`
	Interval  th.Duration `json:"interval"`
	Timeout   th.Duration `json:"timeout"`
}

func (s spec) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.NtpServer, validation.Required, is.Host),
		validation.Field(&s.SrcAddr, validation.Required, is.IPv4),
		validation.Field(&s.Tries, validation.Required),
		validation.Field(&s.Offset, validation.Required),
		validation.Field(&s.Interval, validation.Required),
		validation.Field(&s.Timeout, validation.Required),
	)
}
