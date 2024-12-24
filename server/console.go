package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/pterodactyl/wings/config"
	"github.com/pterodactyl/wings/system"
)

// appName is a local cache variable to avoid having to make expensive copies of
// the configuration every time we need to send output along to the websocket for
// a server.
var appName string
var appNameSync sync.Once

// PublishConsoleOutputFromDaemon sends output to the server console formatted
// to appear correctly as being sent from Wings.
func (s *Server) PublishConsoleOutputFromDaemon(data string) {
    // ANSI escape-коды для цветов
    const (
        Reset      = "\033[0m"
        White      = "\033[97m"
        Purple     = "\033[95m"       // Основной фиолетовый
        LightPink  = "\033[38;5;219m" // Розовый (для "System")
        LightPurple = "\033[38;5;225m" // Светло-фиолетовый (для data)
    )

    // Форматирование строки с использованием ANSI-кодов
    formattedOutput := fmt.Sprintf(
        "%s[%sCrystallSpace %sSystem%s] %s%s%s%s",
        White, Purple, LightPink, White, LightPurple, data, Reset, White,
    )

    // Публикация события с отформатированным текстом
    s.Events().Publish(
        ConsoleOutputEvent,
        formattedOutput,
    )
}


// Throttler returns the throttler instance for the server or creates a new one.
func (s *Server) Throttler() *ConsoleThrottle {
	s.throttleOnce.Do(func() {
		throttles := config.Get().Throttles
		period := time.Duration(throttles.Period) * time.Millisecond

		s.throttler = newConsoleThrottle(throttles.Lines, period)
		s.throttler.strike = func() {
			s.PublishConsoleOutputFromDaemon("Сервер выводит данные на консоль слишком быстро — ограничение скорости...")
		}
	})
	return s.throttler
}

type ConsoleThrottle struct {
	limit  *system.Rate
	lock   *system.Locker
	strike func()
}

func newConsoleThrottle(lines uint64, period time.Duration) *ConsoleThrottle {
	return &ConsoleThrottle{
		limit: system.NewRate(lines, period),
		lock:  system.NewLocker(),
	}
}

// Allow checks if the console is allowed to process more output data, or if too
// much has already been sent over the line. If there is too much output the
// strike callback function is triggered, but only if it has not already been
// triggered at this point in the process.
//
// If output is allowed, the lock on the throttler is released and the next time
// it is triggered the strike function will be re-executed.
func (ct *ConsoleThrottle) Allow() bool {
	if !ct.limit.Try() {
		if err := ct.lock.Acquire(); err == nil {
			if ct.strike != nil {
				ct.strike()
			}
		}
		return false
	}
	ct.lock.Release()
	return true
}

// Reset resets the console throttler internal rate limiter and overage counter.
func (ct *ConsoleThrottle) Reset() {
	ct.limit.Reset()
}
