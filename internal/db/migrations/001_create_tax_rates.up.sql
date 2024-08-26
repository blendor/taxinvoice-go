CREATE TABLE IF NOT EXISTS tax_rates (
    id SERIAL PRIMARY KEY,
    state VARCHAR(2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    effective_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_tax_rates_state_date ON tax_rates (state, effective_date);