---
 CREATE OR REPLACE FUNCTION clear_test_data()
 RETURNS void 
 AS $$
 BEGIN
    IF (SELECT (SELECT count(current_database) FROM current_database() WHERE current_database LIKE '%-test%') = 0) THEN
        RETURN;
    END IF;

    DELETE FROM goal;
 END;
 $$ LANGUAGE plpgsql;
 