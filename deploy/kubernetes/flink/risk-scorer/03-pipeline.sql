INSERT INTO risk_flag
SELECT
    CONCAT('identity:risk:auto:', user_uid),
    CONCAT(
        '{"sub":"', user_uid,
        '","level":"restrict","reason":"auto: score ', CAST(score AS STRING),
        ' (', CAST(total_requests AS STRING), ' req, ',
        CAST(error_requests AS STRING), ' errors, ',
        CAST(distinct_paths AS STRING), ' paths in 5m)",',
        '"flagged_by":"risk-scorer",',
        '"flagged_at":"', DATE_FORMAT(window_end, 'yyyy-MM-dd''T''HH:mm:ss''Z'''), '",',
        '"expires_at":"', DATE_FORMAT(window_end + INTERVAL '1' SECOND * ${AUTO_FLAG_TTL_SECONDS}, 'yyyy-MM-dd''T''HH:mm:ss''Z'''), '"}'
    )
FROM (
    SELECT
        user_uid,
        window_end,
        COUNT(*) AS total_requests,
        SUM(CASE WHEN status >= 400 THEN 1 ELSE 0 END) AS error_requests,
        COUNT(DISTINCT `path`) AS distinct_paths,
        COUNT(*)
            + 4 * SUM(CASE WHEN status >= 400 THEN 1 ELSE 0 END)
            + 2 * COUNT(DISTINCT `path`) AS score
    FROM TABLE(HOP(TABLE access_log, DESCRIPTOR(kafka_ts), INTERVAL '1' MINUTE, INTERVAL '5' MINUTE))
    WHERE user_uid IS NOT NULL AND user_uid <> ''
    GROUP BY user_uid, window_start, window_end
)
WHERE score >= ${RISK_SCORE_THRESHOLD};
