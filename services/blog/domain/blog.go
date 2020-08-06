package domain

import (
	"context"
	"time"
)

// BlogID はブログにユニークに割り当てられる ID
type BlogID uint64

// Blog はブログを表す
type Blog struct {
	ID          BlogID `db:"id"`
	UserID      UserID `db:"user_id"`
	Path        string `db:"path"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

// CreateBlogInput はブログ作成時の入力
type CreateBlogInput struct {
	UserID      UserID
	Path        string
	Title       string
	Description string
}

// UpdateBlogInput はブログ更新時の入力
type UpdateBlogInput struct {
	Title       string
	Description string
}

// BlogRepository はブログのリポジトリ
type BlogRepository interface {
	Create(ctx context.Context, input *CreateBlogInput) (*Blog, error)
	FindByID(ctx context.Context, id BlogID) (*Blog, error)
	FindByPath(ctx context.Context, path string) (*Blog, error)
	List(ctx context.Context, limit, offset int) ([]*Blog, error)
	ListByUserID(ctx context.Context, userID UserID, limit, offset int) ([]*Blog, error)
	Update(ctx context.Context, id BlogID, input *UpdateBlogInput) (*Blog, error)
	Delete(ctx context.Context, id BlogID) error
}

// Edit はブログのタイトルや説明文を更新する
func (b Blog) Edit(title, description string) func(ctx context.Context, r Repository) (*Blog, error) {
	return func(ctx context.Context, r Repository) (*Blog, error) {
		blog, err := r.Blog().Update(ctx, b.ID, &UpdateBlogInput{
			Title:       title,
			Description: description,
		})
		if err != nil {
			return nil, err
		}
		return blog, nil
	}
}

// Delete はブログを削除する
func (b Blog) Delete() func(ctx context.Context, r Repository) error {
	return func(ctx context.Context, r Repository) error {
		return r.Blog().Delete(ctx, b.ID)
	}
}

// PublishEntry は新規エントリを公開する
func (b Blog) PublishEntry(title, body string, publishedAt time.Time) func(ctx context.Context, r Repository, br BodyRenderer) (*Entry, error) {
	return func(ctx context.Context, r Repository, br BodyRenderer) (*Entry, error) {
		bodyHTML, err := br.Render(ctx, body)
		if err != nil {
			return nil, err
		}
		entry, err := r.Entry().Create(ctx, &CreateEntryInput{
			BlogID:      b.ID,
			Title:       title,
			Body:        body,
			BodyHTML:    bodyHTML,
			PublishedAt: publishedAt,
			EditedAt:    publishedAt,
		})
		if err != nil {
			return nil, err
		}
		return entry, nil
	}
}
