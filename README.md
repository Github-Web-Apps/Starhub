<p align="center">
  <img alt="Starhub Logo" src="https://github.com/Github-Web-Apps/Starhub/raw/master/static/logo/logo-128.png" />
  <h1 align="center">Starhub</h1>
  <h3 align="center">https://starhub.be/</h3>
  <p align="center">All about your Github account, public and private activity, stars count, release download count, who followed/unfollowed and starred/unstarred your Github repositories plus daily email notification about changes and much more.</p>
</p>

---

# Features:

- My-Starhub: display total stars for all repository
- My-Starhub: display total public repos
- My-Starhub: public and private activity history listing
- My-Starhub: activity history listing filter/search
- My-Starhub: watch all repos (forked one as well)
- My-Starhub: no github write access required
- My-Starhub: daily email notifciation (only on changes)
- My-Starhub: notification for new followers
- My-Starhub: notification for new unfollower
- My-Starhub: notification for new starred repository
- My-Starhub: notification for new unstarred repository
- My-Starhub: global user statistics on activities and used languages
- My-Starhub: main user organization statistic
- My-Starhub: github profile center
- My-Starhub: starhub iframe app loader
- Statistics: stars analytics over time with a graph
- Statistics: repo commmits, issue and activity analytics (under stars statics section)
- Statistics: releases downloads couter & stats
- Tools: github applications selection and listing
- Tools: user mail finder
- Tools: username converter
- Tools: user-id converter
- Tools: site preview (git html file preview)
- Tools: git folder downloader
- Profiler: github profile for any github user
- Profiler: various statistics (activities/repos/langs)
- Profiler: work with (url) get method for an easy access (starhub.be/username)
- Engine: encrypted requests (ssl https)
- And much more...

# Screenshot

<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/1.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/2.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/3.png" width="700px"></img> 
</div>

# Running it locally

**Cloning**

For Go projects to work they have to be cloned on the right places.

Let's assume `~/Code/Go` as our default Go projects folder.

So:

```sh
git clone git@github.com:Github-Web-Apps/Starhub.git
cd Starhub
```

**Dependencies**

Now, install Go 1.11+ and run:

```sh
make setup
```

To install the other project's dependencies.

**Building**

Just run:

```sh
make build
```

**Running the tests**

Just run:

```sh
make test
```

**Database setup**

Start up postgres and run:

```sh
createdb watchub
for sql in ./migrations/*; do psql watchub -f $sql; done
```

**Tunnel with ngrok**

To test the entire flow, you'll need to install ngrok.

Install it, then just run:

```sh
ngrok http 3000
```

Then, create an application on [github](https://github.com/settings/applications/new).

Fill it like this:

1. Application name: `Starhub-Dev-Username`
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
