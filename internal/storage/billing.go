package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

var _ billing.Repository = (*SQLiteStore)(nil)

func (s *SQLiteStore) CreateCustomer(ctx context.Context, c billing.Customer) (billing.Customer, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO customers (id, email, name, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`, c.ID, c.Email, c.Name, encodeMap(c.Metadata), encodeTime(c.CreatedAt), encodeTime(c.CreatedAt)); err != nil {
		return billing.Customer{}, err
	}
	return s.GetCustomer(ctx, c.ID)
}

func (s *SQLiteStore) GetCustomer(ctx context.Context, id string) (billing.Customer, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, email, name, metadata, created_at FROM customers WHERE id = ?`, id)
	c, err := scanCustomer(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.Customer{}, billing.ErrNotFound
	}
	return c, err
}

func (s *SQLiteStore) ListCustomers(ctx context.Context) ([]billing.Customer, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, email, name, metadata, created_at FROM customers ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateCustomer(ctx context.Context, id string, in billing.Customer) (billing.Customer, error) {
	current, err := s.GetCustomer(ctx, id)
	if err != nil {
		return billing.Customer{}, err
	}
	if in.Email != "" {
		current.Email = in.Email
	}
	if in.Name != "" {
		current.Name = in.Name
	}
	if in.Metadata != nil {
		current.Metadata = in.Metadata
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE customers SET email = ?, name = ?, metadata = ?, updated_at = ? WHERE id = ?`,
		current.Email, current.Name, encodeMap(current.Metadata), encodeTime(time.Now().UTC()), id); err != nil {
		return billing.Customer{}, err
	}
	return s.GetCustomer(ctx, id)
}

func (s *SQLiteStore) CreateProduct(ctx context.Context, p billing.Product) (billing.Product, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO products (id, name, description, active, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, p.ID, p.Name, p.Description, boolInt(p.Active), encodeMap(p.Metadata), encodeTime(p.CreatedAt), encodeTime(p.CreatedAt)); err != nil {
		return billing.Product{}, err
	}
	return s.GetProduct(ctx, p.ID)
}

func (s *SQLiteStore) GetProduct(ctx context.Context, id string) (billing.Product, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, name, description, active, metadata, created_at FROM products WHERE id = ?`, id)
	p, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.Product{}, billing.ErrNotFound
	}
	return p, err
}

