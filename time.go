// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package time provides functionality for measuring and displaying time.
//
// The calendrical calculations always assume a Gregorian calendar, with
// no leap seconds.
//
// Monotonic Clocks
//
// Operating systems provide both a “wall clock,” which is subject to
// changes for clock synchronization, and a “monotonic clock,” which is
// not. The general rule is that the wall clock is for telling time and
// the monotonic clock is for measuring time. Rather than split the API,
// in this package the Time returned by time.Now contains both a wall
// clock reading and a monotonic clock reading; later time-telling
// operations use the wall clock reading, but later time-measuring
// operations, specifically comparisons and subtractions, use the
// monotonic clock reading.
//
// For example, this code always computes a positive elapsed time of
// approximately 20 milliseconds, even if the wall clock is changed during
// the operation being timed:
//
//	start := time.Now()
//	... operation that takes 20 milliseconds ...
//	t := time.Now()
//	elapsed := t.Sub(start)
//
// Other idioms, such as time.Since(start), time.Until(deadline), and
// time.Now().Before(deadline), are similarly robust against wall clock
// resets.
//
// The rest of this section gives the precise details of how operations
// use monotonic clocks, but understanding those details is not required
// to use this package.
//
// The Time returned by time.Now contains a monotonic clock reading.
// If Time t has a monotonic clock reading, t.Add adds the same duration to
// both the wall clock and monotonic clock readings to compute the result.
// Because t.AddDate(y, m, d), t.Round(d), and t.Truncate(d) are wall time
// computations, they always strip any monotonic clock reading from their results.
// Because t.In, t.Local, and t.UTC are used for their effect on the interpretation
// of the wall time, they also strip any monotonic clock reading from their results.
// The canonical way to strip a monotonic clock reading is to use t = t.Round(0).
//
// If Times t and u both contain monotonic clock readings, the operations
// t.After(u), t.Before(u), t.Equal(u), and t.Sub(u) are carried out
// using the monotonic clock readings alone, ignoring the wall clock
// readings. If either t or u contains no monotonic clock reading, these
// operations fall back to using the wall clock readings.
//
// On some systems the monotonic clock will stop if the computer goes to sleep.
// On such a system, t.Sub(u) may not accurately reflect the actual
// time that passed between t and u.
//
// Because the monotonic clock reading has no meaning outside
// the current process, the serialized forms generated by t.GobEncode,
// t.MarshalBinary, t.MarshalJSON, and t.MarshalText omit the monotonic
// clock reading, and t.Format provides no format for it. Similarly, the
// constructors time.Date, time.Parse, time.ParseInLocation, and time.Unix,
// as well as the unmarshalers t.GobDecode, t.UnmarshalBinary.
// t.UnmarshalJSON, and t.UnmarshalText always create times with
// no monotonic clock reading.
//
// Note that the Go == operator compares not just the time instant but
// also the Location and the monotonic clock reading. See the
// documentation for the Time type for a discussion of equality
// testing for Time values.
//
// For debugging, the result of t.String does include the monotonic
// clock reading if present. If t != u because of different monotonic clock readings,
// that difference will be visible when printing t.String() and u.String().
//
package time32

import (
	_ "unsafe" // for go:linkname
)

// Provided by platform specific assembly code

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

//go:noescape
//go:linkname walltime runtime.walltime
func walltime() (sec int64, nsec int32)

// runtimeNano returns the current value of the runtime clock in nanoseconds.
var runtimeNano = nanotime
var now = time_now

