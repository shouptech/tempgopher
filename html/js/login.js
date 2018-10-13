// Retrieve username and password from fields, store a token, and redirect to main app
function processLogin() {
    var username = $("#loginName").val();
    var password = $("#loginPassword").val();
    window.localStorage.setItem("authtoken", btoa(username + ":" + password));
    window.location.replace(jsconfig.baseurl + "/app/");
};

// If the login page is displayed, we need to remvoe any existing auth tokens
window.localStorage.removeItem("authtoken");
