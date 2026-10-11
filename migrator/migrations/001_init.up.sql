CREATE TABLE cars (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    num VARCHAR(20) NOT NULL,
    color VARCHAR(30) NOT NULL,
    mark VARCHAR(50) NOT NULL,
    is_foreign INT NOT NULL,
    CONSTRAINT chk_cars_is_foreign CHECK (is_foreign IN (0, 1)),
    CONSTRAINT uq_cars_num UNIQUE (num)
);

CREATE TABLE masters (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE services (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cost_our NUMERIC(10, 2) NOT NULL,
    cost_foreign NUMERIC(10, 2) NOT NULL,
    CONSTRAINT chk_services_cost_our CHECK (cost_our >= 0),
    CONSTRAINT chk_services_cost_foreign CHECK (cost_foreign >= 0)
);

CREATE TABLE works (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    date_work DATE NOT NULL,
    master_id INT NOT NULL,
    car_id INT NOT NULL,
    service_id INT NOT NULL,
    CONSTRAINT fk_works_cars FOREIGN KEY (car_id) REFERENCES cars(id) ON DELETE CASCADE,
    CONSTRAINT fk_works_masters FOREIGN KEY (master_id) REFERENCES masters(id) ON DELETE CASCADE,
    CONSTRAINT fk_works_services FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    CONSTRAINT chk_users_role CHECK (role IN ('admin', 'operator'))
);