// A Time represents an instant in time with nanosecond precision.
//
// Programs using times should typically store and pass them as values,
// not pointers. That is, time variables and struct fields should be of
// type time.Time, not *time.Time.
//
// A Time value can be used by multiple goroutines simultaneously except
// that the methods GobDecode, UnmarshalBinary, UnmarshalJSON and
// UnmarshalText are not concurrency-safe.
//
// Time instants can be compared using the Before, After, and Equal methods.
// The Sub method subtracts two instants, producing a Duration.
// The Add method adds a Time and a Duration, producing a Time.
//
// The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC.
// As this time is unlikely to come up in practice, the IsZero method gives
// a simple way of detecting a time that has not been initialized explicitly.
//
// Each Time has associated with it a Location, consulted when computing the
// presentation form of the time, such as in the Format, Hour, and Year methods.
// The methods Local, UTC, and In return a Time with a specific location.
// Changing the location in this way changes only the presentation; it does not
// change the instant in time being denoted and therefore does not affect the
// computations described in earlier paragraphs.
//
// Representations of a Time value saved by the GobEncode, MarshalBinary,
// MarshalJSON, and MarshalText methods store the Time.Location's offset, but not
// the location name. They therefore lose information about Daylight Saving Time.
//
// In addition to the required “wall clock” reading, a Time may contain an optional
// reading of the current process's monotonic clock, to provide additional precision
// for comparison or subtraction.
// See the “Monotonic Clocks” section in the package documentation for details.
//
// Note that the Go == operator compares not just the time instant but also the
// Location and the monotonic clock reading. Therefore, Time values should not
// be used as map or database keys without first guaranteeing that the
// identical Location has been set for all values, which can be achieved
// through use of the UTC or Local method, and that the monotonic clock reading
// has been stripped by setting t = t.Round(0). In general, prefer t.Equal(u)
// to t == u, since t.Equal uses the most accurate comparison available and
// correctly handles the case when only one of its arguments has a monotonic
// clock reading.
//
type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64
}

/*
Define our own time unit which will always hold epoch time
in millis. Example: 1588228661
4294967295
int64 size: 8 bytes
uint32 size: 4 bytes

x2 size reduction changing data size
*/
type Time32 uint32

func (t Time32) AddDate(days int) Time32 {
	v := int(t) + (days * 86400)
	return Time32(v)
}

const (
	hasMonotonic = 1 << 63
	maxWall      = wallToInternal + (1<<33 - 1) // year 2157
	minWall      = wallToInternal               // year 1885
	nsecMask     = 1<<30 - 1
	nsecShift    = 30
)

// These helpers for manipulating the wall and monotonic clock readings
// take pointer receivers, even when they don't modify the time,
// to make them cheaper to call.

// nsec returns the time's nanoseconds.
func (t *Time) nsec() int32 {
	return int32(t.wall & nsecMask)
}

// sec returns the time's seconds since Jan 1 year 1.
func (t *Time) sec() int64 {
	if t.wall&hasMonotonic != 0 {
		return wallToInternal + int64(t.wall<<1>>(nsecShift+1))
	}
	return t.ext
}

// unixSec returns the time's seconds since Jan 1 1970 (Unix time).
func (t *Time) unixSec() int64 { return t.sec() + internalToUnix }

// addSec adds d seconds to the time.
func (t *Time) addSec(d int64) {
	if t.wall&hasMonotonic != 0 {
		sec := int64(t.wall << 1 >> (nsecShift + 1))
		dsec := sec + d
		if 0 <= dsec && dsec <= 1<<33-1 {
			t.wall = t.wall&nsecMask | uint64(dsec)<<nsecShift | hasMonotonic
			return
		}
		// Wall second now out of range for packed field.
		// Move to ext.
		t.stripMono()
	}

	// TODO: Check for overflow.
	t.ext += d
}

// stripMono strips the monotonic clock reading in t.
func (t *Time) stripMono() {
	if t.wall&hasMonotonic != 0 {
		t.ext = t.sec()
		t.wall &= nsecMask
	}
}

// setMono sets the monotonic clock reading in t.
// If t cannot hold a monotonic clock reading,
// because its wall time is too large,
// setMono is a no-op.
func (t *Time) setMono(m int64) {
	if t.wall&hasMonotonic == 0 {
		sec := t.ext
		if sec < minWall || maxWall < sec {
			return
		}
		t.wall |= hasMonotonic | uint64(sec-minWall)<<nsecShift
	}
	t.ext = m
}

