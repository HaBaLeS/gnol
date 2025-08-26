package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func migrate(db *sqlx.DB) {
	var version int

	//version = getVersion
	noVersion := db.Get(&version, CURRENT_VERSION)
	if noVersion == sql.ErrNoRows {
		version = 0
	}

	if version < 1 {
		db.MustExec(schema_1)
		db.MustExec(UPDATE_VERSION, 1)
	}

	if version < 2 {
		db.MustExec(schema_2)
		db.MustExec(UPDATE_VERSION, 2)
	}

	if version < 3 {
		db.MustExec(schema_3)
		db.MustExec(UPDATE_VERSION, 3)
	}

	if version < 4 {
		db.MustExec(schema_4)
		db.MustExec(UPDATE_VERSION, 4)
	}

	if version < 5 {
		db.MustExec(schema_5)
		db.MustExec(UPDATE_VERSION, 5)
	}

	if version < 6 {
		db.MustExec(schema_6)
		db.MustExec(UPDATE_VERSION, 6)
	}

	if version < 7 {
		db.MustExec(DEFAULT_SERIES)
		db.MustExec(UPDATE_VERSION, 7)
	}

	if version < 8 {
		db.MustExec(schema_8)
		db.MustExec(UPDATE_VERSION, 8)
	}

	if version < 9 {
		db.MustExec(schema_9)
		db.MustExec(UPDATE_VERSION, 9)
	}

	if version < 10 {
		db.MustExec(schema_10)
		db.MustExec(UPDATE_VERSION, 10)
	}

	if version < 11 {
		db.MustExec(schema_11)
		db.MustExec(UPDATE_VERSION, 11)
	}

	if version < 12 {
		db.MustExec(schema_12)
		db.MustExec(UPDATE_VERSION, 12)
	}

	if version < 13 {
		db.MustExec(schema_13)
		db.MustExec(UPDATE_VERSION, 13)
	}
	if version < 14 {
		db.MustExec(schema_14)
		db.MustExec(UPDATE_VERSION, 14)
	}

	if version < 15 {
		db.MustExec(schema_15)
		db.MustExec(UPDATE_VERSION, 15)
	}

}

var schema_1 = `

DROP TABLE IF EXISTS "schema_version";
DROP TABLE IF EXISTS "gnoluser";
DROP TABLE IF EXISTS "comic";
DROP TABLE IF EXISTS "series";
DROP TABLE IF EXISTS "user_to_comic";

CREATE TABLE "schema_version" (
    id SERIAL  PRIMARY KEY  NOT NULL,
    version INTEGER NOT NULL,
  	migration_date timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "gnoluser"(
   id SERIAL  PRIMARY KEY  NOT NULL,
   name text NOT NULL,
   password_hash bytea,
   salt bytea,
   webauthn bool DEFAULT false                  	
);

CREATE TABLE "comic"(
   id SERIAL  PRIMARY KEY  NOT NULL,
   name text NOT NULL,
   series_id INTEGER,
   nsfw bool,
   cover_image_base64 TEXT,
   num_pages INTEGER,
   file_path TEXT
);

CREATE TABLE "series"(
   id SERIAL  PRIMARY KEY  NOT NULL,
   name text NOT NULL,
   cover_image_base64 TEXT NOT NULL
);

CREATE TABLE "user_to_comic"(
   user_id int NOT NULL,
   comic_id int NOT NULL,
   CONSTRAINT user_to_comic_pkey PRIMARY KEY (user_id,comic_id)
);


CREATE UNIQUE INDEX gnoluser_name_key ON "gnoluser"(name);

`

var schema_2 = `

DROP TABLE IF EXISTS "gnoljobs";
CREATE TABLE "gnoljobs" (
	id SERIAL PRIMARY KEY  NOT NULL,
	job_status INTEGER NOT NULL DEFAULT 0,
	user_id int NOT NULL,
	job_type int NOT NULL,
	input_data TEXT
);

`

var schema_3 = `

DROP TABLE IF EXISTS "webauthn_authenticator";
CREATE TABLE "webauthn_authenticator" (
    id SERIAL PRIMARY KEY  NOT NULL,
    aagu_id  bytea NOT NULL,
    signcount bigint NOT NULL,
    clonewarning bool default false
);

DROP TABLE IF EXISTS "webauthn_credential";
CREATE TABLE "webauthn_credential" (
    id bytea PRIMARY KEY NOT NULL,
    publicKey bytea NOT NULL,
    attestationType TEXT,
    authenticator_id int NOT NULL,
    user_id int NOT NULL
);
`

var schema_4 = `
drop table  if exists apitoken;
create table apitoken (
    id SERIAL PRIMARY KEY  NOT NULL,
    token TEXT not null,
    user_id INTEGER UNIQUE not null
);
`

var schema_5 = `
DROP TABLE IF EXISTS tags;
create table tags (
    Id SERIAL PRIMARY KEY ,
    Tag TEXT UNIQUE NOT NULL
);

DROP TABLE IF EXISTS tag_to_comic;
create table tag_to_comic(
    comic_id integer,
    tag_id integer,
    UNIQUE(tag_id, comic_id)
);
insert into tags (Tag) values ('Nsfw');
insert into tags (Tag) values ('Marvel');
insert into tags (Tag) values ('DC');
`

var schema_6 = "alter table user_to_comic add column last_page integer default 0;"

var schema_8 = "alter table comic add column ownerID integer default 1;"

var schema_9 = "alter table comic add column ordernum INTEGER default 100;"

var schema_10 = "alter table comic add column sha256sum text default '';"

var schema_11 = `
alter table series add nsfw boolean default false, add ordernum integer default 100 ;
create table tag_to_series(
    comic_id integer,
    tag_id integer,
    UNIQUE(tag_id, comic_id)
);
`
var schema_12 = "alter table series add ownerid int default 1;"

var schema_13 = `
drop table if exists gnol_session;
create table gnol_session (
	session_id text,
	valid_until timestamp,
	user_id int
);
`

var schema_14 = `
alter table gnoluser  add column nsfw bool default false;
`

var schema_15 = `
	alter table user_to_comic add column finished boolean default false;
	alter table user_to_comic add column finished_at timestamp;
	alter table user_to_comic add column favorite boolean default false;
	alter table user_to_comic add column did_not_read boolean default false;
`
