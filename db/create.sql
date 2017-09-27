DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
  id          SERIAL PRIMARY KEY,
  username    VARCHAR(20)  NOT NULL,
  email       VARCHAR(255) NOT NULL,
  password    VARCHAR(60)  NOT NULL,
  bio         TEXT,
  enabled     BOOLEAN      NOT NULL DEFAULT FALSE,
  can_login   BOOLEAN      NOT NULL DEFAULT FALSE,
  joined_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
  last_login  TIMESTAMP,
  last_access TIMESTAMP,
  uploaded    BIGINT       NOT NULL DEFAULT 0,
  downloaded  BIGINT       NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX users_username_uindex
  ON users (username);
CREATE UNIQUE INDEX users_email_uindex
  ON users (email);

DROP TABLE IF EXISTS user_passkeys CASCADE;
CREATE TABLE user_passkeys
(
  uid        INT PRIMARY KEY,
  passkey    VARCHAR(64) NOT NULL,
  created_at TIMESTAMP   NOT NULL,
  valid      BOOLEAN     NOT NULL,
  CONSTRAINT user_passkeys_users_id_fk FOREIGN KEY (uid) REFERENCES users (id)
);
CREATE UNIQUE INDEX user_passkeys_passkey_uindex
  ON user_passkeys (passkey);

DROP TABLE IF EXISTS privileges CASCADE;
CREATE TABLE privileges
(
  id        SERIAL PRIMARY KEY,
  privilege VARCHAR(128) NOT NULL
);
CREATE UNIQUE INDEX privileges_privilege_uindex
  ON privileges (privilege);

DROP TABLE IF EXISTS users_privileges CASCADE;
CREATE TABLE users_privileges
(
  uid       INT NOT NULL,
  privilege INT NOT NULL,
  PRIMARY KEY (uid, privilege)
);

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

DROP TABLE IF EXISTS release_roles CASCADE;
CREATE TABLE release_roles
(
  id   SERIAL PRIMARY KEY,
  role VARCHAR(100) NOT NULL
);
CREATE UNIQUE INDEX release_roles_role_uindex
  ON release_roles (role);

DROP TABLE IF EXISTS release_tags CASCADE;
CREATE TABLE release_tags
(
  id  SERIAL PRIMARY KEY,
  tag VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX release_tags_tag_uindex
  ON release_tags (tag);

DROP TABLE IF EXISTS release_properties CASCADE;
CREATE TABLE release_properties
(
  id       SERIAL PRIMARY KEY,
  property VARCHAR(100) NOT NULL
);
CREATE UNIQUE INDEX release_properties_property_uindex
  ON release_properties (property);

DROP TABLE IF EXISTS media CASCADE;
CREATE TABLE media
(
  id     SERIAL PRIMARY KEY,
  medium VARCHAR(20) NOT NULL
);
CREATE UNIQUE INDEX medias_media_uindex
  ON media (medium);

DROP TABLE IF EXISTS leech_types CASCADE;
CREATE TABLE leech_types
(
  id   SERIAL PRIMARY KEY,
  type VARCHAR(50) NOT NULL
);
CREATE UNIQUE INDEX leech_types_type_uindex
  ON leech_types (type);

DROP TABLE IF EXISTS formats CASCADE;
CREATE TABLE formats
(
  id       SERIAL PRIMARY KEY,
  format   VARCHAR(20) NOT NULL,
  encoding VARCHAR(20) NOT NULL
);
CREATE UNIQUE INDEX formats_format_uindex
  ON formats (format);

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
  id       SERIAL PRIMARY KEY,
  name     VARCHAR(255) NOT NULL,
  bio      TEXT,
  added    TIMESTAMP    NOT NULL,
  added_by INT          NOT NULL,
  CONSTRAINT artists_users_id_fk FOREIGN KEY (added_by) REFERENCES users (id)
);

DROP TABLE IF EXISTS artist_aliases CASCADE;
CREATE TABLE artist_aliases
(
  artist   INT          NOT NULL,
  alias    VARCHAR(255) NOT NULL,
  added    TIMESTAMP    NOT NULL,
  added_by INT          NOT NULL,
  PRIMARY KEY (artist, alias),
  CONSTRAINT artist_aliases_artists_id_fk FOREIGN KEY (artist) REFERENCES artists (id),
  CONSTRAINT artist_aliases_users_id_fk FOREIGN KEY (added_by) REFERENCES users (id)
);

