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

CREATE TABLE IF NOT EXISTS operation_records (
  id BIGSERIAL PRIMARY KEY,
  module TEXT NOT NULL,
  title TEXT NOT NULL,
  member_id BIGINT REFERENCES members(id),
  amount_cents INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'pending',
  due_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_operation_records_module ON operation_records(module);
CREATE INDEX IF NOT EXISTS idx_operation_records_status ON operation_records(status);

INSERT INTO plans (name, price_cents, duration_days, description)
SELECT '月度会员', 9900, 30, '适合短期体验的标准会员套餐'
WHERE NOT EXISTS (SELECT 1 FROM plans);

INSERT INTO plans (name, price_cents, duration_days, description)
SELECT '年度会员', 99900, 365, '适合长期使用的优惠会员套餐'
WHERE (SELECT count(*) FROM plans) = 1;

INSERT INTO operation_records (module, title, amount_cents, status)
SELECT 'booking', '明天 10:30 到店预约', 0, 'pending'
WHERE NOT EXISTS (SELECT 1 FROM operation_records);

INSERT INTO operation_records (module, title, amount_cents, status)
SELECT 'finance', '会员收入日报', 101170, 'done'
WHERE (SELECT count(*) FROM operation_records) = 1;
`
