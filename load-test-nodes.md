# Load test notes

## Scenario
Target endpoint: `GET /rooms/{roomId}/slots/list?date=YYYY-MM-DD`

Reason: this is the hottest endpoint in the task statement.

## Tool
- k6
- script: `deploy/k6/slots.js`

## How to run
```bash
k6 run -e BASE_URL=http://localhost:8080 \
       -e ROOM_ID=<room-id> \
       -e TOKEN=<jwt> \
       deploy/k6/slots.js
```

- Execution: local
- Script: `deploy/k6/slots.js`
- VUs: 20
- Active duration: 30s
- Max duration (with graceful stop): 60s

## Thresholds

- `http_req_duration p(95) < 200ms`: passed (`p(95)=6.66ms`)
- `http_req_failed rate < 0.001`: passed (`rate=0.00%`)

## Total Results

- checks total: `220587`
- checks succeeded: `220587 (100.00%)`
- checks failed: `0 (0.00%)`
- status is 200: passed

## HTTP Metrics

- `http_reqs`: `220587` (`7350.394251/s`)
- `http_req_failed`: `0.00%` (`0 out of 220587`)
- `http_req_duration`:
  - avg: `2.69ms`
  - min: `216us`
  - med: `2.04ms`
  - max: `143.26ms`
  - p(90): `5.15ms`
  - p(95): `6.66ms`

## Execution Metrics

- `iteration_duration`:
  - avg: `2.71ms`
  - min: `237.16us`
  - med: `2.06ms`
  - max: `147.4ms`
  - p(90): `5.18ms`
  - p(95): `6.69ms`
- `iterations`: `220587` (`7350.394251/s`)
- `vus`: `20` (min=20, max=20)
- `vus_max`: `20` (min=20, max=20)

## Network

- data received: `27 MB` (`889 kB/s`)
- data sent: `83 MB` (`2.8 MB/s`)

## Conclusion

Нагрузочный тест ручки `/rooms/{roomId}/slots/list` уверенно прошёл заданные пороги. Полученные результаты соответствуют целевым SLI по времени ответа и по доле ошибок.
