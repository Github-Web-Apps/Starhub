var MAX_SUPPORTED_PAGES_NO_AUTH = 30;


$(document).ready(function() {
	$('#plot-container').hide();
	$('#url-container').hide();

	$('#repo').keyup(function(event) {
		if (event.keyCode == 13) {
			go();
		}
	});

	$(window).scroll(function() {
		if ($(this).scrollTop() > 100) {
			$('.scrollToBottom').fadeOut();
		} else {
			$('.scrollToBottom').fadeIn();
		}
	});

	$('.scrollToBottom').click(function() {
		$('html, body').animate({scrollTop : $(document).height()-$(window).height() });
		return false;
	});

	$('#userandpass_input :input').attr('disabled', false);
	$('#token_input :input').attr('disabled', true);
	
	$('input[type=radio][name=auth_group]').change(function() {
		if (this.value == 'userandpass') {
			$('#userandpass_input :input').attr('disabled', false);
			$('#token_input :input').attr('disabled', true);
		}
		else if (this.value == 'token') {
			$('#userandpass_input :input').attr('disabled', true);
			$('#token_input :input').attr('disabled', false);
		}
	});
	
	NProgress.configure({ 
		showSpinner: false,
		minimum: 0.0001,
		trickle: false
	});
	
	parseUrlParams();

	initAuthDetails();
});

function showPlot(data, user, repo) {
	var plot = document.getElementById('plot-container');
	if ($('#add_or_replace').val() == "replace") {
		Plotly.newPlot( plot, [{
		x: data['xaxis'],
		y: data['yaxis'],
		name: user + '/' + repo	}], {
		margin: { t: 0 } } );
	}
	else
	{
		Plotly.plot( plot, [{
		x: data['xaxis'],
		y: data['yaxis'],
		name: user + '/' + repo }], {
		margin: { t: 0 } } );
	}
}

function showStats(stats, user, repo) {
	$('#divPub').show();
    
	var tableHTML = '<table><thead><th colspan="2">Stats for ' + user + '/' + repo + '</th></thead><tbody>';
	$.each(stats, function(i, item) {
		tableHTML += '<tr><td>' + item['text'] + '</td><td>' + item['data'] + '</td></tr>';
	});
	if ($('#add_or_replace').val() == "replace") {
		$('#stats-table tbody').remove();
	}

	tableHTML += '</tbody></table>';
	if ($('#add_or_replace').val() == "replace") {
		$('#stats-container').empty();
	}
	
	$('#stats-container').append(tableHTML);
}

function buildUrl(user, repo) {
	var base_url = '';
	var cur_repo = 'u=' + user + '&r=' + repo;
	if ($('#add_or_replace').val() == "replace") {
		base_url = window.location.href.split('?')[0] + '?';
	}
	else {
		base_url = $('#url-box').val() + '&';
	}

	$('#url-box').val(base_url + cur_repo);
}

function calcStats(data) {
	if (data['xaxis'].length == 0)
		return;

	dateFirst = new Date(data['xaxis'][0]);
	dateLast = new Date(data['xaxis'][data['xaxis'].length-1]);
	var numOfDays = Math.floor((new Date(dateLast - dateFirst))/1000/60/60/24);

	result = [];

	result.push({
		'text': 'Number of stars',
		'data': data['yaxis'][data['yaxis'].length-1]
		});

	result.push({
		'text': 'Number of days',
		'data': numOfDays
		});

	result.push({
		'text': 'Average stars per day',
		'data': (data['yaxis'].length / numOfDays).toFixed(3)
		});

	result.push({
		'text': 'Average days per star',
		'data': (numOfDays / data['yaxis'].length).toFixed(3)
		});

	var daysWithoutStars = 0;
	var maxStarsPerDay = 0;
	var curSameDays = 1;
	var startDate = Math.floor(new Date(0)/1000/60/60/24);
	var prevDate = startDate;
	$.each(data['xaxis'], function(i, date) {
		curDate = Math.floor(new Date(date)/1000/60/60/24);

		if (curDate == prevDate) {
			curSameDays += 1;
		}
		else {
			if (prevDate != startDate) {
				daysWithoutStars += curDate - prevDate - 1;
			}

			if (curSameDays > maxStarsPerDay) {
				maxStarsPerDay = curSameDays;
			}

			curSameDays = 1;
		}

		prevDate = curDate;
	});

	if (curSameDays > maxStarsPerDay) {
		maxStarsPerDay = curSameDays;
	}

	result.push({
		'text': 'Max stars in one day',
		'data': maxStarsPerDay
		});


	result.push({
		'text': 'Days with stars',
		'data': numOfDays - daysWithoutStars
		});

	result.push({
		'text': 'Days with no stars',
		'data': daysWithoutStars
		});

	return result;
}

