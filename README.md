# techlog
Библиотека парсинга Технологического журнала 1С Предприятие

## Функционал

* `techlog.Read(file)` Чтение файла технологического журнала в массив
* `techlog.ReadAt(file, offset)` Чтение файла технологического журнала в массив с опеределнного места
* `techlog.StreamRead(file, 10, offset)` Чтение файла технологического журнала в поток
* `techlog.StreamReadAt(file, 10, offset)` Чтение файла технологического журнала в поток с опеределнного места

## Примеры

### Чтение файла технологического журнала в массив
```go
package main

import (
	"log"
	"v8platform/techlog"
)

func main() {

	file := "./logs/20100521.log"

	events, err := techlog.Read(file)
    //events, offset, err := techlog.ReadAt(file, 500)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("readed <%d> events", len(events))
}

```

### Чтение файла технологического журнала в поток
```go
package main

import (
	"log"
	"v8platform/techlog"
)

func main() {


	file := "./logs/20100521.log"

	events, err := techlog.StreamRead(file, 10)
	//events, offset, err := techlog.StreamReadAt(file, 10, 500)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	for event := range events {
		//pp.Println(event) // Не ракомендую использовать на больших объемах
		count++
	}

	log.Printf("readed <%d> events", count)
}

```
