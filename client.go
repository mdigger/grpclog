package grpclog

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mdigger/grpclog/api/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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

var (
	levelShortnames = map[Level]string{
		TRACE:   "TRC",
		DEBUG:   "DBG",
		INFO:    "INF",
		WARNING: "WRN",
		ERROR:   "ERR",
		FATAL:   "FTL",
	}
	newLineReplacer = strings.NewReplacer("\n", "\\n", "\r", "\\r")
)

// Print выводит в форматированном виде запись лога в os.Stdout.
func Print(logmsg *Item) {
	fmt.Printf("%v %4d %3s %s",
		time.Now().Local().Format("15:04:05"), // локальное время
		logmsg.Id,                             // счётчик
		levelShortnames[logmsg.Level],         // уровень сообщения
		logmsg.Message)                        // текст сообщения

	// дополнительные поля
	keys := lo.Keys(logmsg.Fields) // получаем список ключей
	sort.Strings(keys)             // сортируем его

	// перебираем ключи дополнительного поля
	for _, k := range keys {
		// название ключа
		if strings.ContainsAny(k, " \t\r\n\"") {
			k = strconv.Quote(newLineReplacer.Replace(k))
		}
		fmt.Printf(" %s=", k)

		// значение ключа
		v := logmsg.Fields[k]
		if strings.ContainsAny(v, " \t\r\n\"") {
			v = strconv.Quote(newLineReplacer.Replace(v))
		}
		fmt.Print(v)
	}

	fmt.Println() // перевод на новую строку в конце
}
