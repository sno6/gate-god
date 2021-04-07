package config

type AppConfig struct {
	RelayPinMCU   int      `json:"relay_pin_mcu" validate:"required"`
	AllowedPlates []string `json:"allowed_plates" validate:"required"`
	Token         string   `env:"PLATE_RECOGNIZER_TOKEN" validate:"required"`
	User          string   `env:"FTP_USER" validate:"required"`
	Password      string   `env:"FTP_PASS" validate:"required"`
}
