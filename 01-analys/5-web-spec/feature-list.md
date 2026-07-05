# Фича-лист веб-приложения скалодрома

**Версия:** 1.0.0  
**Дата:** 2026-07-05  
**Платформа:** Web

---

## 1. Назначение

Веб-приложение для самостоятельной записи на тренировки в скалодроме. Заменяет ручную запись через WhatsApp и бумажную тетрадь, устраняя двойные брони и путаницу с местами.

**Скоуп приложения — только роль «Клиент».** Инструктор и Администратор работают через существующую инфраструктуру/админку и в приложение **не входят**.

**Источники:** [Бриф](../0-customer-brief/customer-brief.md) · [Бизнес-требования](../2-requirements/business-requirements.md) · [Функциональные требования](../2-requirements/functional-requirements.md) · [Нефункциональные требования](../2-requirements/non-functional-requirements.md)

---

## 2. Карта навигации

```mermaid
graph TD
    Start([Открытие сайта]) --> Auth{Авторизован?}
    Auth -->|Нет| Landing[SCR-000 Лендинг]
    Auth -->|Да| Schedule[SCR-001 Расписание]
    
    Landing -->|Войти| Login[Модальное окно входа]
    Login --> Schedule
    
    Schedule -->|Клик по карточке| SlotCard[SCR-003 Карточка слота]
    Schedule -->|Клик «Записаться»| Booking[SCR-004 Оформление записи]
    Schedule -->|Клик «Фильтры»| Filters[BS-001 Фильтры]
    
    SlotCard -->|«Записаться»| Booking
    SlotCard -->|Карта| Map[BS-004 Карта маршрута]
    
    Booking -->|Успех| Success[BS-002 Подтверждение записи]
    Booking -->|Ошибка| BookingError[Сообщение об ошибке]
    
    Success -->|«Мои записи»| History[SCR-005 Мои бронирования]
    Success -->|«Готово»| Schedule
    
    Schedule -->|«Мои записи»| History
    History -->|Клик по брони| Details[SCR-006 Детали брони]
    Details -->|«Отменить»| Cancel[BS-003 Подтверждение отмены]
    Cancel -->|Подтвердить| Details
    
    Schedule -->|«Профиль»| Profile[SCR-007 Профиль]
    Profile -->|«Выйти»| Landing