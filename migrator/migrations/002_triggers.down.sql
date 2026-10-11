DROP TRIGGER IF EXISTS trigger_unique_car_num ON cars;
DROP TRIGGER IF EXISTS trigger_masters_count ON masters;
DROP TRIGGER IF EXISTS trigger_master_daily_works ON works;
DROP TRIGGER IF EXISTS trigger_date_change ON works;
DROP TRIGGER IF EXISTS trigger_prevent_car_delete ON cars;

DROP FUNCTION IF EXISTS check_unique_car_num();
DROP FUNCTION IF EXISTS check_masters_count();
DROP FUNCTION IF EXISTS check_master_daily_works();
DROP FUNCTION IF EXISTS check_date_change();
DROP FUNCTION IF EXISTS prevent_car_delete();