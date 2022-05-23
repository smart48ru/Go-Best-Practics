package configuration

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Configuration interface {
	PrintHelp()
	Helper() bool
	MaxDepth() int64
	JsonLog() bool
	FileExt() string
	LogLevel() zerolog.Level
	String() (result string)
}

type configuration struct {
	helper     bool
	maxDepth   int64         `yaml:"max_depth"`
	jsonOutput bool          `yaml:"json_log"`
	fileExt    string        `yaml:"file_ext"`
	logLevel   zerolog.Level `yaml:"log_level"`
}

const (
	DefaultMaxDepth = 2
	DefaultNeedHelp = false
	DefaultJSONLog  = false
	DefaultLogLevel = "info"
	DefaultFileExt  = ".go"
)

func flagsRead() (c configuration, err error) {
	var ll string

	flag.BoolVar(&c.helper, "h", DefaultNeedHelp, "Print help")
	flag.Int64Var(&c.maxDepth, "d", DefaultMaxDepth, "Max depth")
	flag.StringVar(&c.fileExt, "e", DefaultFileExt, "File extension")
	flag.BoolVar(&c.jsonOutput, "j", DefaultJSONLog, "JSON log format")
	flag.StringVar(&ll, "l", DefaultLogLevel, "Log level")
	flag.Parse()

	c.configureLogger() // костыль, что бы вывести красиво ошибку, как сделать лучше я не разобрался :(
	c.logLevel, err = zerolog.ParseLevel(ll)
	if err != nil {
		c.logLevel = 1
		log.Error().Msgf("Ошибка установки уровня логирования | Установлен уровень логирования по умолчанию = %s", DefaultLogLevel)
		//return
	}
	return
}

func (c *configuration) loadFromFlags(cf configuration) {
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

func New() configuration {
	flagConf, err := flagsRead()
	if err != nil {
		log.Error().Err(err)
	}
	cfg := configuration{}
	cfg.loadFromFlags(flagConf)
	cfg.configureLogger()
	return cfg
}

func (c *configuration) Helper() bool {
	return c.helper
}

func (c *configuration) MaxDepth() int64 {
	return c.maxDepth
}

func (c *configuration) JsonLog() bool {
	return c.jsonOutput
}

func (c *configuration) FileExt() string {
	return c.fileExt
}

func (c *configuration) LogLevel() zerolog.Level {
	return c.logLevel
}

func (c configuration) String() (result string) {
	return fmt.Sprintf("%#v", c)
}

func (c *configuration) PrintHelp() {
	fmt.Printf(`
Программа Filescanner V 0.0.1:
Сканирует каталог и подкаталоги на заданную глубину для поиска файлов с определенным расширением

Ключи запуска:
-h вывод подсказки
-d максимальная глубина поиска
-e расширение файла 
-j вывод лога в формате JSON
-l уровень логирования :
			panic (PanicLevel, 5)
			fatal (FatalLevel, 4)
			error (ErrorLevel, 3)
			warn (WarnLevel, 2)
			info (InfoLevel, 1)
			debug (DebugLevel, 0)
			trace (TraceLevel, -1)

Пример использования:
bin/filescanner -m 5 ## Глубина поиска 5 (default = %d)
bin/filescanner -m j ## Выводит лог в формате JSON (default = Stdout)
иin/filescanner -e .go ## Выводит лог в формате JSON (default = %s)
bin/filescanner -l info ## Устанавливает уровень логирования info (default = %s)
`, DefaultMaxDepth, DefaultFileExt, DefaultLogLevel)
}

func (c *configuration) configureLogger() {
	zerolog.SetGlobalLevel(c.logLevel)

	if c.jsonOutput {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		return
	}

	output := zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: false,
		FormatTimestamp: func(i interface{}) string {
			parse, err := time.Parse(time.RFC3339Nano, i.(string))
			if err != nil {
				log.Error().Msgf("%s", err)
			}
			return parse.Format("2006-01-02 15:04:05 Z07:00")
		},
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
