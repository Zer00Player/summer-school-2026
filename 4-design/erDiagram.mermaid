erDiagram
    %% СУЩНОСТИ (СПРАВОЧНИКИ) - ТОЛЬКО ЧТЕНИЕ
    Route {
        uuid id PK
        string name
        string description
        enum type "novice|experienced"
        int capacity_cap "≤8 нович/≤12 опыт"
        int duration_min
        polyline geometry
    }
    
    Instructor {
        uuid id PK
        string name
    }
    
    Slot {
        uuid id PK
        uuid route_id FK
        uuid instructor_id FK
        datetime start_at "UTC"
        int total_seats
        int free_seats
        int free_rental_boards
        money price "за место"
        money rental_price "за доску"
        string meeting_point
        float meeting_point_lat
        float meeting_point_lng
        enum status "scheduled|cancelled"
    }

    %% СУЩНОСТЬ (РАБОЧАЯ) - ЧТЕНИЕ + ЗАПИСЬ
    Client {
        uuid id PK
        string name
        string phone UK
        datetime created_at
    }

    %% СУЩНОСТЬ (РАБОЧАЯ) - СОЗДАЕТСЯ/ОТМЕНЯЕТСЯ
    Booking {
        uuid id PK
        uuid slot_id FK
        uuid client_id FK
        int seats_count "1..3"
        int rental_count "0..seats_count"
        enum status "active|cancelled|late_cancel|club_cancelled"
        money price_total "read-only, расчет сервера"
        datetime created_at
        datetime cancelled_at
    }

    %% СВЯЗИ
    Route ||--o{ Slot : "определяет"
    Instructor ||--o{ Slot : "ведет"
    Client ||--o{ Booking : "создает"
    Slot ||--o{ Booking : "содержит"

    %% ЛЕГЕНДА
    Route ||--o{ Slot : "read-only (справочник)"
    Instructor ||--o{ Slot : "read-only (справочник)"
    Slot ||--o{ Booking : "read-only (проекция)"
    Client ||--o{ Booking : "чтение + запись"
    Booking : "чтение + запись (создание/отмена)"