function buildData(jsonData) {
	var starCount = 0;
	var xaxis = [];
	var yaxis = [];
	for (key in jsonData) {
		starCount = starCount + 1;
		xaxis.push(jsonData[key]['starred_at']);
		yaxis.push(starCount);
	};

	return {
		'xaxis': xaxis,
		'yaxis': yaxis
		};
}

function startLoading() {
	$('#gobtn').hide();
	$('#loading').show();
}

function stopLoading() {
	$('#gobtn').show();
	$('#loading').hide();
}

function finishLoading() {
  stopLoading();
  $('#plot-container').show();
  $('#url-container').show();
  $('#add_or_replace').prop('disabled', false);
  $('.scrollToBottom').fadeIn();
}

function findLastPage(linkHeader) {
	if (linkHeader == null || linkHeader.length == 0) {
		return 0;
	}

	// Split parts by comma
	var parts = linkHeader.split(',');

	// Parse each part into a named link
	for (i in parts) {
		var section = parts[i].split(';');
		if (section.length != 2) {
			continue;
		}

		var url = section[0].replace(/<(.*)>/, '$1').trim();
		var name = section[1].replace(/rel="(.*)"/, '$1').trim();

		// if name is 'last' then extract page and return it
		if (name == 'last') {
			return url.replace(/(.*)&page=(.*)/, '$2').trim();
		}
	}
}


var stargazersData = [];
var done = false;
var loadError = false;
var lastPage = 0;

function checkGithubAuth(userpass, token, callback_on_success, callback_on_fail) {
	url = 'https://api.github.com';
	
	var auth_success = false;
	
	$.ajax({
		beforeSend: function(request) {
						if (token != undefined) {
							var auth_string = 'token ' + token;
							request.setRequestHeader('Authorization', auth_string);
						}
						else if (userpass != undefined) {
							var auth_string = 'Basic ' + userpass;
							request.setRequestHeader('Authorization', auth_string);
						}
					},
		url: url,
		success: function(data, textStatus, xhr) {
			if (xhr.status == 200) {
				auth_success = true;
				callback_on_success();
			} else {
				callback_on_fail();
			}
		},
		
		complete: function(xhr, textStatus) {
			if (auth_success == false && xhr.status != 200) {
				callback_on_fail();
			}
		} 
	});
}

function showMessageBox(message, title, callback) {
	dialog = $('<div>' + message + '</div>').dialog({
		modal: true,
		title: title,
		buttons: {
			Ok: function() {
				dialog.dialog('close');
				callback();
			}
		}
	});
}

function openGithubAuthDialog(succes_auth_callback) {
    dialog = $('#github_auth_dialog').dialog({
		modal: true,
		height: 580,
		width: 500,
		buttons: {
			Ok: function() {
				var auth_option = $('input[name=auth_group]:checked').attr('value');
				if (auth_option == 'userandpass' && $('#github_username').val() != '' && $('#github_password').val() != '') {
					getGitHubAuthObject().setUserAndPass($('#github_username').val().trim(), $('#github_password').val().trim());
					checkGithubAuth(
						getGitHubAuthObject().getUserAndPass(), 
						null, 
						function() {
							getGitHubAuthObject().useLocalStorage($('#save_in_local_storage_chkbox').is(':checked'));
							updateGithubAuthUI();
							if (succes_auth_callback != undefined) {
								succes_auth_callback();
							}
						},
						function() {
							getGitHubAuthObject().cleanUserAndPass();
							showMessageBox(
								'GitHub authentication failed, please try again',
								'Error',
								function() {
									openGithubAuthDialog(succes_auth_callback);
								});
						});
				}
				else if (auth_option == 'token' && $('#github_token').val() != '') {
					getGitHubAuthObject().setAccessToken($('#github_token').val().trim());
					checkGithubAuth(
						null, 
						getGitHubAuthObject().getAccessToken(),
						function() {
							getGitHubAuthObject().useLocalStorage($('#save_in_local_storage_chkbox').is(':checked'));
							updateGithubAuthUI();
							if (succes_auth_callback != undefined) {
								succes_auth_callback();
							}
						},
						function() {
							getGitHubAuthObject().cleanAccessToken();
							showMessageBox(
								'GitHub authentication failed, please try again',
								'Error',
								function() {
									openGithubAuthDialog(succes_auth_callback);
								});
						});
				}
				
				dialog.dialog('close');
			},
			Cancel: function() {
				dialog.dialog('close');
			}
		}
	});
}

