package berkeley

import (
	"github.com/sirupsen/logrus"
	"reflect"
)

// LoggerManager es responsable de gestionar los loggers para las clases que lo invocan.
type LoggerManager struct{}

// GetLogger obtiene el logger correspondiente a la clase que invoca este método.
func (LoggerManager) GetLogger(clazz interface{}) *logrus.Logger {
	// Se usa el nombre de la clase (tipo) como parte del nombre del logger.
	logger := logrus.New()

	// Se asigna el nombre del tipo (nombre de la clase)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	className := reflect.TypeOf(clazz).Name()
	logger.Infof("Logger for %s initialized.", className)

	return logger
}
