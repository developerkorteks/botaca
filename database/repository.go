// Package database - repository untuk operasi database auto promote
package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository interface untuk operasi database
type Repository interface {
	// Auto Promote Groups
	GetAutoPromoteGroup(groupJID string) (*AutoPromoteGroup, error)
	CreateAutoPromoteGroup(groupJID string) (*AutoPromoteGroup, error)
	UpdateAutoPromoteGroup(group *AutoPromoteGroup) error
	GetActiveGroups() ([]AutoPromoteGroup, error)
	
	// Promote Templates
	GetAllTemplates() ([]PromoteTemplate, error)
	GetActiveTemplates() ([]PromoteTemplate, error)
	GetTemplateByID(id int) (*PromoteTemplate, error)
	CreateTemplate(template *PromoteTemplate) error
	UpdateTemplate(template *PromoteTemplate) error
	DeleteTemplate(id int) error
	
	// Promote Logs
	CreateLog(log *PromoteLog) error
	GetLogsByGroup(groupJID string, limit int) ([]PromoteLog, error)
	
	// Stats
	UpdateStats(date string, totalGroups, totalMessages, successMessages, failedMessages int) error
	GetStats(date string) (*PromoteStats, error)
}

// SQLiteRepository implementasi repository untuk SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository membuat repository baru
func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{db: db}
}

// === AUTO PROMOTE GROUPS ===

func (r *SQLiteRepository) GetAutoPromoteGroup(groupJID string) (*AutoPromoteGroup, error) {
	query := `SELECT id, group_jid, is_active, started_at, last_promote_at, created_at, updated_at 
			  FROM auto_promote_groups WHERE group_jid = ?`
	
	row := r.db.QueryRow(query, groupJID)
	
	var group AutoPromoteGroup
	var startedAt, lastPromoteAt sql.NullTime
	
	err := row.Scan(&group.ID, &group.GroupJID, &group.IsActive, 
		&startedAt, &lastPromoteAt, &group.CreatedAt, &group.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Group tidak ditemukan
		}
		return nil, err
	}
	
	if startedAt.Valid {
		group.StartedAt = &startedAt.Time
	}
	if lastPromoteAt.Valid {
		group.LastPromoteAt = &lastPromoteAt.Time
	}
	
	return &group, nil
}

