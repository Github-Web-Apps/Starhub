<p align="center">
  <img alt="Watchub Logo" src="https://github.com/caarlos0/watchub/raw/master/static/apple-touch-icon-144x144.png" />
  <h3 align="center">Watchub - Mirror</h3>
  <h3 align="center">https://watchub.duckdns.org/</h3>
  <p align="center">Be notified of people who followed/unfollowed you and starred/unstarred your repositories.</p>
  <p align="center">
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/caarlos0/watchub"><img alt="Travis" src="https://img.shields.io/travis/caarlos0/watchub.svg?style=flat-square"></a>
    <a href="https://codecov.io/gh/caarlos0/watchub"><img alt="Codecov branch master" src="https://img.shields.io/codecov/c/github/caarlos0/watchub/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/caarlos0/watchub"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/caarlos0/watchub?style=flat-square"></a>
    <a href="https://godoc.org/github.com/caarlos0/watchub"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://saythanks.io/to/caarlos0"><img alt="SayThanks.io" src="https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square"></a>
  </p>
</p>

---

**Watchub-Mirro:** This is a mirror of the upstream project with these key differences:

- Watch all repos (forked one as well)
- No write access required
- SSL HTTPS

**Watchub:** watches your account and notify you via email (at most one email/day) of the following events:

- Someone followed you
- Someone unfollowed you
- Someone starred a repository you have access to
- Someone unstarred a repository you have access to

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
1. Homepage URL: the ngrok http forwarding URL, e.g. `https://6f7ca783.ngrok.io`
1. Application description: empty
1. Authorization callback URL: same as homepage url, but with a `/login/callback`
suffix. e.g.: `https://6f7ca783.ngrok.io/login/callback`

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

# Contributor Covenant Code of Conduct

## Our Pledge

In the interest of fostering an open and welcoming environment, we as
contributors and maintainers pledge to making participation in our project and
our community a harassment-free experience for everyone, regardless of age, body
size, disability, ethnicity, gender identity and expression, level of experience,
nationality, personal appearance, race, religion, or sexual identity and
orientation.

## Our Standards

Examples of behavior that contributes to creating a positive environment
include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or
advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic
  address, without explicit permission
* Other conduct which could reasonably be considered inappropriate in a
  professional setting

## Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable
behavior and are expected to take appropriate and fair corrective action in
response to any instances of unacceptable behavior.

Project maintainers have the right and responsibility to remove, edit, or
reject comments, commits, code, wiki edits, issues, and other contributions
that are not aligned to this Code of Conduct, or to ban temporarily or
permanently any contributor for other behaviors that they deem inappropriate,
threatening, offensive, or harmful.

## Scope

This Code of Conduct applies both within project spaces and in public spaces
when an individual is representing the project or its community. Examples of
representing a project or community include using an official project e-mail
address, posting via an official social media account, or acting as an appointed
representative at an online or offline event. Representation of a project may be
further defined and clarified by project maintainers.

## Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported by contacting the project team at root@carlosbecker.com. All
complaints will be reviewed and investigated and will result in a response that
is deemed necessary and appropriate to the circumstances. The project team is
obligated to maintain confidentiality with regard to the reporter of an incident.
Further details of specific enforcement policies may be posted separately.

Project maintainers who do not follow or enforce the Code of Conduct in good
faith may face temporary or permanent repercussions as determined by other
members of the project's leadership.

## Attribution

This Code of Conduct is adapted from the [Contributor Covenant][homepage], version 1.4,
available at [https://contributor-covenant.org/version/1/4][version]

[homepage]: https://contributor-covenant.org
[version]: https://contributor-covenant.org/version/1/4/


