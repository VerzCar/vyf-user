BEGIN;

alter table companies
    drop constraint fk_companies_owner cascade;

drop table users cascade;

drop table locales cascade;

drop table companies cascade;

drop table contacts cascade;

drop table addresses cascade;

drop table countries cascade;

DROP TYPE gender;

COMMIT;