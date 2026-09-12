UPDATE users
SET roles = jsonb_build_object(
        'values',
        COALESCE(
                (SELECT jsonb_agg(DISTINCT CASE upper(role_name)
                                               WHEN 'PROFILE' THEN 'USER'
                                               WHEN 'ORDER' THEN 'USER'
                                               WHEN 'VIEW' THEN 'USER'
                                               WHEN 'ADMIN' THEN 'MANAGE'
                                               ELSE upper(role_name)
                                           END)
                 FROM jsonb_array_elements_text(users.roles -> 'values') AS role_name),
                '["USER"]'::jsonb))
WHERE jsonb_typeof(roles -> 'values') = 'array'
  AND jsonb_array_length(roles -> 'values') > 0;
