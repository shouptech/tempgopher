function celsiusToFahrenheit(degree) {
    return degree * 1.8 + 32;
}

function renderThermostats() {
    $.ajax({
        url: jsconfig.baseurl + "/api/status/"
    }).then(function(data) {
        $("#thermostats").empty();
        for (var key in data) {
            // Title of thermostat
            var titleh = $("<h4></h4>").text(data[key].alias);
            var titlediv = $("<div></div>").addClass("row").append(titleh);

            // Thermostat status
            var rowdiv = $("<div></div>");
            rowdiv.addClass("row");

            // Display temperature
            if (jsconfig.fahrenheit) {
                var temp = celsiusToFahrenheit(parseFloat(data[key].temp)).toFixed(1) + "°F";
            } else {
                var temp = parseFloat(data[key].temp).toFixed(1) + "°C";
            }
            var temph = $("<h2></h2>").text(temp);
            var tempdiv = $("<div></div>").addClass("two columns").append(temph);
            rowdiv.append(tempdiv);

            // Display status
            if (data[key].cooling) {
                var statustext = "Cooling"
            } else if (data[key].heating) {
                var statustext = "Heating"
            } else {
                var statustext = "Idle"
            }
            var statusp = $("<p></p>").html(statustext);
            var statusdiv = $("<div></div>").addClass("three columns").append(statusp);
            rowdiv.append(statusdiv);

            // Display sensor config
            $.ajax({
                url: jsconfig.baseurl + "/api/config/sensors/" + data[key].alias
            }).then(function(configData){
                if (jsconfig.fahrenheit) {
                    var hightemp = celsiusToFahrenheit(parseFloat(configData.hightemp)).toFixed(1) + "°F";
                    var lowtemp = celsiusToFahrenheit(parseFloat(configData.lowtemp)).toFixed(1) + "°F";
                } else {
                    var hightemp = parseFloat(configData.hightemp).toFixed(1) + "°C";
                    var lowtemp = parseFloat(configData.lowtemp).toFixed(1) + "°C";
                }
                configText = "Chills for " + configData.coolminutes + " minutes when > " + hightemp + ".<br />";
                configText += "Heats for " + configData.heatminutes + " minutes when < " + lowtemp + ".";

                var configp = $("<p></p>").html(configText);
                var configdiv = $("<div></div>").addClass("seven columns").append(configp);
                rowdiv.append(configdiv);
            });

            // Add things back to the thermostat list
            $("#thermostats").append(titlediv);
            $("#thermostats").append(rowdiv);
        };
    });
};

function renderVersion() {
    $.ajax({
        url: jsconfig.baseurl + "/api/version"
    }).then(function(data) {
        var versionText = "TempGopher © 2018 Mike Shoup | Version: " + data.version;
        $("#version").text(versionText);
    });
};
$(document).ready(renderVersion);

$(document).ready(renderThermostats);
setInterval(renderThermostats, 60000)