// mono returns t's monotonic clock reading.
// It returns 0 for a missing reading.
// This function is used only for testing,
// so it's OK that technically 0 is a valid
// monotonic clock reading as well.
func (t *Time) mono() int64 {
	if t.wall&hasMonotonic == 0 {
		return 0
	}
	return t.ext
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	if t.wall&u.wall&hasMonotonic != 0 {
		return t.ext > u.ext
	}
	ts := t.sec()
	us := u.sec()
	return ts > us || ts == us && t.nsec() > u.nsec()
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool {
	if t.wall&u.wall&hasMonotonic != 0 {
		return t.ext < u.ext
	}
	return t.sec() < u.sec() || t.sec() == u.sec() && t.nsec() < u.nsec()
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t Time) Equal(u Time) bool {
	if t.wall&u.wall&hasMonotonic != 0 {
		return t.ext == u.ext
	}
	return t.sec() == u.sec() && t.nsec() == u.nsec()
}

// Computations on time.
//
// The zero value for a Time is defined to be
//	January 1, year 1, 00:00:00.000000000 UTC
// which (1) looks like a zero, or as close as you can get in a date
// (1-1-1 00:00:00 UTC), (2) is unlikely enough to arise in practice to
// be a suitable "not set" sentinel, unlike Jan 1 1970, and (3) has a
// non-negative year even in time zones west of UTC, unlike 1-1-0
// 00:00:00 UTC, which would be 12-31-(-1) 19:00:00 in New York.
//
// The zero Time value does not force a specific epoch for the time
// representation. For example, to use the Unix epoch internally, we
// could define that to distinguish a zero value from Jan 1 1970, that
// time would be represented by sec=-1, nsec=1e9. However, it does
// suggest a representation, namely using 1-1-1 00:00:00 UTC as the
// epoch, and that's what we do.
//
// The Add and Sub computations are oblivious to the choice of epoch.
//
// The presentation computations - year, month, minute, and so on - all
// rely heavily on division and modulus by positive constants. For
// calendrical calculations we want these divisions to round down, even
// for negative values, so that the remainder is always positive, but
// Go's division (like most hardware division instructions) rounds to
// zero. We can still do those computations and then adjust the result
// for a negative numerator, but it's annoying to write the adjustment
// over and over. Instead, we can change to a different epoch so long
// ago that all the times we care about will be positive, and then round
// to zero and round down coincide. These presentation routines already
// have to add the zone offset, so adding the translation to the
// alternate epoch is cheap. For example, having a non-negative time t
// means that we can write
//
//	sec = t % 60
//
// instead of
//
//	sec = t % 60
//	if sec < 0 {
//		sec += 60
//	}
//
// everywhere.
//
// The calendar runs on an exact 400 year cycle: a 400-year calendar
// printed for 1970-2369 will apply as well to 2370-2769. Even the days
// of the week match up. It simplifies the computations to choose the
// cycle boundaries so that the exceptional years are always delayed as
// long as possible. That means choosing a year equal to 1 mod 400, so
// that the first leap year is the 4th year, the first missed leap year
// is the 100th year, and the missed missed leap year is the 400th year.
// So we'd prefer instead to print a calendar for 2001-2400 and reuse it
// for 2401-2800.
//
// Finally, it's convenient if the delta between the Unix epoch and
// long-ago epoch is representable by an int64 constant.
//
// These three considerations—choose an epoch as early as possible, that
// uses a year equal to 1 mod 400, and that is no more than 2⁶³ seconds
// earlier than 1970—bring us to the year -292277022399. We refer to
// this year as the absolute zero year, and to times measured as a uint64
// seconds since this year as absolute times.
//
// Times measured as an int64 seconds since the year 1—the representation
// used for Time's sec field—are called internal times.
//
// Times measured as an int64 seconds since the year 1970 are called Unix
// times.
//
// It is tempting to just use the year 1 as the absolute epoch, defining
// that the routines are only valid for years >= 1. However, the
// routines would then be invalid when displaying the epoch in time zones
// west of UTC, since it is year 0. It doesn't seem tenable to say that
// printing the zero time correctly isn't supported in half the time
// zones. By comparison, it's reasonable to mishandle some times in
// the year -292277022399.
//
// All this is opaque to clients of the API and can be changed if a
// better implementation presents itself.

const (
	// The unsigned zero year for internal calculations.
	// Must be 1 mod 400, and times before it will not compute correctly,
	// but otherwise can be changed at will.
	absoluteZeroYear = -292277022399

	// The year of the zero Time.
	// Assumed by the unixToInternal computation below.
	internalYear = 1

	// Offsets to convert between internal and absolute or Unix times.
	absoluteToInternal int64 = (absoluteZeroYear - internalYear) * 365.2425 * secondsPerDay
	internalToAbsolute       = -absoluteToInternal

	unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	internalToUnix int64 = -unixToInternal

	wallToInternal        int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
	internalToWall        int64 = -wallToInternal
	unixToInternalminWall       = unixToInternal - minWall
)

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return t.sec() == 0 && t.nsec() == 0
}

