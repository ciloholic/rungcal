SELECT
    e.id,
    se.uuid,
    e.project,
    CONCAT(se.project, '/', se.group_path, '/', se.job_name)                                  AS job,
    e.date_started,
    e.date_completed,
    CONCAT_WS(' ', se.seconds, se.minute, se.hour, se.day_of_month, se.month, se.day_of_week) AS schedule,
    TIMESTAMPDIFF(SECOND, e.date_started, e.date_completed)                                   AS execution_time
FROM
    scheduled_execution AS se
    JOIN execution AS e ON se.id = e.scheduled_execution_id
WHERE
    e.execution_type = 'scheduled'
    AND e.status = 'succeeded'
    AND (e.date_started BETWEEN ? AND ?)
ORDER BY e.date_started;
