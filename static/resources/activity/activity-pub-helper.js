/*! Starhub Activity Stream - v2.1.0 * Copyright (c) 2019 Intika - starhub.be * MIT License */

var pageBrowsePub=1;
var pageFilterPub=0;
var pageUsernamePub='';

function pubfilterList (who) {	
	pageFilterPub=who;
	var pubtextDescId 	= 'pubtextFilter' + who;
	var pubtext 		= document.getElementById(pubtextDescId).innerHTML;
	document.getElementById('pubfkingFilterText').innerHTML = pubtext;
	pubLoadActi(pageBrowsePub, '1', pageFilterPub, pageUsernamePub);
	//Warning the inner div trigger showFker() a second time because of the onclick
}

function pubintelliHideFker() {
	if (document.getElementById('pubfkingMenuList').clientHeight == '0' || document.getElementById('pubfkingMenuList').style.visibility == 'hidden') {
		document.getElementById('pubfkingMenuList').style.visibility = 'hidden';
	} else {
		setTimeout(pubintelliHideFker, 500);
	}
}

function pubhideFker() {
	document.getElementById('pubfkingMenuList').style.visibility = 'hidden';
}

//Warning this is triggered on menu click and sub click 
function pubshowFker() {
	if (document.getElementById('pubfkingMenuList').style.visibility == 'hidden') {
		document.getElementById('pubfkingMenuList').style.visibility = 'visible';
		pubintelliHideFker();
	} else {
		pubhideFker();
	}
}

function loadPreviousPub() {
	if (pageBrowsePub < 2 ) {
		document.getElementById('pubPrevButton').classList.remove('primary');
	} else {
		pageBrowsePub--;
		document.getElementById('pubPrevButton').classList.add('primary');
		document.getElementById('pubNextButton').classList.add('primary');
		pubLoadActi(pageBrowsePub, '1', pageFilterPub, pageUsernamePub);
	}
}

function loadNextPub() {
	if (pageBrowsePub > 8 ) {
		document.getElementById('pubNextButton').classList.remove('primary');
	} else {
		pageBrowsePub++;
		document.getElementById('pubNextButton').classList.add('primary');
		document.getElementById('pubPrevButton').classList.add('primary');
		pubLoadActi(pageBrowsePub, '1', pageFilterPub, pageUsernamePub);
	}
}

function pubLoadActi(page, type, filter, usernameParam) {
	if (usernameParam != '') {
		pageUsernamePub = usernameParam;
		var urlSrc="static/resources/activity/activity.html?u=" + pageUsernamePub + "&p=" + page + "&t=" + type +"&f=" + filter;
		document.getElementById("pIframeContainerActPub").style.height = "";
		document.getElementById("divIframeActPubLoader").style.height = "";
		document.getElementById("divIframeActPubLoader").style.marginTop = "40px";
		document.getElementById("divIframeActPubLoader").style.visibility = "visible";
		document.getElementById("iframeActPub").style.visibility = "visible";
		document.getElementById("iframeActPub").src=urlSrc;
		document.getElementById("iframeActPub").onload = function(){
			document.getElementById("iframeActPub").height = "610px";
			document.getElementById("iframeActPub").style.height = "610px";
			document.getElementById("divIframeActPubLoader").style.height = "0px";
			document.getElementById("divIframeActPubLoader").style.marginTop = "0px";
			document.getElementById("divIframeActPubLoader").style.visibility = "hidden";
		};	
	}
}
