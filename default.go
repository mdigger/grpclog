package grpclog

import (
	"github.com/mdigger/grpclog/api/v1"
	"google.golang.org/grpc"
)

// сервис логов по умолчанию
var service = new(LogService)

// Register регистрирует сервис отправки логов по умолчанию в качестве сервиса gRPC.
func Register(srv *grpc.Server) {
	api.RegisterLogServiceServer(srv, service)
}

// Send отвечает за формирование и отправку записи лога.
func Send(level int8, message string, fields ...any) {
	service.Send(api.LogItem_Level(level), message, fields...)
}

// Info отправляет информационное сообщение лога.
func Info(message string, fields ...any) {
	service.Info(message, fields...)
}

// Debug отправляет отладочное сообщение лога.
func Debug(message string, fields ...any) {
	service.Debug(message, fields...)
}

// Error отправляет сообщение лога об ошибке, если ошибка задана.
// В противном случае отправляется информационное сообщение.
func Error(err error, message string, fields ...any) {
	service.Error(err, message, fields...)
}

// Warn отправляет предупреждающее сообщение лога.
func Warn(message string, fields ...any) {
	service.Warn(message, fields...)
}
