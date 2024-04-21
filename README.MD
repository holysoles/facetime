# holysoles-facetime

A Go webserver that wraps the FaceTime application on macOS for interoperability using AppleScript.

## Overview

In macOS 12, Apple [introduced the ability](https://support.apple.com/guide/facetime/join-a-call-from-a-facetime-link-fctm2cd42547/5.0/mac/13.0) to create a link to a FaceTime call, that both others who have Apple devices can join in the native FaceTime apps, or other users via Chromium web browsers. Since this functionality also works on mobile, those who own a Mac but not an iPhone aren't able to initiate FaceTime calls to friends and family.

This runs on a Mac to wrap the native FaceTime application and its interaction with FaceTime links, using AppleScript and its ability to interact with applications.

## Setup

- Support has only been verified on macOS Ventura (13.6+) but may work on others.
- After downloading, ensure you set execute permissions on the executable: `chmod  u=rwx ./facetime_macos_arm64`
- the environment variable `TRUSTED_PROXIES` should be set with the IP addresses/ranges that should be treated as proxy sources for Gin.

## API Call Examples

The below details API calls that can be used to initialize calls from a http client.

### Get webserver status information
```bash
curl --location 'http://localhost:8080/status'
```
### Get list of links for unused/upcoming FaceTime Calls
```bash
curl --location 'http://localhost:8080/link'
```
### Request a new line for a FaceTime Call
```bash
curl --location --request POST 'http://localhost:8080/link/new'
```
### Have the server start a call
Note: User(s) must have already requested to join the call.
```bash
curl --location --request POST 'http://localhost:8080/link/join'
```
### Have the server user leave the call
```bash
curl --location --request POST 'http://localhost:8080/link/leave'
```
### Delete a call from the list
```bash
curl --location --request DELETE 'http://localhost:8080/link' \
--header 'Content-Type: application/json' \
--data '{"link":"https://facetime.apple.com/join#v=1\u0026p=1bv0//lFFe6uYY7be8h5vw\u0026k=fkA7tz4FRASGFLc5345_pkFJ_rCQJAdJH6F8z2OGrIZ"}'
```

## Contributing
- Commits should follow [Angular Commit Guidelines](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines) for the Github Tag workflow's parsing.

## Credits
Special thanks to the following projects for providing essential resources:
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [UI Browser v3](https://latenightsw.com/freeware/ui-browser/)

And a thanks to the following projects for inspiration:
- [mautrix/imessage](https://github.com/mautrix/imessage)
