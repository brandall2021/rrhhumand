package logger

import "go.uber.org/zap"

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}
