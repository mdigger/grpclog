package grpclog

import (
	"sync"
	"sync/atomic"

	"github.com/mdigger/grpclog/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LogService реализует сервис для публикации логов.
type LogService struct {
	counter                           uint64   // счётчик записей лога
	channels                          sync.Map // спиосок каналов для публикации логов подписчикам
	api.UnimplementedLogServiceServer          // поддержка интерфейсов сервиса логов
}

// синоним для определения канала публикации логов
type logChannel = chan *api.LogItem

// Logs поддерживает отправку потока логов подписчику.
func (s *LogService) Logs(_ *emptypb.Empty, stream api.LogService_LogsServer) error {
	ch := make(logChannel, 100) // на случай задержек с отправкой делаем его кешируемым
	defer close(ch)

	s.channels.Store(ch, struct{}{})
	defer s.channels.Delete(ch)

	for logItem := range ch {
		if err := stream.Send(logItem); err != nil {
			return err
		}
	}

	return nil
}

// send публикует новую запись лога и отправляет её всем текущим подписчикам сервиса.
func (s *LogService) send(item *api.LogItem) {
	item.Id = atomic.AddUint64(&s.counter, 1) // увеличиваем счётчик записей логов

	// отправляем всем получателям
	s.channels.Range(func(key, _ any) bool {
		key.(logChannel) <- item
		return true
	})
}

// Send формирует запись лога и отправляет всем подписчикам.
//
// FIXME: я не стал особо заморачиваться по поводу разбора различных вариатов перечисления полей, поэтому сделал
// как удобнее мне самому для использования!
//
// Поля с ошибками преобразуются в текст с именем error. Именованные списки (map) автоматически преобразуются
// к списку полей с такими же именами. Если имя поля не может быть преобразовано к строке, то оно игнорируется
// вместе с его значением. Последнее поле, не имеющего в пару значения, тоже игнорируется.
func (s *LogService) Send(level api.LogItem_Level, message string, fields ...any) {
	s.send(&api.LogItem{
		Message: message,
		Level:   api.LogItem_Level(level),
		Fields:  fields2map(fields),
	})
}

// Info отправляет информационное сообщение лога.
func (s *LogService) Info(message string, fields ...any) {
	s.Send(api.LogItem_LEVEL_INFO, message, fields...)
}

// Debug отправляет отладочное сообщение лога.
func (s *LogService) Debug(message string, fields ...any) {
	s.Send(api.LogItem_LEVEL_DEBUG, message, fields...)
}

// Error отправляет сообщение лога об ошибке, если ошибка задана.
// В противном случае отправляется информационное сообщение.
func (s *LogService) Error(err error, message string, fields ...any) {
	level := api.LogItem_LEVEL_INFO
	if err != nil {
		level = api.LogItem_LEVEL_ERROR
		fields = append(fields, err)
	}

	s.Send(level, message, fields...)
}

// Warn отправляет предупреждающее сообщение лога.
func (s *LogService) Warn(message string, fields ...any) {
	s.Send(api.LogItem_LEVEL_WARNING, message, fields...)
}
