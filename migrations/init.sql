-- Создаём таблицу подписок
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,                                 -- Уникальный идентификатор записи
    service_name VARCHAR(255) NOT NULL,                 -- Название сервиса
    price INTEGER NOT NULL,                             -- Стоимость в рублях (целое число)
    user_id UUID NOT NULL,                              -- ID пользователя (UUID)
    start_date DATE NOT NULL,                           -- Дата начала (первое число месяца)
    end_date DATE                                       -- Дата окончания (опционально)
);