// Nanosecond returns the nanosecond offset within the second specified by t,
// in the range [0, 999999999].
func (t Time) Nanosecond() int {
	return int(t.nsec())
}

// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64

const (
	minDuration Duration = -1 << 63
	maxDuration Duration = 1<<63 - 1
)

// Common durations. There is no definition for units of Day or larger
// to avoid confusion across daylight savings time zone transitions.
//
// To count the number of units in a Duration, divide:
//	second := time.Second
//	fmt.Print(int64(second/time.Millisecond)) // prints 1000
//
// To convert an integer number of units to a Duration, multiply:
//	seconds := 10
//	fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
//
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0s.
func (d Duration) String() string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			return "0s"
		case u < uint64(Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 's'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros. It omits the decimal
// point too when the fraction is 0. It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (d Duration) Nanoseconds() int64 { return int64(d) }

// Microseconds returns the duration as an integer microsecond count.
func (d Duration) Microseconds() int64 { return int64(d) / 1e3 }

// Milliseconds returns the duration as an integer millisecond count.
func (d Duration) Milliseconds() int64 { return int64(d) / 1e6 }

// These methods return float64 because the dominant
// use case is for printing a floating point number like 1.5s, and
// a truncation to integer would make them not useful in those cases.
// Splitting the integer and fraction ourselves guarantees that
// converting the returned float64 to an integer rounds the same
// way that a pure integer conversion would have, even in cases
// where, say, float64(d.Nanoseconds())/1e9 would have rounded
// differently.

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	sec := d / Second
	nsec := d % Second
	return float64(sec) + float64(nsec)/1e9
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	min := d / Minute
	nsec := d % Minute
	return float64(min) + float64(nsec)/(60*1e9)
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	hour := d / Hour
	nsec := d % Hour
	return float64(hour) + float64(nsec)/(60*60*1e9)
}

// Truncate returns the result of rounding d toward zero to a multiple of m.
// If m <= 0, Truncate returns d unchanged.
func (d Duration) Truncate(m Duration) Duration {
	if m <= 0 {
		return d
	}
	return d - d%m
}

// lessThanHalf reports whether x+x < y but avoids overflow,
// assuming x and y are both positive (Duration is signed).
func lessThanHalf(x, y Duration) bool {
	return uint64(x)+uint64(x) < uint64(y)
}

