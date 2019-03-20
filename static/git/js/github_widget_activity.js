;(function($, window, document, undefined) {

  'use strict';

  var defaults = {
    username: 'octocat',
    timeout: 300,
    debug: true
  };

  var views = (function() {
    var helpers = {
      capitalize: function(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
      },
      date: function(str) {
        var d = new Date(str);

        // If a fancy date library is available, use it to generate fuzzy relative dates
        // otherwise use the US-centric "month-day-year" format
        if (window.moment) {
          return window.moment(d).fromNow(); // moment.js
        } else if ($.timeago) {
          return $.timeago(d); // jquery.timeago.js
        } else {
          return (d.getMonth() + 1) + '-' + (d.getDate()) + '-' + (d.getFullYear());
        }
      },
      // Determine if an issue is a Pull Request or just a regular Issue
      issueType: function(issue) {
        if (issue.pull_request && issue.pull_request.html_url &&
            issue.pull_request.diff_url && issue.pull_request.patch_url) {
          return 'pull request';
        } else {
          return 'issue';
        }
      },
      // Get the last piece of a slash separated string (such as an URL or Git ref)
      tail: function(str) {
        return str.split("/").pop();
      }
    };

    var templates = {
      boilerplate: function() {
        return ''+
        '<div class="activity">'+
          '<div class="activity-head"></div>'+
          '<ul class="activity-body"></ul>'+
          '<div class="activity-foot">'+
            '<a href="https://github.com/smuyyh">Powered by smuyyh</a>'+
          '</div>'+
        '</div>';
      },
      profileWithName: function(data) {
        var avatar = data.avatar_url,
            name   = data.name,
            login  = data.login;

        return ''+
        '<a href="https://github.com/'+login+'"><img src="'+avatar+'"></a>'+
        '<h4><a href="https://github.com/'+login+'">'+name+'</a></h4>'+
        '<small><a href="https://github.com/'+login+'">'+login+'</a></small>'+
        '<div class="clearfix"></div>';
      },
      profileWithoutName: function(data) {
        var avatar = data.avatar_url,
            login  = data.login;

        return ''+
        '<a href="https://github.com/'+login+'"><img src="'+avatar+'"></a>'+
        '<h4><a href="https://github.com/'+login+'">'+login+'</a></h4>'+
        '<div class="clearfix"></div>';
      },
      commentOnCommit: function(data) {
        var id         = data.id,
            date       = helpers.date(data.created_at),
            repoName   = data.repo.name,
            commentURL = data.payload.comment.html_url,
            commitID   = data.payload.comment.commit_id.substring(0,7);

        return ''+
        '<li id="'+id+'">'+
          'Commented on commit <a href="'+commentURL+'">'+commitID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createRepository: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            ref      = data.payload.ref;

        return ''+
        '<li id="'+id+'">'+
          'Created repository <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createBranch: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            ref      = data.payload.ref;

        return ''+
        '<li id="'+id+'">'+
          'Created branch <a href="https://github.com/'+repoName+'/tree/'+ref+'">'+ref+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createTag: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            ref      = data.payload.ref;

        return ''+
        '<li id="'+id+'">'+
          'Created tag <a href="https://github.com/'+repoName+'/tree/'+ref+'">'+ref+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      deleteBranch: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            ref      = data.payload.ref;

        return ''+
        '<li id="'+id+'">'+
          'Deleted branch '+ref+' '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      deleteTag: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            ref      = data.payload.ref;

        return ''+
        '<li id="'+id+'">'+
          'Deleted tag '+ref+' '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createDownload: function(data) {
        var id           = data.id,
            date         = helpers.date(data.created_at),
            repoName     = data.repo.name,
            downloadName = data.payload.download.name,
            downloadURL  = data.payload.download.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Created download <a href="'+downloadURL+'">'+downloadName+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      followUser: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            userName = data.payload.target.login,
            userURL  = data.payload.target.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Started following <a href="'+userURL+'">'+userName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      applyPatch: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name;

        return ''+
        '<li id="'+id+'">'+
          'Applied a patch '+
          'to <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      forkRepository: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            forkee   = data.payload.forkee.full_name;

        return ''+
        '<li id="'+id+'">'+
          'Forked <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          'to <a href="https://github.com/'+forkee+'">'+forkee+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createGist: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            gistID   = data.payload.gist.id,
            gistURL  = data.payload.gist.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Created gist <a href="'+gistURL+'">'+gistID+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      updateGist: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            gistID   = data.payload.gist.id,
            gistURL  = data.payload.gist.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Updated gist <a href="'+gistURL+'">'+gistID+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      editWiki: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name;

        return ''+
        '<li id="'+id+'">'+
          'Edited the wiki '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      commentOnIssue: function(data) {
        var id         = data.id,
            date       = helpers.date(data.created_at),
            repoName   = data.repo.name,
            issueID    = data.payload.issue.number,
            commentURL = data.payload.comment.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Commented on issue <a href="'+commentURL+'">#'+issueID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      commentOnPullRequest: function(data) {
        var id         = data.id,
            date       = helpers.date(data.created_at),
            repoName   = data.repo.name,
            issueID    = data.payload.issue.number,
            commentURL = data.payload.issue.pull_request.html_url+'#issuecomment-'+data.payload.comment.id;

        return ''+
        '<li id="'+id+'">'+
          'Commented on pull request <a href="'+commentURL+'">#'+issueID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      openIssue: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            issueID  = data.payload.issue.number,
            issueURL = data.payload.issue.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Opened issue <a href="'+issueURL+'">#'+issueID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      closeIssue: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            issueID  = data.payload.issue.number,
            issueURL = data.payload.issue.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Closed issue <a href="'+issueURL+'">#'+issueID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      reopenIssue: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            issueID  = data.payload.issue.number,
            issueURL = data.payload.issue.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Reopened issue <a href="'+issueURL+'">#'+issueID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      addUserToRepository: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            userName = data.payload.member.login,
            userURL  = data.payload.member.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Added <a href="'+userURL+'">'+userName+'</a> '+
          'to <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      openSourceRepository: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name;

        return ''+
        '<li id="'+id+'">'+
          'Open sourced repository <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      openPullRequest: function(data) {
        var id            = data.id,
            date          = helpers.date(data.created_at),
            repoName      = data.repo.name,
            pullRequestID = data.payload.number;

        return ''+
        '<li id="'+id+'">'+
          'Opened pull request <a href="https://github.com/'+repoName+'/pull/'+pullRequestID+'">#'+pullRequestID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      closePullRequest: function(data) {
        var id            = data.id,
            date          = helpers.date(data.created_at),
            repoName      = data.repo.name,
            pullRequestID = data.payload.number;

        return ''+
        '<li id="'+id+'">'+
          'Closed pull request <a href="https://github.com/'+repoName+'/pull/'+pullRequestID+'">#'+pullRequestID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      reopenPullRequest: function(data) {
        var id            = data.id,
            date          = helpers.date(data.created_at),
            repoName      = data.repo.name,
            pullRequestID = data.payload.number;

        return ''+
        '<li id="'+id+'">'+
          'Reopened pull request <a href="https://github.com/'+repoName+'/pull/'+pullRequestID+'">#'+pullRequestID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      synchronizePullRequest: function(data) {
        var id            = data.id,
            date          = helpers.date(data.created_at),
            repoName      = data.repo.name,
            pullRequestID = data.payload.number;

        return ''+
        '<li id="'+id+'">'+
          'Synchronized pull request <a href="https://github.com/'+repoName+'/pull/'+pullRequestID+'">#'+pullRequestID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      commentOnPullRequestDiff: function(data) {
        var id            = data.id,
            date          = helpers.date(data.created_at),
            repoName      = data.repo.name,
            commentURL    = data.payload.comment.html_url,
            pullRequestID = helpers.tail(data.payload.comment.pull_request_url);

        return ''+
        '<li id="'+id+'">'+
          'Commented on pull request <a href="'+commentURL+'">#'+pullRequestID+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      pushToBranch: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name,
            count    = (data.payload.size === 1) ? '1 commit ' : data.payload.size+' commits ',
            refTail  = helpers.tail(data.payload.ref);

        return ''+
        '<li id="'+id+'">'+
          'Pushed '+count+
          'to <a href="https://github.com/'+repoName+'/tree/'+refTail+'">'+refTail+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      createRelease: function(data) {
        var id          = data.id,
            date        = helpers.date(data.created_at),
            repoName    = data.repo.name,
            releaseName = data.payload.release.name,
            releaseURL  = data.payload.release.html_url;

        return ''+
        '<li id="'+id+'">'+
          'Created release <a href="'+releaseURL+'">'+releaseName+'</a> '+
          'at <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      },
      starRepository: function(data) {
        var id       = data.id,
            date     = helpers.date(data.created_at),
            repoName = data.repo.name;

        return ''+
        '<li id="'+id+'">'+
          'Starred repository <a href="https://github.com/'+repoName+'">'+repoName+'</a> '+
          '<small>'+date+'</small>'+
        '</li>';
      }
    };

    var render = function(view, data) {

      // Create widget skeleton
      if (view === 'boilerplate') {
        return templates.boilerplate();
      }

      // Populate widget header with user info
      if (view === 'profile') {
        return (data.name) ? templates.profileWithName(data) : templates.profileWithoutName(data);
      }

      // Append event info to the widget body
      if (view === 'event') {

        // Comment on a Commit
        if (data.type === 'CommitCommentEvent') {
          return templates.commentOnCommit(data);
        }

        // Create a Repository
        if (data.type === 'CreateEvent' && data.payload.ref_type === 'repository') {
          return templates.createRepository(data);
        }

        // Create a Branch
        if (data.type === 'CreateEvent' && data.payload.ref_type === 'branch') {
          return templates.createBranch(data);
        }

        // Create a Tag
        if (data.type === 'CreateEvent' && data.payload.ref_type === 'tag') {
          return templates.createTag(data);
        }

        // Delete a Branch
        if (data.type === 'DeleteEvent' && data.payload.ref_type === 'branch') {
          return templates.deleteBranch(data);
        }

        // Delete a Tag
        if (data.type === 'DeleteEvent' && data.payload.ref_type === 'tag') {
          return templates.deleteTag(data);
        }

        // Create a Download
        if (data.type === 'DownloadEvent') {
          return templates.createDownload(data);
        }

        // Follow a User
        if (data.type === 'FollowEvent') {
          return templates.followUser(data);
        }

        // Apply a patch in the Fork Queue
        if (data.type === 'ForkApplyEvent') {
          return templates.applyPatch(data);
        }

        // Fork a Repository
        if (data.type === 'ForkEvent') {
          return templates.forkRepository(data);
        }

        // Create a Gist
        if (data.type === 'GistEvent' && data.payload.action === 'create') {
          return templates.createGist(data);
        }

        // Update a Gist
        if (data.type === 'GistEvent' && data.payload.action === 'update') {
          return templates.updateGist(data);
        }

        // Edit a Wiki
        if (data.type === 'GollumEvent') {
          return templates.editWiki(data);
        }

        // Comment on an Issue
        if (data.type === 'IssueCommentEvent' && helpers.issueType(data.payload.issue) === 'issue') {
          return templates.commentOnIssue(data);
        }

        // Comment on a Pull Request
        if (data.type === 'IssueCommentEvent' && helpers.issueType(data.payload.issue) === 'pull request') {
          return templates.commentOnPullRequest(data);
        }

        // Open an Issue
        if (data.type === 'IssuesEvent' && data.payload.action === 'opened') {
          return templates.openIssue(data);
        }

        // Close an Issue
        if (data.type === 'IssuesEvent' && data.payload.action === 'closed') {
          return templates.closeIssue(data);
        }

        // Reopen an Issue
        if (data.type === 'IssuesEvent' && data.payload.action === 'reopened') {
          return templates.reopenIssue(data);
        }

        // Add a User to a Repository
        if (data.type === 'MemberEvent') {
          return templates.addUserToRepository(data);
        }

        // Open Source a Repository
        if (data.type === 'PublicEvent') {
          return templates.openSourceRepository(data);
        }

        // Open a Pull Request
        if (data.type === 'PullRequestEvent' && data.payload.action === 'opened') {
          return templates.openPullRequest(data);
        }

        // Close a Pull Request
        if (data.type === 'PullRequestEvent' && data.payload.action === 'closed') {
          return templates.closePullRequest(data);
        }

        // Reopen a Pull Request
        if (data.type === 'PullRequestEvent' && data.payload.action === 'reopened') {
          return templates.reopenPullRequest(data);
        }

        // Synchronize a Pull Request
        if (data.type === 'PullRequestEvent' && data.payload.action === 'synchronized') {
          return templates.synchronizePullRequest(data);
        }

        // Comment on the Unified Diff of a Pull Request
        if (data.type === 'PullRequestReviewCommentEvent') {
          return templates.commentOnPullRequestDiff(data);
        }

        // Push to a Branch
        if (data.type === 'PushEvent') {
          return templates.pushToBranch(data);
        }

        // Create a Release
        if (data.type === 'ReleaseEvent') {
          return templates.createRelease(data);
        }

        // Star a Repository
        if (data.type === 'WatchEvent') {
          return templates.starRepository(data);
        }
      }
    };

    return {
      boilerplate: function() { return render('boilerplate'); },
      profile: function(data) { return render('profile', data); },
      event: function(data) { return render('event', data); }
    };
  })();

  function ActivityPlugin(element, options) {
    this.container   = $(element);
    this.widgetHead  = null;
    this.widgetBody  = null;
    this.settings    = $.extend({}, defaults, options);
    this.realTimeout = this.settings.timeout;

    this.init();
  }

  $.extend(ActivityPlugin.prototype, {
    /**
     * Initialize the widget and start the polling process
     */
    init: function() {
      this.container.append(views.boilerplate());

      this.widgetHead = this.container.find('.activity-head');
      this.widgetBody = this.container.find('.activity-body');

      if (this.settings.debug) {
        console.log('Settings: ', this.settings);
      }

      this.poll();
    },

    /**
     * Poll for new data
     * Uses setTimeout to recursively loop indefinitely
     */
    poll: function() {
      var obj = this;

      obj.fetchProfile()
        .then(obj.fetchActivity)
        .always(function() {
          setTimeout($.proxy(obj.poll, obj), obj.realTimeout*1000);
        });
    },

    /**
     * Get user profile and if successful update DOM
     * Returns a promise
     */
    fetchProfile: function() {
      var url, promise;

      url = 'https://api.github.com/users/'+this.settings.username;

      promise = $.ajax({
        url: url,
        headers: {'Accept': 'application/vnd.github.v3+json'},
        dataType: 'json',
        ifModified: true,
        context: this
      });

      // Success, if new data update DOM
      promise.done(function(data, status, xhr) {
        if (this.settings.debug) {
          console.log('Fetch '+url+' ('+xhr.status+' '+xhr.statusText+')');
        }

        if (data) {
          if (this.settings.debug) {
            console.log(data);
          }

          // Populate widget head
          this.widgetHead.html(views.profile(data));
        }
      });

      // Failure, ...
      promise.fail(function(xhr, status, error) {
        if (this.settings.debug) {
          console.log('Fetch '+url+' ('+xhr.status+' '+xhr.statusText+')');
        }
      });

      return promise;
    },

    /**
     * Get user activity and if successful update DOM
     * Returns a promise
     */
    fetchActivity: function() {
      var obj, master, fetch;

      obj = this;

      master = new $.Deferred();

      fetch = function(url) {
        var firstPageURL, promise;

        firstPageURL = 'https://api.github.com/users/'+obj.settings.username+'/events/public';

        if (url === undefined) {
          url = firstPageURL;
        }

        promise = $.ajax({
          url: url,
          headers: {'Accept': 'application/vnd.github.v3+json'},
          dataType: 'json',
          ifModified: true,
          context: obj
        });

        // Success, if new data update DOM
        promise.done(function(data, status, xhr) {
          var content, matches, i;

          if (this.settings.debug) {
            console.log('Fetch '+url+' ('+xhr.status+' '+xhr.statusText+')');
          }

          // Read X-Poll-Interval header to get the polite minimum timeout
          this.setRealTimeout(xhr.getResponseHeader('X-Poll-Interval'));

          if (data) {
            if (this.settings.debug) {
              console.log(data);
            }

            // Build-up content by passing the data through the view templates
            content = '';
            for (i=0; i<data.length; i++) {
              content += views.event(data[i]);
            }

            // Populate the widget body
            if (url === firstPageURL) {
              this.widgetBody.html('').append(content);
            } else {
              this.widgetBody.append(content);
            }

            // If Link header contains a "next" url then fetch it, otherwise we're done
            matches = /<(\S+)>; rel="next"/.exec(xhr.getResponseHeader('Link'));
            if (matches) {
              fetch(matches[1]);
            } else {
              master.resolve();
            }
          } else {
            // No data, probably got "304 Not Modified"
            master.resolve();
          }
        });

        // Failure, ...
        promise.fail(function(xhr, status, error) {
          if (this.settings.debug) {
            console.log('Fetch '+url+' ('+xhr.status+' '+xhr.statusText+')');
          }

          // Read X-Poll-Interval header to get the polite minimum timeout
          this.setRealTimeout(xhr.getResponseHeader('X-Poll-Interval'));

          master.reject();
        });
      };

      fetch();

      return master.promise();
    },

    /**
     * Make sure realTimeout is never less than minTimeout
     * (minTimeout comes from the X-Poll-Interval response header)
     */
    setRealTimeout: function(minTimeout) {
      this.realTimeout = this.settings.timeout;

      if (minTimeout > this.realTimeout) {
        this.realTimeout = minTimeout;
      }
    }
  });

  /* Create a jQuery plugin */
  $.fn.activity = function(options) {
    this.each(function() {
      if (!$.data(this, "activity_plugin")) {
        $.data(this, "activity_plugin", new ActivityPlugin(this, options));
      }
    });

    return this;
  };

})(jQuery, window, document);