DROP TABLE IF EXISTS release_groups CASCADE;
CREATE TABLE release_groups
(
  id           SERIAL PRIMARY KEY,
  name         VARCHAR(255) NOT NULL,
  added        TIMESTAMP    NOT NULL,
  added_by     INT          NOT NULL,
  type         INT          NOT NULL,
  release_date TIMESTAMP    NOT NULL,
  CONSTRAINT release_groups_users_id_fk FOREIGN KEY (added_by) REFERENCES users (id),
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

DROP TABLE IF EXISTS torrents CASCADE;
CREATE TABLE torrents
(
  id        SERIAL PRIMARY KEY,
  release   INT       NOT NULL,
  uploaded  TIMESTAMP NOT NULL,
  uploader  INT       NOT NULL,
  info_hash BYTEA     NOT NULL,
  format    INT       NOT NULL,
  size      BIGINT    NOT NULL,
  comment   VARCHAR(255),
  CONSTRAINT torrents_releases_id_fk FOREIGN KEY (release) REFERENCES releases (id),
  CONSTRAINT torrents_users_id_fk FOREIGN KEY (uploader) REFERENCES users (id),
  CONSTRAINT torrents_formats_id_fk FOREIGN KEY (format) REFERENCES formats (id)
);
CREATE UNIQUE INDEX torrents_info_hash_uindex
  ON torrents (info_hash);

DROP TABLE IF EXISTS torrent_trackerdata CASCADE;
CREATE TABLE torrent_trackerdata
(
  torrent          INT PRIMARY KEY,
  leech_type       INT    NOT NULL,
  seeders          INT    NOT NULL DEFAULT 0,
  leechers         INT    NOT NULL DEFAULT 0,
  snatches         INT    NOT NULL DEFAULT 0,
  total_uploaded   BIGINT NOT NULL DEFAULT 0,
  total_downloaded BIGINT NOT NULL DEFAULT 0,
  CONSTRAINT torrent_metadata_torrents_id_fk FOREIGN KEY (torrent) REFERENCES torrents (id),
  CONSTRAINT torrent_metadata_leech_types_id_fk FOREIGN KEY (leech_type) REFERENCES leech_types (id)
);

DROP TABLE IF EXISTS user_stat_changes CASCADE;
CREATE TABLE user_stat_changes
(
  uid              INT PRIMARY KEY,
  torrent          INT       NOT NULL,
  reported_at      TIMESTAMP NOT NULL,
  uploaded_delta   INT       NOT NULL DEFAULT 0,
  downloaded_delta INT       NOT NULL DEFAULT 0,
  CONSTRAINT user_stat_changes_users_id_fk FOREIGN KEY (uid) REFERENCES users (id),
  CONSTRAINT user_stat_changes_torrents_id_fk FOREIGN KEY (torrent) REFERENCES torrents (id)
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

DROP TABLE IF EXISTS release_properties_releases CASCADE;
CREATE TABLE release_properties_releases
(
  release  INT NOT NULL,
  property INT NOT NULL,
  value    VARCHAR(100),
  PRIMARY KEY (release, property),
  CONSTRAINT release_properties_releases_releases_id_fk FOREIGN KEY (release) REFERENCES releases (id),
  CONSTRAINT release_properties_releases_release_properties_id_fk FOREIGN KEY (property) REFERENCES release_properties (id)
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
INSERT INTO privileges (id, privilege) VALUES
  (0, 'get_blogs'),
  (1, 'post_blog'),
  (2, 'update_blog'),
  (3, 'update_blog_set_posted_at'),
  (4, 'update_blog_set_posted_by'),
  (5, 'update_blog_set_tags');
ALTER SEQUENCE privileges_id_seq RESTART WITH 6;

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
ALTER SEQUENCE release_group_types_id_seq RESTART WITH 9;

INSERT INTO media (id, medium) VALUES
  (0, 'CD'),
  (1, 'DVD'),
  (2, 'Vinyl'),
  (3, 'WEB'),
  (4, 'Blu-Ray'),
  (5, 'SACD'),
  (6, 'Cassette'),
  (7, 'DAT');
ALTER SEQUENCE media_id_seq RESTART WITH 8;

INSERT INTO release_roles (id, role) VALUES
  (0, 'Main'),
  (1, 'Guest'),
  (2, 'Composer'),
  (3, 'Conductor'),
  (4, 'Remixer'),
  (5, 'Producer');
ALTER SEQUENCE release_roles_id_seq RESTART WITH 6;

INSERT INTO formats (id, format, encoding) VALUES
  (0, 'FLAC', 'Lossless'),
  (1, 'FLAC/24bit', 'Lossless'),
  (2, 'MP3/320', 'Lossy'),
  (3, 'MP3/V0', 'Lossy'),
  (4, 'MP3/V2', 'Lossy');
ALTER SEQUENCE formats_id_seq RESTART WITH 5;

INSERT INTO release_properties (id, property) VALUES
  (0, 'LossyMasterApproved'),
  (1, 'LossyWebApproved'),
  (2, 'CassetteApproved');
ALTER SEQUENCE release_properties_id_seq RESTART WITH 3;

INSERT INTO leech_types (id, type) VALUES
  (0, 'Normal'),
  (1, 'Freeleech'),
  (2, 'Neutral'),
  (3, 'DoubleUp'),
  (4, 'DoubleDown');
ALTER SEQUENCE leech_types_id_seq RESTART WITH 5;

INSERT INTO users (id, username, email, password, bio, enabled, can_login, joined_at, last_login, last_access, uploaded, downloaded)
VALUES
  (0, 'boiling', 'boiling@boiling.rip', '', 'The one', TRUE, FALSE,
      '2000-01-01 00:00',
      '2000-01-01 00:00', '2000-01-01 00:00', 0,
   0),
  (1, 'test', 'test@boiling.rip',
      '$2a$14$2v2YkEAjBx9ZEYZdYQgDR.H4r.CmdOTI.10cmqnKvQ7Ucq60prUGm', '',
      TRUE, TRUE, '2000-01-01 00:00', '2000-01-01 00:00',
      '2000-01-01 00:00', 0, 0); --password is test
ALTER SEQUENCE users_id_seq RESTART WITH 2;

INSERT INTO users_privileges(uid, privilege) SELECT 1,id FROM privileges;