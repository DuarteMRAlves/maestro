package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net"
)

type Options struct {
	// address to listen
	addr string
}

func NewCmdRoot() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "maestro",
		Short: "maestro is a server to orchestrate grpc services into pipelines.",
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := logs.DefaultProductionLogger()
			// Should never happen
			if err != nil {
				panic(err)
			}

			db, err := kv.NewDb()
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
	cmd.Flags().StringVar(&o.addr, "addr", "0.0.0.0:50051", "address to listen")
}
