<p align="center">
  <img alt="Starhub Logo" src="https://github.com/Github-Web-Apps/Starhub/raw/master/static/logo/logo-128.png" />
  <h1 align="center">Starhub</h1>
  <h3 align="center">https://starhub.be/</h3>
  <h5 align="center">https://starhub.be/YOUR-GITHUB-LOGIN</h5>
  <p align="center">All about your Github account, public and private activity, stars count, release download count, who followed/unfollowed and starred/unstarred your Github repositories plus daily email notification about changes and much more.</p>
</p>

---

# Features:

- My-Starhub: display total repos, stars and followers
- My-Starhub: public and private activity history listing plus filter
- My-Starhub: watch all repos for changes
- My-Starhub: daily email notification *(only on changes)*
- My-Starhub: notification for new followers, un-follower, stars, un-star
- My-Starhub: global user statistics on activities and used languages
- My-Starhub: main user organization statistic
- My-Starhub: list search and filter starred repos
- Statistics: stars and releases downloads counter analytics
- Tools: github applications selection and listing
- Tools: mail finder, username and user-id converter, site preview and git downloader
- Profiler: github profile for any github user with various statistics
- Profiler: direct access (starhub.be/github-user-name)
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
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/4.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/5.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/6.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/7.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/8.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/9.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/10.png" width="700px"></img> 
</div>
</br>
<div align="center">
    <img src="https://raw.githubusercontent.com/Github-Web-Apps/Starhub/master/screenshot/11.png" width="700px"></img> 
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
