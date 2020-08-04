CREATE OR REPLACE FUNCTION create_goal (in_user_id goal.user_id%TYPE
                                        , in_type goal.type%TYPE
                                        , in_category goal.category%TYPE
                                        , in_value goal.value%TYPE
                                        , in_name goal.name%TYPE
                                        , in_description goal.description%TYPE
                                        , in_starts_at goal.starts_at%TYPE
                                        , in_expires_at goal.expires_at%TYPE)
RETURNS VOID
AS $$
DECLARE
    v_now TIMESTAMP = now() AT TIME ZONE 'utc';
    v_goal_count INTEGER;
BEGIN
    -- validate and insert a weekly personal goal
    IF in_category = 'weekly' THEN
        -- check for max goals
        SELECT COUNT(1)
        FROM goal
        WHERE user_id = in_user_id
        AND type = in_type
        AND category = in_category
        AND EXTRACT('week' FROM starts_at) = EXTRACT('week' FROM in_starts_at)
        INTO v_goal_count;

        IF in_type = 'personal' AND v_goal_count >= 3 THEN
            PERFORM throw_max_goals_reached();
        END IF;

        IF in_type = 'point' AND v_goal_count >= 1 THEN
            PERFORM throw_max_goals_reached(); 
        END IF;
    ELSE
        PERFORM throw_unknown_goal_type();
    END IF;

    INSERT INTO goal 
    (
        user_id
        , type
        , category
        , value
        , name
        , description
        , starts_at
        , expires_at
        , created_at
    )
    VALUES
    (
        in_user_id
        , in_type
        , in_category
        , in_value
        , in_name
        , in_description
        , in_starts_at
        , in_expires_at
        , v_now
    );
END;
$$
LANGUAGE plpgsql;
