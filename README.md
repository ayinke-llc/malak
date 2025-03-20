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
- [Self hosting](#self-hosting)

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

Detailed guide at <https://docs.malak.vc/self-hosting/installation>

## LICENSE

Malak is a commercial open source product, which means some parts of this
open source repository require a commercial license. The concept is
called "Open Core" where the core software (99%) is fully open source,
licensed under AGPLv3 and the last 1% is covered under a
commercial license.
([/internal/integrations Enterprise/Commercial Edition](https://github.com/ayinke-llc/malak/tree/main/internal/integrations)) which we
believe is entirely relevant for larger
organisations that require those features.

Our philosophy is simple, all features are open-source under AGPLv3.
But 3rd party integrations and auto-syncing are under a commercial/Enterprise license.

[See details of commercial/Enterprise license here](https://github.com/ayinke-llc/malak/tree/main/internal/integrations#readme)
