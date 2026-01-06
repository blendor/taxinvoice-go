ALTER TABLE invoices
ADD COLUMN status VARCHAR(20) DEFAULT 'unpaid',
ADD COLUMN due_date DATE,
ADD COLUMN paid_at TIMESTAMP;

CREATE INDEX idx_invoices_status ON invoices(status);
