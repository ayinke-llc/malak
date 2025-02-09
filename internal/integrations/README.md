# Self hosted Commercial/Enterprise Edition of Malak

Welcome to the Self hosted Commercial/Enterprise Edition of Malak

Our philosophy is simple, all the source code under AGPLv3 and free to self host yourself.
All auto syncing features from 3rd party systems are under a commercial license.

The [/integrations](https://github.com/ayinke-llc/malak/tree/main/internal/integrations) subfolder is the place
for all the **Self hosted Enterprise Edition** features from our [hosted](https://malak.vc/prigin)
plan such as Mercury, Brex, Stripe,Mixpanel, Paystack, Flutterwave and many more

> [!WARNING]
> This repository / tree is copyrighted (unlike the [main repo and rest of the code](https://github.com/ayinke-llc/malak)).
> You are not allowed to use this code to host your own version of <https://app.malak.vc>
> to include automatic data syncing
> without obtaining a proper license first. Shoot an email to <lanre@ayinke.ventures>

## Syncing integrations

The binary comes with a `integrations sync` command, it will auto sync the required data
from the integration and insert into the database.

```sh
malak integrations sync
```
