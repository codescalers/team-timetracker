# Timetracker

TimeTracker is a robust command-line interface (CLI) and server application designed to help teams to efficiently track time spent on various URLs, such as GitHub issues, project management tasks, or any web-based activities. By seamlessly integrating with your workflow, TimeTracker allows you to start and stop time entries, view detailed logs, and export data in multiple formats.

## Getting started

### Download the Binaries

Get the binaries from the releases or download/clone the source code and build it manually

### Manual Building

### Building the server

Executing `make build-server` is enough, a binary named `timetrackerd` should be created in the `bin` directory

### Building the CLI

Executing `make build-client` will create a binary named `timetracker` in the `bin` directory

### Running the server

#### Sample server configuration

```json
{
    "server": {
        "addr": "0.0.0.0:8080",
        "allowed_origins": ["*"]
    },
    "database": {
        "driver": "sqlite",
        "data_source": "timetracker.db"
    }
}
```

the command to run the server is `./bin/timetrackerd -config ./team-timetracker-server.sample.json`

### Running the client

The client has 3 main commands `start`, `stop`, `entries`

#### Client configurations

The client requires configuration either passed as `-config` flag or in `~/.config/team-timetracker.json` that has the username and backend URL

```json
{
    "username": "xmonader",
    "backend_url": "http://localhost:8080"
}
```

#### Start time entry

```
~> ./bin/timetracker start github.com/theissues/issue/142 "helllooo starting" 
Started tracking ID 9 for URL 'github.com/theissues/issue/142'.
```

It takes the URL as an argument and a description and returns you a tracking ID
> Note: you can't start an already started entry.

#### Stop time entry

```
~> ./bin/timetracker stop github.com/theissues/issue/142   

Stopped tracking for URL 'github.com/theissues/issue/142'. Duration: 35 minutes.
```

#### Querying the time entries

> You can generate CSV data by specifying --format=csv

##### Querying all of the entries

```
~> ./bin/timetracker-cli entries                                       
Time Entries:
------------------------------
ID: 9
Username: azmy
URL: github.com/theissues/issue/154
Description: helllooo starting
Start Time: Sun, 15 Sep 2024 22:08:01 EEST
End Time: Sun, 15 Sep 2024 22:08:14 EEST
Duration: 3 minutes
------------------------------
ID: 8
Username: xmonader
URL: github.com/theissues/issue/112
Description: hello 142
Start Time: Sun, 15 Sep 2024 22:07:10 EEST
End Time: Sun, 15 Sep 2024 22:07:12 EEST
Duration: 42 minutes
------------------------------
ID: 7
Username: xmonader
URL: github.com/theissues/issue/52111
Description: the descrpition
Start Time: Sun, 15 Sep 2024 22:06:55 EEST
End Time: Sun, 15 Sep 2024 22:06:58 EEST
Duration: 11 minutes
------------------------------
ID: 6
Username: guido
URL: github.com/theissues/issue/4
Description: alright done
Start Time: Sun, 15 Sep 2024 22:03:59 EEST
End Time: Sun, 15 Sep 2024 22:04:09 EEST
Duration: 20 minutes
------------------------------
ID: 5
Username: guido
URL: github.com/rfs/issue/87
Description: that's tough
Start Time: Sun, 15 Sep 2024 22:03:00 EEST
End Time: Sun, 15 Sep 2024 22:03:04 EEST
Duration: 15 minutes
------------------------------
ID: 4
Username: xmonader
URL: github.com/zos/issue/414
Description: woh
Start Time: Sun, 15 Sep 2024 22:02:26 EEST
End Time: Sun, 15 Sep 2024 22:02:41 EEST
Duration: 33 minutes
------------------------------
ID: 3
Username: xmonader
URL: github.com/thesites/www_what_io
Description: another desc1123
Start Time: Sun, 15 Sep 2024 22:00:13 EEST
End Time: ---
Duration: ---
------------------------------
ID: 2
Username: ahmed
URL: https://github.com/xmonader/team-timetracker/issues/121
Description: el desc
Start Time: Sun, 15 Sep 2024 20:03:56 EEST
End Time: Sun, 15 Sep 2024 20:06:05 EEST
Duration: 24 minutes
------------------------------
ID: 1
Username: ahmed
URL: https://github.com/xmonader/team-timetracker/issues/1
Description: another desc
Start Time: Sun, 15 Sep 2024 20:03:24 EEST
End Time: Sun, 15 Sep 2024 20:03:50 EEST
Duration: 23 minutes
------------------------------
```

###### Generating CSV data

```
~> ./bin/timetracker-cli entries --format=csv
ID,Username,URL,Description,Start Time,End Time,Duration (minutes),Created At,Updated At
9,azmy,github.com/theissues/issue/154,helllooo starting,2024-09-15T22:08:01+03:00,2024-09-15T22:08:14+03:00,3,2024-09-15T22:08:01+03:00,2024-09-15T22:08:14+03:00
8,xmonader,github.com/theissues/issue/112,hello 142,2024-09-15T22:07:10+03:00,2024-09-15T22:07:12+03:00,42,2024-09-15T22:07:10+03:00,2024-09-15T22:07:12+03:00
7,xmonader,github.com/theissues/issue/52111,the descrpition,2024-09-15T22:06:55+03:00,2024-09-15T22:06:58+03:00,11,2024-09-15T22:06:55+03:00,2024-09-15T22:06:58+03:00
6,guido,github.com/theissues/issue/4,alright done,2024-09-15T22:03:59+03:00,2024-09-15T22:04:09+03:00,20,2024-09-15T22:03:59+03:00,2024-09-15T22:04:09+03:00
5,guido,github.com/rfs/issue/87,that's tough,2024-09-15T22:03:00+03:00,2024-09-15T22:03:04+03:00,15,2024-09-15T22:03:00+03:00,2024-09-15T22:03:04+03:00
4,xmonader,github.com/zos/issue/414,woh,2024-09-15T22:02:26+03:00,2024-09-15T22:02:41+03:00,33,2024-09-15T22:02:26+03:00,2024-09-15T22:02:41+03:00
3,xmonader,github.com/thesites/www_what_io,another desc1123,2024-09-15T22:00:13+03:00,N/A,5,2024-09-15T22:00:13+03:00,2024-09-15T22:00:13+03:00
2,ahmed,https://github.com/xmonader/team-timetracker/issues/121,el desc,2024-09-15T20:03:56+03:00,2024-09-15T20:06:05+03:00,24,2024-09-15T20:03:56+03:00,2024-09-15T20:06:05+03:00
1,ahmed,https://github.com/xmonader/team-timetracker/issues/1,another desc,2024-09-15T20:03:24+03:00,2024-09-15T20:03:50+03:00,23,2024-09-15T20:03:24+03:00,2024-09-15T20:03:50+03:00

```

##### Querying entries by username

```
------------------------------
~> ./bin/timetracker-cli entries --username="azmy"    
Time Entries:
------------------------------
ID: 9
Username: azmy
URL: github.com/theissues/issue/154
Description: helllooo starting
Start Time: Sun, 15 Sep 2024 22:08:01 EEST
End Time: Sun, 15 Sep 2024 22:08:14 EEST
Duration: 3 minutes
------------------------------
```

##### Querying entries by URL

```
~> ./bin/timetracker-cli entries --url="github.com/theissues/issue/154"
Time Entries:
------------------------------
ID: 9
Username: azmy
URL: github.com/theissues/issue/154
Description: helllooo starting
Start Time: Sun, 15 Sep 2024 22:08:01 EEST
End Time: Sun, 15 Sep 2024 22:08:14 EEST
Duration: 3 minutes
------------------------------
```
