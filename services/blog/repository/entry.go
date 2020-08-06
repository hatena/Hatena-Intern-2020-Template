package repository

import (
	"context"
	"database/sql"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/jmoiron/sqlx"
)

// EntryRepository は domain.EntryRepository に対するデータベースを使った実装
type EntryRepository struct {
	db DB
}

func newEntryRepository(db DB) *EntryRepository {
	return &EntryRepository{db}
}

// Create は新規エントリを作成し, リポジトリに保存する
func (r *EntryRepository) Create(ctx context.Context, input *domain.CreateEntryInput) (*domain.Entry, error) {
	id, err := generateID(r.db)
	if err != nil {
		return nil, err
	}
	entry := &domain.Entry{
		ID:          domain.EntryID(id),
		BlogID:      input.BlogID,
		Title:       input.Title,
		Body:        input.Body,
		BodyHTML:    input.BodyHTML,
		PublishedAt: input.PublishedAt,
		EditedAt:    input.EditedAt,
	}
	_, err = r.db.ExecContext(
		ctx,
		`
			INSERT INTO entries (id, blog_id, title, body, body_html, published_at, edited_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		entry.ID, entry.BlogID, entry.Title, entry.Body, entry.BodyHTML, entry.PublishedAt, entry.EditedAt,
	)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// FindByID はリポジトリから ID でエントリを検索する
func (r *EntryRepository) FindByID(ctx context.Context, id domain.EntryID) (*domain.Entry, error) {
	var entry domain.Entry
	err := sqlx.GetContext(
		ctx,
		r.db,
		&entry,
		`
			SELECT id, blog_id, title, body, body_html, published_at, edited_at FROM entries
				WHERE id = ? LIMIT 1
		`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &entry, nil
}

// ListByBlogID はリポジトリからブログの ID でエントリを検索する
func (r *EntryRepository) ListByBlogID(ctx context.Context, blogID domain.BlogID, limit, offset int) ([]*domain.Entry, error) {
	entries := make([]*domain.Entry, 0, limit)
	err := sqlx.SelectContext(
		ctx,
		r.db,
		&entries,
		`
			SELECT id, blog_id, title, body, body_html, published_at, edited_at FROM entries
				WHERE blog_id = ?
				ORDER BY published_at DESC LIMIT ? OFFSET ?
		`,
		blogID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// Update はエントリを更新する
func (r *EntryRepository) Update(ctx context.Context, id domain.EntryID, input *domain.UpdateEntryInput) (*domain.Entry, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			UPDATE entries SET title = ?, body = ?, body_html = ?, published_at = ?, edited_at = ?
				WHERE id = ?
		`,
		input.Title, input.Body, input.BodyHTML, input.PublishedAt, input.EditedAt, id,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

// Delete はエントリをリポジトリから削除する
func (r *EntryRepository) Delete(ctx context.Context, id domain.EntryID) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM entries WHERE id = ?
		`,
		id,
	)
	return err
}