func (r *SQLiteRepository) CreateAutoPromoteGroup(groupJID string) (*AutoPromoteGroup, error) {
	query := `INSERT INTO auto_promote_groups (group_jid, is_active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?)`
	
	now := time.Now()
	result, err := r.db.Exec(query, groupJID, false, now, now)
	if err != nil {
		return nil, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return &AutoPromoteGroup{
		ID:        int(id),
		GroupJID:  groupJID,
		IsActive:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *SQLiteRepository) UpdateAutoPromoteGroup(group *AutoPromoteGroup) error {
	query := `UPDATE auto_promote_groups 
			  SET is_active = ?, started_at = ?, last_promote_at = ?, updated_at = ? 
			  WHERE id = ?`
	
	group.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query, group.IsActive, group.StartedAt, 
		group.LastPromoteAt, group.UpdatedAt, group.ID)
	
	return err
}

func (r *SQLiteRepository) GetActiveGroups() ([]AutoPromoteGroup, error) {
	query := `SELECT id, group_jid, is_active, started_at, last_promote_at, created_at, updated_at 
			  FROM auto_promote_groups WHERE is_active = true`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var groups []AutoPromoteGroup
	
	for rows.Next() {
		var group AutoPromoteGroup
		var startedAt, lastPromoteAt sql.NullTime
		
		err := rows.Scan(&group.ID, &group.GroupJID, &group.IsActive,
			&startedAt, &lastPromoteAt, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, err
		}
		
		if startedAt.Valid {
			group.StartedAt = &startedAt.Time
		}
		if lastPromoteAt.Valid {
			group.LastPromoteAt = &lastPromoteAt.Time
		}
		
		groups = append(groups, group)
	}
	
	return groups, nil
}

// === PROMOTE TEMPLATES ===

func (r *SQLiteRepository) GetAllTemplates() ([]PromoteTemplate, error) {
	query := `SELECT id, title, content, category, is_active, created_at, updated_at 
			  FROM promote_templates ORDER BY created_at DESC`
	
	return r.queryTemplates(query)
}

func (r *SQLiteRepository) GetActiveTemplates() ([]PromoteTemplate, error) {
	query := `SELECT id, title, content, category, is_active, created_at, updated_at 
			  FROM promote_templates WHERE is_active = true ORDER BY created_at DESC`
	
	return r.queryTemplates(query)
}

func (r *SQLiteRepository) queryTemplates(query string, args ...interface{}) ([]PromoteTemplate, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var templates []PromoteTemplate
	
	for rows.Next() {
		var template PromoteTemplate
		err := rows.Scan(&template.ID, &template.Title, &template.Content,
			&template.Category, &template.IsActive, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	
	return templates, nil
}

func (r *SQLiteRepository) GetTemplateByID(id int) (*PromoteTemplate, error) {
	query := `SELECT id, title, content, category, is_active, created_at, updated_at 
			  FROM promote_templates WHERE id = ?`
	
	row := r.db.QueryRow(query, id)
	
	var template PromoteTemplate
	err := row.Scan(&template.ID, &template.Title, &template.Content,
		&template.Category, &template.IsActive, &template.CreatedAt, &template.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &template, nil
}

func (r *SQLiteRepository) CreateTemplate(template *PromoteTemplate) error {
	query := `INSERT INTO promote_templates (title, content, category, is_active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now
	
	result, err := r.db.Exec(query, template.Title, template.Content, 
		template.Category, template.IsActive, template.CreatedAt, template.UpdatedAt)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	template.ID = int(id)
	return nil
}

func (r *SQLiteRepository) UpdateTemplate(template *PromoteTemplate) error {
	query := `UPDATE promote_templates 
			  SET title = ?, content = ?, category = ?, is_active = ?, updated_at = ? 
			  WHERE id = ?`
	
	template.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(query, template.Title, template.Content, 
		template.Category, template.IsActive, template.UpdatedAt, template.ID)
	
	return err
}

func (r *SQLiteRepository) DeleteTemplate(id int) error {
	query := `DELETE FROM promote_templates WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// === PROMOTE LOGS ===

func (r *SQLiteRepository) CreateLog(log *PromoteLog) error {
	query := `INSERT INTO promote_logs (group_jid, template_id, content, sent_at, success, error_msg) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := r.db.Exec(query, log.GroupJID, log.TemplateID, 
		log.Content, log.SentAt, log.Success, log.ErrorMsg)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	log.ID = int(id)
	return nil
}

func (r *SQLiteRepository) GetLogsByGroup(groupJID string, limit int) ([]PromoteLog, error) {
	query := `SELECT id, group_jid, template_id, content, sent_at, success, error_msg 
			  FROM promote_logs WHERE group_jid = ? 
			  ORDER BY sent_at DESC LIMIT ?`
	
	rows, err := r.db.Query(query, groupJID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var logs []PromoteLog
	
	for rows.Next() {
		var log PromoteLog
		var errorMsg sql.NullString
		
		err := rows.Scan(&log.ID, &log.GroupJID, &log.TemplateID,
			&log.Content, &log.SentAt, &log.Success, &errorMsg)
		if err != nil {
			return nil, err
		}
		
		if errorMsg.Valid {
			log.ErrorMsg = &errorMsg.String
		}
		
		logs = append(logs, log)
	}
	
	return logs, nil
}

// === STATS ===

func (r *SQLiteRepository) UpdateStats(date string, totalGroups, totalMessages, successMessages, failedMessages int) error {
	query := `INSERT OR REPLACE INTO promote_stats 
			  (date, total_groups, total_messages, success_messages, failed_messages, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, date, totalGroups, totalMessages, 
		successMessages, failedMessages, time.Now())
	
	return err
}

func (r *SQLiteRepository) GetStats(date string) (*PromoteStats, error) {
	query := `SELECT id, date, total_groups, total_messages, success_messages, failed_messages, created_at 
			  FROM promote_stats WHERE date = ?`
	
	row := r.db.QueryRow(query, date)
	
	var stats PromoteStats
	err := row.Scan(&stats.ID, &stats.Date, &stats.TotalGroups,
		&stats.TotalMessages, &stats.SuccessMessages, &stats.FailedMessages, &stats.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &stats, nil
}

// === UTILITY FUNCTIONS ===

// InitializeDatabase menginisialisasi database dan menjalankan migrasi
func InitializeDatabase(dbPath string) (*sql.DB, Repository, error) {
	// Buka koneksi database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %v", err)
	}
	
	// Test koneksi
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %v", err)
	}
	
	// Jalankan migrasi
	if err := RunMigrations(db); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	
	// Buat repository
	repo := NewSQLiteRepository(db)
	
	return db, repo, nil
}