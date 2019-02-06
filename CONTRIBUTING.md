# Running it locally

## Cloning

For Go projects to work they have to be cloned on the right places.

Let's assume `~/Code/Go` as our default Go projects folder.

So:

```sh
git clone git@github.com:caarlos0/watchub.git
cd watchub
```

## Dependencies

Now, install Go 1.11+ and run:

```sh
make setup
```

To install the other project's dependencies.

## Linting

Just run:

```sh
make lint
```

## Running the tests

Just run:

```sh
make test
```

## Database setup

Start up postgres and run:

```sh
createdb watchub
for sql in ./migrations/*; do psql watchub -f $sql; done
```

## Tunnel with ngrok

To test the entire flow, you'll need to install ngrok.

Install it, then just run:

```sh
ngrok http 3000
```

Then, create an application on [github](https://github.com/settings/applications/new).

Fill it like this:

1. Application name: `Watchub dev`
1. Homepage URL: the ngrok http forwarding URL, e.g. `http://6f7ca783.ngrok.io`
1. Application description: empty
1. Authorization callback URL: same as homepage url, but with a `/login/callback`
suffix. e.g.: `http://6f7ca783.ngrok.io/login/callback`

GitHub will then give you a Client ID and a Client Secret.

Export them like this:

```sh
export GITHUB_CLIENT_ID="your client id"
export GITHUB_CLIENT_SECRET="your client secret"
```

And then just run the app:

```sh
go run main.go
```
