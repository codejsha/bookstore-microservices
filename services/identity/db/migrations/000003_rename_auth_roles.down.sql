UPDATE users
SET roles = jsonb_build_object(
        'values',
        COALESCE(
                (SELECT jsonb_agg(DISTINCT reverted_role)
                 FROM jsonb_array_elements_text(users.roles -> 'values') AS role_name,
                      LATERAL unnest(CASE upper(role_name)
                                         WHEN 'USER' THEN ARRAY ['PROFILE', 'ORDER', 'VIEW']
                                         ELSE ARRAY [upper(role_name)]
                                     END) AS reverted_role),
                '["PROFILE","ORDER","VIEW"]'::jsonb))
WHERE jsonb_typeof(roles -> 'values') = 'array'
  AND jsonb_array_length(roles -> 'values') > 0;
