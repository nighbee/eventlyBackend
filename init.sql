-- Инициализация базы данных Evently
-- Этот файл автоматически выполняется при первом запуске PostgreSQL контейнера

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'user'
);

-- Создание таблицы событий
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255) NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    price BIGINT NOT NULL,
    capacity INTEGER NOT NULL
);

-- Создание таблицы бронирований
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL,
    event_id BIGINT NOT NULL,
    seat_number INTEGER NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    paid BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(32) NOT NULL DEFAULT 'reserved'
);

-- Создание индексов для производительности
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_events_starts_at ON events(starts_at);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings(event_id);

-- Уникальный индекс для предотвращения двойного бронирования
CREATE UNIQUE INDEX IF NOT EXISTS uniq_event_seat_active ON bookings(event_id, seat_number, active);

-- Внешние ключи для целостности данных
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_bookings_user_id') THEN
        ALTER TABLE bookings ADD CONSTRAINT fk_bookings_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_bookings_event_id') THEN
        ALTER TABLE bookings ADD CONSTRAINT fk_bookings_event_id FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Создание пользователя-админа (пароль: admin123)
-- В реальном приложении используйте bcrypt для хеширования
INSERT INTO users (name, email, password_hash, role) 
VALUES ('Admin User', 'admin@evently.com', '$2a$10$example_hash_here', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Создание тестового события
INSERT INTO events (title, description, location, starts_at, price, capacity) 
VALUES (
    'Тестовый концерт', 
    'Описание тестового концерта для демонстрации функционала', 
    'Москва, Кремлёвский дворец', 
    '2025-09-01 19:00:00+03', 
    5000, 
    100
) ON CONFLICT DO NOTHING;

-- Создание триггера для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создание триггеров для всех таблиц
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_events_updated_at ON events;
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Логирование успешной инициализации
DO $$
BEGIN
    RAISE NOTICE 'База данных Evently успешно инициализирована!';
    RAISE NOTICE 'Админ: admin@evently.com (роль: admin)';
    RAISE NOTICE 'Тестовое событие создано';
END $$;