func (s *SQLiteStore) ListProducts(ctx context.Context) ([]billing.Product, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, active, metadata, created_at FROM products ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateProduct(ctx context.Context, id string, in billing.Product) (billing.Product, error) {
	current, err := s.GetProduct(ctx, id)
	if err != nil {
		return billing.Product{}, err
	}
	if in.Name != "" {
		current.Name = in.Name
	}
	if in.Description != "" {
		current.Description = in.Description
	}
	current.Active = in.Active
	if in.Metadata != nil {
		current.Metadata = in.Metadata
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE products SET name = ?, description = ?, active = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		current.Name, current.Description, boolInt(current.Active), encodeMap(current.Metadata), id); err != nil {
		return billing.Product{}, err
	}
	return s.GetProduct(ctx, id)
}

func (s *SQLiteStore) CreatePrice(ctx context.Context, p billing.Price) (billing.Price, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO prices (id, product_id, currency, unit_amount, lookup_key, recurring_interval, recurring_interval_count, active, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, p.ID, p.ProductID, strings.ToLower(p.Currency), p.UnitAmount, p.LookupKey, p.RecurringInterval, p.RecurringIntervalCount, boolInt(p.Active), encodeMap(p.Metadata), encodeTime(p.CreatedAt), encodeTime(p.CreatedAt)); err != nil {
		return billing.Price{}, err
	}
	return s.GetPrice(ctx, p.ID)
}

func (s *SQLiteStore) GetPrice(ctx context.Context, id string) (billing.Price, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, product_id, currency, unit_amount, lookup_key, recurring_interval, recurring_interval_count, active, metadata, created_at FROM prices WHERE id = ?`, id)
	p, err := scanPrice(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.Price{}, billing.ErrNotFound
	}
	return p, err
}

func (s *SQLiteStore) ListPrices(ctx context.Context) ([]billing.Price, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, product_id, currency, unit_amount, lookup_key, recurring_interval, recurring_interval_count, active, metadata, created_at FROM prices ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Price
	for rows.Next() {
		p, err := scanPrice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdatePrice(ctx context.Context, id string, in billing.Price) (billing.Price, error) {
	current, err := s.GetPrice(ctx, id)
	if err != nil {
		return billing.Price{}, err
	}
	if in.ProductID != "" {
		current.ProductID = in.ProductID
	}
	if in.Currency != "" {
		current.Currency = strings.ToLower(in.Currency)
	}
	if in.UnitAmount > 0 {
		current.UnitAmount = in.UnitAmount
	}
	if in.LookupKey != "" {
		current.LookupKey = in.LookupKey
	}
	if in.RecurringInterval != "" {
		current.RecurringInterval = in.RecurringInterval
	}
	if in.RecurringIntervalCount > 0 {
		current.RecurringIntervalCount = in.RecurringIntervalCount
	}
	current.Active = in.Active
	if in.Metadata != nil {
		current.Metadata = in.Metadata
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE prices SET product_id = ?, currency = ?, unit_amount = ?, lookup_key = ?, recurring_interval = ?, recurring_interval_count = ?, active = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		current.ProductID, current.Currency, current.UnitAmount, current.LookupKey, current.RecurringInterval, current.RecurringIntervalCount, boolInt(current.Active), encodeMap(current.Metadata), id); err != nil {
		return billing.Price{}, err
	}
	return s.GetPrice(ctx, id)
}

func (s *SQLiteStore) CreateCheckoutSession(ctx context.Context, cs billing.CheckoutSession) (billing.CheckoutSession, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO checkout_sessions (id, customer_id, mode, line_items, success_url, cancel_url, status, payment_status, allow_promotion_codes, trial_period_days, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, cs.ID, cs.CustomerID, cs.Mode, encodeLineItems(cs.LineItems), cs.SuccessURL, cs.CancelURL, cs.Status, cs.PaymentStatus, boolInt(cs.AllowPromotionCodes), cs.TrialPeriodDays, encodeTime(cs.CreatedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if err := s.insertTimeline(ctx, nil, timelineCreate(cs.ID, "checkout_session.created", "Checkout session created", billing.ObjectCheckoutSession, cs.ID, cs.CustomerID, cs.ID, "", "", "", nil, cs.CreatedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	return s.GetCheckoutSession(ctx, cs.ID)
}

func (s *SQLiteStore) GetCheckoutSession(ctx context.Context, id string) (billing.CheckoutSession, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, mode, line_items, success_url, cancel_url, status, payment_status, allow_promotion_codes, trial_period_days, subscription_id, invoice_id, payment_intent_id, created_at, completed_at FROM checkout_sessions WHERE id = ?`, id)
	cs, err := scanCheckoutSession(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.CheckoutSession{}, billing.ErrNotFound
	}
	return cs, err
}

