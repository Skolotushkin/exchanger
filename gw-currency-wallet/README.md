# Разработка микросервиса для управления кошельком и обмена валют (gw-currency-wallet)

Микросервис поддерживает регистрацию и авторизацию пользователей, пополнение счета, вывод средств,  
получение курсов валют и обмен валют. В качестве базы данных используется **PostgreSQL**.  
Для взаимодействия с внешним сервисом курсов валют используется **gRPC**.

**Стек технологий:**  
_Gin (или другой HTTP-фреймворк), JWT для авторизации, gRPC для получения курсов валют и обмена_

---

## Описание работы

Курс валют осуществляется по данным сервиса `exchange`. Если в течение короткого времени был запрос курса валют (`/api/v1/exchange`) до обмена, то курс берется из кэша.  
Если же запроса курса валют не было или он устарел, выполняется gRPC-вызов к внешнему сервису, который предоставляет актуальные курсы валют.  

При обмене проверяется наличие средств пользователя, выполняется пересчет и обновляется баланс.  

---

## Требования к сервису

1. **Безопасность:** Все запросы защищены JWT-токенами.  
2. **Производительность:** Время отклика сервиса не должно превышать 200 мс для большинства операций.  
3. **Логирование:** Все операции логируются для анализа и отладки.  
4. **Тестирование:** Необходимо покрыть юнит-тестами все основные функции.  
5. **Документация:** API документируется с помощью Swagger или аналогичного инструмента.  
6. **Конфигурация:** Сервис должен читать переменные окружения из `config.env` для локального запуска.  

---

## API

### 1. Регистрация пользователя

**Метод:** POST  
**URL:** `/api/v1/auth/register`  

**Тело запроса:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Ответ:**

Успех: 201 Created
```json
{ 
  "message": "User registered successfully"
}
```

Ошибка: 400 Bad Request
```json
{
  "error": "Invalid request"
}
```

### 2. Авторизация пользователя

**Метод:** POST
URL: /api/v1/auth/login

**Тело запроса:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Ответ:**

Успех: 200 OK
```json
{
  "token": "JWT_TOKEN"
}
```

Ошибка: 401 Unauthorized
```json
{
  "error": "invalid email or password"
}
```

### 3. Получение баланса пользователя

**Метод:** GET
URL: /api/v1/wallet/balance
Заголовки: Authorization: Bearer JWT_TOKEN

**Ответ:**

Успех: 200 OK
```json
{
  "balance": {
    "USD": "float",
    "RUB": "float",
    "EUR": "float"
  }
}
```

### 4. Пополнение счета

**Метод:** POST
URL: /api/v1/wallet/deposit
Заголовки: Authorization: Bearer JWT_TOKEN

**Тело запроса:**
```json
{
  "amount": "100.00",
  "currency": "USD"
}
```

**Ответ:**

Успех: 200 OK
```json
{
  "message": "Deposit successful",
  "new_balance": {
    "USD": "float",
    "RUB": "float",
    "EUR": "float"
  }
}
```

Ошибка: 400 Bad Request
```json
{
  "error": "Invalid amount"
}
```

### 5. Вывод средств

**Метод:** POST
URL: /api/v1/wallet/withdraw
Заголовки: Authorization: Bearer JWT_TOKEN

**Тело запроса:**
```json
{
  "amount": "50.00",
  "currency": "USD"
}
```

Ответ:

Успех: 200 OK
```json
{
  "message": "Withdrawal successful",
  "new_balance": {
    "USD": "float",
    "RUB": "float",
    "EUR": "float"
  }
}
```

Ошибка: 400 Bad Request
```json
{
  "error": "Insufficient funds or invalid amount"
}
```

### 6. Получение курса валют

**Метод:** GET
URL: /api/v1/exchange/rates
Заголовки: Authorization: Bearer JWT_TOKEN

Ответ:

Успех: 200 OK
```json
{
  "rates": {
    "USD": "float",
    "RUB": "float",
    "EUR": "float"
  }
}
```

Ошибка: 500 Internal Server Error
```json
{
  "error": "Failed to retrieve exchange rates"
}
```

### 7. Обмен валют

**Метод:** POST
URL: /api/v1/exchange
Заголовки: Authorization: Bearer JWT_TOKEN

**Тело запроса:**
```json
{
  "from_currency": "USD",
  "to_currency": "EUR",
  "amount": "100.00"
}
```

Ответ:

Успех: 200 OK
```json
{
  "message": "Exchange successful",
  "exchanged_amount": "85.00",
  "new_balance": {
    "USD": "0.00",
    "EUR": "85.00"
  }
}
```

Ошибка: 400 Bad Request
```json
{
  "error": "Insufficient funds or invalid currencies"
}
```
