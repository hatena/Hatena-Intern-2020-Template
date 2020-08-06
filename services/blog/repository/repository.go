package repository

import (
	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/jmoiron/sqlx"
)

// DB はデータベースのインターフェース
type DB interface {
	sqlx.Execer
	sqlx.ExecerContext
	sqlx.Queryer
	sqlx.QueryerContext
}

// Repository は domain.Repository に対するデータベースを使った実装
type Repository struct {
	user    *UserRepository
	session *SessionRepository
	blog    *BlogRepository
	entry   *EntryRepository
}

// NewRepository は Repository を作成する
func NewRepository(db DB) *Repository {
	return &Repository{
		user:    newUserRepository(db),
		session: newSessionRepository(db),
		blog:    newBlogRepository(db),
		entry:   newEntryRepository(db),
	}
}

// User はユーザーに対するリポジトリを返す
func (r *Repository) User() domain.UserRepository {
	return r.user
}

// Session はセッションに対するリポジトリを返す
func (r *Repository) Session() domain.SessionRepository {
	return r.session
}

// Blog はブログに対するリポジトリを返す
func (r *Repository) Blog() domain.BlogRepository {
	return r.blog
}

// Entry はエントリに対するリポジトリを返す
func (r *Repository) Entry() domain.EntryRepository {
	return r.entry
}

func generateID(db DB) (uint64, error) {
	var id uint64
	err := sqlx.Get(db, &id, "SELECT UUID_SHORT()")
	return id, err
}
