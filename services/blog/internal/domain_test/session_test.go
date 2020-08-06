package domain

import (
	"testing"
	"time"

	"github.com/hatena/Hatena-Intern-2020/services/blog/domain"
	"github.com/stretchr/testify/assert"
)

func Test_Session_IsExpired(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	sess := &domain.Session{
		ID:        domain.SessionID(0),
		UserID:    domain.UserID(0),
		Key:       "sessionkey",
		ExpiresAt: expiresAt,
	}
	assert.False(t, sess.IsExpired(now))
	assert.False(t, sess.IsExpired(now.Add(24*time.Hour)))
	assert.True(t, sess.IsExpired(now.Add(24*time.Hour+time.Second)))
	assert.True(t, sess.IsExpired(now.Add(48*time.Hour)))
}
