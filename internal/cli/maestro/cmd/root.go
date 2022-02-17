package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/kv"
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
			logger, err := zap.NewProduction()
			// Should never happen
			if err != nil {
				panic(err)
			}
			sugar := logger.Sugar()

			db, err := kv.NewDb()
			// Should never happen
			if err != nil {
				panic(err)
			}

			lis, err := net.Listen("tcp", o.addr)
			if err != nil {
				sugar.Fatal("Failed to listen.", "err", err)
			}
			sugar.Infof("Server listening at: %v", lis.Addr())

			s, err := server.NewBuilder().
				WithGrpc().
				WithLogger(logger).
				WithDb(db).
				Build()
			if err != nil {
				sugar.Fatalf("build server: %v", err)
			}

			if err := s.ServeGrpc(lis); err != nil {
				sugar.Fatalf("Failed to serve: %s", err)
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

func (o *Options) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.addr, "addr", "0.0.0.0:50051", "address to listen")
}
