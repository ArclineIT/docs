package store

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// planCatalog mirrors billing's PlanCatalog — key → display name.
var planCatalog = map[string]string{
	"shared_starter":   "Shared Starter",
	"shared_pro":       "Shared Pro",
	"shared_business":  "Shared Business",
	"wp_starter":       "WordPress Starter",
	"wp_pro":           "WordPress Pro",
	"wp_business":      "WordPress Business",
	"vps_1":            "VPS Tier 1",
	"vps_2":            "VPS Tier 2",
	"vps_3":            "VPS Tier 3",
	"vps_4":            "VPS Tier 4",
}

// PlanOptions returns plan catalog entries for admin UI dropdowns.
func PlanOptions() []PlanOption {
	order := []string{
		"shared_starter", "shared_pro", "shared_business",
		"wp_starter", "wp_pro", "wp_business",
		"vps_1", "vps_2", "vps_3", "vps_4",
	}
	out := make([]PlanOption, 0, len(order))
	for _, k := range order {
		out = append(out, PlanOption{Key: k, Name: planCatalog[k]})
	}
	return out
}

type PlanOption struct {
	Key  string
	Name string
}

// Customer holds the fields docs needs from billing's customers table.
type Customer struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
}

// Subscription holds the active subscription for a customer (may be nil).
type Subscription struct {
	PriceID  string
	PlanName string
	Status   string
}

// Page is a client or admin-managed doc page stored in docs.db.
type Page struct {
	ID           int64
	Title        string
	Slug         string
	Description  string
	Section      string
	Content      string // raw markdown
	Visibility   string // "public" | "plan:key" | "customer:id"
	DisplayOrder int
	CreatedAt    string
	UpdatedAt    string
}

// Store wraps the billing DB (read-only) and docs DB (read-write).
type Store struct {
	billing *sql.DB
	docs    *sql.DB
}

// New opens both databases and migrates the docs schema.
// If the billing database file does not exist the store starts without it
// (auth and admin features will be unavailable, but public docs still serve).
func New(billingPath, docsPath string) (*Store, error) {
	var billing *sql.DB

	if _, err := os.Stat(billingPath); err == nil {
		billing, err = sql.Open("sqlite", billingPath+"?_journal_mode=WAL&mode=ro")
		if err != nil {
			return nil, fmt.Errorf("open billing db: %w", err)
		}
		if err := billing.Ping(); err != nil {
			billing.Close()
			billing = nil
			slog.Warn("billing db ping failed, auth/admin disabled", "path", billingPath, "err", err)
		}
	} else {
		slog.Warn("billing db not found, auth/admin disabled", "path", billingPath)
	}

	// Ensure the docs directory exists before opening.
	if err := os.MkdirAll(filepath.Dir(docsPath), 0o755); err != nil {
		if billing != nil {
			billing.Close()
		}
		return nil, fmt.Errorf("create docs dir: %w", err)
	}

	docs, err := sql.Open("sqlite", docsPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		if billing != nil {
			billing.Close()
		}
		return nil, fmt.Errorf("open docs db: %w", err)
	}
	if err := docs.Ping(); err != nil {
		if billing != nil {
			billing.Close()
		}
		docs.Close()
		return nil, fmt.Errorf("ping docs db: %w", err)
	}

	if err := migrate(docs); err != nil {
		if billing != nil {
			billing.Close()
		}
		docs.Close()
		return nil, fmt.Errorf("migrate docs db: %w", err)
	}

	return &Store{billing: billing, docs: docs}, nil
}

