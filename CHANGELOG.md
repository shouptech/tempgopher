# TempGopher Changelog

## 0.4.0

Release 2018-11-01

* Enable/disable heating or cooling via UI
* Improved unit tests!
* Fixed bug where updating/clearing changes would cause multiple refreshes in a row [#17](https://gitlab.com/shouptech/tempgopher/issues/17)

## 0.3.1

Release: 2018-10-24

* Fixes bug where UI would not display two sensors correctly [#15](https://gitlab.com/shouptech/tempgopher/issues/15)
* Fixes a typo during the CLI configuration [#14](https://gitlab.com/shouptech/tempgopher/issues/14)

## 0.3.0

Release: 2018-10-20

* You can now supply a list of users for simple authentication
* Will write data to an Influx DB if configured
* Adds the ability to selectively disable heating or cooling
* Checks for the existence of a config file before generating a new one

## 0.2.0

Release: 2018-10-11

* You can now update the configuration using the UI or API
* A script to help install has been added
* The binary will now generate a usable configuration

## 0.1.1

Release: 2018-10-07

* Changes temperature logic. See #8. Fixes a situation where temperature 'floats' at the threshold and the switch is rapidly cycled.

## 0.1.0

Released: 2018-10-04

* Added single page Web UI, packaged into app
* Moved to GitLab, uses GitLab CI to run builds

## 0.0.1

Released: 2018-10-02

* Initial Release
* API only
* Supports multiple thermostats
