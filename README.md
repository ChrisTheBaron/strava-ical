# strava-ical

Allows you to import your Strava activies into your Google calendar (or any calendar that supports ICAL).

[How to add to calendar](http://lmgtfy.com/?q=add+webcal+to+calendar).

Uses a lot of the same code as [ury-ical](https://github.com/UniversityRadioYork/ury-ical).

## Preparation

- Get a Strava API key. [https://strava.github.io/api/#access](https://strava.github.io/api/#access)
- Get your Strava athlete ID. _It's the number at the end of the URL when you visit your profile._
- Copy `config.example.toml` to `config.toml`
- Customise your calendar entry name and description in `config.toml` using the fields of [`strava.ActivitySummary`](https://github.com/strava/go.strava/blob/master/activities.go)
- There is also a field available called `DistanceKm` which is rounded to 2dp.

## Running

```bash
$ go build
$ ./strava-ical -c /path/to/config.toml
```

## Endpoints

## `/strava.ics`

All activities for the athlete.

# TODO

- [ ] Filter by activity type.
- [ ] Get description of activities.
- [ ] "Follow" other athletes.
