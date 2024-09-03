package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	app "github.com/Arcadian-Sky/datakkeeper/internal/app/server"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/handler"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/router"
	"github.com/sirupsen/logrus"
)

func main() {

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ap, err := app.NewApp(&ctx, &cancel)
	if err != nil {
		ap.Logger.Fatal(err)
	}
	defer ap.DBPG.Close()

	go handleSignals(ap.CncF, ap.Logger)

	//set user repo
	repo := repository.NewUserRepository(ap.DBPG, ap.Logger)
	ap.SetUserRepo(repo)

	//set data repo
	repod := repository.NewDataRepository(ap.DBPG, ap.Logger)
	ap.SetDataRepo(repod)

	err = ap.MigrateDBPG()
	if err != nil {
		ap.Logger.Fatal(err)
	}
	// mgrepo := repository.NewDataRepository(ap.DBMG, ap.Logger)
	// ap.SetDBMGRepo(mgrepo)

	frepo := repository.NewFileRepository(ap.Storage, ap.Logger, &ap.Ctx)
	ap.SetDBFileRepo(frepo)

	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	server, err := router.InitGRPCServer(
		ap.Flags,
		ap.Logger,
		ap.GetFileRepo(),
		ap.GetUserRepo(),
		ap.GetDataRepo(),
	)

	//set handlers
	handler.NewHandler(ap)

	// ap.Logger.WithField("Flags", *ap.Flags).Info("App init")

	// go func() {
	// 	ap.Logger.Info("Start ListenAndServe")
	// 	ap.Logger.Fatal(http.ListenAndServe(ap.Flags.Endpoint, router.InitRouter(*vhandler)))
	// }()

	go func() {
		ap.Logger.Info("Start ListenAndServe")
		// start the server
		if err = server.Start(); err != nil {
			ap.Logger.Fatal(err)
		}
		// ap.Logger.Fatal(http.ListenAndServe(ap.Flags.Endpoint, router.InitRouter(*vhandler)))
	}()

	<-ap.Ctx.Done()
	if err = server.ShutDown(); err != nil {
		// ошибки закрытия Listener
		ap.Logger.Printf("gRPC server shutdown err: %v", err)
	}
	ap.Logger.Info("Server stopped gracefully")
}

func handleSignals(cancel context.CancelFunc, logger *logrus.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал
	sig := <-sigs
	logger.Info("Received signal: ", sig)
	cancel()
}
