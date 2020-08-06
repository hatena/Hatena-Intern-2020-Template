package domain

// Repository は各エンティティのリポジトリを束ねる
type Repository interface {
	User() UserRepository
}
