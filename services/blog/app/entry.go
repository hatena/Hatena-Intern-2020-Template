package app

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/repository"
)

// ListEntriesByBlog はブログのエントリを検索する
func (a *App) ListEntriesByBlog(ctx context.Context, blog *domain.Blog, page, limit int) ([]*domain.Entry, bool, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	repo := repository.NewRepository(a.db)
	entries, err := repo.Entry().ListByBlogID(ctx, blog.ID, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	if len(entries) > limit {
		return entries[:limit], true, nil
	}
	return entries, false, nil
}

// PublishEntry は新規エントリを投稿する
func (a *App) PublishEntry(ctx context.Context, user *domain.User, blog *domain.Blog, title, body string) (*domain.Entry, error) {
	if blog.UserID != user.ID {
		return nil, ErrPermissionDenied
	}
	if utf8.RuneCountInString(title) > 500 {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	return blog.PublishEntry(title, body, time.Now())(ctx, repo, a)
}

// FindEntryByID は ID でエントリを検索する
func (a *App) FindEntryByID(ctx context.Context, blog *domain.Blog, entryID domain.EntryID) (*domain.Entry, error) {
	repo := repository.NewRepository(a.db)
	entry, err := repo.Entry().FindByID(ctx, entryID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if entry.BlogID != blog.ID {
		return nil, ErrNotFound
	}
	return entry, nil
}

// EditEntry はエントリを編集する
func (a *App) EditEntry(ctx context.Context, user *domain.User, blog *domain.Blog, entry *domain.Entry, title, body string) (*domain.Entry, error) {
	if blog.UserID != user.ID {
		return nil, ErrPermissionDenied
	}
	if entry.BlogID != blog.ID {
		return nil, ErrNotFound
	}
	if utf8.RuneCountInString(title) > 500 {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	return entry.Edit(title, body, time.Now())(ctx, repo, a)
}

// UnpublishEntry はエントリを削除する
func (a *App) UnpublishEntry(ctx context.Context, user *domain.User, blog *domain.Blog, entry *domain.Entry) error {
	if blog.UserID != user.ID {
		return ErrPermissionDenied
	}
	if entry.BlogID != blog.ID {
		return ErrNotFound
	}
	repo := repository.NewRepository(a.db)
	return entry.Unpublish()(ctx, repo)
}
