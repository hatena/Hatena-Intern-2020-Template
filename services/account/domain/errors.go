package domain

import "errors"

// ErrNotFound はリポジトリにエンティティが見つからなかったときに返される
var ErrNotFound = errors.New("not found")

// ErrAlreadyExists はリポジトリに既に重複するエンティティが存在したときに返される
var ErrAlreadyExists = errors.New("already exists")
