package actioninfo

import (
	"fmt"
	"log"
)

// DataParser объединяет разбор данных и формирование отчёта об активности.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info выводит отчёты; ошибки отдельных записей не прерывают обработку.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Printf("ошибка разбора данных %q: %v", data, err)
			continue
		}
		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("ошибка формирования отчёта для %q: %v", data, err)
			continue
		}
		fmt.Println(info)
	}
}