function GitHubAuth() {

	var static_vars = {
		STORAGE_USERANDPASS_KEY : 'starttrack-js-userandpass',
		STORAGE_TOKEN_KEY : 'starttrack-js-token',
	};

	var vars = {
		save_in_local_storage : false,
	};

	var getStorage = function() {
		if (vars.save_in_local_storage) {
			return localStorage;
		}
		else {
			return sessionStorage;
		}
	}

	this.construct = function() {
		var userandpass = localStorage.getItem(static_vars.STORAGE_USERANDPASS_KEY);
		var token = localStorage.getItem(static_vars.STORAGE_TOKEN_KEY);
		if (userandpass != undefined || token != undefined) {
			vars.save_in_local_storage = true;
		}
	};

	this.getUserAndPass = function() {
		return getStorage().getItem(static_vars.STORAGE_USERANDPASS_KEY);
	};

	this.getUser = function() {
		var userandpass = this.getUserAndPass();
		if (userandpass == undefined) {
			return;
		}

		return window.atob(userandpass).split(':')[0];
	}

	this.setUserAndPass = function(user, passwd) {
		if (user == undefined || passwd == undefined) {
			return;
		}
		var userandpass = window.btoa(user.trim() + ':' + passwd.trim());
		getStorage().setItem(static_vars.STORAGE_USERANDPASS_KEY, userandpass);
		this.cleanAccessToken();
	}

	this.setUserAndPassOneWord = function(userandpass) {
		if (userandpass == undefined) {
			return;
		}

		getStorage().setItem(static_vars.STORAGE_USERANDPASS_KEY, userandpass);
		this.cleanAccessToken();
	}

	this.cleanUserAndPass = function() {
		getStorage().removeItem(static_vars.STORAGE_USERANDPASS_KEY);
	}

	this.getAccessToken = function() {
		if (getStorage().getItem(static_vars.STORAGE_TOKEN_KEY) != undefined) {
			return window.atob(getStorage().getItem(static_vars.STORAGE_TOKEN_KEY));
		} else {
			return undefined;
		}
	};

	this.setAccessToken = function(token) {
		if (token == undefined) {
			return;
		}
		var token_btoa = window.btoa(token);
		getStorage().setItem(static_vars.STORAGE_TOKEN_KEY, token_btoa);
		this.cleanUserAndPass();
	}

	this.cleanAccessToken = function() {
		getStorage().removeItem(static_vars.STORAGE_TOKEN_KEY);
	}

	this.cleanStorage = function() {
		this.cleanUserAndPass();
		this.cleanAccessToken();
	}

	this.useLocalStorage = function(val) {
		if (vars.save_in_local_storage != val) {
			var userandpass = this.getUserAndPass();
			var token = this.getAccessToken();
			this.cleanStorage();
			vars.save_in_local_storage = val;
			this.setUserAndPassOneWord(userandpass);
			this.setAccessToken(token);
		}

		vars.save_in_local_storage = val;
	}

	this.construct();
}

var gitHubAuth = undefined;

function getGitHubAuthObject() {
	if (gitHubAuth == undefined) {
		gitHubAuth = new GitHubAuth();
	}

	return gitHubAuth;
}

function initAuthDetails() {
	getGitHubAuthObject();
	updateGithubAuthUI();
}

function updateGithubAuthUI() {
	if (getGitHubAuthObject().getAccessToken() != undefined) {
		$('#github-auth').text('GitHub Access Token: ' + getGitHubAuthObject().getAccessToken().trim().substring(0, 6));
		$('#logout-btn').show();
	} else if (getGitHubAuthObject().getUserAndPass() != undefined) {
		$('#github-auth').text('GitHub User: ' + getGitHubAuthObject().getUser().trim());
		$('#logout-btn').show();
	} else {
		$('#github-auth').text('Authentication');
		$('#logout-btn').hide();
	}
}

