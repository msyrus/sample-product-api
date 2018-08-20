package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/msyrus/simple-product-inv/config"
	"github.com/msyrus/simple-product-inv/infra/pgsql"
	"github.com/msyrus/simple-product-inv/log"
	"github.com/msyrus/simple-product-inv/repo"
	"github.com/msyrus/simple-product-inv/service"
	"github.com/msyrus/simple-product-inv/web"
)

var cfgPath string

// srvCmd is the serve sub command to start the api server
var srvCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves the api server",
	RunE:  serve,
}

func init() {
	srvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yml", "config file path")
	rootCmd.AddCommand(srvCmd)
}

func serve(cmd *cobra.Command, args []string) error {
	lgr := log.DefaultOutputLogger
	errorLgr := log.DefaultOutputLogger
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}

	cfg, err := config.Parse(f)
	if err != nil {
		return err
	}

	addr := cfg.Host
	if cfg.Port != 0 {
		addr = addr + ":" + strconv.Itoa(cfg.Port)
	}

	db, err := sql.Open("postgres", cfg.Postgres.URI)
	if err != nil {
		return err
	}

	pg := pgsql.NewDB(db, nil)

	ratSvc := service.NewRating(repo.NewCritic("ratings", pg))
	pdtSvc := service.NewProduct(repo.NewChef("products", pg), ratSvc)
	sysSvc := service.NewSystem()

	r := chi.NewMux()
	r.Mount("/api/v1", web.NewRouter(web.NewProductController(pdtSvc), web.NewSystemController(sysSvc)))

	srvr := http.Server{
		Addr:         addr,
		Handler:      r,
		ErrorLog:     errorLgr,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		lgr.Println("Server Listening on", addr)

		if err := srvr.ListenAndServe(); err != nil {
			lgr.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)
	<-stop

	lgr.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulWait)
	defer cancel()

	srvr.Shutdown(ctx)

	lgr.Println("Server shutteddown gracefully")
	return nil
}
