CREATE TABLE emoji (
  id         TEXT PRIMARY KEY,
  name       TEXT,
  category   TEXT,
  "group"    TEXT,
  html_code  TEXT[],
  unicode    TEXT[]
);

CREATE TABLE favorites (
  session_id TEXT,
  emoji_id   TEXT REFERENCES emoji(id),
  created_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (session_id, emoji_id)
);

CREATE TABLE llm_cache (
  emoji_id TEXT PRIMARY KEY,
  mood     TEXT,
  updated  TIMESTAMP DEFAULT NOW()
);
