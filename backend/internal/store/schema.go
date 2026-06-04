package store

const schemaSQL = `
CREATE TABLE IF NOT EXISTS plans (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
  duration_days INTEGER NOT NULL CHECK (duration_days > 0),
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS members (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  phone TEXT NOT NULL,
  email TEXT NOT NULL,
  plan_id BIGINT NOT NULL REFERENCES plans(id),
  status TEXT NOT NULL DEFAULT 'pending',
  start_at TIMESTAMPTZ,
  expire_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_members_status ON members(status);
CREATE INDEX IF NOT EXISTS idx_members_phone ON members(phone);

INSERT INTO plans (name, price_cents, duration_days, description)
SELECT '月度会员', 9900, 30, '适合短期体验的标准会员套餐'
WHERE NOT EXISTS (SELECT 1 FROM plans);

INSERT INTO plans (name, price_cents, duration_days, description)
SELECT '年度会员', 99900, 365, '适合长期使用的优惠会员套餐'
WHERE (SELECT count(*) FROM plans) = 1;
`
