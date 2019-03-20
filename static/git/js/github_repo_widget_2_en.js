$(function() {
    var $gitpane = $(".gitinfo");
    $.ajax({
        url: "https://api.github.com/repos" + git_name,
        dataType: "jsonp",
        success: function(results) {
            var gitres = results.data;
            var pushed_at = gitres.pushed_at.substr(0, 10);
            var $gitdiv = $(' <div class="github-box repo"> <div class="github-box-title"><h3>    <a class="owner" href="' + gitres.owner.url.replace("api.", "").replace("users/", "") + '" target="_blank">' + gitres.owner.login + '</a>    /     <a class="repo" href="' + gitres.url.replace("api.", "").replace("repos/", "") + '" target="_blank">' + gitres.name + '</a></h3><div class="github-stats">Star<a class="watchers" title="Watchers" href="' + gitres.url.replace("api.", "").replace("repos/", "") + '/watchers" target="_blank">' + gitres.watchers + '</a>Fork<a class="forks" title="Forks" href="' + gitres.url.replace("api.", "").replace("repos/", "") + '/network" target="_blank">' + gitres.forks + '</a></div></div><div class="github-box-content"><p class="description"><strong>Description：</strong>' + gitres.description + ' &mdash; <a href="' + gitres.url.replace("api.", "").replace("repos/", "") + '#readme"  target="_blank">More...</a></p><p class="link"><br><a href="' + gitres.homepage + '">' + (gitres.homepage == null ? "": gitres.homepage) + '</a></p> <br><table class="issues" width="100%"></table> <div style="height:1px; margin-top:-1px;clear: both;overflow:hidden;"></div>    </div><div class="github-box-download"><p class="updated"><a href="' + gitres.url.replace("api.", "").replace("repos/", "") + '/tree/master" target="_blank"><strong>master</strong>branch</a> Latest Commit：' + pushed_at + '</p> <p class="updated">Language:' + gitres.language + '</p><a class="download" href="' + gitres.url.replace("api.", "").replace("repos/", "") + '/zipball/master">Download Zip</a></div> </div>');
            $gitdiv.appendTo($gitpane);
            var userg = git_name.split("/")[1];
            var nameg = git_name.split("/")[2];
            try {} catch(e) {}
            if (gitres.has_issues && gitres.open_issues > 0) {
                $.ajax({
                    url: "https://api.github.com/repos" + git_name + "/issues?state=open&per_page=5&page=1&sort=updated",
                    dataType: "jsonp",
                    success: function(results) {
                        var issues = results.data;
                        var $issues_table = $(".github-box-content table");
                        $issues_table.append('<tr><td colspan="2"><a href="' + gitres.html_url + '/issues" target="_blank"><strong>Issues：</strong></a></td></tr>');
                        for (var i = 0; i < issues.length; i++) {
                            var updated_at = issues[i].updated_at.substr(0, 10);
                            $issues_table.append('<tr><td class="number" width="30"> #' + issues[i].number + '</td><td class="info"><a href="' + issues[i].html_url + '" target="_blank">' + issues[i].title + '</a></td><td class="author" width="250" align="right">From <a href="' + issues[i].user.url.replace("api.", "").replace("users/", "") + '" target="_blank">' + issues[i].user.login + '</a>&nbsp;&nbsp;<em style="font-size: 8pt;font-family: Candara,arial;color: #666;-webkit-text-size-adjust: none;">' + updated_at + "</em></td></tr>")
                        }
                    }
                })
            }
        }
    })
});