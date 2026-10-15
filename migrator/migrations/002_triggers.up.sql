CREATE OR REPLACE FUNCTION check_unique_car_num() 
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM cars WHERE cars.num = NEW.num) THEN
        RAISE EXCEPTION 'Автомобиль с номером % уже существует!', NEW.num;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_unique_car_num
BEFORE INSERT ON cars
FOR EACH ROW
EXECUTE FUNCTION check_unique_car_num();



CREATE OR REPLACE FUNCTION check_masters_count()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF (SELECT COUNT(*) FROM masters) >= 10 THEN
        RAISE EXCEPTION 'Нельзя добавить мастера: их уже больше 10!';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_masters_count
BEFORE INSERT ON masters
FOR EACH ROW
EXECUTE FUNCTION check_masters_count();



CREATE OR REPLACE FUNCTION check_master_daily_works()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM works 
        WHERE works.master_id = NEW.master_id 
          AND works.date_work = NEW.date_work
          AND works.id != NEW.id
    ) THEN
        RAISE EXCEPTION 'Мастер уже выполнил работу в этот день!';
    END IF;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trigger_master_daily_works ON works;
CREATE TRIGGER trigger_master_daily_works
BEFORE INSERT OR UPDATE ON works
FOR EACH ROW
EXECUTE FUNCTION check_master_daily_works();



CREATE OR REPLACE FUNCTION check_date_change() 
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    IF ABS(NEW.date_work - OLD.date_work) > 1 THEN
        RAISE EXCEPTION 'Нельзя изменить дату работы более чем на 1 день! Было: %, стало: %', 
            OLD.date_work, NEW.date_work;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_date_change
BEFORE UPDATE ON works
FOR EACH ROW
WHEN (NEW.date_work IS DISTINCT FROM OLD.date_work)
EXECUTE FUNCTION check_date_change();



CREATE OR REPLACE FUNCTION prevent_car_delete() 
RETURNS TRIGGER 
LANGUAGE plpgsql 
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM works WHERE works.car_id = OLD.id) THEN
        RAISE EXCEPTION 'Нельзя удалить автомобиль %, так как по нему есть работы!', OLD.num;
    END IF;
    RETURN OLD;
END;
$$;

CREATE TRIGGER trigger_prevent_car_delete
BEFORE DELETE ON cars
FOR EACH ROW
EXECUTE FUNCTION prevent_car_delete();