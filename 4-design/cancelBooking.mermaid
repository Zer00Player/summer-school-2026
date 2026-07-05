sequenceDiagram
    actor User as Клиент
    participant App as Приложение (UI)
    participant API as API Бэкенд
    participant DB as База данных

    Note over App: SCR-06: бронь active<br/>start_at в будущем
    User->>App: Тап «Отменить запись»
    App-->>User: BS-003: «Подтверждение отмены»<br/>«Если до старта ≥2ч — место освободится.<br/>Если <2ч — место не освободится.»

    User->>App: Подтверждает отмену (целиком)

    App->>API: POST /bookings/{bookingId}/cancel<br/>Headers: Authorization: Bearer {token}

    Note over API: Шаг 1: Проверка токена
    alt Токен невалиден / истек
        API-->>App: 401 Unauthorized
        App-->>User: Переход на SCR-01 (Вход)
    end

    Note over API,DB: Шаг 2: Проверка существования брони
    API->>DB: SELECT * FROM bookings<br/>WHERE id = {bookingId}

    alt Бронь не найдена
        API-->>App: 404 Not Found
        App-->>User: «Бронь не найдена»
    end

    alt Бронь принадлежит другому клиенту
        API-->>App: 403 Forbidden
        App-->>User: «У вас нет прав на эту бронь»
    end

    Note over API: Шаг 3: Проверка текущего статуса
    alt Бронь уже отменена (status != active)
        API-->>App: 409 Conflict {code: already_cancelled}
        App-->>User: «Бронь уже отменена»
        App->>App: Актуализировать статус
    end

    Note over API: Шаг 4: Проверка времени старта
    API->>DB: SELECT start_at FROM slots<br/>WHERE id = {slot_id}

    alt Слот уже стартовал (start_at < NOW())
        API-->>App: 422 Unprocessable {code: slot_started}
        App-->>User: «Отмена недоступна: тренировка уже началась»
    end

    Note over API: Шаг 5: Расчет типа отмены<br/>(источник истины — start_at в UTC)
    API->>API: hours_until_start = (start_at - NOW()) / 3600

    alt hours_until_start >= 2
        Note over API: РАННЯЯ ОТМЕНА (≥2ч)
        API->>DB: BEGIN TRANSACTION
        
        API->>DB: UPDATE bookings<br/>SET status = 'cancelled', cancelled_at = NOW()<br/>WHERE id = {bookingId}
        
        API->>DB: UPDATE slots<br/>SET free_seats = free_seats + {seats_count},<br/>    free_rental_boards = free_rental_boards + {rental_count}<br/>WHERE id = {slot_id}
        
        API->>DB: COMMIT
        
        API-->>App: 200 OK<br/>Body: {id, status: 'cancelled', cancelled_at}
        
        App-->>User: «Бронь отменена. Место освобождено»<br/>+ снек-сообщение
    else hours_until_start < 2
        Note over API: ПОЗДНЯЯ ОТМЕНА (<2ч)
        Note over API: Места НЕ возвращаются в слот
        API->>DB: BEGIN TRANSACTION
        
        API->>DB: UPDATE bookings<br/>SET status = 'late_cancel', cancelled_at = NOW()<br/>WHERE id = {bookingId}
        
        API->>DB: COMMIT
        
        API-->>App: 200 OK<br/>Body: {id, status: 'late_cancel', cancelled_at}
        
        App-->>User: «Поздняя отмена: место не освобождено.<br/>Штраф не взимается.»
    end

    App->>App: Обновить SCR-06