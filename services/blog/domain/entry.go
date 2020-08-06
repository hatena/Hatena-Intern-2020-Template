package domain

import (
	"context"
	"strconv"
	"time"
)

// EntryID はブログのエントリにユニークに割り当てられる ID
type EntryID uint64

func (id EntryID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// ParseEntryID は文字列の EntryID をパースする
func ParseEntryID(str string) (EntryID, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return EntryID(0), err
	}
	return EntryID(id), nil
}

// Entry はブログのエントリを表す
type Entry struct {
	ID          EntryID   `db:"id"`
	BlogID      BlogID    `db:"blog_id"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	BodyHTML    string    `db:"body_html"`
	PublishedAt time.Time `db:"published_at"`
	EditedAt    time.Time `db:"edited_at"`
}

// CreateEntryInput はエントリ作成時の入力
type CreateEntryInput struct {
	BlogID      BlogID
	Title       string
	Body        string
	BodyHTML    string
	PublishedAt time.Time
	EditedAt    time.Time
}

// UpdateEntryInput はエントリ更新時の入力
type UpdateEntryInput struct {
	Title       string
	Body        string
	BodyHTML    string
	PublishedAt time.Time
	EditedAt    time.Time
}

// EntryRepository はブログのエントリのリポジトリ
type EntryRepository interface {
	Create(ctx context.Context, input *CreateEntryInput) (*Entry, error)
	FindByID(ctx context.Context, id EntryID) (*Entry, error)
	ListByBlogID(ctx context.Context, blogID BlogID, limit, offset int) ([]*Entry, error)
	Update(ctx context.Context, id EntryID, input *UpdateEntryInput) (*Entry, error)
	Delete(ctx context.Context, id EntryID) error
}

// BodyRenderer はエントリの本文を HTML に変換する
type BodyRenderer interface {
	Render(ctx context.Context, body string) (string, error)
}

// Edit はエントリを編集する
func (e Entry) Edit(title, body string, editedAt time.Time) func(ctx context.Context, r Repository, t BodyRenderer) (*Entry, error) {
	return func(ctx context.Context, r Repository, br BodyRenderer) (*Entry, error) {
		bodyHTML, err := br.Render(ctx, body)
		if err != nil {
			return nil, err
		}
		entry, err := r.Entry().Update(ctx, e.ID, &UpdateEntryInput{
			Title:       title,
			Body:        body,
			BodyHTML:    bodyHTML,
			PublishedAt: e.PublishedAt,
			EditedAt:    editedAt,
		})
		if err != nil {
			return nil, err
		}
		return entry, nil
	}
}

// Unpublish はエントリを削除する
func (e Entry) Unpublish() func(ctx context.Context, r Repository) error {
	return func(ctx context.Context, r Repository) error {
		return r.Entry().Delete(ctx, e.ID)
	}
}
