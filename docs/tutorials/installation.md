# Installation

These instructions explain how to install and run Musings on a linux server.
The steps to install Musings on other operating systems are similar,
but these are left to the reader to discern.

This document does not describe how to connect Musings to the outside world.
For that, see the documentation for your web server.

## 1. Obtain Musings

### Option A: Download a release

The easiest way to obtain Musings is to download a release.

Releases are listed on the [GitHub repo](https://github.com/subtributary/musings/releases).
Download the latest release, then extracte the files.

### Option B: Build from source

Developers wishing to [customize](customization.md) Musings will want to build it from source.

First, clone the repo, then build and publish it.

```shell
git clone https://github.com/subtributary/musings.git
cd musings

# For linux:
./scripts/build.sh release
./scripts/publish.sh # copies the needed files to ./publish/
```

On Windows, open the ".sh" files and run the equivalent commands in your environment.

The files needed for the remaining steps will be in the "/publish" folder.

## 2. Copy files to destination

If the destination directory on the server has not been prepared, then set that up now.

If you downloaded a release, copy all of those extracted files to the destination.
If you built from source code, copy all the files in "./publish" to the destination.
After copying, the destination should contain the "bin", "content", "scripts", and "web" directories.

Finally, ensure that "./bin/server" at the destination can be executed.

## 3. Configure Musings

Musings is most easily configured with environment variables.
For a list of these, see the "./.env" file.

Storing the configuration in "./.env" is not required,
but it is recommended for an easy installation.
The included helper scripts load settings from this file.
If you run the `./bin/server` directly, you must set environment variables yourself. 

Most of the presets are sensible,
but the binding address will likely need changed:

```dotenv
MUSINGS_BIND_ADDRESS=localhost:8080
```

Set this to the address Musings should listen on.

## 4. Run Musings

To easily run Musings with the settings in "./.env", run this script from the installation directory:

```shell
./scripts/run.sh
```

Alternatively, set the environment variables, then run "./bin/server":

```shell
source .env # or however you want to set your variables
./bin/server
```

At this point, you may want to configure your web server (e.g., nginx) to route traffic to Musings.
You may also choose to run Musings as a system service.
These tasks are out of scope for this document.

When finished, navigate in your browser to the URL that you configured for Musings.
You should see your new website.
