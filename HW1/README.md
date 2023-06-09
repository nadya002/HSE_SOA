# Сериализация и десериализация данных

## Запуск

Проект запускается командой
```bash
docker-compouse up --build
```
Эта команда строит 6 докер контейнеров для каждого протокола, а так же прокси сервер, который слушает на UDP порту 2000.

Обращаться к этому порту можно из терминала с помощью команды 
```bash
nc -u 127.0.0.1 2000
```
Можно отправлять серверу запрос вида get_result {название формата}, где варианты формата могут быть такими:
1) json
2) xml
3) msgpack
4) avro
5) yaml
6) protobuf

Пример взаимодействия:
```bash
$ nc -u 127.0.0.1 2000
get_result avro
avro - 66176 - 173.61µs - 798.267µs
get_result json
json - 78063 - 1.07782ms - 1.859386ms
get_result
Wrong number of params
```

Вывод результатов осуществляется в следующем виде:\
{Формат сериализации} – {Размер сериализованной структуры в байтах} – {Время сериализации} – {Время десериализации}

## Метод решения

Создаются 6 докер контейнеров, которые принимают соединения на своем порту и обрабатывают запросы вида get_result. Для каждого такого запроса сервер 1000 раз сериализует и десериализует структуру данных (из числа, строки, массива, словаря и дробного числа), считает среднее время сериализации и десериализации, а также размер получившийся структуры в байтах, полученный ответ возвращается клиенту.

Прокси сервер слушает порт 2000, обрабатывает запросы вида get_result {название формата}, и перенаправляет запросы к правильному докер контейнеру, а затем возвращает ответ клиенту.

