package initializer

import (
	"context"
	"fmt"
	"time"
)

type (
	InitFunc = func(string) (interface{}, error)
	CallFunc = func() (interface{}, error)
)

type Exponential struct {
	timeOut, interval int
}

func NewExponential(timeOut, interval int) Exponential {
	return Exponential{
		timeOut, interval,
	}
}

func (e Exponential) Run(ctx context.Context, initFunc InitFunc, source string) (interface{}, error) {
	isConnDB := func() (interface{}, error) {
		conn, err := initFunc(source)
		if err != nil {
			return nil, fmt.Errorf("unable to init source: %s", err)
		}
		return conn, nil
	}

	conn, ok := RetryIncrementallyUntil(ctx,
		time.Duration(e.interval)*time.Second,
		time.Duration(e.timeOut)*time.Second,
		isConnDB,
	)

	if !ok {
		return nil, fmt.Errorf("can't init source: %s", source)
	}
	return conn, nil
}

// RetryIncrementallyUntil retries passed function call until it end with
// success (returns true), else repeat and increments retrying interval
// up to 2 hour until function call succeeded (returns true) or until until
// time reached (returns false).
func RetryIncrementallyUntil(ctx context.Context, interval time.Duration, until time.Duration, call CallFunc) (interface{}, bool) {
	var counter int64
	untilTimer := time.NewTimer(until)

	counter++
	conn, err := call()
	if err == nil {
		return conn, true
	}
	intervalTimer := time.NewTimer(incrementInterval(interval, counter))

	for {
		select {
		case <-untilTimer.C:
			return nil, false
		case <-intervalTimer.C:
			counter++
			conn, err = call()
			if err == nil {
				return conn, true
			}
			intervalTimer.Reset(incrementInterval(interval, counter))
		case <-ctx.Done():
			return nil, false
		}
	}
}

func incrementInterval(interval time.Duration, counter int64) time.Duration {
	quartile := int64(interval) / 4
	addendum := quartile * counter
	result := interval + time.Duration(addendum)

	if result > 2*time.Hour {
		result = 2 * time.Hour
	}
	return result
}
