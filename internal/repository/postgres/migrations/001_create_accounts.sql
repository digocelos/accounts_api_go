CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY,
  document text NOT NULL UNIQUE,
  name text NOT NULL,
  email text NULL,
  version integer NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_accounts_document ON accounts(document);