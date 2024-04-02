package app

import (
	"context"
	"net"
	"sync"

	"github.com/pg-sharding/spqr/pkg/spqrlog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/pg-sharding/spqr/coordinator"
	"github.com/pg-sharding/spqr/coordinator/provider"
	"github.com/pg-sharding/spqr/pkg/config"
	protos "github.com/pg-sharding/spqr/pkg/protos"

	"golang.org/x/sync/semaphore"
)

type App struct {
	coordinator coordinator.Coordinator
	sem         *semaphore.Weighted
}

const (
	maxWorkers = 50
)

func NewApp(c coordinator.Coordinator) *App {
	return &App{
		coordinator: c,
		sem:         semaphore.NewWeighted(int64(maxWorkers)),
	}
}

func (app *App) Run(withPsql bool) error {
	spqrlog.Zero.Info().Msg("running coordinator app")

	app.coordinator.RunCoordinator(context.TODO(), !withPsql)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		if err := app.ServeGrpcApi(wg); err != nil {
			spqrlog.Zero.Error().Err(err).Msg("")
		}
	}(wg)
	if withPsql {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if err := app.ServeCoordinator(wg); err != nil {
				spqrlog.Zero.Error().Err(err).Msg("")
			}
		}(wg)
	}

	wg.Wait()

	spqrlog.Zero.Debug().Msg("exit coordinator app")
	return nil
}

func (app *App) ServeCoordinator(wg *sync.WaitGroup) error {
	defer wg.Done()

	var lwg sync.WaitGroup

	listen := []string{
		"localhost:7002",
		net.JoinHostPort(config.CoordinatorConfig().Host, config.CoordinatorConfig().CoordinatorPort),
	}

	lwg.Add(len(listen))

	for _, l := range listen {
		go func(address string) {
			defer lwg.Done()

			listener, err := net.Listen("tcp", address)
			if err != nil {
				spqrlog.Zero.Error().
					Err(err).
					Msg("error serve coordinator console")
				return
			}
			spqrlog.Zero.Info().
				Str("address", address).
				Msg("serve coordinator console")

			for {
				conn, err := listener.Accept()
				if err != nil {
					spqrlog.Zero.Error().Err(err).Msg("")
					return err
				}

				if err := app.sem.Acquire(context.Background(), 1); err != nil {
					return err
				}

				go func() {
					defer app.sem.Release(1)

					err := app.coordinator.ProcClient(context.TODO(), conn)
					if err != nil {
						spqrlog.Zero.Error().Err(err).Msg("failed to serve client")
					}
				}()
			}
		}(l)
	}
	lwg.Wait()
	return nil
}

func (app *App) ServeGrpcApi(wg *sync.WaitGroup) error {
	defer wg.Done()

	serv := grpc.NewServer()
	reflection.Register(serv)

	krServ := provider.NewKeyRangeService(app.coordinator)
	rrServ := provider.NewRouterService(app.coordinator)
	topServ := provider.NewTopologyService(app.coordinator)
	shardServ := provider.NewShardServer(app.coordinator)
	dsServ := provider.NewDistributionServer(app.coordinator)
	tasksServ := provider.NewTasksServer(app.coordinator)
	protos.RegisterKeyRangeServiceServer(serv, krServ)
	protos.RegisterRouterServiceServer(serv, rrServ)
	protos.RegisterTopologyServiceServer(serv, topServ)
	protos.RegisterShardServiceServer(serv, shardServ)
	protos.RegisterDistributionServiceServer(serv, dsServ)
	protos.RegisterTasksServiceServer(serv, tasksServ)

	address := net.JoinHostPort(config.CoordinatorConfig().Host, config.CoordinatorConfig().GrpcApiPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		spqrlog.Zero.Error().
			Err(err).
			Msg("error serve grpc coordinator service")
		return err
	}

	spqrlog.Zero.Info().
		Str("address", address).
		Msg("serve grpc coordinator service")

	return serv.Serve(listener)
}
