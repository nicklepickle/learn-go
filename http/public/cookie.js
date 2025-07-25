const Cookie = {
    setCookie: function (key, value) {
        let d = new Date();
        let days= 365;
        d.setTime(d.getTime() + (days*24*60*60*1000));
        document.cookie = (key + '=' + btoa(value) + '; expires=' + d.toUTCString() + "; path=/");

    },
    getCookie: function (key) {
        var cookies = document.cookie.split(';');
        var val = null;
        for(var i=0; i < cookies.length;i++) {
            var c = cookies[i];
            while (c.charAt(0)==' ') {
                c = c.substring(1,c.length);
            }
            if (c.indexOf(key + '=') == 0) {
                val = c.substring(key.length + 1,c.length);
                break;
            }
        }
        return val;                  
    }
};

//export default Cookie;