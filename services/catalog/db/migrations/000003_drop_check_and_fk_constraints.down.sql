ALTER TABLE edition
    ADD CONSTRAINT chk_edition_isbn10_length
        CHECK (isbn10 IS NULL OR CHAR_LENGTH(isbn10) = 10);
ALTER TABLE edition
    ADD CONSTRAINT chk_edition_isbn13_length
        CHECK (isbn13 IS NULL OR CHAR_LENGTH(isbn13) = 13);

ALTER TABLE edition
    ADD CONSTRAINT fk_edition__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE RESTRICT;
ALTER TABLE edition
    ADD CONSTRAINT fk_edition__publisher
        FOREIGN KEY (publisher_id) REFERENCES publisher (id) ON DELETE SET NULL;

ALTER TABLE work_author_mapping
    ADD CONSTRAINT fk_work_author__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE CASCADE;
ALTER TABLE work_author_mapping
    ADD CONSTRAINT fk_work_author__author
        FOREIGN KEY (author_id) REFERENCES author (id) ON DELETE CASCADE;

ALTER TABLE work_subject_mapping
    ADD CONSTRAINT fk_work_subject__work
        FOREIGN KEY (work_id) REFERENCES work (id) ON DELETE CASCADE;
ALTER TABLE work_subject_mapping
    ADD CONSTRAINT fk_work_subject__subject
        FOREIGN KEY (subject_id) REFERENCES subject (id) ON DELETE CASCADE;
