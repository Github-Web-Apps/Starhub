/*
    Minimal client-side Javascript API library for GitHub v3.
    Sebastian Hanula - http://www.hanula.com

    Date: 2015-01-09
*/

var GitHub = {

    // Parse 'Link' response headers and returns and object with rel as keys.
    // For example: {'next': 'https://...', 'last': '...'}
    _parseLinkHeader: function(xhr) {

        var link = xhr.getResponseHeader('Link');
        if(link) {

            var parsed = {},
                reg = new RegExp('<(.*?)>; rel="(.*)"');

            link.split(',').forEach(function(part) {
                var r = reg.exec(part);
                parsed[r[2]] = r[1];
            });
            return parsed;
        }
        return {};
    },

    // Fetch JSON contents from given github's API path or full url.
    // returns an $.ajax() promise object.
    //
    // When additional pagination data is found (via Link response header)
    // then it loads it automatically and merges the results.
    fetch: function(path, _previousData) {

        if($.isArray(path)) {
            path = 'https://api.github.com/' + path.join('/');
        }

        return $.ajax({
                url:  path,
                dataType: 'json'
            }).then(function(data, status, xhr) {
                var links = this._parseLinkHeader(xhr);
                if(_previousData) {
                    data = $.merge(_previousData, data);
                }
                if(links.next) {
                    return this.fetch(links.next, data);
                }
                return data;
            }.bind(this));
    },

    // Get user information by given username.
    getUser: function(username) {
        return this.fetch(['users', username]);
    },

    // Get starrted repository list for given user.
    getUserStars: function(username) {
        return this.fetch(['users', username, 'starred']);
    }
}
