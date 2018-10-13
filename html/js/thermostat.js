// Set the auth header if necessary
function authHeaders(xhr) {
    if (window.localStorage.getItem("authtoken") !== null) {
        var authToken = window.localStorage.getItem("authtoken");
        xhr.setRequestHeader("Authorization", "Basic " + authToken);
    }
}

// Redirect if not authorized
function redirectIfNotAuthorized() {
    $.ajax({
        url: jsconfig.baseurl + "/api/version",
        beforeSend: authHeaders,
        statusCode: {
            401: function() {
                window.location.replace(jsconfig.baseurl + "/app/login.html");
            },
            403: function() {
                window.location.replace(jsconfig.baseurl + "/app/login.html");
            }
        }
    });
};

// Call the redirect function
redirectIfNotAuthorized();

// Display version at bottom of page
function renderVersion() {
    $.ajax({
        url: jsconfig.baseurl + "/api/version",
        beforeSend: authHeaders
    }).then(function(data) {
        var versionText = "TempGopher © 2018 Mike Shoup | Version: " + data.version;
        $("#version").text(versionText);
    });
};
$(document).ready(renderVersion);

function displayLogoutButton() {
    if (window.localStorage.getItem("authtoken") !== null) {
        // Display a logout button
        var logoutButton = $("<button>")
            .text('Logout')
            .click(function() {
                window.localStorage.removeItem("authtoken");
                window.location.replace(jsconfig.baseurl + "/app/login.html");
            });
        $("#logoutDiv").append(logoutButton);
    };
}

$(document).ready(displayLogoutButton);

function celsiusToFahrenheit(degree) {
    return degree * 1.8 + 32;
}

function fahrenheitToCelsius(degree) {
    return (degree - 32) * 5 / 9;
};

function renderThermostats() {
    $.ajax({
        url: jsconfig.baseurl + "/api/status/",
        beforeSend: authHeaders
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
                url: jsconfig.baseurl + "/api/config/sensors/" + data[key].alias,
                beforeSend: authHeaders
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

                var cmIn = $("<input>").attr("id", "cm" + configData.alias).val(configData.coolminutes).attr("size", "2").attr("pattern", rp).on('input', function(){window.clearInterval(rtHandle)});
                var htIn = $("<input>").attr("id", "ht" + configData.alias).val(hightemp).attr("size", "4").attr("pattern", rp).on('input', function(){window.clearInterval(rtHandle)});
                var hmIn = $("<input>").attr("id", "hm" + configData.alias).val(configData.heatminutes).attr("size", "2").attr("pattern", rp).on('input', function(){window.clearInterval(rtHandle)});
                var ltIn = $("<input>").attr("id", "lt" + configData.alias).val(lowtemp).attr("size", "4").attr("pattern", rp).on('input', function(){window.clearInterval(rtHandle)});

                var configp = $("<p></p>").text("Chills for ").append(cmIn).append(" minutes when &gt; ").append(htIn).append(degUnit).append($("<br>"));
                configp.append("Heats for ").append(hmIn).append(" minutes when &lt; ").append(ltIn).append(degUnit);

                var configdiv = $("<div></div>").addClass("five columns").append(configp);
                rowdiv.append(configdiv);

                var yesButton = $("<button></button>").addClass("button button-primary").text("✔").css("margin-right", "5px").click(function() {
                    if (jsconfig.fahrenheit) {
                        var newHT = fahrenheitToCelsius(parseFloat(htIn.val()));
                        var newLT = fahrenheitToCelsius(parseFloat(ltIn.val()));
                    } else {
                        var newHT = parseFloat(htIn.val());
                        var newLT = parseFloat(ltIn.val());
                    }
                    $.ajax({
                        type: "POST",
                        url: jsconfig.baseurl + "/api/config/sensors",
                        beforeSend: authHeaders,
                        data: JSON.stringify([{
                            "id": configData.id,
                            "alias": configData.alias,
                            "hightemp": newHT,
                            "lowtemp": newLT,
                            "heatgpio": configData.heatgpio,
                            "heatinvert": configData.heatInvert,
                            "heatminutes": parseFloat(hmIn.val()),
                            "coolgpio": configData.coolgpio,
                            "coolinvert": configData.coolinvert,
                            "coolminutes": parseFloat(cmIn.val()),
                            "verbose": configData.verbose
                        }])
                    })
                    window.setInterval(renderThermostats, 60000);
                    renderThermostats();
                });

                var noButton = $("<button></button>").addClass("button").text("✘").click(function() {
                    window.setInterval(renderThermostats, 60000);
                    renderThermostats();
                });

                var buttonDiv = $("<div></div>").addClass("three columns").append(yesButton).append(noButton);
                rowdiv.append(buttonDiv);
                //var confForm = $("<form></form>").append(rowdiv);

                // Add things back to the thermostat list
                $("#thermostats").append(titlediv);
                $("#thermostats").append(rowdiv);
            });
        };
    });
};

$(document).ready(renderThermostats);
var rtHandle = window.setInterval(renderThermostats, 60000);