func (s *SQLiteStore) ListCheckoutSessions(ctx context.Context) ([]billing.CheckoutSession, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, mode, line_items, success_url, cancel_url, status, payment_status, allow_promotion_codes, trial_period_days, subscription_id, invoice_id, payment_intent_id, created_at, completed_at FROM checkout_sessions ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.CheckoutSession
	for rows.Next() {
		cs, err := scanCheckoutSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) RecordCheckoutCompletion(ctx context.Context, c billing.CheckoutCompletion) (billing.CheckoutSession, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return billing.CheckoutSession{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO subscriptions (id, customer_id, status, items, current_period_start, current_period_end, cancel_at_period_end, canceled_at, latest_invoice_id, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Subscription.ID, c.Subscription.CustomerID, c.Subscription.Status, encodeLineItems(c.Subscription.Items), encodeTime(c.Subscription.CurrentPeriodStart), encodeTime(c.Subscription.CurrentPeriodEnd), boolInt(c.Subscription.CancelAtPeriodEnd), encodeOptionalTime(c.Subscription.CanceledAt), c.Subscription.LatestInvoiceID, encodeMap(c.Subscription.Metadata)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO invoices (id, customer_id, subscription_id, status, currency, subtotal, total, amount_due, amount_paid, attempt_count, next_payment_attempt, payment_intent_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Invoice.ID, c.Invoice.CustomerID, c.Invoice.SubscriptionID, c.Invoice.Status, c.Invoice.Currency, c.Invoice.Subtotal, c.Invoice.Total, c.Invoice.AmountDue, c.Invoice.AmountPaid, c.Invoice.AttemptCount, encodeOptionalTime(c.Invoice.NextPaymentAttempt), c.Invoice.PaymentIntentID, encodeTime(c.Invoice.CreatedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO payment_intents (id, customer_id, invoice_id, amount, currency, status, failure_code, failure_decline_code, failure_message, payment_method_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.PaymentIntent.ID, c.PaymentIntent.CustomerID, c.PaymentIntent.InvoiceID, c.PaymentIntent.Amount, c.PaymentIntent.Currency, c.PaymentIntent.Status, c.PaymentIntent.FailureCode, c.PaymentIntent.DeclineCode, c.PaymentIntent.FailureMessage, c.PaymentIntent.PaymentMethodID, encodeTime(c.PaymentIntent.CreatedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE checkout_sessions
		SET status = ?, payment_status = ?, subscription_id = ?, invoice_id = ?, payment_intent_id = ?, completed_at = ?
		WHERE id = ?`,
		firstNonEmpty(c.SessionStatus, "complete"), paymentStatus(c), c.Subscription.ID, c.Invoice.ID, c.PaymentIntent.ID, encodeTime(c.CompletedAt), c.SessionID); err != nil {
		return billing.CheckoutSession{}, err
	}
	checkoutEvent := firstNonEmpty(c.CheckoutEvent, "checkout.session.completed")
	if err := s.insertTimeline(ctx, tx, timelineCreate(c.SessionID+"_"+c.Outcome, checkoutEvent, checkoutMessage(checkoutEvent, c.Outcome), billing.ObjectCheckoutSession, c.SessionID, c.Subscription.CustomerID, c.SessionID, c.Subscription.ID, c.Invoice.ID, c.PaymentIntent.ID, map[string]string{"outcome": c.Outcome}, c.CompletedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if err := s.insertTimeline(ctx, tx, timelineCreate(c.Subscription.ID, "customer.subscription.created", "Subscription created from checkout", billing.ObjectSubscription, c.Subscription.ID, c.Subscription.CustomerID, c.SessionID, c.Subscription.ID, c.Invoice.ID, c.PaymentIntent.ID, map[string]string{"status": c.Subscription.Status}, c.CompletedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if err := s.insertTimeline(ctx, tx, timelineCreate(c.Invoice.ID, invoiceEvent(c.Invoice.Status, c.PaymentIntent.Status), "Invoice "+c.Invoice.Status, billing.ObjectInvoice, c.Invoice.ID, c.Invoice.CustomerID, c.SessionID, c.Subscription.ID, c.Invoice.ID, c.PaymentIntent.ID, map[string]string{"status": c.Invoice.Status}, c.CompletedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}
	if err := s.insertTimeline(ctx, tx, timelineCreate(c.PaymentIntent.ID, paymentIntentEvent(c.PaymentIntent.Status), "Payment intent "+c.PaymentIntent.Status, billing.ObjectPaymentIntent, c.PaymentIntent.ID, c.PaymentIntent.CustomerID, c.SessionID, c.Subscription.ID, c.Invoice.ID, c.PaymentIntent.ID, map[string]string{"status": c.PaymentIntent.Status}, c.CompletedAt)); err != nil {
		return billing.CheckoutSession{}, err
	}

	if err := tx.Commit(); err != nil {
		return billing.CheckoutSession{}, err
	}
	return s.GetCheckoutSession(ctx, c.SessionID)
}

func (s *SQLiteStore) GetSubscription(ctx context.Context, id string) (billing.Subscription, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, status, items, current_period_start, current_period_end, cancel_at_period_end, canceled_at, latest_invoice_id, metadata FROM subscriptions WHERE id = ?`, id)
	sub, err := scanSubscription(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.Subscription{}, billing.ErrNotFound
	}
	return sub, err
}

func (s *SQLiteStore) ListSubscriptions(ctx context.Context) ([]billing.Subscription, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, status, items, current_period_start, current_period_end, cancel_at_period_end, canceled_at, latest_invoice_id, metadata FROM subscriptions ORDER BY current_period_start DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListSubscriptionsByCustomer(ctx context.Context, customerID string) ([]billing.Subscription, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, status, items, current_period_start, current_period_end, cancel_at_period_end, canceled_at, latest_invoice_id, metadata
		FROM subscriptions WHERE customer_id = ? ORDER BY current_period_start DESC, id DESC`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Subscription
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateSubscription(ctx context.Context, sub billing.Subscription, timeline []billing.TimelineEntry) (billing.Subscription, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return billing.Subscription{}, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `UPDATE subscriptions
		SET status = ?, items = ?, current_period_start = ?, current_period_end = ?, cancel_at_period_end = ?, canceled_at = ?, latest_invoice_id = ?, metadata = ?
		WHERE id = ?`,
		sub.Status, encodeLineItems(sub.Items), encodeTime(sub.CurrentPeriodStart), encodeTime(sub.CurrentPeriodEnd), boolInt(sub.CancelAtPeriodEnd), encodeOptionalTime(sub.CanceledAt), sub.LatestInvoiceID, encodeMap(sub.Metadata), sub.ID)
	if err != nil {
		return billing.Subscription{}, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return billing.Subscription{}, err
	}
	if changed == 0 {
		return billing.Subscription{}, billing.ErrNotFound
	}
	for _, entry := range timeline {
		if err := s.insertTimeline(ctx, tx, entry); err != nil {
			return billing.Subscription{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return billing.Subscription{}, err
	}
	return s.GetSubscription(ctx, sub.ID)
}

func (s *SQLiteStore) GetInvoice(ctx context.Context, id string) (billing.Invoice, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, subscription_id, status, currency, subtotal, total, amount_due, amount_paid, attempt_count, next_payment_attempt, payment_intent_id, created_at FROM invoices WHERE id = ?`, id)
	inv, err := scanInvoice(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.Invoice{}, billing.ErrNotFound
	}
	return inv, err
}

func (s *SQLiteStore) ListInvoices(ctx context.Context) ([]billing.Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, subscription_id, status, currency, subtotal, total, amount_due, amount_paid, attempt_count, next_payment_attempt, payment_intent_id, created_at FROM invoices ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListInvoicesFiltered(ctx context.Context, filter billing.InvoiceFilter) ([]billing.Invoice, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.CustomerID != "" {
		clauses = append(clauses, "customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.SubscriptionID != "" {
		clauses = append(clauses, "subscription_id = ?")
		args = append(args, filter.SubscriptionID)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, subscription_id, status, currency, subtotal, total, amount_due, amount_paid, attempt_count, next_payment_attempt, payment_intent_id, created_at
		FROM invoices WHERE `+strings.Join(clauses, " AND ")+` ORDER BY created_at DESC, id DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetPaymentIntent(ctx context.Context, id string) (billing.PaymentIntent, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, customer_id, invoice_id, amount, currency, status, failure_code, failure_decline_code, failure_message, payment_method_id, created_at FROM payment_intents WHERE id = ?`, id)
	pi, err := scanPaymentIntent(row)
	if errors.Is(err, sql.ErrNoRows) {
		return billing.PaymentIntent{}, billing.ErrNotFound
	}
	return pi, err
}

func (s *SQLiteStore) ListPaymentIntents(ctx context.Context) ([]billing.PaymentIntent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, invoice_id, amount, currency, status, failure_code, failure_decline_code, failure_message, payment_method_id, created_at FROM payment_intents ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.PaymentIntent
	for rows.Next() {
		pi, err := scanPaymentIntent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pi)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) ListPaymentIntentsFiltered(ctx context.Context, filter billing.PaymentIntentFilter) ([]billing.PaymentIntent, error) {
	if filter.InvoiceIDs != nil && len(filter.InvoiceIDs) == 0 {
		return []billing.PaymentIntent{}, nil
	}
	clauses := []string{"1=1"}
	args := []any{}
	if filter.CustomerID != "" {
		clauses = append(clauses, "customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if len(filter.InvoiceIDs) > 0 {
		placeholders := make([]string, 0, len(filter.InvoiceIDs))
		for _, invoiceID := range filter.InvoiceIDs {
			placeholders = append(placeholders, "?")
			args = append(args, invoiceID)
		}
		clauses = append(clauses, "invoice_id IN ("+strings.Join(placeholders, ",")+")")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, customer_id, invoice_id, amount, currency, status, failure_code, failure_decline_code, failure_message, payment_method_id, created_at
		FROM payment_intents WHERE `+strings.Join(clauses, " AND ")+` ORDER BY created_at DESC, id DESC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.PaymentIntent
	for rows.Next() {
		pi, err := scanPaymentIntent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pi)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) Timeline(ctx context.Context, filter billing.TimelineFilter) ([]billing.TimelineEntry, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.CustomerID != "" {
		clauses = append(clauses, "customer_id = ?")
		args = append(args, filter.CustomerID)
	}
	if filter.CheckoutSessionID != "" {
		clauses = append(clauses, "checkout_session_id = ?")
		args = append(args, filter.CheckoutSessionID)
	}
	if filter.SubscriptionID != "" {
		clauses = append(clauses, "subscription_id = ?")
		args = append(args, filter.SubscriptionID)
	}
	if filter.InvoiceID != "" {
		clauses = append(clauses, "invoice_id = ?")
		args = append(args, filter.InvoiceID)
	}
	if filter.PaymentIntentID != "" {
		clauses = append(clauses, "payment_intent_id = ?")
		args = append(args, filter.PaymentIntentID)
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, action, message, object_type, object_id, customer_id, checkout_session_id, subscription_id, invoice_id, payment_intent_id, data, created_at
		FROM timeline_entries WHERE `+strings.Join(clauses, " AND ")+` ORDER BY created_at ASC, id ASC`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []billing.TimelineEntry
	for rows.Next() {
		entry, err := scanTimelineEntry(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) RecordTimeline(ctx context.Context, entry billing.TimelineEntry) error {
	return s.insertTimeline(ctx, nil, entry)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCustomer(row scanner) (billing.Customer, error) {
	var c billing.Customer
	var metadata, createdAt string
	if err := row.Scan(&c.ID, &c.Email, &c.Name, &metadata, &createdAt); err != nil {
		return c, err
	}
	c.Object = billing.ObjectCustomer
	c.Metadata = decodeMap(metadata)
	c.CreatedAt = decodeTime(createdAt)
	return c, nil
}

func scanProduct(row scanner) (billing.Product, error) {
	var p billing.Product
	var active int
	var metadata, createdAt string
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &active, &metadata, &createdAt); err != nil {
		return p, err
	}
	p.Object = billing.ObjectProduct
	p.Active = active != 0
	p.Metadata = decodeMap(metadata)
	p.CreatedAt = decodeTime(createdAt)
	return p, nil
}

func scanPrice(row scanner) (billing.Price, error) {
	var p billing.Price
	var active int
	var metadata, createdAt string
	if err := row.Scan(&p.ID, &p.ProductID, &p.Currency, &p.UnitAmount, &p.LookupKey, &p.RecurringInterval, &p.RecurringIntervalCount, &active, &metadata, &createdAt); err != nil {
		return p, err
	}
	p.Object = billing.ObjectPrice
	p.Active = active != 0
	p.Metadata = decodeMap(metadata)
	p.CreatedAt = decodeTime(createdAt)
	return p, nil
}

func scanCheckoutSession(row scanner) (billing.CheckoutSession, error) {
	var cs billing.CheckoutSession
	var items, createdAt string
	var allowPromotionCodes int
	var completedAt, subscriptionID, invoiceID, paymentIntentID sql.NullString
	if err := row.Scan(&cs.ID, &cs.CustomerID, &cs.Mode, &items, &cs.SuccessURL, &cs.CancelURL, &cs.Status, &cs.PaymentStatus, &allowPromotionCodes, &cs.TrialPeriodDays, &subscriptionID, &invoiceID, &paymentIntentID, &createdAt, &completedAt); err != nil {
		return cs, err
	}
	cs.Object = billing.ObjectCheckoutSession
	cs.LineItems = decodeLineItems(items)
	cs.URL = "/checkout/" + cs.ID
	cs.AllowPromotionCodes = allowPromotionCodes != 0
	cs.SubscriptionID = subscriptionID.String
	cs.InvoiceID = invoiceID.String
	cs.PaymentIntentID = paymentIntentID.String
	cs.CreatedAt = decodeTime(createdAt)
	if completedAt.Valid {
		t := decodeTime(completedAt.String)
		cs.CompletedAt = &t
	}
	return cs, nil
}

func scanSubscription(row scanner) (billing.Subscription, error) {
	var sub billing.Subscription
	var items, start, end, metadata string
	var cancelAtPeriodEnd int
	var canceledAt sql.NullString
	if err := row.Scan(&sub.ID, &sub.CustomerID, &sub.Status, &items, &start, &end, &cancelAtPeriodEnd, &canceledAt, &sub.LatestInvoiceID, &metadata); err != nil {
		return sub, err
	}
	sub.Object = billing.ObjectSubscription
	sub.Items = decodeLineItems(items)
	sub.CurrentPeriodStart = decodeTime(start)
	sub.CurrentPeriodEnd = decodeTime(end)
	sub.CancelAtPeriodEnd = cancelAtPeriodEnd != 0
	if canceledAt.Valid {
		t := decodeTime(canceledAt.String)
		sub.CanceledAt = &t
	}
	sub.Metadata = decodeMap(metadata)
	return sub, nil
}

func scanInvoice(row scanner) (billing.Invoice, error) {
	var inv billing.Invoice
	var nextPaymentAttempt sql.NullString
	var createdAt string
	if err := row.Scan(&inv.ID, &inv.CustomerID, &inv.SubscriptionID, &inv.Status, &inv.Currency, &inv.Subtotal, &inv.Total, &inv.AmountDue, &inv.AmountPaid, &inv.AttemptCount, &nextPaymentAttempt, &inv.PaymentIntentID, &createdAt); err != nil {
		return inv, err
	}
	inv.Object = billing.ObjectInvoice
	if nextPaymentAttempt.Valid {
		t := decodeTime(nextPaymentAttempt.String)
		inv.NextPaymentAttempt = &t
	}
	inv.CreatedAt = decodeTime(createdAt)
	return inv, nil
}

func scanPaymentIntent(row scanner) (billing.PaymentIntent, error) {
	var pi billing.PaymentIntent
	var createdAt string
	if err := row.Scan(&pi.ID, &pi.CustomerID, &pi.InvoiceID, &pi.Amount, &pi.Currency, &pi.Status, &pi.FailureCode, &pi.DeclineCode, &pi.FailureMessage, &pi.PaymentMethodID, &createdAt); err != nil {
		return pi, err
	}
	pi.Object = billing.ObjectPaymentIntent
	pi.CreatedAt = decodeTime(createdAt)
	return pi, nil
}

func scanTimelineEntry(row scanner) (billing.TimelineEntry, error) {
	var e billing.TimelineEntry
	var data, createdAt string
	if err := row.Scan(&e.ID, &e.Action, &e.Message, &e.ObjectType, &e.ObjectID, &e.CustomerID, &e.CheckoutSessionID, &e.SubscriptionID, &e.InvoiceID, &e.PaymentIntentID, &data, &createdAt); err != nil {
		return e, err
	}
	e.Object = billing.ObjectTimelineEntry
	e.Data = decodeMap(data)
	e.CreatedAt = decodeTime(createdAt)
	return e, nil
}

func timelineCreate(seed, action, message, objectType, objectID, customerID, checkoutSessionID, subscriptionID, invoiceID, paymentIntentID string, data map[string]string, at time.Time) billing.TimelineEntry {
	return billing.TimelineEntry{
		ID:                "tl_" + sanitizeID(seed),
		Object:            billing.ObjectTimelineEntry,
		Action:            action,
		Message:           message,
		ObjectType:        objectType,
		ObjectID:          objectID,
		CustomerID:        customerID,
		CheckoutSessionID: checkoutSessionID,
		SubscriptionID:    subscriptionID,
		InvoiceID:         invoiceID,
		PaymentIntentID:   paymentIntentID,
		Data:              data,
		CreatedAt:         at,
	}
}

func (s *SQLiteStore) insertTimeline(ctx context.Context, tx *sql.Tx, e billing.TimelineEntry) error {
	query := `INSERT INTO timeline_entries (id, action, message, object_type, object_id, customer_id, checkout_session_id, subscription_id, invoice_id, payment_intent_id, data, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	args := []any{e.ID, e.Action, e.Message, e.ObjectType, e.ObjectID, e.CustomerID, e.CheckoutSessionID, e.SubscriptionID, e.InvoiceID, e.PaymentIntentID, encodeMap(e.Data), encodeTime(e.CreatedAt)}
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func encodeMap(m map[string]string) string {
	if m == nil {
		return "{}"
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func decodeMap(raw string) map[string]string {
	if raw == "" {
		return nil
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func encodeLineItems(items []billing.LineItem) string {
	if items == nil {
		return "[]"
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodeLineItems(raw string) []billing.LineItem {
	var out []billing.LineItem
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func encodeTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func encodeOptionalTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return encodeTime(*t)
}

func decodeTime(raw string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func paymentStatus(c billing.CheckoutCompletion) string {
	if c.PaymentStatus != "" {
		return c.PaymentStatus
	}
	if c.Subscription.Status == "trialing" && c.Invoice.Total == 0 {
		return "no_payment_required"
	}
	if c.Outcome == "success" {
		return "paid"
	}
	return "unpaid"
}

func invoiceEvent(status string, paymentIntentStatus string) string {
	if status == "paid" {
		return "invoice.payment_succeeded"
	}
	if status == "void" {
		return "invoice.voided"
	}
	if paymentIntentStatus == "processing" {
		return "invoice.finalized"
	}
	return "invoice.payment_failed"
}

func paymentIntentEvent(status string) string {
	switch status {
	case "succeeded":
		return "payment_intent.succeeded"
	case "processing":
		return "payment_intent.processing"
	case "canceled":
		return "payment_intent.canceled"
	}
	return "payment_intent.payment_failed"
}

func checkoutMessage(eventType string, outcome string) string {
	if eventType == "checkout.session.expired" {
		return "Checkout expired with " + outcome
	}
	return "Checkout completed with " + outcome
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func sanitizeID(raw string) string {
	replacer := strings.NewReplacer("/", "_", ".", "_", " ", "_")
	return replacer.Replace(raw)
}
