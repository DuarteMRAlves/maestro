package execute

import (
	"context"
	"time"

	"github.com/DuarteMRAlves/maestro/internal"
)

type offlineUnary struct {
	name    internal.StageName
	address internal.Address

	input  <-chan offlineState
	output chan<- offlineState

	clientBuilder internal.UnaryClientBuilder

	logger Logger
}

func newOfflineUnary(
	name internal.StageName,
	input <-chan offlineState,
	output chan<- offlineState,
	address internal.Address,
	clientBuilder internal.UnaryClientBuilder,
	logger Logger,
) Stage {
	return &offlineUnary{
		name:          name,
		input:         input,
		output:        output,
		address:       address,
		clientBuilder: clientBuilder,
		logger:        logger,
	}
}

func (s *offlineUnary) Run(ctx context.Context) error {
	var (
		in, out offlineState
		more    bool
	)
	client, err := s.clientBuilder(s.address)
	if err != nil {
		return err
	}
	defer client.Close()
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
		s.logger.Debugf("'%s': recv msg: %#v\n", s.name, in.msg)
		req := in.msg
		rep, err := s.call(ctx, client, req)
		if err != nil {
			return err
		}

		out = newOfflineState(rep)
		s.logger.Debugf("'%s': send msg: %#v\n", s.name, out.msg)
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
	ctx context.Context,
	client internal.UnaryClient,
	req internal.Message,
) (internal.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return client.Call(ctx, req)
}

type onlineUnary struct {
	name    internal.StageName
	address internal.Address

	input  <-chan onlineState
	output chan<- onlineState

	clientBuilder internal.UnaryClientBuilder

	logger Logger
}

func newOnlineUnary(
	name internal.StageName,
	input <-chan onlineState,
	output chan<- onlineState,
	address internal.Address,
	clientBuilder internal.UnaryClientBuilder,
	logger Logger,
) Stage {
	return &onlineUnary{
		name:          name,
		input:         input,
		output:        output,
		address:       address,
		clientBuilder: clientBuilder,
		logger:        logger,
	}
}

func (s *onlineUnary) Run(ctx context.Context) error {
	var (
		in, out onlineState
		more    bool
	)
	client, err := s.clientBuilder(s.address)
	if err != nil {
		return err
	}
	defer client.Close()
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
		s.logger.Debugf("'%s': recv id: %d, msg: %#v\n", s.name, in.id, in.msg)
		req := in.msg
		rep, err := s.call(ctx, client, req)
		if err != nil {
			return err
		}

		out = fromOnlineState(in, rep)
		s.logger.Debugf("'%s': send id: %d, msg: %#v\n", s.name, out.id, out.msg)
		select {
		case s.output <- out:
		case <-ctx.Done():
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
	}
}

func (s *onlineUnary) call(
	ctx context.Context,
	client internal.UnaryClient,
	req internal.Message,
) (internal.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return client.Call(ctx, req)
}
