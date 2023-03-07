package backoff

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

func EqualJitter(b Backoff) Backoff {
	return NextDelayFunc(func(i int) (int64, error) {
		dur, err := b.NextDelay(i)
		if err != nil {
			return 0, err
		}

		dur = dur / 2

		if dur <= 0 {
			return dur, nil
		}

		jitter, err := rand.Int(rand.Reader, big.NewInt(int64(dur+1)))
		if err != nil {
			panic(err)
		}

		return dur + jitter.Int64(), nil
	})
}

func FullJitter(b Backoff) Backoff {
	return NextDelayFunc(func(i int) (int64, error) {
		dur, err := b.NextDelay(i)
		if err != nil {
			return 0, err
		}

		if dur <= 0 {
			return dur, nil
		}

		jitter, err := rand.Int(rand.Reader, big.NewInt(int64(dur+1)))
		if err != nil {
			panic(err)
		}

		return jitter.Int64(), nil
	})
}

var ErrMaxAttempts = errors.New("max attempts reached")

func MaxAttemps(attempts int, b Backoff) Backoff {
	return NextDelayFunc(func(i int) (int64, error) {
		if i >= attempts {
			return 0, ErrMaxAttempts
		}

		return b.NextDelay(i)
	})
}

func Chain(backoff Backoff, items ...func(Backoff) Backoff) Backoff {
	var wrappedbackoff Backoff

	for _, item := range items {
		wrappedbackoff = item(backoff)
	}

	return wrappedbackoff
}

func Initialdelay(dur time.Duration, b Backoff) Backoff {
	return NextDelayFunc(func(i int) (int64, error) {
		nextDelay, err := b.NextDelay(i)
		if err != nil {
			return 0, err
		}

		return int64(dur) * nextDelay, nil
	})
}
