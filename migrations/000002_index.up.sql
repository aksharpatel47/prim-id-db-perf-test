create index address_person_id_bigserial on addresses_bigserial (person_id);
create index address_person_id_random_int on addresses_random_int (person_id);
create index address_person_id_date_random_int on addresses_date_random_int (person_id);
create index address_person_id_uuid on addresses_uuid (person_id);
create index addresses_person_id_uuidv7 on addresses_uuidv7 (person_id);