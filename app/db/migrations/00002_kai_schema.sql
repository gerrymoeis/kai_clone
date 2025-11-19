-- +goose Up

-- KAI clone core schema
CREATE TABLE IF NOT EXISTS stations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100),
    lat DECIMAL(9,6),
    lon DECIMAL(9,6),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_stations_code ON stations(code);

CREATE TABLE IF NOT EXISTS trains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    operator VARCHAR(100),
    class_support VARCHAR(50)
);
CREATE INDEX IF NOT EXISTS idx_trains_code ON trains(code);

CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_code VARCHAR(40) UNIQUE NOT NULL,
    origin_station_id UUID NOT NULL REFERENCES stations(id) ON DELETE RESTRICT,
    dest_station_id UUID NOT NULL REFERENCES stations(id) ON DELETE RESTRICT,
    distance_km DECIMAL(7,2)
);
CREATE INDEX IF NOT EXISTS idx_routes_origin ON routes(origin_station_id);
CREATE INDEX IF NOT EXISTS idx_routes_dest ON routes(dest_station_id);

CREATE TABLE IF NOT EXISTS service_calendars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    weekday_flags INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    exceptions JSONB
);
CREATE INDEX IF NOT EXISTS idx_svc_cal_train ON service_calendars(train_id);

CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    service_date DATE NOT NULL,
    depart_time TIME NOT NULL,
    arrive_time TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    base_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_trips_route ON trips(route_id);
CREATE INDEX IF NOT EXISTS idx_trips_train ON trips(train_id);
CREATE INDEX IF NOT EXISTS idx_trips_date ON trips(service_date);
CREATE INDEX IF NOT EXISTS idx_trips_route_date ON trips(route_id, service_date);

CREATE TABLE IF NOT EXISTS coaches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    coach_no INT NOT NULL,
    class VARCHAR(20) NOT NULL,
    layout_code VARCHAR(20) NOT NULL,
    rows INT NOT NULL,
    cols INT NOT NULL,
    UNIQUE (trip_id, coach_no)
);
CREATE INDEX IF NOT EXISTS idx_coaches_trip ON coaches(trip_id);

CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    coach_no INT NOT NULL,
    seat_no VARCHAR(16) NOT NULL,
    class VARCHAR(20) NOT NULL,
    is_accessible BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (trip_id, coach_no, seat_no)
);
CREATE INDEX IF NOT EXISTS idx_seats_trip ON seats(trip_id);
CREATE INDEX IF NOT EXISTS idx_seats_trip_coach ON seats(trip_id, coach_no);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(24) UNIQUE NOT NULL,
    user_ref VARCHAR(200) NOT NULL,
    trip_id UUID NOT NULL REFERENCES trips(id) ON DELETE RESTRICT,
    total_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_bookings_trip ON bookings(trip_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_ref ON bookings(user_ref);

CREATE TABLE IF NOT EXISTS booking_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id UUID NOT NULL REFERENCES seats(id) ON DELETE RESTRICT,
    price DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'held',
    UNIQUE (booking_id, seat_id)
);
CREATE INDEX IF NOT EXISTS idx_booking_items_booking ON booking_items(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_items_seat ON booking_items(seat_id);

-- +goose Down
DROP TABLE IF EXISTS booking_items;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS coaches;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS service_calendars;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS trains;
DROP TABLE IF EXISTS stations;
