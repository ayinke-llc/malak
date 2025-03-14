# Malak

[![Codecov](https://codecov.io/gh/ayinke-llc/malak/graph/badge.svg?token=J1AVNTOCVY)](https://codecov.io/gh/ayinke-llc/malak)
[![Go Report Card](https://goreportcard.com/badge/github.com/ayinke-llc/malak)](https://goreportcard.com/report/github.com/ayinke-llc/malak)

Meet Malak, an opensource performance insights management tool for startups and their investors.

Angelist Raise, Carta and other tools are awesome to communicate with investors. However, these tools are very
limited in terms of control and customization. Or can cost thousands of $ to run and still be
limited to certain features subsets. Or have your data being used [to cross-sell to other businesses for
their own gain](https://x.com/karrisaarinen/status/1743398553500971331)

That's where Malak comes in. Self-hosted on your own infra or hosted by us.
White-label by design. Ready to be deployed on your own domain.
Own full control of your data and company state if you decide to selfhost.

- [Features](#features)
- [License](#license)
- [FAQs](#faqs)
- [Self hosting](#self-hosting)
  - [Frontend](#frontend)
    - [Dashboard app](#dashboard)
    - [Deck viewer](#deck-viewer)
  - [Migrations](#migrations)
  - [Backend](#backend)

## Features

> [!NOTE]
> Hosted version available by mid March :)

- 100% Opensource. Own your own data, metrics, and story
- Send investors' updates to anyone via email
- Manage metrics automatically from many sources such as Mercury,
  Brex, Stripe and more. [Full list](https://malak.vc/integrations)
- Create, manage and share your decks.
- Manage your fundraising pipeline
- Own your captable

### Built With

- Golang
- Postgresql
- Redis
- NextJS
- Tailwind CSS
- Stripe ( optional )

## Self-hosting

### Migrations

```sh
malak migrate
```

> [!IMPORTANT]
> Everytime you upgrade the backend or re-download the binary/docker image,
> it makes sense to run the migrations again as migrations are usually added to
> support newer features or enhance existing ones

### Backend

You can either download the raw binary or use our docker image. You can view
a list of all available releases on [Github](https://github.com/ayinke-llc/malak/releases).

The docker image is also available as `docker pull ghcr.io/ayinke-llc/malak:version`
where `version` can be the version you want. The version can be in two formats:

- semver version number e.g: `docker pull ghcr.io/ayinke-llc/malak:0.4.2`
- commit hash e.g: `docker pull ghcr.io/ayinke-llc/malak:101a434d`

At the very minimum, your configuration needs to look like this:

> [!IMPORTANT]
> You can also use environment values to configure this in cloud environments
> using the syntax `MALAK_`. So for `auth.google.client_id`, you can use `MALAK_AUTH_GOOGLE_CLIENT_ID`

Your `config.yml` below:

```yml
auth:
  google:
    client_id: "GOOGLE_CLIENT_ID"
    client_secret: "GOOGLE_CLIENT_SECRET"
    redirect_uri: "GOOGLE_REDIRECT_URI"
    is_enabled: true

  jwt:
    key: YOUR_SECRET_KEY

uploader:
  s3:
    access_key: YOUR_ACCESS_KEY
    access_secret: YOUR_SECRET_KEY
    region: lagos-1
    endpoint: http://localhost:9000
    use_tls: false

email:
  provider: smtp # OR resend
  sender: email@yourdomain
  sender_name: Malak
  resend:
    api_key: RESEND_API_KEY
  smtp:
    host: localhost
    password: random
    username: random
    use_tls: false
    port: 9125
```

> [!IMPORTANT]
> The docker image runs the http server by default but if
> you are using the binary, you can run it with `malak http`

Plans are extremely important in Malak. Even though you are self hosting
and not taking payments, you can still limit certain features for users
on your instance. Every company/workspace must be have a plan.

You can view a complete list of all available configuration here.

#### Listing plans

```sh
malak plans list
```

#### Create a new plan

```sh
malak plans create
```

#### Set a default plan for newly created workspaces

> every newly created user inherits this plan by default

```sh
malak plans set-default plan_id
```

#### Cron to send updates to investors

> [!IMPORTANT]
> You need `google-chrome` in your path as we use Headless Chrome to render and download the image we
> share in the email if your updates includes bar or pie charts

```sh
malak cron updates
```

> You can probably run this cron every 5 minutes or 1 hour or daily depending on what you want

#### Deck analytics cron

```sh
malak cron analytics
```

### Frontend

#### Dashboard

```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=
NEXT_PUBLIC_MALAK_POSTHOG_KEY=
NEXT_PUBLIC_MALAK_ENABLE_POSTHOG=
NEXT_PUBLIC_MALAK_POSTHOG_HOST=
NEXT_PUBLIC_MALAK_TERMS_CONDITION_LINK=
NEXT_PUBLIC_MALAK_PRIVACY_POLICY_LINK=
NEXT_PUBLIC_SENTRY_DSN=
NEXT_PUBLIC_DECKS_DOMAIN=https://deck.yourdowmain
NEXT_PUBLIC_SUPPORT_EMAIL=support@yourdomain
```

#### Deck viewer

The deck viewer is hosted in another repo <https://github.com/ayinke-llc/malak-deck> .
You need to set this up if you want users to view and share decks with investors and track
them via emails and time spent

## LICENSE

Malak is a commercial open source product, which means some parts of this
open source repository require a commercial license. The concept is
called "Open Core" where the core software (99%) is fully open source,
licensed under AGPLv3 and the last 1% is covered under a
commercial license
([/internal/integrations Enterprise/Commercial Edition](https://github.com/ayinke-llc/malak/tree/main/internal/integrations)) which we
believe is entirely relevant for larger organisations that require those features.

Our philosophy is simple, all features are open-source under AGPLv3.
But 3rd party integrations and auto-syncing are under a commercial/Enterprise license.

[See details of commercial/Enterprise license here](https://github.com/ayinke-llc/malak/tree/main/internal/integrations#readme)

## FAQs

### Using Minio locally for S3

#### Buckets

You need to create two buckets:

- `malak` : for publicly available objects
- `deck`: Make sure this bucket is private. Items here should not be publicly available

> [!NOTE]
> This bucket names are not hard-coded. They are just defaults, you can
> name your buckets however you choose to. Just add the bucket names correctly
> in your configuration

#### Images are not showing correctly in the editor

Please make sure the image bucket is publicly available. This allows images to be read correctly

#### Configure mc client

```sh
mc alias set malak http://localhost:9000 access secret
## this assumes you created a bucket in Minio called malak
mc stat malak/malak
mc share download --recursive malak/malak
```

This should return something like:

```txt
URL: http://localhost:9000/malak/575ca9e2-9782-439e-9a74-c0a8923c7e1e
Expire: 7 days 0 hours 0 minutes 0 seconds
Share: http://localhost:9000/malak/575ca9e2-9782-439e-9a74-c0a8923c7e1e?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=yy9Od9rfMjCX
I553bZKp%2F20240921%2Flagos-1%2Fs3%2Faws4_request&X-Amz-Date=20240921T161000Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Sign
ature=b96a68d4d4b048d7e8092ff3c9241880f3d32d78ae85023e0abd68b681ed33bd

URL: http://localhost:9000/malak/CleanShot 2024-09-17 at 20.22.15.png
Expire: 7 days 0 hours 0 minutes 0 seconds
Share: http://localhost:9000/malak/CleanShot%202024-09-17%20at%2020.22.15.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=yy9Od9
rfMjCXI553bZKp%2F20240921%2Flagos-1%2Fs3%2Faws4_request&X-Amz-Date=20240921T161000Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Am
z-Signature=6cc3487799e07e2de2681c024a842f81b189fc5968f1ed79b8e7760ef1c3019e

```

### Viewing swagger docs

You need to first enable the dev swagger UI with the following config:

```yml
http:
  swagger:
    port: 9999
    ui_enabled: true
```

After which you can visit <http://localhost:9999/swagger/>
