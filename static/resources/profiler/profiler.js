/* Intika - Starhub - 2019*/

function tabClickSwitch (){
	var activateMeTab 	  = 'tabTitle1';
	var activateMeDiv     = 'tabDiv1';
	var desactivateMeTab  = 'tabTitle2';
	var desactivateMeDiv  = 'tabDiv2';
	document.getElementById(activateMeTab).classList.toggle('active');
	document.getElementById(activateMeDiv).classList.toggle('active');
	document.getElementById(desactivateMeTab).classList.toggle('active');
	document.getElementById(desactivateMeDiv).classList.toggle('active');
}

function showDiv(div){
	document.getElementById(div).style.height = "";
	document.getElementById(div).style.padding = "";
	document.getElementById(div).style.margin = "";
	document.getElementById(div).style.visibility = "visible";
}

function hideDiv(div){
	document.getElementById(div).style.height = "0px";
	document.getElementById(div).style.padding = "0px";
	document.getElementById(div).style.margin = "0px";
	document.getElementById(div).style.visibility = "hidden";
}

const myNotification = window.createNotification({
	closeOnClick: true,
	displayCloseButton: false,
	//nfc-top-right
	//nfc-top-left
	//nfc-bottom-right
	//nfc-bottom-left
	positionClass: 'nfc-top-left',
	onclick: false,
	showDuration:5000,
	// success, info, warning, error, and none
	theme: 'info'
});

function checkUserIDandProceed (login) {
	getIdFor(login, function(username, id) {
		//alert(id);
		//alert(login);
		loadGlobal(login);
	})
}

function getIdFor (username, callback) {
	$.getJSON('https://api.github.com/users/' + username + "?callback=?", {},
	function(json) {
		var id = json["data"]["id"]
		if (id) {callback(username, id)} else {
			//$('#div2').hide();
			//$('#div5').hide();
			showDiv("titleBox");
			showDiv("searchBox");
			myNotification({message: 'Unknown github user, or limit reached',});
		}
	});
}

var usernameProf = window.location.pathname;
usernameProf = usernameProf[0] == '/' ? usernameProf.substr(1) : usernameProf;

$(document).ready(function() {

	if ((usernameProf) && (usernameProf != 'profile')) {
		checkUserIDandProceed(usernameProf);
	} else {
		//$('#div2').hide();
		//$('#div5').hide();
		showDiv("titleBox");
		showDiv("searchBox");
	}
	
});

function profile() {
	var userInput = $("#username").val();
	if (userInput) {
		checkUserIDandProceed(userInput);
	}
}

function loadScript(src){
	var script = document.createElement('script');
	script.onload = function () {
	};
	script.src = src;
	document.head.appendChild(script);
}

function loadGlobal(user) {

	showDiv("activityDiv");
	//showDiv("div2");
	//showDiv("div5");
	//$('#div2').show();
	//$('#div5').show();
	showDiv("div3");
	showDiv("div4");
	showDiv("div6");
	document.getElementById("activityDiv").style.height = "auto";
	document.getElementById("activityDiv").style.visibility = "visible";
	
	document.getElementById("github-language-distribution").innerHTML = "";
	document.getElementById("github-contributions").innerHTML = "<div id='divGitStatLoader' class='ui active centered inline loader' style='height:0px; visibility:hidden;'></div>";
	document.getElementById("gitWid").innerHTML = "";
	
	document.getElementById("grapHref").href = "https://github.com/" + user;
	document.getElementById("ghuserLink").href = "https://ghuser.io/" + user;

	document.getElementById("iframeProfile").height = "350px";
	document.getElementById("iframeProfile").src = "static/resources/github-id/?q=" + user;	
	document.getElementById("divGitStatLoader").style.height = "";
	document.getElementById("divGitStatLoader").style.visibility = "visible";
	
	document.getElementById("gitWid").innerHTML = "<div class='github-widget' style='width:550px; margin:auto;' data-username='" + user + "'></div>";

	loadScript("static/resources/widgets/widget.js");
	
	setTimeout(privLoadActi(pageBrowsePriv, '2', pageFilterPriv, user), 2000);
	setTimeout(pubLoadActi(pageBrowsePub, '1', pageFilterPub, user), 2000);
		
	(async () => {	
		const GITHUB_USERNAME = user;
		const COMMITS_CONTAINER = '#github-contributions';
		const LANGUAGES_CONTAINER = '#github-language-distribution';
		
		const githubStats = await GithubStats(GITHUB_USERNAME);
		
		let githubCommits = document.querySelector(COMMITS_CONTAINER);
		/* Render SVG for commit contributions */
		let commitsContribSVG = githubStats.commitsContribSVG({
		   rows: 7,
		   space: 2,
		   rectWidth: 8,
		levelColors: [
		    {
		        minCommits: 0,
		        color: '#ebedf0'
		    },
		    {
		        minCommits: 1,
		        color: '#c6e48b'
		    },
		    {
		        minCommits: 9,
		        color: '#7bc96f'
		    },
		    {
		        minCommits: 17,
		        color: '#239a3b'
		    },
		    {
		        minCommits: 26,
		        color: '#196127'
		    }
		]
		});
		githubCommits.appendChild(commitsContribSVG);
		
		    let githubLanguageDistribution = document.querySelector(LANGUAGES_CONTAINER);
		    /* Render SVG for language contributions */
		    let languageContribSVG = githubStats.languagesContribSVG({
		        barHeight: 20,
		        barWidth: githubLanguageDistribution.offsetWidth/2,
		        lineSpacing: 4,
		        languageNameWidth: 100,
		        fontSize: 14
		    });
		    githubLanguageDistribution.appendChild(languageContribSVG);
		    document.getElementById("divGitStatLoader").style.height = "0px";
			document.getElementById("divGitStatLoader").style.visibility = "hidden";
	})();
}
