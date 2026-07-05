sequenceDiagram
    actor User as Клиент
    participant App as Приложение (UI)
    participant API as API Бэкенд
    participant DB as База данных

    Note over App: SCR-04: выбраны<br/>seats_count=2, rental_count=1<br/>цена показана из slot.price_total

    User->>App: Тап «Записаться»
    App->>App: Генерирует Idempotency-Key (UUID)
    
    App->>API: POST /bookings<br/>Headers: Authorization: Bearer {token}<br/>          Idempotency-Key: {uuid}<br/>Body: {slot_id, seats_count, rental_count}
    
    Note over API: Шаг 1: Проверка токена
    alt Токен невалиден / истек
        API-->>App: 401 Unauthorized
        App-->>User: Переход на SCR-01 (Вход)
    end

    Note over API: Шаг 2: Валидация запроса
    alt Невалидные данные
        API-->>App: 400 Bad Request / 422 Unprocessable
        App-->>User: Подсказка по полям
    end

    Note over API,DB: Шаг 3: Атомарная проверка<br/>(блокировка строки слота)
    API->>DB: SELECT * FROM slots WHERE id = {slot_id} FOR UPDATE
    
    alt Слот не найден
        API-->>App: 404 Not Found
        App-->>User: «Слот не найден»
    end

    alt Слот отменен (status = cancelled)
        API-->>App: 410 Gone {code: slot_cancelled}
        App-->>User: «Прогулка отменена, запись недоступна»
    end

    alt Слот уже стартовал (start_at < NOW())
        API-->>App: 422 Unprocessable {code: slot_started}
        App-->>User: «Запись на прошедшую тренировку невозможна»
    end

    Note over API: Шаг 4: Проверка свободных мест
    alt Недостаточно мест (free_seats < seats_count)
        API-->>App: 409 Conflict {code: slot_full, available_seats}
        App-->>User: «Свободно только X мест из Y»
    end

    alt Недостаточно прокатных досок (free_rental_boards < rental_count)
        API-->>App: 409 Conflict {code: rental_unavailable, available_rental_boards}
        App-->>User: «Доступно только X досок из Y»
    end

    Note over API: Шаг 5: Проверка Idempotency-Key (защита от дублей)
    alt Бронь с таким ключом уже существует
        API-->>App: 409 Conflict {code: double_booking}
        App-->>User: «Вы уже записаны на эту тренировку»
    end

    Note over API,DB: Шаг 6: Создание брони (транзакция)
    API->>DB: BEGIN TRANSACTION
    
    API->>DB: INSERT INTO bookings<br/>(id, slot_id, client_id, seats_count,<br/> rental_count, status, price_total, created_at)<br/>VALUES (..., 'active', calculated_price, NOW())
    
    API->>DB: UPDATE slots<br/>SET free_seats = free_seats - {seats_count},<br/>    free_rental_boards = free_rental_boards - {rental_count}<br/>WHERE id = {slot_id}
    
    API->>DB: INSERT INTO idempotency_keys<br/>(key, booking_id, created_at)<br/>VALUES ({idempotency_key}, {booking_id}, NOW())
    
    API->>DB: COMMIT

    Note over API: Шаг 7: Возврат успешного ответа
    API-->>App: 201 Created<br/>Body: {<br/>  id,<br/>  status: 'active',<br/>  price_total,<br/>  created_at,<br/>  slot: { ... }<br/>}
    
    App-->>User: BS-002: «Запись оформлена!»<br/>+ сводка брони
    App->>App: Переход на SCR-06 (Детали брони)
    
    Note over App: Дополнительно: запрос на разрешение push-уведомлений (при первой записи)