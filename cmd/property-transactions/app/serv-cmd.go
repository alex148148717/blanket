package app

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"property_transactions/cmd/property-transactions/app/property_transactions_bl"
	"property_transactions/cmd/property-transactions/app/property_transactions_db"

	"syscall"
	"time"
)

var (
	port               int
	clickhouseAddr     string
	clickhouseDatabase string
	clickhouseUsername string
	clickhousePassword string
)

var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		clickhouseOptions := clickhouse.Options{
			Addr: []string{clickhouseAddr},
			Auth: clickhouse.Auth{
				Database: clickhouseDatabase,
				Username: clickhouseUsername,
				Password: clickhousePassword,
			},
			DialTimeout: time.Second * 10,
			Debug:       false,
		}

		propertyTransactionsDBClient, err := property_transactions_db.New(ctx, property_transactions_db.Config{ClickhouseOptions: clickhouseOptions})
		if err != nil {
			return err
		}
		propertyTransactionsClient, err := property_transactions_bl.New(propertyTransactionsDBClient)
		if err != nil {
			return err
		}
		r := chi.NewRouter()

		if _, err = New(ctx, r, propertyTransactionsClient); err != nil {
			return err
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		}

		go func() {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

			for sig := range signalChan {
				switch sig {
				case syscall.SIGHUP:
					log.Println("Reloading the service...")
					if err := srv.Shutdown(context.Background()); err != nil {
						log.Fatalf("Failed to shutdown server: %v", err)
					}
					os.Exit(0) // Exit to let systemd or a service manager restart with the updated binary
				case syscall.SIGINT, syscall.SIGTERM:
					log.Println("Shutting down the service...")
					if err := srv.Shutdown(context.Background()); err != nil {
						log.Fatalf("Failed to shutdown server: %v", err)
					}
					return
				}
			}
		}()

		if err := srv.ListenAndServe(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	ServCmd.Flags().IntVar(&port, "port", 80, "Port to run the server on")
	ServCmd.Flags().StringVar(&clickhouseAddr, "clickhouseAddr", "localhost:9000", "server Addr like HOST:PORT")
	ServCmd.Flags().StringVar(&clickhouseDatabase, "clickhouseDatabase", "default", "database")
	ServCmd.Flags().StringVar(&clickhouseUsername, "clickhouseUsername", "myuser", "user name ")
	ServCmd.Flags().StringVar(&clickhousePassword, "clickhousePassword", "mypassword", "password")

}
