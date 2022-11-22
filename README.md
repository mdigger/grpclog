# gRPC сервис для получения логов

Вспомогательный сервис для получения логов от gRPC сервиса в реальном времени.

```golang
conn, err := grpc.Dial("localhost:3000",
    grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    panic(err)
}
defer conn.Close()

stream, err := grpclog.Receiver(context.Background(), conn)
if err != nil {
    panic(err)
}

for {
    logMsg, err := stream()
    if err == io.EOF {
        break
    }
    if err != nil {
        panic(err)
    }

    grpclog.Print(logMsg)
}
```
