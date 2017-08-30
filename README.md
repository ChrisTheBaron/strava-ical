# strava-ical

Allows users to import their Strava activies into any calendar that supports ICAL.

[How to add to calendar](http://lmgtfy.com/?q=add+webcal+to+calendar).

Uses a lot of the same code as [ury-ical](https://github.com/UniversityRadioYork/ury-ical).

## Running

```bash
$ go build
$ ./strava-ical -c /path/to/config.toml
```

## Endpoints

### `GET /`

Show Static Content.

### `GET /login`

 - Check for cookie
 - If cookie present, then `/calendars`
 - Else `/auth`

### `GET /auth`

  - Redirect to Strava login.
    - On success.
    - Store token.
    - Generate JWT.
    - Store JWT as cookie.
 - On failure
   - Show static content.

### `GET /calendars`

  - If no cookie, redirect to `/login`
  - List users calendars

### `POST /calendars/`

 - Create new calendar for user.
 - Redirect to `/calendar/{id}`

### `GET /calendar/{id}`

 - Show calendar settings.

### `GET /calendar/{id}.ics`

 - Get calendar as ICS.

### `DELETE /calendar/{id}`

 - Delete calendar.
 - Redirect to `/calendars`
