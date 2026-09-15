UPDATE users
SET roles = jsonb_build_object(
        'values',
        (SELECT jsonb_agg(DISTINCT CASE upper(role_name)
                                       WHEN 'MANAGER' THEN 'MANAGE'
                                       ELSE upper(role_name)
                                   END)
         FROM jsonb_array_elements_text(users.roles -> 'values') AS role_name))
WHERE jsonb_typeof(roles -> 'values') = 'array'
  AND roles -> 'values' @> '["MANAGER"]'::jsonb;