// Close closes both database connections.
func (s *Store) Close() {
	if s.billing != nil {
		s.billing.Close()
	}
	s.docs.Close()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pages (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		title         TEXT    NOT NULL,
		slug          TEXT    NOT NULL UNIQUE,
		description   TEXT    NOT NULL DEFAULT '',
		section       TEXT    NOT NULL DEFAULT 'client',
		content       TEXT    NOT NULL DEFAULT '',
		visibility    TEXT    NOT NULL DEFAULT 'public',
		display_order INTEGER NOT NULL DEFAULT 0,
		created_at    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
		updated_at    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now'))
	)`)
	return err
}

// ── Auth (reads billing DB) ────────────────────────────────────────────────────

// GetCustomerBySession looks up a billing session token and returns the
// associated customer. Returns nil, nil if the token is missing, expired,
// or billing db is not available.
func (s *Store) GetCustomerBySession(token string) (*Customer, error) {
	if s.billing == nil {
		return nil, nil
	}
	var customerID int64
	var expiresAt string
	err := s.billing.QueryRow(
		`SELECT customer_id, expires_at FROM sessions WHERE token = ?`, token,
	).Scan(&customerID, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query session: %w", err)
	}

	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil || time.Now().UTC().After(exp) {
		return nil, nil
	}

	var c Customer
	err = s.billing.QueryRow(
		`SELECT id, email, first_name, last_name FROM customers WHERE id = ?`, customerID,
	).Scan(&c.ID, &c.Email, &c.FirstName, &c.LastName)
	if err != nil {
		return nil, fmt.Errorf("query customer: %w", err)
	}
	return &c, nil
}

// GetSubscription returns the most recent active/cancelling subscription for
// the given customer, or nil if they have none.
func (s *Store) GetSubscription(customerID int64) (*Subscription, error) {
	if s.billing == nil {
		return nil, nil
	}
	var sub Subscription
	err := s.billing.QueryRow(`
		SELECT stripe_price_id, plan_name, status
		FROM   subscriptions
		WHERE  customer_id = ?
		  AND  status IN ('active','cancelling')
		ORDER  BY created_at DESC
		LIMIT  1`,
		customerID,
	).Scan(&sub.PriceID, &sub.PlanName, &sub.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query subscription: %w", err)
	}
	return &sub, nil
}

// ListCustomers returns all billing customers (for admin visibility picker).
func (s *Store) ListCustomers() ([]Customer, error) {
	rows, err := s.billing.Query(
		`SELECT id, email, first_name, last_name FROM customers ORDER BY email`,
	)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var out []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Email, &c.FirstName, &c.LastName); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ── Visibility ────────────────────────────────────────────────────────────────

// CanSee reports whether a customer (with optional subscription) can view a
// page with the given visibility string.
func CanSee(c *Customer, sub *Subscription, visibility string) bool {
	if visibility == "public" {
		return true
	}
	if c == nil {
		return false
	}
	if strings.HasPrefix(visibility, "customer:") {
		tail := strings.TrimPrefix(visibility, "customer:")
		return fmt.Sprintf("%d", c.ID) == tail
	}
	if strings.HasPrefix(visibility, "plan:") && sub != nil {
		key := strings.TrimPrefix(visibility, "plan:")
		name, ok := planCatalog[key]
		if !ok {
			return false
		}
		return sub.PlanName == name && (sub.Status == "active" || sub.Status == "cancelling")
	}
	return false
}

// VisibilityLabel returns a human-readable label for a visibility string.
func VisibilityLabel(visibility string) string {
	switch {
	case visibility == "public":
		return "Public"
	case strings.HasPrefix(visibility, "plan:"):
		key := strings.TrimPrefix(visibility, "plan:")
		if name, ok := planCatalog[key]; ok {
			return "Plan: " + name
		}
		return "Plan: " + key
	case strings.HasPrefix(visibility, "customer:"):
		return "Customer #" + strings.TrimPrefix(visibility, "customer:")
	}
	return visibility
}

// ── Pages (reads/writes docs DB) ──────────────────────────────────────────────

func (s *Store) ListPages() ([]*Page, error) {
	rows, err := s.docs.Query(`
		SELECT id, title, slug, description, section, content,
		       visibility, display_order, created_at, updated_at
		FROM   pages
		ORDER  BY section, display_order, title`)
	if err != nil {
		return nil, fmt.Errorf("list pages: %w", err)
	}
	defer rows.Close()

	var out []*Page
	for rows.Next() {
		var p Page
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.Section,
			&p.Content, &p.Visibility, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

// ListVisiblePages returns pages the given customer can see, ordered for display.
func (s *Store) ListVisiblePages(c *Customer, sub *Subscription) ([]*Page, error) {
	all, err := s.ListPages()
	if err != nil {
		return nil, err
	}
	var out []*Page
	for _, p := range all {
		if CanSee(c, sub, p.Visibility) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (s *Store) GetPageBySlug(slug string) (*Page, error) {
	var p Page
	err := s.docs.QueryRow(`
		SELECT id, title, slug, description, section, content,
		       visibility, display_order, created_at, updated_at
		FROM   pages WHERE slug = ?`, slug,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.Section,
		&p.Content, &p.Visibility, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get page by slug: %w", err)
	}
	return &p, nil
}

func (s *Store) GetPageByID(id int64) (*Page, error) {
	var p Page
	err := s.docs.QueryRow(`
		SELECT id, title, slug, description, section, content,
		       visibility, display_order, created_at, updated_at
		FROM   pages WHERE id = ?`, id,
	).Scan(&p.ID, &p.Title, &p.Slug, &p.Description, &p.Section,
		&p.Content, &p.Visibility, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get page by id: %w", err)
	}
	return &p, nil
}

func (s *Store) CreatePage(p *Page) (int64, error) {
	res, err := s.docs.Exec(`
		INSERT INTO pages (title, slug, description, section, content, visibility, display_order)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.Title, p.Slug, p.Description, p.Section, p.Content, p.Visibility, p.DisplayOrder,
	)
	if err != nil {
		return 0, fmt.Errorf("create page: %w", err)
	}
	return res.LastInsertId()
}

func (s *Store) UpdatePage(p *Page) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.docs.Exec(`
		UPDATE pages
		SET    title = ?, slug = ?, description = ?, section = ?, content = ?,
		       visibility = ?, display_order = ?, updated_at = ?
		WHERE  id = ?`,
		p.Title, p.Slug, p.Description, p.Section, p.Content,
		p.Visibility, p.DisplayOrder, now, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update page: %w", err)
	}
	return nil
}

func (s *Store) DeletePage(id int64) error {
	_, err := s.docs.Exec(`DELETE FROM pages WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete page: %w", err)
	}
	return nil
}
