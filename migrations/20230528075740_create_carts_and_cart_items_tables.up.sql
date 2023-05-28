CREATE TABLE carts (
  id BIGSERIAL PRIMARY KEY,
  uid TEXT NOT NULL,
  quantity INT NOT NULL,
  total_price TEXT NOT NULL,
  total_price_value BIGINT NOT NULL,
  total_weight TEXT NOT NULL,
  total_weight_value NUMERIC(10, 2) NOT NULL,
  user_id BIGINT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL,

  FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TABLE cart_items (
  id BIGSERIAL PRIMARY KEY,
  uid TEXT NOT NULL,
  quantity INT NOT NULL,
  total_price TEXT NOT NULL,
  total_price_value BIGINT NOT NULL,
  total_weight TEXT NOT NULL,
  total_weight_value NUMERIC(10, 2) NOT NULL,
  product_name TEXT NOT NULL,
  product_slug TEXT NOT NULL,
  product_image TEXT NOT NULL,
  product_weight TEXT NOT NULL,
  product_weight_value NUMERIC(10, 2) NOT NULL,
  base_price TEXT NOT NULL,
  base_price_value INT NOT NULL,
  offer_price TEXT NOT NULL,
  offer_price_value INT NOT NULL,
  discount SMALLINT NOT NULL,
  cart_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL,

  FOREIGN KEY(product_id)
    REFERENCES products(id)
    ON DELETE CASCADE,
  FOREIGN KEY(cart_id)
    REFERENCES carts(id)
    ON DELETE CASCADE
);