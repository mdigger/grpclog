package grpclog

import "fmt"

// F используется как синомим для списка именованных полей.
type F = map[string]any

// fields2map конвертирует список последовательностей в именнованные значения.
// Используется для разбора дополнительных полей лога.
func fields2map(fields []any) map[string]string {
	// отдельно обрабатываем пустые списки
	if len(fields) == 0 || fields[0] == nil {
		return nil
	}

	// формируем список именованных полей с их значениями
	attrs := make(map[string]string)
	for i := 0; i < len(fields); i++ {
		// в зависимости от типа информации преобразуем её в имя и значение
		switch v := fields[i].(type) {
		case string:
			attrs[v] = fmt.Sprint(fields[i+1])
		case fmt.Stringer:
			attrs[v.String()] = fmt.Sprint(fields[i+1])
		case error:
			if v != nil {
				attrs["error"] = v.Error()
			}
			continue
		case map[string]string:
			for k, v := range v {
				attrs[k] = v
			}
			continue
		case map[string]any:
			for k, v := range v {
				attrs[k] = fmt.Sprint(v)
			}
			continue
		}
		i++ // увеличиваем счётчик на прочитанное значение
	}

	return attrs
}
