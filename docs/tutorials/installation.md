# Installation

These instructions explain how to install and run Musings on a linux server.
The steps to install Musings on other operating systems are similar,
but these are left to the reader to discern.

This document does not describe how to connect Musings to the outside world.
For that, see the documentation of your web server software.

## 1. Obtain Musings

### Option A: Download a release

Download the latest version of Musings from the releases on GitHub,
then extract the files.

### Option B: Build from source

Clone the Musings source code, then run these commands:

```shell
# For linux:
./scripts/build.sh release
./scripts/publish.sh # copies the needed files to ./publish/
```

## 2. Copy files to destination

If the destination directory on the server has not been prepared, 
then set that up now.

If you downloaded a release, copy all of those extracted files to the destination.

If you built from source code, copy all the files in "./publish" to the destination.

Finally, ensure that "./bin/server" at the destination can be executed.

## 3. Configure Musings

Musings is configured with environment variables.
For a list of these, see the "./.env" file.

Storing the configuration in "./.env" is not required,
but it is recommended for an easy installation.
The included helper scripts load settings from this file.

Most of the presets are sensible,
but the binding address will likely need changed:

```dotenv
MUSINGS_BIND_ADDRESS=localhost:8080
```

Set this to whatever address you want Musings to listen on.

## 4. Run musings

To easily run Musings with the settings in "./.env", run this script from the installation directory:

```shell
./scripts/run.sh
```

Alternatively, set the environment variables, then run "./bin/server":

```shell
source .env # or however you want to set your variables
./bin/server
```

At or near this point, you may want to configure your other web server software to route traffic to the address that Musings is listening on.
You may also wish to configure a web service to more easily manage how Musings runs.
These tasks are out of scope of this document as they pertain to other tools.

When finished, navigate in your browser to the URL that you configured for Musings.
You should see your new website.
