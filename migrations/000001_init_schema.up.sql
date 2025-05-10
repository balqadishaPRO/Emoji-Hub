CREATE TABLE emoji (
  id         UUID PRIMARY KEY,
  name       TEXT,
  category   TEXT,
  "group"    TEXT,
  html_code  TEXT[],
  unicode    TEXT[]
);

CREATE TABLE favorites (
  session_id UUID,
  emoji_id   UUID REFERENCES emoji(id),
  PRIMARY KEY (session_id, emoji_id)
);

CREATE TABLE llm_cache (
  emoji_id UUID PRIMARY KEY,
  mood     TEXT,
  updated  TIMESTAMP DEFAULT NOW()
);
