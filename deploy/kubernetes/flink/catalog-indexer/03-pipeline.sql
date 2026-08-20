INSERT INTO book_search_index
SELECT w.uid                                                                  AS uid,
       w.title                                                                AS title,
       w.description                                                          AS description,
       w.ol_key                                                               AS ol_key,
       w.first_publish_date                                                   AS first_publish_date,
       JSON_QUERY(w.cover_uids, '$' RETURNING ARRAY<STRING>)                  AS cover_uids,
       CAST(w.created_at AS TIMESTAMP_LTZ(3))                                 AS created_at,
       CAST(w.updated_at AS TIMESTAMP_LTZ(3))                                 AS updated_at,
       CAST(authors_agg.items AS ARRAY<ROW<
            uid STRING, NAME STRING, bio STRING, birth_date STRING,
            death_date STRING, ol_key STRING, alternate_names ARRAY<STRING>
                >>)  AS authors,
       CAST(subjects_agg.items AS ARRAY<ROW<
            uid STRING, NAME STRING
                >>) AS subjects,
       CAST(editions_agg.items AS ARRAY<ROW<
            uid STRING, title STRING, isbn10 STRING, isbn13 STRING,
            publish_date STRING, languages ARRAY<STRING>, physical_format STRING,
            description STRING,
            publisher ROW<uid STRING, NAME STRING>
                >>) AS editions
FROM work w
    LEFT JOIN (
                  SELECT wam.work_id,
                         ARRAY_AGG(ROW(
                                 a.uid,
                                 a.name,
                                 a.bio,
                                 a.birth_date,
                                 a.death_date,
                                 a.ol_key,
                                 JSON_QUERY(a.alternate_names, '$' RETURNING ARRAY<STRING>)
                                   )) AS items
                  FROM work_author_mapping wam
                      JOIN author a ON a.id = wam.author_id AND a.deleted_at IS NULL
                  GROUP BY wam.work_id
                  ) AS authors_agg ON authors_agg.work_id = w.id
    LEFT JOIN (
                  SELECT wsm.work_id,
                         ARRAY_AGG(ROW(
                                 s.uid,
                                 s.name
                                   )) AS items
                  FROM work_subject_mapping wsm
                      JOIN subject s ON s.id = wsm.subject_id AND s.deleted_at IS NULL
                  GROUP BY wsm.work_id
                  ) AS subjects_agg ON subjects_agg.work_id = w.id
    LEFT JOIN (
                  SELECT e.work_id,
                         ARRAY_AGG(ROW(
                                 e.uid,
                                 e.title,
                                 e.isbn10,
                                 e.isbn13,
                                 e.publish_date,
                                 JSON_QUERY(e.languages, '$' RETURNING ARRAY<STRING>),
                                 e.physical_format,
                                 e.description,
                                 ROW(
                                         p.uid,
                                         p.name
                                 )
                                   )) AS items
                  FROM edition e
                      LEFT JOIN publisher p ON p.id = e.publisher_id AND p.deleted_at IS NULL
                  WHERE e.deleted_at IS NULL
                  GROUP BY e.work_id
                  ) AS editions_agg ON editions_agg.work_id = w.id
WHERE w.deleted_at IS NULL;
