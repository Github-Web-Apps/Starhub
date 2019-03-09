{{if .NewFollowers }}
## Users that are now **following** you:

<ul>
{{ range $index, $element := .NewFollowers }}
<li><a href="https://github.com/{{ . }}">@{{ . }}</a></li>
{{ end }}
</ul>
{{ end }}

{{if .Unfollowers }}
## Users that are **not following** you anymore:

<ul>
{{ range $index, $element := .Unfollowers }}
<li><a href="https://github.com/{{ . }}">@{{ . }}</a></li>
{{ end }}
</ul>
{{ end }}

{{if .NewStars}}
## New stargazers on public repos you have access to:

{{ range $index, $element := .NewStars }}
### [{{ .Repo }}](https://github.com/{{ .Repo }})

<ul>
{{ range $index, $element := .Users }}
<li><a href="https://github.com/{{ . }}">@{{ . }}</a></li>
{{ end }}
</ul>
{{ end }}
{{ end }}

{{if .Unstars}}
## Users that unstarred public repos you have access to:

{{ range $index, $element := .Unstars }}
### [{{ .Repo }}](https://github.com/{{ .Repo }})

<ul>
{{ range $index, $element := .Users }}
<li><a href="https://github.com/{{ . }}">@{{ . }}</a></li>
{{ end }}
</ul>
{{ end }}
{{ end }}

---

You now have:

- {{ .Followers }} followers;
- {{ .Repos }} repositories;
- {{ .Stars }} stars across all those repositories.
