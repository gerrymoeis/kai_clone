-- +goose Up

-- Seed minimal KAI-like data for development

-- Stations
INSERT INTO stations (id, code, name, city, lat, lon) VALUES
    (gen_random_uuid(),'GMR','Gambir','Jakarta',-6.176655,106.830583),
    (gen_random_uuid(),'BD','Bandung','Bandung',-6.914744,107.609810),
    (gen_random_uuid(),'YK','Yogyakarta','Yogyakarta',-7.789133,110.363022),
    (gen_random_uuid(),'SMT','Semarang Tawang','Semarang',-6.967222,110.427222),
    (gen_random_uuid(),'SGU','Surabaya Gubeng','Surabaya',-7.265357,112.750000)
ON CONFLICT (code) DO NOTHING;

-- Trains
INSERT INTO trains (id, code, name, operator, class_support) VALUES
    (gen_random_uuid(),'AP','Argo Parahyangan','KAI','economy,business,executive'),
    (gen_random_uuid(),'TAK','Taksaka','KAI','business,executive'),
    (gen_random_uuid(),'ARW','Argo Wilis','KAI','business,executive')
ON CONFLICT (code) DO NOTHING;

-- Routes
INSERT INTO routes (id, route_code, origin_station_id, dest_station_id, distance_km)
SELECT gen_random_uuid(), 'GMR-BD', s1.id, s2.id, 150.0
FROM stations s1, stations s2
WHERE s1.code='GMR' AND s2.code='BD'
ON CONFLICT (route_code) DO NOTHING;

INSERT INTO routes (id, route_code, origin_station_id, dest_station_id, distance_km)
SELECT gen_random_uuid(), 'BD-GMR', s1.id, s2.id, 150.0
FROM stations s1, stations s2
WHERE s1.code='BD' AND s2.code='GMR'
ON CONFLICT (route_code) DO NOTHING;

-- Trips (two sample days, two directions)
INSERT INTO trips (id, route_id, train_id, service_date, depart_time, arrive_time, status, base_price)
SELECT gen_random_uuid(), r.id, t.id, (current_date + interval '1 day'), time '07:00', time '10:00', 'scheduled', 150000
FROM routes r JOIN trains t ON t.code='AP' WHERE r.route_code='GMR-BD'
ON CONFLICT DO NOTHING;

INSERT INTO trips (id, route_id, train_id, service_date, depart_time, arrive_time, status, base_price)
SELECT gen_random_uuid(), r.id, t.id, (current_date + interval '1 day'), time '17:00', time '20:00', 'scheduled', 150000
FROM routes r JOIN trains t ON t.code='AP' WHERE r.route_code='BD-GMR'
ON CONFLICT DO NOTHING;

INSERT INTO trips (id, route_id, train_id, service_date, depart_time, arrive_time, status, base_price)
SELECT gen_random_uuid(), r.id, t.id, (current_date + interval '2 day'), time '07:00', time '10:00', 'scheduled', 150000
FROM routes r JOIN trains t ON t.code='AP' WHERE r.route_code='GMR-BD'
ON CONFLICT DO NOTHING;

INSERT INTO trips (id, route_id, train_id, service_date, depart_time, arrive_time, status, base_price)
SELECT gen_random_uuid(), r.id, t.id, (current_date + interval '2 day'), time '17:00', time '20:00', 'scheduled', 150000
FROM routes r JOIN trains t ON t.code='AP' WHERE r.route_code='BD-GMR'
ON CONFLICT DO NOTHING;

-- Coaches (6 per trip, 2-2 layout, 15x4)
INSERT INTO coaches (id, trip_id, coach_no, class, layout_code, rows, cols)
SELECT gen_random_uuid(), tr.id, gs, 'economy', '2-2', 15, 4
FROM trips tr
CROSS JOIN generate_series(1,6) AS gs
WHERE tr.service_date IN (current_date + interval '1 day', current_date + interval '2 day');

-- Seats (15x4 per coach)
INSERT INTO seats (id, trip_id, coach_no, seat_no, class, is_accessible)
SELECT gen_random_uuid(), c.trip_id, c.coach_no, concat(lpad(row::string,2,'0'),'-',col::string), c.class, FALSE
FROM coaches c
CROSS JOIN generate_series(1,15) AS row
CROSS JOIN generate_series(1,4) AS col
WHERE c.trip_id IN (
    SELECT id FROM trips WHERE service_date IN (current_date + interval '1 day', current_date + interval '2 day')
);

-- +goose Down
-- Remove seeded trips/coaches/seats/routes/trains/stations inserted above
DELETE FROM seats WHERE trip_id IN (
    SELECT id FROM trips WHERE service_date IN (current_date + interval '1 day', current_date + interval '2 day')
);
DELETE FROM coaches WHERE trip_id IN (
    SELECT id FROM trips WHERE service_date IN (current_date + interval '1 day', current_date + interval '2 day')
);
DELETE FROM trips WHERE service_date IN (current_date + interval '1 day', current_date + interval '2 day');
DELETE FROM routes WHERE route_code IN ('GMR-BD','BD-GMR');
DELETE FROM trains WHERE code IN ('AP','TAK','ARW');
DELETE FROM stations WHERE code IN ('GMR','BD','YK','SMT','SGU');
