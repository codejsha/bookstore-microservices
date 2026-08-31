ALTER TABLE edition DROP CONSTRAINT IF EXISTS chk_edition_isbn10_length;
ALTER TABLE edition DROP CONSTRAINT IF EXISTS chk_edition_isbn13_length;

ALTER TABLE edition DROP CONSTRAINT IF EXISTS fk_edition__work;
ALTER TABLE edition DROP CONSTRAINT IF EXISTS fk_edition__publisher;

ALTER TABLE work_author_mapping DROP CONSTRAINT IF EXISTS fk_work_author__work;
ALTER TABLE work_author_mapping DROP CONSTRAINT IF EXISTS fk_work_author__author;

ALTER TABLE work_subject_mapping DROP CONSTRAINT IF EXISTS fk_work_subject__work;
ALTER TABLE work_subject_mapping DROP CONSTRAINT IF EXISTS fk_work_subject__subject;
