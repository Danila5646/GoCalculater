# Калькулятор выражений

## Описание
Этот проект предоставляет API для расчёта математических выражений. Можно отправить выражение в виде строки, и сервер вернёт результат вычисления. Поддерживаются четыре основные операции: сложение, вычитание, умножение, деление, а также скобки для управления приоритетом операций.

## Использование

### Пример использования с помощью `curl`

#### Успешный запрос
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Ответ:
```json
{
  "result": "6.000000"
}
```

#### Ошибка 422 (некорректное выражение)
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+"
}'
```
Ответ:
```text
HTTP/1.1 422 Unprocessable Entity
Некорректное выражение!
```

#### Ошибка 500 (деление на ноль)
```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "10/0"
}'
```
Ответ:
```text
HTTP/1.1 500 Internal Server Error
Деление на ноль!
```

## Инструкция для запуска

1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/Danila5646/GoCalculater
   cd Calculater_project
   ```
2. Запустите сервер:
   ```bash
   go run ./cmd/calculator_project/...
   ```
3. Сервер будет доступен по адресу `http://localhost:8080`.

## Описание API
### POST /api/v1/calculate

**Описание:** Рассчитывает математическое выражение.

**Тело запроса:**
```json
{
  "expression": "Cтрока с математическим выражением"
}
```

**Коды ответа:**
- `200 OK`: Запрос выполнен успешно.
- `422 Unprocessable Entity`: Некорректное математическое выражение.
- `500 Internal Server Error`: Ошибка сервера.
