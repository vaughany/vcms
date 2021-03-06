# VCMS - Vaughany's Computer Monitoring System

![Screenshot of v0.0.5](./screenshots/v0.0.5.png)

[![Go Report Card](https://goreportcard.com/badge/vaughany/vcms)](https://goreportcard.com/report/vaughany/vcms)
[![codecov](https://codecov.io/gh/vaughany/vcms/branch/main/graph/badge.svg?token=RC1PUVBO78)](https://codecov.io/gh/vaughany/vcms)

## What is VCMS?

> **tl;dr:**  VCMS is a work-in-progress cross-platform zero-config computer health monitoring tool (for many computers to report their status to a central point) and is a side-project I am working on to learn [Go](https://golang.org/).  Be kind!

I've used [Munin](https://munin-monitoring.org/) for years and [M/Monit](https://mmonit.com/) a lot also.  Both are excellet montoring solutions.  I wanted to try my hand at something similar, as a learning exercise.  *This project doesn't (yet) do anything that Munin or M/Monit don't already do,* but the upside is that it's really, really simple.  I also recently discovered Go project and enterprise-level monitoring solution [Sensu](https://sensu.io/): hopefully one day this project will be a tenth as good as that! :)

VCMS is two programs: one (called Collector) sits on a computer and monitors it's state, periodically sending the information to the other (called Receiver) which makes a nice live web page out of it.  You can have as many collectors on as many computers as you like.

This project is the result of a much larger, private project in which I attempted too much too soon, with minimal Go knowledge.  The scope crept faster than my knowledge and ability to implement it, and it got unwieldly and complex, so this is me 'starting over', one feature at a time.  

This project is by no means feature complete and probably won't be for a while, and is being written *Linux-first*. I know less about the inner workings of Windows than I used to due to simply not using it regularly for years, but it's my intention to make all reported details work for Linux, Windows and macOS in time.

### Metrics reported

Version 0.0.1:

* hostname 
* username
* IP address
* first seen
* last seen

Added in version 0.0.2:

* uptime
* operating system
* if a reboot is required (Linux)
* load averages (Linux)
* memory, total and free (Linux)
* swap memory, total and free (Linux)
* disk space, total and free (Linux)

Added in version 0.0.3:

* CPU number and clock speed (Linux)

> **Reminder:** this project is in no way feature-complete!

### Tested with

I aim to test both Receiver and Collector programs on as many OSes and architectures as I reasonably can, but it's going to be the amd64/x86_64 versions of OSes I can easily download and use via [Vagrant](https://www.vagrantup.com/) in the first instance.

* Debian:
  * 8 / Jessie
  * 9 / Stretch
  * 10 / Buster
  * 11 / Bullseye
* Ubuntu (recent LTS versions and supported non-LTS versions since the last LTS version):
  * 16.04 LTS / Xenial Xerus
  * 18.04 LTS / Bionic Beaver
  * 20.04 LTS / Focal Fossa
  * 20.10 / Groovy Gorilla
  * 21.04 / Hirsute Hippo
  * Kubuntu versions as Ubuntu versions
  * Lubuntu versions as Ubuntu versions
* openSUSE:
  * Leap 15
  * Tumbleweed
* Red Hat Enterprise Linux:
  * 7 / Maipo
  * 8 / Ootpa
* Mint
  * 20
  * ~~19~~
* CentOS:
  * ~~Linux 7~~
  * Linux 8
  * Stream 8
* Fedora:
  * 34 (Cloud Edition)
  * 33 (Cloud Edition)
* Elementary
  * 6 / Odin
* ~~Gentoo~~
* ~~Mandriva~~
* ~~Turbolinux~~
* ~~Xandros~~
* Manjaro
* Oracle:
  * 7
  * 8
* Kali:
  * Rolling
* Zorin
* Solaris:
  * 11
* BSD:
  * FreeBSD
    * ~~14~~
    * 13
    * 12
  * OpenBSD
    * 6.9
* ~~Plan 9~~
* ~~IBM OS / 2~~
* ~~macOS~~
* Windows:
  * 10

---

## Basic Linux Installation

1. Download the binaries ([from here](https://github.com/vaughany/vcms/releases)).  They're not large.  You'll need both the Collector and Receiver for your platform.

2. Open two terminals, `cd` to the download location in each.

3. Run the Receiver in one terminal:

   `./receiver`

   > **Note:** you might have to make it executable before you can run it:
   > 
   > `chmod +x receiver`

   You should see output similar to the following:

   ```
   VCMS - Receiver v0.0.7 (2021-08-08), go1.17.
   Receives data from the Collector apps, creates a web page.
   Loading nodes from persistent storage.
   nodes.json could not be read, so could not load node data.
   Running web server on http://127.0.0.1:8080.
   To connect a Collector, run: './collector -r http://127.0.0.1:8080'.
   ```

4. Run the Collector in the other terminal:

   `./collector`

   > **Note:** you might have to make it executable before you can run it:
   > 
   > `chmod +x collector`

   You should see output similar to the following:

   ```
   VCMS - Collector v0.0.7 (2021-08-08), go1.17.
   Collects information about the computer, and sends it to the Receiver app.
   Sending data to http://127.0.0.1:8080/api/announce
   Response: 200 OK
   ```

   > **Note:** If that `200 OK` line is not there, there's a problem somewhere.

   In the other terminal running the Receiver, you should see the following:

   ```
   Received data from 127.0.0.1:8080
   ```

5. I've tried to make both programs user-friendly, so you can specify the following flags:

   `-d` for 'debug mode', which basically makes the output more verbose.

   `-v` for version info.

   `-h` for help.  More options will most likely be added over time.

6. For testing purposes, so I could run multiple Collectors on just one computer, I added a testing mode which generates random hostname, IP address and so on.  For each extra Collector that you want to run, open a new terminal for it and run it with the `-t` flag:

   ```
   ./collector -t
   ```

   Random hostnames are in the format _adjective-color-animal_, IP addresses always start with the number '512', and the username is a Norse god.

---

## Better Linux installation

By default, the Receiver runs the web page and API on `127.0.0.1` ('localhost') on port `8080`.  These are fair defaults, but won't work outside your computer.  Both programs allow you to specify the IP address and port to use with the `-r` flag.  

As of v0.0.10 you can force the use of a pre-shared key, so that data sent without it is rejected.  There's also a helper function to generate one for you.

1. To run Receiver on a different IP address and/or port:

   ```
   ./receiver -r 192.168.0.100:8081
   ./receiver -r mywebsite.com:8081
   ```

   > **Note:** do not specify the protocol.  We're using `http` for now.

   Receiver will complain if it cannot use the IP address or port you specify.  Ideally, you should be using the IP address of your computer as shown in `ifconfig`.

2. Make the Collector use a different IP address and/or port:

   ```
   ./collector -r http://192.168.0.100:8081
   ```

   > **Note:** Some TCP/IP ports, such as the standard web port 80, are restricted and cannot be used unless you run the program with elevated privileges:
   >
   >`./receiver -r 192.168.0.100:80 // reserved port: fails.`
   >
   >`sudo ./receiver -r 192.168.0.100:80 // works.`

3. To force the use of a pre-shared key for added security:

   If the Receiver specifies a pre-shared (API) key, then any Collector sending data is required to send the same key within the JSON, or the data will be rejected. 

   > **Note:**  The Receiver is not yet using HTTPS so data is not encrypted in transit, nor is the API key encrypted in any way in the JSON.

   Run the Receiver as follows (continuing the example from above):

   ```
   ./receiver -r 192.168.0.100:8081 --apikey long-string-of-letters-and-numbers
   ```

   > **Note:** You can't specify **no** API key (e.g. `--apikey`) but you can specify an **empty** API key (e.g. `--apikey ""`), which is treated as no API key. This is silly, don't do it.

   At startup, the helper text shown to run a correcly-configured Collector will show the new command:

   ```
   To connect a Collector, run: './collector -r http://192.168.0.100:8081 --apikey long-string-of-letters-and-numbers'.
   ```

   If a Collector has not been similarly configured, you'll get an error message and HTTP status code 403:

   ```
   Response: 403 Forbidden
   ERROR: 192.168.0.100:8081: Collector API key '' does not match Receiver API key (not shown). Ignoring data.
   ```

   ...so configure the Collector as described above.

4. Getting a random pre-shared key:

   Both Receiver and Collector apps can be run with `-k` as the only argument, and this will generate a handful of keys of varying lengths in both hex and base64 formats, then quit.  It's just a convenience helper, nothing more.  Run it multiple times for more random keys.

   ```
   $ ./collector -k
   Hex:
    16 chars: 1869abd649331542 / 1869ABD649331542
    32 chars: 2c282ff203bb30879a5505493cbc2d13 / 2C282FF203BB30879A5505493CBC2D13
    48 chars: bb12c3223fa26ed0b04f7767758f48a3c42bfbc07c07d826 / BB12C3223FA26ED0B04F7767758F48A3C42BFBC07C07D826
   Base64:
    16 chars: i0L2oZt32HSIKWhC
    32 chars: Pjb3niW4La7OsTcuJeKI8PPYeinBjhBn
    48 chars: nBqtfU-gdvZ74XZP3phTv7CE1jezCFjVfkajfYofU7V8B2Bs
    64 chars: eS5uz5UyV5x31T_aRYEHWHRC3Jj5dtDErPmjN3SqOOI9MVPKUBmiQ4HYazy37Nef
   ```

---

## Use

When you run the Receiver, it mentions a URL (`http://127.0.0.1:8080` by default).  Click it, or copy-paste it into your web browser.  You should see a nice, if basic web page detailing the hosts in alphabetical order of host (as per the screenshot at the top).  Web page accesses are logged.  The page will automatically refresh every 60 seconds.

A 'ping' API endpoint is provided so you can check if the Receiver is responding:

`http://127.0.0.1:8080/api/ping`

Accessing the above should return the following JSON response:

`{"result":"pong"}`

Pings are logged.

> **Note:** You can of course access the above from a web browser as it's a simple 'GET' request, or 'curl' it if you prefer:
>
> `curl http://127.0.0.1:8080/api/ping`

In v0.0.9 I added logging to a file. Both commands will create and use a folder called `logs` in the current working directory, but each command's log file is named differently and datestamped with the date the command was run, e.g.:

```
$ ls -hl logs/
-rwxr-xr-x 1 paulvaughan paulvaughan 258K Sep 10 23:12 vcms-collector_2021-09-10.log
-rwxr-xr-x 1 paulvaughan paulvaughan 167K Sep 10 23:13 vcms-receiver_2021-09-10.log
```

There's nothing in the commands which culls old logs, or logs over a certain size: they will continue to grow in size and number over time.  This is acceptable to me at this time, but I may come up with a different solution (e.g. configurable log settings) in the future.  For now, just keep an eye on them and delete occasionally.  And don't use debugging (`-d`) unless absolutely required.

---

## Building from source

If you have Go installed, clone the repo and build it yourself. 

### Linux:

1. Clone the repository, e.g.:

   `git clone git@github.com:vaughany/vcms.git`

2. Change directory into the folder:

   `cd vcms`

3. Run the following helper script to run a variety of tests on the code (you may have to install some Go tools):

   `./format.sh`

4. Run the following command to build the binaries for your architecture (Go figures it out) and put into `bin/`:

   `go build -trimpath -ldflags "-s -w" -o bin/ ./cmd/...`

5. Alternatively, run the helper script `./build.sh` to build Linux and Windows binaries.

   Add or amend as you see fit to build as many different binaries as required.  To see a list of possible OSes and architectures, run `go tool dist list`, and then add new lines, e.g.

   ```
   echo -e "\e[1mBuilding macOS/Arm64...'\e[0m"
   env GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o bin/ ./cmd/...
   ```

6. Run the Collector:

   `./bin/collector`
   
7. Open a new terminal and (assuming the same folder) run the Receiver:
   
   `./bin/receiver`

---

## Troubleshooting and Gotchas

1. The Collector **needs** to be able to connect to the Receiver on the chosen URL, or the Receiver will never receive any data.  You might not be able to communicate across VLANs, for example.

2. ~~Data the Receiver receives is stored only in memory at this time.  If you quit the Receiver, everything is lost.~~ As of version 0.0.4, data is persisted to disk regularly and when the program quits, and is re-loaded on startup.  

3. We currently only use HTTP, which means your data is not encrypted in transit.  For my intended use, this is not a huge concern, but it might be for you.  Addressing this is on the to-do list.

4. Internally, data is stored with the hostname as the primary key, so if you have two servers with the same hostname, they'll only appear as one, and the details of each will be constantly overwritten by the other.  I doubt this will be a problem during regular use, but it caught me out during testing.

---

## To do

There's a lot I want to do:

* Receiver to persist data to database
* Store historical data for comparison, graphs etc.
* Choice of HTTP or HTTPS.
* Create a proper web-app with secure login.
* Create groups to group Collectors by, e.g. 'dev', 'production' etc.
* More active monitoring, with alerts should something go awry.
* Web-based, Email and Slack notifications of alerts.
* Monitor software deployed by the Ruby gem Capistrano.
* Make web page refresh configurable.
* Report any data-collecting errors in the Collector through to the Receiver.
* Export to:
  * HTML (web page saved locally, regularly)
  * XML, CSV etc.
* Web app:
  * Page listing just the distributions
  * Page listing just one host
* Consider using github.com/shirou/gopsutil for OS details.
* Look into using gRPC / protocol buffers. Unsure if they have any advantage over simple JSON.
* Send a 'pause', 'restart' or 'force send' from the Receiver.
* Controls on the Receiver web page to increase / decrease the speed of individual Collector or pause it for a bit / a lot.
* Ability to monitor the Receiver's 'logs' folder, and the individual Collector command's 'logs' folders, warning of too many / too large files, and delete them.

Completed to-do's:

* Data is persisted to disk (as JSON) regularly and when the program quits, and reloaded on startup (v0.0.4).
* Added operating system icon, where one can be derived from the OS's name (v0.0.4).
* Ability to remove a node from the Receiver (mostly this was for testing, so I didn't have to quit and restart it) (v0.0.5).
* Export as JSON (v0.0.7).
* Page listing just the hosts (v0.0.8).
* Logging to a file as well as stdout (v0.0.9).
* Authentication of Collectors by pre-shared key (v0.0.10).

---

## History

* **2021-07-01**, v0.0.1.  Initial release.  Collector registers basic info with the Receiver.
* **2021-07-12**, v0.0.2.  Collects more information, but only on Linux.  Change to struct to allow 'meta' data such as app version, errors.  Version check: Receiver will reject data if the Collector is not the same version.
* **2021-08-02**, v0.0.3.  Bug fixes: using host address, not remote address; hung on failed data send.  Ensured changing data is logged on first attempt: it's anti-DRY, but the web page doesn't look like it's broken now.  Added basic Windows version string.  Tested Collector on many versions of Linux, and Windows 10.  Added CPU count and clock speed.  Added more future to-do's to readme.
* **2021-08-05**, v0.0.4.  Upgraded to Go v1.16.7.  Added operating system images for many systems, to made identification easier (I found out that although they're marketed differently, Lubuntu and Kubuntu identify themselves as Ubuntu, and 'MX' as Debian).  More testing, including Solaris (had to add some 'not yet implemented' checks) requiring new binaries and new build config in the 'build.sh' helper script (only the binaries for Linux are kept in the repo).  Data in the Receiver is persisted to disk at regular intervals and when the program is terminated (additionally created a shutdown handler to achieve this).  Data is then loaded again at startup.
* **2021-08-06**, v0.0.5.  Added ability to remove a node's data from the Receiver.
* **2021-08-07**, v0.0.6.  Added timestamp to the disk-persisted data; added macOS to build file.
* **2021-08-08**, v0.0.7.  Added ability to export all data as JSON.
* **2021-09-07**, v0.0.8.  Lots of changes based on comments from the nice people at the Gophers Slack, including: removing init(), removing globals (some: work in progress), cleaning up the readme, adding OS logos, using GoReleaser and more linters, basic tests and benchmarks (also work in progress), new lighter dashboard and hosts web pages, and tidied up the struct used to marshal/unmarshal data.
* **2021-09-09**, v0.0.9.  Added logging to a file. Creates folder in cwd, then creates a named and dated file, logs to it and stdout.
* **2021-09-10**, v0.0.10.  Added a basic form of authentication by the use of pre-shared keys, using the `--apikey` flag. Collectors sending data with a non-matching or empty key will have their data rejected.  Generate random keys with the `-k` helper flag.

---

## Packages used

* [cloudfoundry/bytefmt](https://github.com/cloudfoundry/bytefmt/) for human-readable byte formatter.
* [hako/durafmt](https://github.com/hako/durafmt) human-readable formatting of time.Duration.

---

## Contributing and reporting issues

If you want to contribute, you're more than welcome.  This is a learning project for me, but I'm happy to consider a pull request that fixes an issue, or adds something useful (such as instructions for installation, use, and building from source on Windows or macOS).

Feel free to file a [bug report or feature request](https://github.com/vaughany/vcms/issues/new/choose) if you find bugs or think of cool features.  Happier with a pull request, but happy to learn from a bug or feature request also. 

---

## About

I'm Paul 'Vaughany' Vaughan, and I've been doing dev-ops for a while and web development for _ages_.  I'm learning Go, and finding interesting little projects to create with it, usually trying to do a better job than just a Bash script run via Cron.  This is one such project.  I'm publishing it online as I develop it, rather than as a finished product, so please remember this if you use this software or get in touch.  If you want to contact me, my email address is _paulieboo at gmail dot com_ and on my [GitHub profile](https://github.com/vaughany/).

---

## Licence

[GNU GPLv3](https://choosealicense.com/licenses/gpl-3.0/):

"_Permissions of this strong copyleft license are conditioned on making available complete source code of licensed works and modifications, which include larger works using a licensed work, under the same license. Copyright and license notices must be preserved. Contributors provide an express grant of patent rights._"
