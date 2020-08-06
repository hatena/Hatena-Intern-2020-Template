package repository

import (
	"context"
	"database/sql"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/jmoiron/sqlx"
)

// BlogRepository は domain.BlogRepository に対するデータベースを使った実装
type BlogRepository struct {
	db DB
}

func newBlogRepository(db DB) *BlogRepository {
	return &BlogRepository{db}
}

// Create は新規ブログを作成し, リポジトリに保存する
func (r *BlogRepository) Create(ctx context.Context, input *domain.CreateBlogInput) (*domain.Blog, error) {
	id, err := generateID(r.db)
	if err != nil {
		return nil, err
	}
	blog := &domain.Blog{
		ID:          domain.BlogID(id),
		UserID:      input.UserID,
		Path:        input.Path,
		Title:       input.Title,
		Description: input.Description,
	}
	_, err = r.db.ExecContext(
		ctx,
		`
			INSERT INTO blogs (id, user_id, path, title, description)
				VALUES (?, ?, ?, ?, ?)
		`,
		blog.ID, blog.UserID, blog.Path, blog.Title, blog.Description,
	)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

// FindByID はリポジトリからブログを ID で検索する
func (r *BlogRepository) FindByID(ctx context.Context, id domain.BlogID) (*domain.Blog, error) {
	var blog domain.Blog
	err := sqlx.GetContext(
		ctx,
		r.db,
		&blog,
		`
			SELECT id, user_id, path, title, description FROM blogs
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
	return &blog, nil
}

// FindByPath はリポジトリからブログをパスで検索する
func (r *BlogRepository) FindByPath(ctx context.Context, path string) (*domain.Blog, error) {
	var blog domain.Blog
	err := sqlx.GetContext(
		ctx,
		r.db,
		&blog,
		`
			SELECT id, user_id, path, title, description FROM blogs
				WHERE path = ? LIMIT 1
		`,
		path,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &blog, nil
}

// List はリポジトリから全ユーザーのブログを検索する
func (r *BlogRepository) List(ctx context.Context, limit, offset int) ([]*domain.Blog, error) {
	blogs := make([]*domain.Blog, 0, limit)
	err := sqlx.SelectContext(
		ctx,
		r.db,
		&blogs,
		`
			SELECT id, user_id, path, title, description FROM blogs
				ORDER BY created_at DESC LIMIT ? OFFSET ?
		`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// ListByUserID はリポジトリからユーザーのブログを検索する
func (r *BlogRepository) ListByUserID(ctx context.Context, userID domain.UserID, limit, offset int) ([]*domain.Blog, error) {
	blogs := make([]*domain.Blog, 0, limit)
	err := sqlx.SelectContext(
		ctx,
		r.db,
		&blogs,
		`
			SELECT id, user_id, path, title, description FROM blogs
				WHERE user_id = ?
				ORDER BY title DESC LIMIT ? OFFSET ?
		`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// Update はブログを更新する
func (r *BlogRepository) Update(ctx context.Context, id domain.BlogID, input *domain.UpdateBlogInput) (*domain.Blog, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			UPDATE blogs SET title = ?, description = ?
				WHERE id = ?
		`,
		input.Title, input.Description, id,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

// Delete はブログを削除する
func (r *BlogRepository) Delete(ctx context.Context, id domain.BlogID) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM blogs WHERE id = ?
		`,
		id,
	)
	return err
}
