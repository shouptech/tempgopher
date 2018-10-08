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
            var statusdiv = $("<div></div>").addClass("two columns").append(statusp);
            rowdiv.append(statusdiv);

            // Display sensor config
            $.ajax({
                url: jsconfig.baseurl + "/api/config/sensors/" + data[key].alias
            }).then(function(configData){
                if (jsconfig.fahrenheit) {
                    var degUnit = "°F";
                    var hightemp = celsiusToFahrenheit(parseFloat(configData.hightemp)).toFixed(1);
                    var lowtemp = celsiusToFahrenheit(parseFloat(configData.lowtemp)).toFixed(1);
                } else {
                    var hightemp = parseFloat(configData.hightemp).toFixed(1);
                    var lowtemp = parseFloat(configData.lowtemp).toFixed(1);
                }

                rp = '[0-9]+(\.[0-9]+)?'

                configText = "Chills for " +
                    "<input id=\"cm" + configData.alias + "\" value=\"" + configData.coolminutes + "\" size=\"2\" pattern=\"" + rp +"\"> minutes when &gt; " +
                    "<input id=\"ht" + configData.alias + "\" value=\"" + hightemp + "\" size=\"4\" pattern=\"" + rp +"\">" + degUnit + ".<br />";

                configText += "Heats for " +
                    "<input id=\"hm" + configData.alias + "\" value=\"" + configData.heatminutes + "\" size=\"2\" pattern=\"" + rp +"\"> minutes when &lt; " +
                    "<input id=\"lt" + configData.alias + "\" value=\"" + lowtemp + "\" size=\"4\" pattern=\"" + rp +"\">" + degUnit + ".<br />";

                var configp = $("<p></p>").html(configText);
                var configdiv = $("<div></div>").addClass("five columns").append(configp);
                rowdiv.append(configdiv);

                var yesButton = $("<button></button>").attr("type", "submit").addClass("button button-primary").text("✔").css("margin-right", "5px");
                var noButton = $("<button></button>").addClass("button").text("✘");
                var buttonDiv = $("<div></div>").addClass("three columns").append(yesButton).append(noButton);
                rowdiv.append(buttonDiv);
            });

            // Add things back to the thermostat list
            $("#thermostats").append(titlediv);
            $("#thermostats").append($("<form></form>").append(rowdiv));
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
