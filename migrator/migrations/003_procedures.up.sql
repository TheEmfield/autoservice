CREATE OR REPLACE FUNCTION get_cost_by_car_type(
    p_start_date DATE DEFAULT NULL,
    p_end_date DATE DEFAULT NULL
)
RETURNS TABLE(car_type TEXT, total_cost NUMERIC, works_count BIGINT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        CASE WHEN cars.is_foreign = 1 THEN 'Импортные' ELSE 'Отечественные' END::TEXT,
        COALESCE(SUM(
            CASE 
                WHEN cars.is_foreign = 1 THEN services.cost_foreign
                ELSE services.cost_our
            END
        ), 0),
        COUNT(works.id)
    FROM cars
    JOIN works ON works.car_id = cars.id
    JOIN services ON works.service_id = services.id
    WHERE (p_start_date IS NULL OR works.date_work >= p_start_date)
      AND (p_end_date IS NULL OR works.date_work <= p_end_date)
    GROUP BY cars.is_foreign;
END;
$$;

CREATE OR REPLACE FUNCTION get_top_5_masters_by_month(
    p_year INT,
    p_month INT
)
RETURNS TABLE(master_name TEXT, works_count BIGINT, cars_count BIGINT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        masters.name::TEXT,
        COUNT(works.id),
        COUNT(DISTINCT works.car_id)
    FROM masters
    JOIN works ON works.master_id = masters.id
    WHERE EXTRACT(YEAR FROM works.date_work) = p_year
      AND EXTRACT(MONTH FROM works.date_work) = p_month
    GROUP BY masters.id, masters.name
    ORDER BY 2 DESC, 3 DESC
    LIMIT 5;
END;
$$;