function removeGithubAuth() {
	getGitHubAuthObject().cleanStorage();
	updateGithubAuthUI();
}

function loadStargazers(user, repo, on_complete_callback, cur) {
    var stargazersURL = "https://api.github.com/repos/{user}/{repo}/stargazers?per_page=100&page={page}";
	if (typeof(cur) == 'undefined') {
		cur = 1;
		stargazersData = [];
		done = false;
		loadError = false;
	}

	if (done == false) {
		startLoading();
		url = stargazersURL.replace('{page}', cur).replace('{user}', user).replace('{repo}', repo);
		//console.log(url);
		$.ajax({
			beforeSend: function(request) {
							request.setRequestHeader('Accept', 'application/vnd.github.v3.star+json');
							if (getGitHubAuthObject().getAccessToken() != undefined) {
								var auth_string = 'token ' + getGitHubAuthObject().getAccessToken();
								request.setRequestHeader('Authorization', auth_string);
							}
							else if (getGitHubAuthObject().getUserAndPass() != undefined) {
								var auth_string = 'Basic ' + getGitHubAuthObject().getUserAndPass();
								request.setRequestHeader('Authorization', auth_string);
							}
						},
			datatype: 'json',
			url: url,
			success: function(data, textStatus, request) {
						if ($.isEmptyObject(data) == true) {
							done = true;
							return;
						}

						if (cur == 1) {
							lastPage = findLastPage(request.getResponseHeader('Link'));
							if (lastPage > MAX_SUPPORTED_PAGES_NO_AUTH && getGitHubAuthObject().getAccessToken() == undefined && getGitHubAuthObject().getUserAndPass() == undefined) {
								done = true;
								loadError = true;
								stopLoading();
								openGithubAuthDialog(function() {
									loadStargazers(user, repo, on_complete_callback);
								});
								return;
							}
							NProgress.start();
						}

						stargazersData = $.merge(stargazersData, data);
						NProgress.inc(1/lastPage);
					},
			error: function(xhr, ajaxContext, thrownError) {
						if (xhr.responseText != 'undefined') {
							var responseText = $.parseJSON(xhr.responseText);
							showMessageBox('Error occured: '+ responseText['message'], 'Error');
						}
						else {
							showMessageBox('Error occured: ' + thrownError, 'Error');
						}

						stopLoading();
					}
			}).done(function() {
					if (loadError == true) {
						return;
					}
					cur = cur + 1;
					loadStargazers(user, repo, on_complete_callback, cur);
				});
	}
	else if (loadError == false) {
		finishLoading();
		NProgress.done(true);
		xyData = buildData(stargazersData);
		showStats(calcStats(xyData), user, repo);
		showPlot(xyData, user, repo);
		buildUrl(user, repo);
		if (on_complete_callback != undefined)
			on_complete_callback();
	}
}

function loadRepos(repos) {
	if (repos == undefined || repos.length == 0)
		return;

	var pair = repos.shift();
	$('#user').val(pair[0]);
	$('#repo').val(pair[1]);
	go(function() {
		$('#add_or_replace').val('add');
		loadRepos(repos);
	});
}

function parseUrlParams() {
	var params = window.location.href.split('?')[1];
	if (params == undefined)
		return;

	params = params.split('&');
	if (params == undefined)
		return;

	var repos = [];
	var user = undefined;
	var repo = undefined;
	var turn = 'user';
	for (i in params) {
		if (turn == 'user' && params[i].startsWith('u=')) {
			user = params[i].substring(2);
			turn = 'repo';
		}
		else if (turn == 'repo' && params[i].startsWith('r=')) {
			repo = params[i].substring(2);
			turn = 'user';
			repos.push([user, repo]);
		}
		else {
			showMessageBox('Wrong URL parameter: ' + params[i], 'Error');
			return;
		}

	}

	loadRepos(repos);
}

function go(on_complete_callback) {
	if ($('#user').val() == "" || $('#repo').val() == "") {
		showMessageBox('Please enter GitHub username and GitHub repository', 'Error');
	return;
	}

	$('#user').val($('#user').val().trim());
	$('#repo').val($('#repo').val().trim());

	loadStargazers($('#user').val(), $('#repo').val(), on_complete_callback);
}
