# holysoles-facetime
A Go webserver that wraps the FaceTime application on MacOS for interoperability using AppleScript.

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

## Credits
Special thanks to the following projects for providing essential resources:
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [UI Browser v3](https://latenightsw.com/freeware/ui-browser/)

And a thanks to the following projects for inspiration:
- [mautrix/imessage](https://github.com/mautrix/imessage)
