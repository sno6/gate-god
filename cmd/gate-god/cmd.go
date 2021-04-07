package cmd

import (
	"log"

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
			var cfg config.AppConfig
			err := gonfig.NewFromFile(&gonfig.Config{
				Path:        cfgPath,
				Environment: gonfig.Local,
			}, &cfg)
			if err != nil {
				log.Fatal(err)
			}

			relay, err := relay.New(cfg.RelayPinMCU)
			if err != nil {
				log.Fatal(err)
			}

			recognizer := platerecognizer.New(cfg.Token)
			engine := engine.New(recognizer, relay, cfg.AllowedPlates)
			batcher := batch.New(engine)

			s := ftp.New(&ftp.Config{
				User:     cfg.User,
				Password: cfg.Password,
			}, batcher)
			if err = s.Serve(); err != nil {
				log.Fatal(err)
			}
		},
	}

	return rootCmd.Execute()
}
