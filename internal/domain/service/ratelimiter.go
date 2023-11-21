package service

import (
	"time"

	"golang.org/x/time/rate"
)

// RateLimiterWithLastEventTime представляет структуру, которая объединяет ограничитель скорости и
// время последнего события, связанного с использованием ограничителя.
type RateLimiterWithLastEventTime struct {
	rate      *rate.Limiter // Ограничитель скорости из библиотеки "golang.org/x/time/rate".
	LastEvent time.Time     // Время последнего события, связанного с использованием ограничителя.
}

// NewLimiter создает новый экземпляр RateLimiterWithLastEventTime, используя заданный лимит и предел.
func NewLimiter(r rate.Limit, b int) *RateLimiterWithLastEventTime {
	limiter := rate.NewLimiter(r, b) // Создаем ограничитель скорости.
	return &RateLimiterWithLastEventTime{rate: limiter}
}

// Allow проверяет, разрешено ли действие в соответствии с ограничителем скорости.
// Если разрешено, обновляется время последнего события и возвращается true, иначе возвращается false.
func (t *RateLimiterWithLastEventTime) Allow() bool {
	t.LastEvent = time.Now() // Обновляем время последнего события.
	allow := t.rate.Allow()  // Проверяем, разрешено ли действие согласно ограничителю.

	return allow
}
