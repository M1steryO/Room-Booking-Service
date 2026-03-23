# Room Booking Service

[![CI](https://github.com/yourusername/room-booking-service/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/room-booking-service/actions/workflows/ci.yml)

Продакшн-ориентированное решение тестового задания на Go + PostgreSQL.

## Что сделано

- clean architecture: `domain -> usecase -> repository/postgres -> transport/http`
- PostgreSQL с отдельными миграциями
- JWT авторизация через `Authorization: Bearer <token>`
- обязательный `/dummyLogin`
- дополнительные `/register` и `/login`
- mock `Conference Service`
- `/_info`
- unit tests + e2e tests
- `docker-compose.yml` и `docker-compose.e2e.yml`
- `Makefile`
- CI в GitHub Actions
- `.golangci.yaml`
- k6 сценарий для нагрузочного тестирования

## Быстрый старт

```bash
cp local.env .env
go mod tidy
make up
```

Проверка health:

```bash
curl http://localhost:8080/_info
```

Получить тестовый токен:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"admin"}'
```

## Основные архитектурные решения

### 1. Стабильные `slotId`
Требование задания: слоты должны иметь стабильные UUID и храниться в БД, иначе бронирование по `slotId` невозможно. Поэтому `slotId` вычисляется детерминированно из `(room_id, start_at, end_at)` и сохраняется в таблице `slots`. Это даёт:
- одинаковый UUID при повторной генерации;
- возможность безопасного `upsert`;
- устойчивость к повторным запросам списка слотов.

### 2. Как генерируются слоты
При создании расписания сервис предсоздаёт слоты на горизонт `SLOT_HORIZON_DAYS` вперёд.  
При запросе `/rooms/{roomId}/slots/list?date=...` сервис сначала проверяет, есть ли слоты на эту дату, и, если их ещё нет, лениво генерирует их именно на нужный день. Такой подход даёт:
- быстрый hot path для ближайших дат;
- отсутствие "вечной" материализации всех будущих слотов;
- простую модель данных без фоновых джобов.

### 3. Защита от двойного бронирования
В таблице `bookings` используется partial unique index:
```sql
CREATE UNIQUE INDEX bookings_active_slot_uniq
ON bookings (slot_id)
WHERE status = 'active';
```
Это гарантирует, что у одного слота не появятся две активные брони даже при гонках.

### 4. Идемпотентная отмена
`POST /bookings/{bookingId}/cancel` всегда возвращает текущее состояние брони. Если бронь уже была отменена, это не ошибка.

### 5. Конференц-ссылка
Если передан `createConferenceLink=true`, сервис синхронно вызывает mock `Conference Service` и сохраняет ссылку вместе с бронью.

Почему так:
- для тестового задания это проще и прозрачнее;
- поведение детерминировано и легко тестируется.

Что делал бы в проде:
- вынес бы создание конференции в асинхронный процесс через outbox/saga;
- добавил бы retry policy и идемпотентный внешний ключ запроса;
- отдельно отслеживал бы orphaned links при сбоях после внешнего ответа.

## Структура

```text
cmd/api
internal/
  app
  conference
  config
  domain
  platform/
    clock
    identity
    migrations
    security
  repository/postgres
  transport/httpx
  usecase
migrations/
tests/
  unit
  e2e
```

## Команды

```bash
make up         # поднять postgres + app
make down       # остановить окружение
make seed       # наполнить БД тестовыми данными
make test       # все тесты
make test-e2e   # отдельное e2e окружение
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

Пример запуска:
```bash
k6 run -e BASE_URL=http://localhost:8080 \
       -e ROOM_ID=<room-id> \
       -e TOKEN=<jwt> \
       deploy/k6/slots.js
```

В этой среде я не запускал реальный k6 against live service, поэтому честных численных результатов в репозиторий не добавлял. Скрипт и методика запуска приложены.

## Проверка покрытия

Требование задания — покрытие выше 40%.  
В репозитории есть unit и e2e тесты; ожидается, что после `go test ./... -coverprofile=coverage.out` покрытие будет выше порога после локальной сборки с зависимостями.

## Соответствие советам по оформлению

Решение специально построено так, чтобы не повторять типичные ошибки из файла `Solution and advice.md`:
- обработчики разложены по отдельному транспортному слою;
- бизнес-логика вынесена в `usecase`;
- зависимости задаются через интерфейсы на стороне потребителя;
- миграции разнесены по отдельным файлам;
- есть отдельное e2e окружение;
- секреты не коммитятся, используется `.env.example`.
