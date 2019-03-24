/*! Starhub Activity Stream - v2.1.0 * Copyright (c) 2019 Intika - starhub.be * MIT License */

var pageBrowsePriv=1;
var pageFilterPriv=0;
var pageUsernamePriv='';

function privfilterList (who) {	
	pageFilterPriv=who;
	var privtextDescId 	= 'privtextFilter' + who;
	var privtext 		= document.getElementById(privtextDescId).innerHTML;
	document.getElementById('privfkingFilterText').innerHTML = privtext;
	privLoadActi(pageBrowsePriv, '2', pageFilterPriv, pageUsernamePriv);
	//Warning the inner div trigger showFker() a second time because of the onclick
}

function privintelliHideFker() {
	if (document.getElementById('privfkingMenuList').clientHeight == '0' || document.getElementById('privfkingMenuList').style.visibility == 'hidden') {
		document.getElementById('privfkingMenuList').style.visibility = 'hidden';
	} else {
		setTimeout(privintelliHideFker, 500);
	}
}

function privhideFker() {
	document.getElementById('privfkingMenuList').style.visibility = 'hidden';
}

//Warning this is triggered on menu click and sub click 
function privshowFker() {
	if (document.getElementById('privfkingMenuList').style.visibility == 'hidden') {
		document.getElementById('privfkingMenuList').style.visibility = 'visible';
		privintelliHideFker();
	} else {
		privhideFker();
	}
}

function loadPreviousPriv() {
	if (pageBrowsePriv < 2 ) {
		document.getElementById('privPrevButton').classList.remove('primary');
	} else {
		pageBrowsePriv--;
		document.getElementById('privPrevButton').classList.add('primary');
		document.getElementById('privNextButton').classList.add('primary');
		privLoadActi(pageBrowsePriv, '2', pageFilterPriv, pageUsernamePriv);
	}
}

function loadNextPriv() {
	if (pageBrowsePriv > 8 ) {
		document.getElementById('privNextButton').classList.remove('primary');
	} else {
		pageBrowsePriv++;
		document.getElementById('privNextButton').classList.add('primary');
		document.getElementById('privPrevButton').classList.add('primary');
		privLoadActi(pageBrowsePriv, '2', pageFilterPriv, pageUsernamePriv);
	}
}

function privLoadActi(page, type, filter, usernameParam) {
	if (usernameParam != '') {
		pageUsernamePriv = usernameParam;
		var urlSrc="static/resources/activity/activity.html?u=" + pageUsernamePriv + "&p=" + page + "&t=" + type +"&f=" + filter;
		document.getElementById("pIframeContainerActPriv").style.height = "";
		document.getElementById("divIframeActPrivLoader").style.height = "";
		document.getElementById("divIframeActPrivLoader").style.marginTop = "40px";
		document.getElementById("divIframeActPrivLoader").style.visibility = "visible";
		document.getElementById("iframeActPriv").style.visibility = "visible";
		document.getElementById("iframeActPriv").src=urlSrc;
		document.getElementById("iframeActPriv").onload = function(){
			document.getElementById("iframeActPriv").height = "610px";
			document.getElementById("iframeActPriv").style.height = "610px";
			document.getElementById("divIframeActPrivLoader").style.height = "0px";
			document.getElementById("divIframeActPrivLoader").style.marginTop = "0px";
			document.getElementById("divIframeActPrivLoader").style.visibility = "hidden";
		};	
	}
}
