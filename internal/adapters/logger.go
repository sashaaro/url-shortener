package adapters

import "go.uber.org/zap"

var Logger zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer l.Sync()

	Logger = *l.Sugar()
}
