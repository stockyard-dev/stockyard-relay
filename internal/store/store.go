package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type Channel struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Targets       string `json:"targets"`
	Enabled       bool   `json:"enabled"`
	CreatedAt     string `json:"created_at"`
	DeliveryCount int    `json:"delivery_count"`
	LastWebhook   string `json:"last_webhook,omitempty"`
}
type Delivery struct {
	ID         string `json:"id"`
	ChannelID  string `json:"channel_id"`
	Method     string `json:"method"`
	Headers    string `json:"headers,omitempty"`
	Body       string `json:"body,omitempty"`
	SourceIP   string `json:"source_ip,omitempty"`
	Status     string `json:"status"`
	TargetsHit int    `json:"targets_hit"`
	CreatedAt  string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "relay.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS channels(id TEXT PRIMARY KEY,name TEXT NOT NULL,slug TEXT UNIQUE NOT NULL,targets TEXT DEFAULT '',enabled INTEGER DEFAULT 1,created_at TEXT DEFAULT(datetime('now')),last_webhook TEXT DEFAULT '')`,
		`CREATE TABLE IF NOT EXISTS deliveries(id TEXT PRIMARY KEY,channel_id TEXT NOT NULL,method TEXT DEFAULT 'POST',headers TEXT DEFAULT '',body TEXT DEFAULT '',source_ip TEXT DEFAULT '',status TEXT DEFAULT 'received',targets_hit INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`,
		`CREATE INDEX IF NOT EXISTS idx_del_channel ON deliveries(channel_id)`,
	} {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(resource TEXT NOT NULL,record_id TEXT NOT NULL,data TEXT NOT NULL DEFAULT '{}',PRIMARY KEY(resource, record_id))`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) CreateChannel(c *Channel) error {
	c.ID = genID()
	c.CreatedAt = now()
	en := 1
	if !c.Enabled {
		en = 0
	}
	_, err := d.db.Exec(`INSERT INTO channels(id,name,slug,targets,enabled,created_at)VALUES(?,?,?,?,?,?)`, c.ID, c.Name, c.Slug, c.Targets, en, c.CreatedAt)
	return err
}
func (d *DB) GetChannel(id string) *Channel {
	var c Channel
	var en int
	if d.db.QueryRow(`SELECT id,name,slug,targets,enabled,created_at,last_webhook FROM channels WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Slug, &c.Targets, &en, &c.CreatedAt, &c.LastWebhook) != nil {
		return nil
	}
	c.Enabled = en == 1
	d.db.QueryRow(`SELECT COUNT(*) FROM deliveries WHERE channel_id=?`, c.ID).Scan(&c.DeliveryCount)
	return &c
}
func (d *DB) GetBySlug(slug string) *Channel {
	var c Channel
	var en int
	if d.db.QueryRow(`SELECT id,name,slug,targets,enabled,created_at,last_webhook FROM channels WHERE slug=?`, slug).Scan(&c.ID, &c.Name, &c.Slug, &c.Targets, &en, &c.CreatedAt, &c.LastWebhook) != nil {
		return nil
	}
	c.Enabled = en == 1
	return &c
}
func (d *DB) ListChannels() []Channel {
	rows, _ := d.db.Query(`SELECT id,name,slug,targets,enabled,created_at,last_webhook FROM channels ORDER BY name`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Channel
	for rows.Next() {
		var c Channel
		var en int
		rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Targets, &en, &c.CreatedAt, &c.LastWebhook)
		c.Enabled = en == 1
		d.db.QueryRow(`SELECT COUNT(*) FROM deliveries WHERE channel_id=?`, c.ID).Scan(&c.DeliveryCount)
		o = append(o, c)
	}
	return o
}
func (d *DB) DeleteChannel(id string) error {
	d.db.Exec(`DELETE FROM deliveries WHERE channel_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM channels WHERE id=?`, id)
	return err
}
func (d *DB) RecordDelivery(del *Delivery) error {
	del.ID = genID()
	del.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO deliveries VALUES(?,?,?,?,?,?,?,?,?)`, del.ID, del.ChannelID, del.Method, del.Headers, del.Body, del.SourceIP, del.Status, del.TargetsHit, del.CreatedAt)
	d.db.Exec(`UPDATE channels SET last_webhook=? WHERE id=?`, del.CreatedAt, del.ChannelID)
	return err
}
func (d *DB) ListDeliveries(channelID string, limit int) []Delivery {
	if limit <= 0 {
		limit = 50
	}
	rows, _ := d.db.Query(`SELECT id,channel_id,method,headers,body,source_ip,status,targets_hit,created_at FROM deliveries WHERE channel_id=? ORDER BY created_at DESC LIMIT ?`, channelID, limit)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Delivery
	for rows.Next() {
		var dl Delivery
		rows.Scan(&dl.ID, &dl.ChannelID, &dl.Method, &dl.Headers, &dl.Body, &dl.SourceIP, &dl.Status, &dl.TargetsHit, &dl.CreatedAt)
		o = append(o, dl)
	}
	return o
}

type Stats struct {
	Channels   int `json:"channels"`
	Deliveries int `json:"deliveries"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM channels`).Scan(&s.Channels)
	d.db.QueryRow(`SELECT COUNT(*) FROM deliveries`).Scan(&s.Deliveries)
	return s
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
