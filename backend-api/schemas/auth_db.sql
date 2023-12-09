CREATE DATABASE cloud_render_auth;

\c cloud_render_auth;

CREATE TABLE IF NOT EXISTS users(
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    login text UNIQUE NOT NULL,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    refresh_token text
);
