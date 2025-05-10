CREATE TABLE emoji (
  id         UUID PRIMARY KEY,
  name       TEXT,
  category   TEXT,
  "group"    TEXT,
  html_code  TEXT[],
  unicode    TEXT[]
);

-- Для избранного; пока без авторизации
CREATE TABLE favorites (
  id         SERIAL PRIMARY KEY,
  emoji_id   UUID REFERENCES emoji(id)
);
