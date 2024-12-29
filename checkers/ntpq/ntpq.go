package ntpq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/beevik/ntp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/runityru/anycastd/checkers"
)

var (
	_ checkers.Checker = (*ntpq)(nil)

	ErrOffset = errors.New("Offset is too big")
)

type ntpq struct {
	server   string
	srcaddr  string
	tries    uint8
	offset   time.Duration 
	interval time.Duration
	timeout  time.Duration
}

const checkName = "ntpq"

func init() {
	checkers.MustRegister(checkName, NewFromSpec)
}

func New(s spec) (checkers.Checker, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	return &ntpq{
		server:   s.NtpServer,
		srcaddr:  s.SrcAddr,
		tries:    s.Tries,
		offset:   s.Offset.TimeDuration(),
		interval: s.Interval.TimeDuration(),
		timeout:  s.Timeout.TimeDuration(),
	}, nil
}

func NewFromSpec(in json.RawMessage) (checkers.Checker, error) {
	s := spec{}
	if err := json.Unmarshal(in, &s); err != nil {
		return nil, err
	}

	return New(s)
}

func (h *ntpq) Kind() string {
	return checkName
}

func (d *ntpq) Check(ctx context.Context) error {
	var lastErr error
	for i := 0; i < int(d.tries); i++ {
		log.WithFields(log.Fields{
			"check":   checkName,
			"attempt": i + 1,
		}).Tracef("running check")

		if err := d.check(ctx); err != nil {
			lastErr = err
			log.WithFields(log.Fields{
				"check":   checkName,
				"attempt": i + 1,
			}).Infof("error received: %s", err)
		} else {
			return nil
		}

		time.Sleep(d.interval)
	}

	if lastErr != nil {
		return errors.Errorf(
			"check failed: %d tries with %s interval; last error: `%s`",
			d.tries, d.interval, lastErr.Error(),
		)
	}
	return nil
}

func (d *ntpq) check(ctx context.Context) error {
	// defaut timeout is 5s
	options := ntp.QueryOptions{LocalAddress: d.srcaddr, Timeout: d.timeout}
	response, err := ntp.QueryWithOptions(d.server, options)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"check": checkName,
	}).Tracef("Offset: %d, RTT: %d, RefID: %d", response.ClockOffset.Milliseconds(), response.RTT.Milliseconds(), response.ReferenceID)

	if response.ClockOffset > d.offset {
		return ErrOffset
	}

	return nil

}
