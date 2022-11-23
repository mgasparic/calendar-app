# Calendar App

A simulation of randomly generated calendar events that get stored to a database. Emphasis on concurrency and
multi-threading.

## Services and Supported Endpoints

### Generator

| Path                | Method | Params          | Possible Codes                                                                                      | Success Response |
|:--------------------|:-------|:----------------|:----------------------------------------------------------------------------------------------------|:-----------------|
| /events/:user_email | GET    | offset<br>limit | 200 - OK<br>204 - No Content<br>400 - Bad Request<br>404 - Not Found<br>500 - Internal Server Error | List of ics data |

### Consumer

No endpoint

## Simulation

Run `./scripts/setup.sh` on Linux machine with Docker. Or perform equivalent actions on a different machine.
