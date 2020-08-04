---
CREATE OR REPLACE FUNCTION throw_max_goals_reached()
RETURNS VOID 
AS $$
BEGIN
    RAISE EXCEPTION 'Max goals reached' USING ERRCODE = 'HM001';
END;
$$ LANGUAGE plpgsql;

---
CREATE OR REPLACE FUNCTION throw_unknown_goal_type()
RETURNS VOID 
AS $$
BEGIN
    RAISE EXCEPTION 'Unknown goal type' USING ERRCODE = 'HM002';
END;
$$ LANGUAGE plpgsql;

