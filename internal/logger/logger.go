package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет тип логгера.
type Logger = *zap.SugaredLogger

// Config представляет конфигурацию для логгера.
type Config struct {
	Level     string `env:"LOG_LEVEL" envDefault:"info"`        // Уровень логиров��ния
	LogToFile bool   `env:"LOG_TO_FILE" envDefault:"false"`     // Логировать в файл
	FilePath  string `env:"LOG_FILE_PATH" envDefault:"app.log"` // Путь к файлу логов
}

// New создает новый экземпляр логгера на основе конфигурации.
func New(config Config) (Logger, error) {
	var level zapcore.Level
	// Разбираем уровень логирования из конфигурации.
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	// Настраиваем конфигурацию кодировщика.
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",                           // Ключ для времени
		LevelKey:       "level",                        // Ключ для уровня логирования
		NameKey:        "logger",                       // Ключ для имени логгера
		MessageKey:     "message",                      // Ключ для сообщения
		StacktraceKey:  "stacktrace",                   // Ключ для трассировки стека
		CallerKey:      "caller",                       // Ключ для вызывающего
		LineEnding:     zapcore.DefaultLineEnding,      // Окончание строки
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // Кодировщик уровня логирования
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // Кодировщик времени
		EncodeDuration: zapcore.SecondsDurationEncoder, // Кодировщик длит��льности
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Кодировщик вызывающего
	}

	var core zapcore.Core
	// Проверяем, нужно ли логировать в файл.
	if config.LogToFile {
		// Открываем файл для логирования.
		file, err := os.OpenFile(config.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		// Настраиваем ядро логгера для логирования в файл и стандартный вывод.
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file)),
			zap.NewAtomicLevelAt(level),
		)
	} else {
		// Настраиваем ядро логгера для логирования только в стандартный вывод.
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(level),
		)
	}

	// Создаем новый экземпляр логгера с добавлением информации о вызывающем.
	l := zap.New(core, zap.AddCaller())
	// Заменяем глобальный логгер на новый экземпляр.
	zap.ReplaceGlobals(l)
	// Возвращаем экземпляр логгера.
	return l.Sugar(), nil
}
