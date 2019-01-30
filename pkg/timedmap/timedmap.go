package timedmap

import (
	"errors"
	"time"
)

// TimedMap contains a map with all key-value pairs,
// and a timer, which cleans the map in the set
// tick durations from expired keys.
type TimedMap struct {
	cleanupTickTime time.Duration
	container       map[interface{}]*element
	cleaner         *time.Ticker
	cleanerStopChan chan bool
}

type element struct {
	value   interface{}
	expires time.Time
}

// New creates and returns a new instance of TimedMap.
// The passed cleanupTickTime will be passed to the
// cleanup Timer, which iterates through the map and
// deletes expired key-value pairs.
func New(cleanupTickTime time.Duration) *TimedMap {
	tm := &TimedMap{
		container: make(map[interface{}]*element),
	}

	tm.cleaner = time.NewTicker(cleanupTickTime)

	go func() {
		for {
			select {
			case <-tm.cleaner.C:
				tm.cleanUp()
			case <-tm.cleanerStopChan:
				break
			}
		}
	}()

	return tm
}

func (tm *TimedMap) cleanUp() {
	for k, v := range tm.container {
		if time.Now().After(v.expires) {
			delete(tm.container, k)
		}
	}
}

func (tm *TimedMap) get(key interface{}) *element {
	v, ok := tm.container[key]
	if !ok {
		return nil
	}

	if time.Now().After(v.expires) {
		delete(tm.container, key)
		return nil
	}

	return v
}

// Set appends a key-value pair to the mao ir sets the value of
// a key. expiresAfter sets the expire time after the key-value pair
// will automatically be removed from the map.
func (tm *TimedMap) Set(key, value interface{}, expiresAfter time.Duration) {
	tm.container[key] = &element{
		value:   value,
		expires: time.Now().Add(expiresAfter),
	}
}

// GetValue returns an interface of the value of a key in the
// map. The returned value is nil if there is no value to the
// passed key or if the value was expired.
func (tm *TimedMap) GetValue(key interface{}) interface{} {
	v := tm.get(key)
	if v == nil {
		return nil
	}
	return v.value
}

// GetExpires returns the expire time of a key-value pair.
// If the key-value pair does not exist in the map or
// was expired, this will return an error object.
func (tm *TimedMap) GetExpires(key interface{}) (time.Time, error) {
	v := tm.get(key)
	if v == nil {
		return time.Time{}, errors.New("key not found")
	}
	return v.expires, nil
}

// Contains returns true, if the key exists in the map.
// false will be returned, if there is no value to the
// key or if the key-value pair was expired.
func (tm *TimedMap) Contains(key interface{}) bool {
	return tm.get(key) != nil
}

// Remove deletes a key-value pair in the map.
func (tm *TimedMap) Remove(key interface{}) {
	delete(tm.container, key)
}

// Refresh extends the expire time for a key-value pair
// about the passed duration. If there is no value to
// the key passed, this will return an error object.
func (tm *TimedMap) Refresh(key interface{}, d time.Duration) error {
	v := tm.get(key)
	if v == nil {
		return errors.New("key not found")
	}
	v.expires = v.expires.Add(d)
	return nil
}

// Size returns the current number of key-value pairs
// existent in the map.
func (tm *TimedMap) Size() int {
	return len(tm.container)
}

// StopCleaner stops the cleaner go routine and timer.
// This should always be called after exiting a scope
// where TimedMap is used that the data can be cleaned
// up correctly.
func (tm *TimedMap) StopCleaner() {
	go func() {
		tm.cleanerStopChan <- true
	}()
	tm.cleaner.Stop()
}
