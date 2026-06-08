# Deploying

These instructions explain how to install and run Musings on a linux server.
The steps to install Musings on other operating systems are similar but are left to the reader to discern.

This document does not describe how to connect Musings to the outside world.
For that, see the documentation for your web server.

## 1. Obtain Musings

### Option A: Download a release

The easiest way to obtain Musings is to download a release.

Releases are listed on the [GitHub repo](https://github.com/subtributary/musings/releases).
Download the latest release, then extract the files.

### Option B: Build from source

Developers wishing to [customize](customizing.md) Musings will want to build it from source.

Clone the repo, build it, then publish it.

The following commands are for linux. Run the equivalent commands in your environment.

```shell
git clone https://github.com/subtributary/musings.git
cd musings

# For linux:
npm install
npm run build
go build -o bin/server cmd/server
```

Most of the files in the repository will not be needed for deployment.
The files and folders that you will need are:

 * `bin/server`
 * `content/...`
 * `data/...`
 * `web/static/...`
 * `web/templates/...`

## 2. Copy files to destination

If the destination directory on the server has not been prepared, then set that up now.

Copy all the deployment files to the directory on the server.
If you downloaded a release, these will be all the extracted files.
If you built from source code, these will be the files listed above.

Ensure that "./bin/server" at the destination can be executed.

## 3. Run Musings

Run musings directly with the `bin/server` executable.
Use the `--bind` parameter or the `MUSINGS_BIND` environment variable to configure the address at which Musings listens. 

At this point, you may want to configure your web server (e.g., nginx) to route traffic to Musings.
You may also choose to run Musings as a system service.
These tasks are out of scope for this document.

When finished, navigate in your browser to the URL that you configured for Musings.
You should see your new website.
