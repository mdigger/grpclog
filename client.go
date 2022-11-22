package grpclog

import (
	"context"

	"github.com/mdigger/grpclog/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LogItem является синонимом описания записи лога.
type LogItem = api.LogItem

// Receiver возвращает функцию для чтения логов с удалённого сервиса по gRPC.
// Каждый вызов этой функции возвращает либо запись лога, либо ошибку.
// В случае планового закрытия канала получения лога сервером возвращается ошибка io.EOF.
func Receiver(ctx context.Context, cc grpc.ClientConnInterface) (func() (*LogItem, error), error) {
	client := api.NewLogServiceClient(cc)               // инициализируем клиент
	stream, err := client.Logs(ctx, new(emptypb.Empty)) // инициализируем получение логов
	if err != nil {
		return nil, err
	}

	return stream.Recv, nil // возвращаем функцию для чтения логов
}
