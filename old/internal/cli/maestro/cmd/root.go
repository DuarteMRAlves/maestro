package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net"
)

type Options struct {
	// logLvl specifies the level at which the server should log.
	logLvl string
	// address to listen
	addr string
}

func NewCmdRoot() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "maestro",
		Short: "maestro is a server to orchestrate grpc services into pipelines.",
		Run: func(cmd *cobra.Command, args []string) {
			// save err for later log if necessary
			lvl, logErr := o.getZapLevel()

			logger, err := logs.DefaultProductionLogger(lvl)
			// Should never happen
			if err != nil {
				panic(err)
			}
			if logErr != nil {
				logger.Warn(
					"Invalid log level.",
					zap.String("received", o.logLvl),
				)
			}
			logger.Info(
				"Server logging level.",
				zap.String("value", lvl.String()),
			)

			db, err := storage.NewDb()
			// Should never happen
			if err != nil {
				panic(err)
			}

			lis, err := net.Listen("tcp", o.addr)
			if err != nil {
				logger.Fatal("Listen.", zap.Error(err))
			}
			logger.Info("Listen.", zap.String("address", lis.Addr().String()))

			s, err := server.NewBuilder().
				WithGrpc().
				WithLogger(logger).
				WithDb(db).
				Build()
			if err != nil {
				logger.Fatal("Build", zap.Error(err))
			}

			if err := s.ServeGrpc(lis); err != nil {
				logger.Fatal("Serve", zap.Error(err))
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

func (o *Options) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.logLvl, "log", "info", "Zap logging level.")
	cmd.Flags().StringVar(&o.addr, "addr", "0.0.0.0:50051", "Address to listen")
}

func (o *Options) getZapLevel() (zap.AtomicLevel, error) {
	if o.logLvl == "" {
		o.logLvl = "info"
	}
	lvl, err := zap.ParseAtomicLevel(o.logLvl)
	if err != nil {
		return zap.NewAtomicLevelAt(zap.InfoLevel), err
	}
	return lvl, nil
}
