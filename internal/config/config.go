package configuration

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration interface {
	PrintHelp()
	Helper() bool
	MaxDepth() int64
	JSONnLog() bool
	FileExt() string
	LogLevel() zerolog.Level
	String() (result string)
}

type Config struct {
	helper     bool
	jsonOutput bool
	configFile string
	fileExt    string
	logLevel   zerolog.Level
	maxDepth   int64
}

const (
	DefaultMaxDepth   = 2
	DefaultNeedHelp   = false
	DefaultJSONLog    = false
	DefaultLogLevel   = "info"
	DefaultFileExt    = ".go"
	DefaultConfigFile = ""
)

func flagsRead() (c Config, err error) {
	var ll string

	pflag.BoolVar(&c.helper, "h", DefaultNeedHelp, "Print help")
	pflag.Int64Var(&c.maxDepth, "d", DefaultMaxDepth, "Max depth")
	pflag.StringVar(&c.fileExt, "e", DefaultFileExt, "File extension")
	pflag.BoolVar(&c.jsonOutput, "j", DefaultJSONLog, "JSON log format")
	pflag.StringVar(&ll, "l", DefaultLogLevel, "Log level")
	pflag.StringVar(&c.configFile, "c", DefaultConfigFile, "Config file JSON, TOML, YAML, HCL, envfile")
	pflag.Parse()
	c.logLevel, err = zerolog.ParseLevel(ll)
	if err != nil {
		c.logLevel = 1
		log.Error().Msgf("Ошибка установки уровня логирования | Установлен уровень логирования по умолчанию = %s", DefaultLogLevel)
	}

	return
}

// func flagsRead() (c config, err error) {
//	var ll string
//
//	flag.BoolVar(&c.helper, "h", DefaultNeedHelp, "Print help")
//	flag.Int64Var(&c.maxDepth, "d", DefaultMaxDepth, "Max depth")
//	flag.StringVar(&c.fileExt, "e", DefaultFileExt, "File extension")
//	flag.BoolVar(&c.jsonOutput, "j", DefaultJSONLog, "JSON log format")
//	flag.StringVar(&ll, "l", DefaultLogLevel, "Log level")
//	flag.Parse()
//
//	c.configureLogger() // костыль, что бы вывести красиво ошибку, как сделать лучше я не разобрался :(
//	c.logLevel, err = zerolog.ParseLevel(ll)
//	if err != nil {
//		c.logLevel = 1
//		log.Error().Msgf("Ошибка установки уровня логирования | Установлен уровень логирования по умолчанию = %s", DefaultLogLevel)
//		//return
//	}
//	return
//}

func (c *Config) loadFromFlags(cf Config) {
	if cf.helper {
		c.helper = cf.helper
	}
	if cf.jsonOutput {
		c.jsonOutput = cf.jsonOutput
	}
	if cf.logLevel != 0 {
		c.logLevel = cf.logLevel
	}
	if cf.maxDepth != 0 {
		c.maxDepth = cf.maxDepth
	}
	if cf.fileExt != "" {
		c.fileExt = cf.fileExt
	}
}

func (c *Config) loadFromViper(cf Config) {
	log.Info().Msg("Load from viper")
	viper.SetConfigName(c.configFile)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Msgf("Fatal error config file: %v \n", err)
		c.loadFromFlags(*c)

		return
	}
	if cf.helper {
		c.helper = cf.helper
	}
	if viper.GetString("file_ext") != "" {
		cf.fileExt = viper.GetString("file_ext")
		c.fileExt = cf.fileExt
	}
	if viper.GetBool("json_output") {
		cf.jsonOutput = viper.GetBool("json_output")
		c.jsonOutput = cf.jsonOutput
	}
	if viper.GetString("log_level") != "" {
		ll, err := zerolog.ParseLevel(viper.GetString("log_level"))
		if err != nil {
			cf.logLevel = 1
			c.logLevel = cf.logLevel
			log.Error().Msgf("Ошибка установки уровня логирования | Установлен уровень логирования по умолчанию = %s", DefaultLogLevel)
		}
		cf.logLevel = ll
		c.logLevel = cf.logLevel
	}
	if viper.GetInt("max_depth") != 0 {
		cf.maxDepth = viper.GetInt64("max_depth")
		c.maxDepth = cf.maxDepth
	}
}

func New() (cfg Config) {
	flagConf, err := flagsRead()
	if err != nil {
		log.Error().Err(err)
	}
	cfg = Config{} //nolint:exhaustivestruct
	if flagConf.configFile != "" {
		cfg.loadFromViper(flagConf)
		cfg.configureLogger()

		return cfg
	}
	cfg.loadFromFlags(flagConf)
	cfg.configureLogger()

	return cfg
}

func (c *Config) Helper() bool {
	return c.helper
}

func (c *Config) MaxDepth() int64 {
	return c.maxDepth
}

func (c *Config) JSONnLog() bool {
	return c.jsonOutput
}

func (c *Config) FileExt() string {
	return c.fileExt
}

func (c *Config) LogLevel() zerolog.Level {
	return c.logLevel
}

func (c Config) String() (result string) {
	return fmt.Sprintf("%#v", c)
}

func (c *Config) PrintHelp() {
	fmt.Printf(`
Программа Filescanner V 0.0.1:
Сканирует каталог и подкаталоги на заданную глубину для поиска файлов с определенным расширением

Ключи запуска:
--h вывод подсказки
--d максимальная глубина поиска
--e расширение файла
--j вывод лога в формате JSON
--l уровень логирования :
			panic (PanicLevel, 5)
			fatal (FatalLevel, 4)
			error (ErrorLevel, 3)
			warn (WarnLevel, 2)
			info (InfoLevel, 1)
			debug (DebugLevel, 0)
			trace (TraceLevel, -1)
--c имя конфигурационного файла
	поддерживает JSON, TOML, YAML, HCL, envfile

Пример использования:
bin/filescanner --d 5 ## Глубина поиска 5 (default = %d)
bin/filescanner --j ## Выводит лог в формате JSON (default = Stdout)
иin/filescanner --e .go ## Выводит лог в формате JSON (default = %s)
bin/filescanner --l info ## Устанавливает уровень логирования info (default = %s)
`, DefaultMaxDepth, DefaultFileExt, DefaultLogLevel)
}

func (c *Config) configureLogger() {
	zerolog.SetGlobalLevel(c.logLevel)
	if c.jsonOutput {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		return
	}

	output := zerolog.ConsoleWriter{ //nolint:exhaustivestruct
		Out:     os.Stdout,
		NoColor: false,
		FormatTimestamp: func(i interface{}) string {
			parse, err := time.Parse(time.RFC3339, i.(string))
			if err != nil {
				log.Error().Msgf("Time format parse error: %s", err)
			}

			return parse.Format("2006-01-02 15:04:05 Z07:00")
		},
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
