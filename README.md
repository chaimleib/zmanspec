BS"D

# zmanspec

## Summary

zmanspec describes a file format containing instructions for calculating zmanim.

## Structure

A zmanspec contains three top-level objects:

1. `var`
2. `daily`
3. `zman`

### `var`

The `var` object defines values which can be referenced during function calls.
It contains three child objects, organized by value type:

1. `ang` for angles.
2. `dur` for durations.
3. `hc` for DayHour configurations,
   which specify how to calculate halachic hours.

`ang.default` and `hc.default` will be used as default values
for functions where an angle or Hour Configuration are optional.

Certain values are commonly overridden by the consuming program,
and to respect such user overrides, a standard name should be used.

Overridable angle names (when referenced, they are prefixed with `ang.`):

- `havdalah` - sun's angle below the horizon used for Havdalah calculations.
  This is combined with `dur.havdalah`. Default `8.5`.

Overridable duration names (when referenced, they are prefixed with `dur.`):

- `candles` - a negative amount of time,
  specifying how early to light candles relative to sunset.
  Default `-18m`.
- `havdalah` - a positive amount of time, combined with `ang.havdalah`.
  Default `0`.

There are no overridable Hour Configurations (`hc.`).

### `daily`

The `daily` object maps human-readable zman names to function calls.
These are zmanim relevant for every day of the year.
It is the responsibility of the consuming program to sort these for display.

### `zman`

The `zman` object maps machine-readable zman IDs (strings) to function calls.
These zmanim are only relevant on certain days of the year,
or only calculated when requested by the consuming program.
It is up to the consumer to reference a Jewish calendar
and decide which ones must be displayed.
The expected IDs are as follows:

- `sunrise` - Hanetz HaChamah
- `sunset` - Shkiat HaChamah

- `candles` - when to light candles prior to Shabbat and holidays,
  when no other holiday or Shabbat immediately precedes candlelighting.
- `havdalah` - when Shabbat or the holiday ends.
- `candles_tzeit` - when to light candles for a holiday,
  when Shabbat or another holiday immediately precedes candlelighting.
  This is typically the same as havdalah,
  unless extra time is added to havdalah
  to add onto the Shabbat or holiday.

- `chanukah` - the preferred time to light Chanukah candles.
- `no_chametz` - on Erev Pesach, the time when the prohibition
  on possessing chametz begins.

- `fast_begin` - the time in the morning when a daytime-only fast begins.
- `fast_end` - the time in the evening when a daytime-only fast ends.

- `tisha_bav_begin` - the time in the evening when the fast of
  the Ninth of Av begins.
- `tisha_bav_end` - the time in the evening when the fast of
  the Ninth of Av ends.

## Field types

### Duration

Durations are strings like "13m30s" or "-18m".

```
int := "0" | [1-9][0-9]{0,4} ;
duration := "0" |
  [ "-" (
    int "h" [ int "m" ] [ int "s" ] |
    int "m" [ int "s" ] |
    int "s"
  ) ;
```

Allowed units are hours (h), minutes (m) and seconds (s),
and they must appear in that order, if present.
A unit may be omitted, as long as there remains
at least one other unit amount or "0".

A related reference implementation of a durations parser,
which supports a strict superset of our allowed syntax,
can be found in the Go standard library.
https://cs.opensource.google/go/go/+/refs/tags/go1.25.5:src/time/format.go;l=1621

For zmanim purposes only hours (h), minutes (m) and seconds (s) are allowed,
and only with integers up to 5 digits long.
The only prefix allowed is the optional `-` negative sign.

### Angles

Angles are numbers like `8.5`, `0.833` and the like.

These are solar angles of depression below the horizon, measured in degrees.
They may be used for both sunrise and sunset,
depending on which function they are used with.
Negative numbers are allowed, in case a synthetic sunrise or sunset time
is needed as an intermediate value.

### Variable references

A variable reference is a string like `"dur.havdalah"`,
`"hc.mga"` or `"ang.default"`.

It evaluates to the corresponding value in the `var` object.

### Function calls

Function calls are lists like `["Sunrise"]`, `["Sunrise", 10.2]`,
`["Sunset", "ang.default", "72m"]`.

They begin with a string containing the name of the function,
which must begin with a capital letter.
Most functions are variadic (take a variable number of arguments).
A function should respect default values if not specified.

If a function name ends in `0`, sea level elevation
and standard weather conditions must be assumed by the consuming program.
Without the `0` suffix,
additional location and weather information may be considered.

The following functions are defined:

- `["Sunrise", angle, duration, altLocation, dayOffset]` - returns when
  the center of the sun will reach the given angle of depression
  below the horizon in the morning.
  `angle` is degrees below the horizon.
  If not provided, assumes `ang.default`.
  `duration` is a time adjustment away from noon.
  If not provided, assumes 0.
  `altLocation` is a boolean which changes the location
  in which to calculate the time.
  If not provided, assumes `false`.
  This is useful if the specified time cannot be calculated
  in the primary location
  (e.g. near the poles when sunrise/sunset may not occur),
  and it is necessary to use a time calculated for a nearby location.
  `dayOffset` changes the calendar day, where positive numbers
  adjust to later days.
  If not provided, assumes 0.
- `["Sunrise0", angle, duration, altLocation, dayOffset]` - like `Sunrise`,
  at sea level with standard weather.
- `["Sunset", angle, duration, altLocation, dayOffset]` - returns when
  the center of the sun will reach the given angle of depression
  below the horizon in the evening.
  Parameters are like `Sunrise`.
