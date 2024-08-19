package adapters

import "go.uber.org/zap"

// создание логина
func CreateLogger() zap.SugaredLogger {
	l, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	//nolint:errcheck
	defer l.Sync()

	return *l.Sugar()
}
