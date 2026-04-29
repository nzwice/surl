# URL Shortener

## Functional Requirements

1. user inputs url and gets a shorten url
   - optional custom alias
   - optional expiration date
2. user goes to shorten url and being redirected to the original url

## Non-functional Requirements

1. redirection should be fast (<100ms)
2. system suports 100M DAU
   - assume 10 redirects / user -> 100M _ 10 / (24 _ 60 \* 60) = 11.5K QPS read
   - scale to 1B urls
3. ensure uniqueness of the shorten url

## Out of scopes

- Authenticating user
- Analytics on link clicks

## Core Entities

- User
- Original URL
- Shorten URL

## APIs

- FR1

```
POST /api/v1/urls
{
	"originalUrl": "https://myawesomewebsite.com",
	"alias": "awesome", // optional
	"expiratationDate": "2026-02-10", // optional
}

{
	"shortUrl": "https://surl.com/awesome123"
}
```

- FR2

```
GET /awesome123

302
Location: https://myawesomewebsite.com
```

302 is good here because it allows server to specify expired url as needed, also for future analytics

301 is permanent redirect

## Architecture

(draw along with High-Level Design)

![architecture](./surl.drawio.svg)

| URL                       |
| ------------------------- |
| id (PK)                   |
| shortCode (index, unique) |
| originalUrl               |
| expirationAt (index)      |
| createdBy                 |

## High-Level Design

- FR1:
  - steps to shorten urls:
    - is the originalUrl valid, if not abort
    - sanitize originalUrl
    - generate short code:
      - if user provides custom alias, check if it's already exists. If not return shortCode
      - generate shortCode (optionally use alias as prefix)
    - save to DB
    - return shortenUrl
  - cron job to clean up expired url
- FR2:
  - check if the short code exists
  - check if expiration date is before now
  - query the original url
  - return redirect response
- NFR1:
  - cache shortCode:{originalUrl, expirationTime} + TTL + LRU. Mem access time ~= (0.1 to 1)ms + 100K IOPS/sec
- NFR2:
  - scale to 1B urls: db size = (8 bytes + 2KB + 8 bytes + 8 bytes + 8 bytes) \* 1B ~= estimate 2.5TB. Still good with normal DB (10TB)
- NFR3:
  - shortCode = base62encode(now().nano % sizeSpace) // base62 is 0..9a..zA..Z