- `["Sunset0", angle, duration, altLocation, dayOffset]` - like `Sunset`,
  at sea level with standard weather.
- `["Midnight", duration]` - returns the time when the sun
  is lowest in the sky.
  For a given date, `Midnight` precedes that day's `Noon`.
  If provided, `duration` gets added.
- `["Noon", duration]` - returns the time when the sun
  is highest in the sky.
  If provided, `duration` gets added.
- `["Hour", numHours, hourConfig]` - returns when this `numHours`
  of halachic hours into the day will be reached.
  `hourConfig` configures how hours are calculated.
  If `hourConfig` is not provided, `hc.default` is assumed.
- `["Hour0", numHours, hourConfig]` - like `Hour`, at sea level
  with standard weather.
- `["Min", a, b ...]` - returns the earliest of the provided times,
  or the smallest of the provided quantities.
  It must support times, durations, and numbers.
  Nulls are discarded.
- `["Max", a, b ...]` - returns the latest of the provided times,
  or the largest of the provided quantities.
  It must support times, durations, and numbers.
  Nulls are discarded.
- `["Priority", a, b ...]` - returns the first non-null item.

### Hour Configurations

Hour Configurations are lists of values used in computing halachic hours.

The most common form is `[angle, duration]`.
Here, `angle` is the solar angle of depression below the horizon
at during sunrise and sunset.
`duration` is an adjustment added on to the times when the sun
reaches the specified angle.
Positive durations adjust away from solar noon.
When `["DayHour", h, [angle, duration]]` is called,
The time between `["Sunrise", angle, duration]`
and `["Sunset", angle, duration]` is divided into 12 parts
and multiplied by `h`.

`[riseAngle, riseDuration, setAngle, setDuration]` allows
the beginning and end positions of the sun to be asymmetric.

`[startTime, endTime, divisions]` allows dividing a part of the day
into any number of hours.
The `startTime` and `endTime` will be function calls.
`divisions` may be 12, 6, or any other number of parts
to divide the duration with.

## Background

### What are zmanim?

Many Jewish observances are bound to particular times of the day and night.
For example:

 - A certain time range assigned to the morning prayers.
 - Shabbat candles must be lit before sunset.
 - Daytime fasts begin as the night sky begins to brighten.

These times are called zmanim (Hebrew for "times").
They are calculated based on the position of the sun
at a particular location, timezone, and day.
Sometimes fixed-duration offsets are added or subtracted.

The shared rabbinic tradition of all Jews prescribes how to calculate zmanim.
For example, all Jews agree that tallis and tefilin may be worn
starting from when one person can recognize an acquaintance
from four cubits away.
Likewise, all Jews agree that halachic hours are calculated
by dividing the daylight hours into twelve parts.

### Slightly different customs

The difficulty is in translating certain definitions
into concrete numerical values.
At what angle of the sun does it become bright enough
to recognize an acquaintance?
How much familiarity with this acquaintance is too much?
For calculating halachic hours, do we measure sunrise and sunset
using sea level, the elevation of Jerusalem, or the local elevation?
Does the elevation of the east and west horizon matter?

There is also debate about fixed duration offsets.
In the Talmud, these are described as the time it takes a person
to walk a unit of distance called a mil.
Opinions on what this time is range from
18 minutes to 20 minutes to 22.5 minutes to 24 minutes.
It is also debated whether to translate this offset to a solar angle or not,
in order to reproduce the sky brightness
which would have been observed in Jerusalem.

There are also various twilight times described in the Talmud,
and which one is significant for which context is debated.

Aside from the above, it is a practice to add time
to the end of Shabbat and holy days.
Some add time by adjusting the solar angle, and others by adding a fixed offset.

Because of all these factors, calculated times will vary
from rabbi to rabbi, from community to community,
even when calculated for the same date and place.

### Motivation

While Jewish communities all benefit from zmanim calculations,
no zmanim calculator currently serves them all.[^1]

[^1]:
- `hebcal` allows customizing the candle-lighting offset
  and the havdalah solar angle and offset,
  while all other zmanim are produced from hard-coded angles
  and offsets.
- MyZmanim.com is one of the most comprehensive and advanced zmanim web sites.
  It allows customizing elevation of the observer and the east and west horizons.
  It also has settings for temperature and air pressure at sunrise and sunset.
  It provides both a basic and an extended zmanim view,
  but few users will ever need to compare the zmanim from multiple opinions,
  which appear even in the basic view.
- Ou.org/calendar/#daily provides a listing of zmanim,
  but which method is being followed is difficult to figure out.
  Certain zmanim are also missing,
  such as when is the last time to burn chametz before Pesach,
  and the extended view leaves out all special zmanim
  for Shabbat and holidays.
- Chabad.org/zmanim is good for presenting a single opinion's set of zmanim.
  However, that also makes it less useful for those who follow other opinions.
- Kosherjava.com/zmanim-project lists over 150 zmanim types
  that it is capable of calculating,
  including 17 ways to calculate Sof Zman Kriat Shma
  and 28 ways to calculate Tzais HaKochavim.
  The enumeration, however, shows that these are hard-coded algorithms.
  Many of the algorithms have been marked deprecated
  because they were written based on misunderstandings
  of the rabbinic literature.

The objective of this zmanspec is to define the shape and semantics
of a portable file which encodes how to calculate zmanim.

It is not a goal to be the authority on how to calculate zmanim,
but rather to allow users to configure the calculation and display of zmanim
as taught by their rabbi.
