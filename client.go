package grpclog

import (
	"context"

	"github.com/mdigger/grpclog/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Экспортируемые синонимы данных.
type (
	Item  = api.LogItem       // запись лога
	Level = api.LogItem_Level // уровень лога
)

const (
	// информация
	INFO Level = api.LogItem_LEVEL_INFO
	// предупреждение
	WARNING Level = api.LogItem_LEVEL_WARNING
	// ошибка
	ERROR Level = api.LogItem_LEVEL_ERROR
	// фатальная ошибка
	FATAL Level = api.LogItem_LEVEL_FATAL
	// отладочный уровень
	DEBUG Level = api.LogItem_LEVEL_DEBUG
	// уровень трассировки
	TRACE Level = api.LogItem_LEVEL_TRACE
)

// Receiver возвращает функцию для чтения логов с удалённого сервиса по gRPC.
// Каждый вызов этой функции возвращает либо запись лога, либо ошибку.
// В случае планового закрытия канала получения лога сервером возвращается ошибка io.EOF.
func Receiver(ctx context.Context, cc grpc.ClientConnInterface) (func() (*Item, error), error) {
	client := api.NewLogServiceClient(cc)               // инициализируем клиент
	stream, err := client.Logs(ctx, new(emptypb.Empty)) // инициализируем получение логов
	if err != nil {
		return nil, err
	}

	return stream.Recv, nil // возвращаем функцию для чтения логов
}
