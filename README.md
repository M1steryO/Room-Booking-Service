# Room Booking Service

[![CI](https://github.com/avito-internships/test-backend-1-M1steryO/actions/workflows/ci.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-M1steryO/actions/workflows/ci.yml)
[![Classroom](https://github.com/avito-internships/test-backend-1-M1steryO/actions/workflows/classroom.yml/badge.svg)](https://github.com/avito-internships/test-backend-1-M1steryO/actions/workflows/classroom.yml)


## Запуск 

```bash
go mod tidy
make up
```

Проверка работоспособности:

```bash
curl http://localhost:8080/_info
```

Получить тестовый токен:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"admin"}'
```

## Пояснения

### 1 Как генерируются слоты
При создании расписания сервис предсоздаёт слоты на определенное число дней вперёд (задается переменной `SLOT_HORIZON_DAYS`).  
При запросе `/rooms/{roomId}/slots/list?date=...` сервис сначала проверяет, есть ли слоты на эту дату, и, если их ещё нет, лениво генерирует их именно на нужный день. Такой подход даёт:
- быстрый hot path для ближайших дат;
- отсутствие "вечной" материализации всех будущих слотов;
- простую модель данных без фоновых джобов.

### 2. Поведение при сбоях Conference Service

Если при создании брони с флагом createConferenceLink внешний сервис не отвечает, падает по таймауту или отменяется контекст, бронь не создаётся. Клиент получает ошибку, в базу бронь не попадает.

Если ссылку на конференцию уже получили, а сохранить бронь в базе не удалось (например, слот уже занят другим запросом), клиент тоже получает ошибку. В моей реализации отката созданной ссылки во внешнем сервисе понятное дело нет

Если createConferenceLink не передают или false, внешний сервис не вызывается, бронь создаётся без ссылки.

Повторный запрос на создание брони после ошибки — это новый запрос. Две активные брони на один слот база не даст за счёт уникального индекса по slot_id для активных броней.

В продакшене логичнее вынести создание конференции в фон (использовать outbox например), добавить ретраи и идемпотентность к внешнему API и отдельно мониторить лишние ссылки, если запись брони не прошла после успешного ответа конференции.

### 3. Индексы
Индекс `bookings_active_slot_uniq` не даёт создать две активные брони на один и тот же слот даже при одновременных запросах, индекс `slots_room_start_idx` ускоряет самый частый сценарий — получение слотов переговорки на дату по room_id и времени начала, а индекс `bookings_user_status_idx` ускоряет выборки пользовательских броней, в связи с чем база тратит меньше времени на сканирование таблиц и лучше держит нагрузку.


## Команды

```bash
make up         # поднять postgres + app
make down       # остановить окружение
make seed       # наполнить БД тестовыми данными
make test       # тесты
make test-e2e   # e2e тест
make lint       # линтер
make swagger    # генерация swagger по аннотациям (требует swag)
```

## Swagger

В репозитории уже лежит `api/api.yaml` как источник спецификации.  
Дополнительно предусмотрен `make swagger` для генерации документации из аннотаций в коде через `swag`.

Установка:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger
```

## Нагрузочное тестирование

Готовый сценарий лежит в `deploy/k6/slots.js`.

Установка:

[Ссылка](https://grafana.com/docs/k6/latest/set-up/install-k6/)


Пример запуска:
```bash
k6 run -e BASE_URL=http://localhost:8080 \
       -e ROOM_ID=<room-id> \
       -e TOKEN=<jwt> \
       deploy/k6/slots.js
```

Отчет лежит в `load-test-nodes.md`