// Round returns the result of rounding d to the nearest multiple of m.
// The rounding behavior for halfway values is to round away from zero.
// If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration,
// Round returns the maximum (or minimum) duration.
// If m <= 0, Round returns d unchanged.
func (d Duration) Round(m Duration) Duration {
	if m <= 0 {
		return d
	}
	r := d % m
	if d < 0 {
		r = -r
		if lessThanHalf(r, m) {
			return d + r
		}
		if d1 := d - m + r; d1 < d {
			return d1
		}
		return minDuration // overflow
	}
	if lessThanHalf(r, m) {
		return d - r
	}
	if d1 := d + m - r; d1 > d {
		return d1
	}
	return maxDuration // overflow
}

// Add returns the time t+d.
func (t Time) Add(d Duration) Time {
	dsec := int64(d / 1e9)
	nsec := t.nsec() + int32(d%1e9)
	if nsec >= 1e9 {
		dsec++
		nsec -= 1e9
	} else if nsec < 0 {
		dsec--
		nsec += 1e9
	}
	t.wall = t.wall&^nsecMask | uint64(nsec) // update nsec
	t.addSec(dsec)
	if t.wall&hasMonotonic != 0 {
		te := t.ext + int64(d)
		if d < 0 && te > t.ext || d > 0 && te < t.ext {
			// Monotonic clock reading now out of range; degrade to wall-only.
			t.stripMono()
		} else {
			t.ext = te
		}
	}
	return t
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(u Time) Duration {
	if t.wall&u.wall&hasMonotonic != 0 {
		te := t.ext
		ue := u.ext
		d := Duration(te - ue)
		if d < 0 && te > ue {
			return maxDuration // t - u is positive out of range
		}
		if d > 0 && te < ue {
			return minDuration // t - u is negative out of range
		}
		return d
	}
	d := Duration(t.sec()-u.sec())*Second + Duration(t.nsec()-u.nsec())
	// Check for overflow or underflow.
	switch {
	case u.Add(d).Equal(t):
		return d // d is correct
	case t.Before(u):
		return minDuration // t - u is negative out of range
	default:
		return maxDuration // t - u is positive out of range
	}
}

// Since returns the time elapsed since t.
// It is shorthand for time.Now().Sub(t).
func Since(t Time) Duration {
	var now Time
	if t.wall&hasMonotonic != 0 {
		// Common case optimization: if t has monotonic time, then Sub will use only it.
		now = Time{hasMonotonic, runtimeNano() - startNano}
	} else {
		now = Now()
	}
	return now.Sub(t)
}

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
)

// Monotonic times are reported as offsets from startNano.
// We initialize startNano to runtimeNano() - 1 so that on systems where
// monotonic time resolution is fairly low (e.g. Windows 2008
// which appears to have a default resolution of 15ms),
// we avoid ever reporting a monotonic time of 0.
// (Callers may want to use 0 as "time not set".)
var startNano = runtimeNano() - 1

// Now returns the current local time.
func Now() Time {
	sec, nsec, mono := now()
	mono -= startNano
	sec += unixToInternalminWall
	if uint64(sec)>>33 != 0 {
		return Time{uint64(nsec), sec + minWall}
	}
	return Time{hasMonotonic | uint64(sec)<<nsecShift | uint64(nsec), mono}
}

/*
Returns current server epoch millis time without
GC dealing with *loc pointers
*/
func Epoch() Time32 {
	sec, nsec, mono := now()
	mono -= startNano
	var wall uint64
	var ext int64
	var unsec = uint64(nsec)
	sec += unixToInternalminWall
	var usec = uint64(sec)
	if usec >> 33 != 0 {
		wall, ext = unsec, sec+minWall
	} else {
		wall, ext = hasMonotonic|usec<<nsecShift|unsec, mono
	}
	tsec := ext
	if wall&hasMonotonic != 0 {
		tsec = wallToInternal + int64(wall<<1>>(nsecShift+1))
	}
	return Time32(tsec + internalToUnix)
}

