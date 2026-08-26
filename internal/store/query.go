package store

import (
	"database/sql"
	"errors"
	"time"
)

// Count 返回某类记录数量。
func (s *Store) Count(kind string) (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM records WHERE kind=?", kind).Scan(&n)
	return n, err
}

// Exists 判断某个记录是否存在。
func (s *Store) Exists(kind, id string) (bool, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM records WHERE kind=? AND id=?", kind, id).Scan(&n)
	return n > 0, err
}

// UpdatedAt 返回记录的更新时间，供报告服务排序。
func (s *Store) UpdatedAt(kind, id string) (time.Time, error) {
	var raw string
	if err := s.db.QueryRow("SELECT updated_at FROM records WHERE kind=? AND id=?", kind, id).Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, sql.ErrNoRows
		}
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339Nano, raw)
}

// Ping 检查 SQLite 连接是否仍可用。
func (s *Store) Ping() error { return s.db.Ping() }
