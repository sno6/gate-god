package cmd

import (
	"log"

	"github.com/sno6/gate-god/api"

	"github.com/joho/godotenv"
	"github.com/sno6/gate-god/camera/batch"
	"github.com/sno6/gate-god/config"
	"github.com/sno6/gate-god/engine"
	"github.com/sno6/gate-god/recognition/platerecognizer"
	"github.com/sno6/gate-god/relay"
	"github.com/sno6/gate-god/server/ftp"
	"github.com/spf13/cobra"

	gonfig "github.com/sno6/config"
)

const cfgPath = "./config"

func Run() error {
	rootCmd := &cobra.Command{
		Use:   "gate-god",
		Short: "Run the god of gate controllers.",
		Run: func(cmd *cobra.Command, args []string) {
			err := godotenv.Load()
			if err != nil {
				log.Fatal(err)
			}

			var cfg config.AppConfig
			err = gonfig.NewFromFile(&gonfig.Config{
				Path:        cfgPath,
				Environment: gonfig.Local,
			}, &cfg)
			if err != nil {
				log.Fatal(err)
			}

			rly, err := relay.New(cfg.RelayPinMCU)
			if err != nil {
				log.Fatal(err)
			}

			recognizer := platerecognizer.New(cfg.Token)
			eng := engine.New(recognizer, rly, cfg.AllowedPlates)
			batcher := batch.New(eng)

			errChan := make(chan error)

			go func() {
				ftpServer := ftp.New(&ftp.Config{
					User:     cfg.User,
					Password: cfg.Password,
				}, batcher)

				err := ftpServer.Serve()
				if err != nil {
					errChan <- err
				}
			}()

			go func() {
				apiServer := api.New(rly)

				err := apiServer.Serve(cfg.HTTPPort)
				if err != nil {
					errChan <- err
				}
			}()

			// If any of the servers error, fail.
			log.Fatal(<-errChan)
		},
	}

	return rootCmd.Execute()
}