func unixTime(sec int64, nsec int32) Time {
	return Time{uint64(nsec), sec + unixToInternal}
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
// Unix-like operating systems often record time as a 32-bit
// count of seconds, but since the method here returns a 64-bit
// value it is valid for billions of years into the past or future.
func (t Time) Unix() int64 {
	return t.unixSec()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64 (a date before the year
// 1678 or after 2262). Note that this means the result of calling UnixNano
// on the zero Time is undefined. The result does not depend on the
// location associated with t.
func (t Time) UnixNano() int64 {
	return (t.unixSec())*1e9 + int64(t.nsec())
}

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such
// value is 1<<63-1 (the largest int64 value).
func Unix(sec int64, nsec int64) Time {
	if nsec < 0 || nsec >= 1e9 {
		n := nsec / 1e9
		sec += n
		nsec -= n * 1e9
		if nsec < 0 {
			nsec += 1e9
			sec--
		}
	}
	return unixTime(sec, int32(nsec))
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Truncate(d Duration) Time {
	t.stripMono()
	if d <= 0 {
		return t
	}
	_, r := div(t, d)
	return t.Add(-r)
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t Time) Round(d Duration) Time {
	t.stripMono()
	if d <= 0 {
		return t
	}
	_, r := div(t, d)
	if lessThanHalf(r, d) {
		return t.Add(-r)
	}
	return t.Add(d - r)
}

// div divides t by d and returns the quotient parity and remainder.
// We don't use the quotient parity anymore (round half up instead of round to even)
// but it's still here in case we change our minds.
func div(t Time, d Duration) (qmod2 int, r Duration) {
	neg := false
	nsec := t.nsec()
	sec := t.sec()
	if sec < 0 {
		// Operate on absolute value.
		neg = true
		sec = -sec
		nsec = -nsec
		if nsec < 0 {
			nsec += 1e9
			sec-- // sec >= 1 before the -- so safe
		}
	}

	switch {
	// Special case: 2d divides 1 second.
	case d < Second && Second%(d+d) == 0:
		qmod2 = int(nsec/int32(d)) & 1
		r = Duration(nsec % int32(d))

	// Special case: d is a multiple of 1 second.
	case d%Second == 0:
		d1 := int64(d / Second)
		qmod2 = int(sec/d1) & 1
		r = Duration(sec%d1)*Second + Duration(nsec)

	// General case.
	// This could be faster if more cleverness were applied,
	// but it's really only here to avoid special case restrictions in the API.
	// No one will care about these cases.
	default:
		// Compute nanoseconds as 128-bit number.
		sec := uint64(sec)
		tmp := (sec >> 32) * 1e9
		u1 := tmp >> 32
		u0 := tmp << 32
		tmp = (sec & 0xFFFFFFFF) * 1e9
		u0x, u0 := u0, u0+tmp
		if u0 < u0x {
			u1++
		}
		u0x, u0 = u0, u0+uint64(nsec)
		if u0 < u0x {
			u1++
		}

		// Compute remainder by subtracting r<<k for decreasing k.
		// Quotient parity is whether we subtract on last round.
		d1 := uint64(d)
		for d1>>63 != 1 {
			d1 <<= 1
		}
		d0 := uint64(0)
		for {
			qmod2 = 0
			if u1 > d1 || u1 == d1 && u0 >= d0 {
				// subtract
				qmod2 = 1
				u0x, u0 = u0, u0-d0
				if u0 > u0x {
					u1--
				}
				u1 -= d1
			}
			if d1 == 0 && d0 == uint64(d) {
				break
			}
			d0 >>= 1
			d0 |= (d1 & 1) << 63
			d1 >>= 1
		}
		r = Duration(u0)
	}

	if neg && r != 0 {
		// If input was negative and not an exact multiple of d, we computed q, r such that
		//	q*d + r = -t
		// But the right answers are given by -(q-1), d-r:
		//	q*d + r = -t
		//	-q*d - r = t
		//	-(q-1)*d + (d - r) = t
		qmod2 ^= 1
		r = d - r
	}
	return
}
