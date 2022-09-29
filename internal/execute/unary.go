package execute

import (
	"context"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
)

type offlineUnary struct {
	name compiled.StageName

	input  <-chan offlineState
	output chan<- offlineState

	dialer method.Dialer

	logger Logger
}

func newOfflineUnary(
	name compiled.StageName,
	input <-chan offlineState,
	output chan<- offlineState,
	dialer method.Dialer,
	logger Logger,
) Stage {
	return &offlineUnary{
		name:   name,
		input:  input,
		output: output,
		dialer: dialer,
		logger: logger,
	}
}

func (s *offlineUnary) Run(ctx context.Context) error {
	var (
		in, out offlineState
		more    bool
	)
	conn, err := s.dialer.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()
	s.logger.Infof("'%s': started\n", s.name)
	for {
		select {
		case in, more = <-s.input:
		case <-ctx.Done():
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
		// channel is closed
		if !more {
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
		s.logger.Debugf("'%s': recv msg: %v\n", s.name, in.msg)
		req := in.msg
		rep, err := s.call(ctx, conn, req)
		if err != nil {
			return err
		}

		out = newOfflineState(rep)
		s.logger.Debugf("'%s': send msg: %v\n", s.name, out.msg)
		select {
		case s.output <- out:
		case <-ctx.Done():
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
	}
}

func (s *offlineUnary) call(
	ctx context.Context, conn method.Conn, req message.Instance,
) (message.Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return conn.Call(ctx, req)
}
