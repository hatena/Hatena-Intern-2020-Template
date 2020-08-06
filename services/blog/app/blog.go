package app

import (
	"context"
	"regexp"
	"unicode/utf8"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/hatena/Hatena-Intern-2020/services/blog/repository"
)

// ListBlogs は全ユーザーのブログを検索する
func (a *App) ListBlogs(ctx context.Context, page, limit int) ([]*domain.Blog, bool, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	repo := repository.NewRepository(a.db)
	blogs, err := repo.Blog().List(ctx, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	if len(blogs) > limit {
		return blogs[:limit], true, nil
	}
	return blogs, false, nil
}

// ListBlogsByUser はユーザーのブログを検索する
func (a *App) ListBlogsByUser(ctx context.Context, user *domain.User, page, limit int) ([]*domain.Blog, bool, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	repo := repository.NewRepository(a.db)
	blogs, err := repo.Blog().ListByUserID(ctx, user.ID, limit+1, offset)
	if err != nil {
		return nil, false, err
	}
	if len(blogs) > limit {
		return blogs[:limit], true, nil
	}
	return blogs, false, nil
}

var pathRE = regexp.MustCompile(`^[0-9A-Za-z][-_0-9A-Za-z]{2,63}$`)

// CreateBlog は新規ブログを作成する
func (a *App) CreateBlog(ctx context.Context, user *domain.User, path, title, description string) (*domain.Blog, error) {
	if ok := pathRE.MatchString(path); !ok {
		return nil, ErrInvalidArgument
	}
	if utf8.RuneCountInString(title) > 200 {
		return nil, ErrInvalidArgument
	}
	if utf8.RuneCountInString(description) > 500 {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	return user.CreateBlog(path, title, description)(ctx, repo)
}

// FindBlogByPath はパスでブログを検索する
func (a *App) FindBlogByPath(ctx context.Context, path string) (*domain.Blog, error) {
	repo := repository.NewRepository(a.db)
	blog, err := repo.Blog().FindByPath(ctx, path)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return blog, nil
}

// EditBlog はブログの情報を更新する
func (a *App) EditBlog(ctx context.Context, user *domain.User, blog *domain.Blog, title, description string) (*domain.Blog, error) {
	if blog.UserID != user.ID {
		return nil, ErrPermissionDenied
	}
	if utf8.RuneCountInString(title) > 200 {
		return nil, ErrInvalidArgument
	}
	if utf8.RuneCountInString(description) > 500 {
		return nil, ErrInvalidArgument
	}
	repo := repository.NewRepository(a.db)
	return blog.Edit(title, description)(ctx, repo)
}

// DeleteBlog はブログを削除する
func (a *App) DeleteBlog(ctx context.Context, user *domain.User, blog *domain.Blog) error {
	if blog.UserID != user.ID {
		return ErrPermissionDenied
	}
	repo := repository.NewRepository(a.db)
	return blog.Delete()(ctx, repo)
}
