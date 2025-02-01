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
- [FAQs](#faqs)
- [Self hosting](#self-hosting)
  - [Frontend](#frontend)
  - [Migrations](#migrations)
  - [Backend](#backend)

## Features

> [!NOTE]
> Hosted version available by mid March :)

- 100% Opensource. Own your own data, metrics, and story
- Send investors' updates to anyone via email
- Manage metrics automatically from many sources. [See Integrations](https://malak.vc/integrations)
- Create, manage and share your decks.
- Manage your fundraising pipeline

## Built With

- Golang
- Postgresql
- Redis
- NextJS
- Tailwind CSS
- Stripe ( optional )

## Self-hosting

### Frontend

### Migrations

```sh
malak migrate
```

> [!IMPORTANT]
> Everytime you upgrade the backend or re-download the binary/docker image,
> it makes sense to run the migrations again as migrations are usually added to
> support newer
> features or enhance existing ones

### Backend

You can either download the raw binary or use our docker image. You can view
a list of all available releases on [Github](https://github.com/ayinke-llc/malak/releases).

The docker image is also available as `docker pull ghcr.io/ayinke-llc/malak:version`
where `version` can be the version you want. The version can be in two formats:

- semver version number e.g: `docker pull ghcr.io/ayinke-llc/malak:0.4.2`
- commit hash e.g: `docker pull ghcr.io/ayinke-llc/malak:101a434d`

Plans are extremely important in Malak. Even though you are self hosting
and not taking payments, you can still limit certain features for users
on your instance. Every company/workspace must be have a plan.

#### Listing plans

```sh
malak plans list
```

#### Create a new plan

```sh
malak plans create
```

#### Set a default plan for newly created workspaces

```sh
malak plans set-default plan_id
```

## FAQs

### Managing Minio

#### Images are not showing correctly in the editor

Please make sure the buckets are publicly available. This allows images to be read correctly

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
