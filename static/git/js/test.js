show('<img style=" height:15px;" src="http://img.blog.csdn.net/20170122091917400" />&nbsp;加载中，请稍后');
    get("https://raw.githubusercontent.com/racaljk/hosts/master/hosts");

    function show(msg) {
        document.getElementById("loadding").innerHTML = msg;
    }

    function copyUrl2() {
        var Url2 = document.getElementById("ta");
        Url2.select();
        document.execCommand("Copy");
        show("复制成功，赶紧去保存HOSTS吧~");
    }

    function ajaxObject() {
        var xmlHttp;
        try {
            // Firefox, Opera 8.0+, Safari
            xmlHttp = new XMLHttpRequest();
        } catch(e) {
            // IE
            try {
                xmlHttp = new ActiveXObject("Msxml2.XMLHTTP");
            } catch(e) {
                try {
                    xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
                } catch(e) {
                    show('加载失败，换个浏览器试试吧~');
                    return false;
                }
            }
        }
        return xmlHttp;
    }

    function get(url) {
        var ajax = ajaxObject();
        ajax.open("get", url, true);
        ajax.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        ajax.onreadystatechange = function() {
            if (ajax.readyState == 4) {
                if (ajax.status == 200) {
                    document.getElementById('ta').value = ajax.responseText;
                    show('<input type="button" onclick="copyUrl2();" value="复制"/>');
                } else {
                    show('加载失败，刷新试试吧~');
                }
            } else {

			}
        }
        ajax.send();

    }


