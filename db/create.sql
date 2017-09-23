DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
  id          SERIAL PRIMARY KEY,
  username    VARCHAR(20)  NOT NULL,
  email       VARCHAR(255) NOT NULL,
  password    VARCHAR(60)  NOT NULL,
  enabled     BOOLEAN      NOT NULL DEFAULT FALSE,
  can_login   BOOLEAN      NOT NULL DEFAULT FALSE,
  join_date   TIMESTAMP    NOT NULL DEFAULT NOW(),
  last_login  TIMESTAMP,
  last_access TIMESTAMP,
  uploaded    BIGINT       NOT NULL DEFAULT 0,
  downloaded  BIGINT       NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX users_username_uindex
  ON users (username);
CREATE UNIQUE INDEX users_email_uindex
  ON users (email);

DROP TABLE IF EXISTS release_roles CASCADE;
CREATE TABLE release_roles
(
  id   SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);
CREATE UNIQUE INDEX release_roles_name_uindex
  ON release_roles (name);

DROP TABLE IF EXISTS blog_tags CASCADE;
CREATE TABLE blog_tags
(
  id  SERIAL PRIMARY KEY,
  tag VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX blog_tags_tag_uindex
  ON blog_tags (tag);

DROP TABLE IF EXISTS blogs CASCADE;
CREATE TABLE blogs
(
  id        SERIAL PRIMARY KEY,
  title     VARCHAR(255) NOT NULL,
  content   TEXT         NOT NULL,
  author    INT          NOT NULL,
  posted_at TIMESTAMP    NOT NULL,
  CONSTRAINT blogs_users_id_fk FOREIGN KEY (author) REFERENCES users (id)
);

DROP TABLE IF EXISTS blog_tags_blogs CASCADE;
CREATE TABLE blog_tags_blogs
(
  blog INT NOT NULL,
  tag  INT NOT NULL,
  PRIMARY KEY (blog, tag),
  CONSTRAINT blog_tags_blogs_blogs_id_fk FOREIGN KEY (blog) REFERENCES blogs (id),
  CONSTRAINT blog_tags_blogs_blog_tags_id_fk FOREIGN KEY (TAG) REFERENCES blog_tags (id)
);

-- Music-specific tables
DROP TABLE IF EXISTS record_labels CASCADE;
CREATE TABLE record_labels
(
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description TEXT,
  founded     DATE,
  added       TIMESTAMP    NOT NULL,
  added_by    INT          NOT NULL,
  CONSTRAINT record_labels_users_id_fk FOREIGN KEY (added_by) REFERENCES users (id)
);

DROP TABLE IF EXISTS release_tags CASCADE;
CREATE TABLE release_tags
(
  id  SERIAL PRIMARY KEY,
  tag VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX release_tags_tag_uindex
  ON release_tags (tag);

DROP TABLE IF EXISTS media CASCADE;
CREATE TABLE media
(
  id     SERIAL PRIMARY KEY,
  medium VARCHAR(20) NOT NULL
);
CREATE UNIQUE INDEX medias_media_uindex
  ON media (medium);

DROP TABLE IF EXISTS artist_tags CASCADE;
CREATE TABLE artist_tags
(
  id  SERIAL PRIMARY KEY,
  tag VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX artist_tags_tag_uindex
  ON artist_tags (tag);

DROP TABLE IF EXISTS release_group_tags CASCADE;
CREATE TABLE release_group_tags
(
  id  SERIAL PRIMARY KEY,
  tag VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX release_group_tags_tag_uindex
  ON release_group_tags (tag);

DROP TABLE IF EXISTS release_group_types CASCADE;
CREATE TABLE release_group_types
(
  id   SERIAL PRIMARY KEY,
  type VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX release_group_types_type_uindex
  ON release_group_types (type);

DROP TABLE IF EXISTS artists CASCADE;
CREATE TABLE artists
(
  id   SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  TEXT
);

DROP TABLE IF EXISTS artist_aliases CASCADE;
CREATE TABLE artist_aliases
(
  artist INT          NOT NULL,
  alias  VARCHAR(255) NOT NULL,
  PRIMARY KEY (artist, alias),
  CONSTRAINT artist_aliases_artists_id_fk FOREIGN KEY (artist) REFERENCES artists (id)
);

DROP TABLE IF EXISTS release_groups CASCADE;
CREATE TABLE release_groups
(
  id           SERIAL PRIMARY KEY,
  name         VARCHAR(255) NOT NULL,
  type         INT          NOT NULL,
  release_date TIMESTAMP    NOT NULL,
  CONSTRAINT release_groups_release_group_types_id_fk FOREIGN KEY (type) REFERENCES release_group_types (id)
);

DROP TABLE IF EXISTS release_groups_artists CASCADE;
CREATE TABLE release_groups_artists
(
  release_group INT NOT NULL,
  artist        INT NOT NULL,
  role          INT NOT NULL,
  PRIMARY KEY (release_group, artist, role),
  CONSTRAINT release_groups_artists_artists_artist_id_fk FOREIGN KEY (artist) REFERENCES artists (id),
  CONSTRAINT release_groups_artists_release_groups_release_group_id_fk FOREIGN KEY (release_group) REFERENCES release_groups (id),
  CONSTRAINT release_groups_artists_release_roles_role_id_fk FOREIGN KEY (role) REFERENCES release_roles (id)
);

DROP TABLE IF EXISTS releases CASCADE;
CREATE TABLE releases
(
  id               SERIAL PRIMARY KEY,
  edition          VARCHAR(255) NOT NULL DEFAULT '',
  medium           INT          NOT NULL,
  release_group    INT          NOT NULL,
  record_label     INT          NOT NULL,
  added            TIMESTAMP    NOT NULL,
  added_by         INT          NOT NULL,
  release_date     DATE,
  catalogue_number VARCHAR(50),
  original         BOOLEAN,
  CONSTRAINT releases_media_id_fk FOREIGN KEY (medium) REFERENCES media (id),
  CONSTRAINT releases_record_label_id_fk FOREIGN KEY (record_label) REFERENCES record_labels (id),
  CONSTRAINT releases_release_group_id_fk FOREIGN KEY (release_group) REFERENCES release_groups (id),
  CONSTRAINT releases_users_id_fk FOREIGN KEY (added_by) REFERENCES users (id)
);

DROP TABLE IF EXISTS artist_tags_artists CASCADE;
CREATE TABLE artist_tags_artists
(
  artist INT NOT NULL,
  tag    INT NOT NULL,
  PRIMARY KEY (artist, tag),
  CONSTRAINT artist_tags_artists_artists_id_fk FOREIGN KEY (artist) REFERENCES artists (id),
  CONSTRAINT artist_tags_artists_artist_tags_id_fk FOREIGN KEY (tag) REFERENCES artist_tags (id)
);

DROP TABLE IF EXISTS release_group_tags_release_groups CASCADE;
CREATE TABLE release_group_tags_release_groups
(
  release_group INT NOT NULL,
  tag           INT NOT NULL,
  PRIMARY KEY (release_group, tag),
  CONSTRAINT release_group_tags_release_groups_release_groups_id_fk FOREIGN KEY (release_group) REFERENCES release_groups (id),
  CONSTRAINT release_group_tags_release_groups_release_group_tags_id_fk FOREIGN KEY (tag) REFERENCES release_group_tags (id)
);

DROP TABLE IF EXISTS release_tags_releases CASCADE;
CREATE TABLE release_tags_releases
(
  release INT NOT NULL,
  tag     INT NOT NULL,
  PRIMARY KEY (release, tag),
  CONSTRAINT release_tags_releases_releases_id_fk FOREIGN KEY (release) REFERENCES releases (id),
  CONSTRAINT release_tags_releases_release_tags_id_fk FOREIGN KEY (tag) REFERENCES release_tags (id)
);

-- Other
DROP TABLE IF EXISTS api_tokens CASCADE;
CREATE TABLE api_tokens
(
  token      VARCHAR(128) PRIMARY KEY,
  uid        INT       NOT NULL,
  created_at TIMESTAMP NOT NULL,
  CONSTRAINT api_tokens_users_id_fk FOREIGN KEY (uid) REFERENCES users (id)
);

-- Populate
INSERT INTO release_group_types (id, type) VALUES
  (0, 'Album'),
  (1, 'EP'),
  (2, 'Single'),
  (3, 'Compilation'),
  (4, 'Soundtrack'),
  (5, 'Live album'),
  (6, 'Bootleg'),
  (7, 'Mixtape'),
  (8, 'Unknown');

INSERT INTO media (id, medium) VALUES
  (0, 'CD'),
  (1, 'DVD'),
  (2, 'Vinyl'),
  (3, 'WEB'),
  (4, 'Blu-Ray'),
  (5, 'SACD'),
  (6, 'Cassette'),
  (7, 'DAT');

INSERT INTO release_roles (id, name) VALUES
  (0, 'Main'),
  (1, 'Guest'),
  (2, 'Composer'),
  (3, 'Conductor'),
  (4, 'Remixer'),
  (5, 'Producer');

INSERT INTO users (id, username, email, password, enabled, can_login, join_date, last_login, last_access, uploaded, downloaded)
VALUES
  (0, 'boiling', 'boiling@boiling.rip', '', TRUE, FALSE, '2000-01-01 00:00',
      '2000-01-01 00:00', '2000-01-01 00:00', 0,
      0),
  (DEFAULT, 'test', 'test@boiling.rip',
            '$2a$14$2v2YkEAjBx9ZEYZdYQgDR.H4r.CmdOTI.10cmqnKvQ7Ucq60prUGm',
            TRUE,
            TRUE, '2000-01-01 00:00', '2000-01-01 00:00', '2000-01-01 00:00', 0,
            0); --password is test
