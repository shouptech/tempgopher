function renderVersion() {
    $.ajax({
        url: jsconfig.baseurl + "/api/version"
    }).then(function(data) {
        var versionText = "TempGopher © 2018 Mike Shoup | Version: " + data.version;
        $("#version").text(versionText);
    });
};
$(document).ready(renderVersion);
