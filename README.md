# Media Toolkit

A tool for importing and managing media data from GoPro, photo cameras, camcoders, smartphones and other media devices. Contains features to automatize file name and Exif/metadata related routine.

## Motivation

There are two groups of the tasks to be automatized:

- import data from various digital cameras(GoPro, Nikon/Cannon DSLR or Panasonic camcorders).
- prepare proper file names and/or correct Exif data (date creation, removing some tags etc).

Base workflows:

1. Importing photos/videos from camera:
    - connect camera
    - run tool which will ...
    - ...move photo and video content to specific place and organize by folder with `YYYY.MM.DD` naming format
    - ...remove original content from the camera
2. Cleanup names:
    - run tool in directory with wrong/unexpected names OR/AND "copy" name suffixes OR/AND not full Exif data (e.g. photos from WhatsApp)
    - rename files to desired naming schema (e.g. `VID_${TIMESTAMP}.mp4` and `IMG_{TIMESTAMP}.jpg`)
    - ... including handling a "Copy" suffixes.
    - fill/correct missed Exif data
3. Cleanup unnecessary Exif data:
    - run tool in directory with photos to remove unnecessary information about camera, location and other tags, which should not be shared with other users

*⚠️ WARNING* This application may perform destructive actions for media files (move, delete files or change metadata). You use it at your own risk and without any warranties. Author is not responsible for any kind of loss or damage of your data. It is strongly recommended to make data backups before any file operations.


### Commands and Features

List of commands with a short description. Please run `media_tool <command> -h` to get more information.

| Command | Description |
| ------------- | ------------- |
| `clean metadata` | Remove camera, location and other information from images |
| `clean names` | Cleanup file names, remove `Copy` prefixes|
| `fix dates` | Reads date from file name and put it into Exif/QuickTime metadata attributes as well as updating file creation/modification attributes. |
| `import camVideo` | Copy video from Panasonic camcoder (WPD) to disk |
| `import gopro` | Copy images and video from GoPro card (WPD) to disk |
| `import local` | Copy images and video from directory to disk |
| `import sdphotos` | Copy images and video from SD card(s)to disk |
| `completion` | Generate the autocompletion script for media-tool for the specified shell |
| `config save` | Generate example configuration file |
| `config print` | Print configuration |


## Installation

Manual installation:

1. Download one of [`release packages`](https://github.com/redrathnure/media-tool/releases) and put `media-tool.exe`/`media-tool` to some folder (preferably in `$PATH` locations).
2. Install [ExifTool by Phil Harvey](https://Exiftool.org/) which is used to perform files and metadata manipulations.
    - Windows: `Exiftool.exe` should be placed into `APP_DIR\Exiftool` dir OR into any `$PATH` location.
    - Linux: something like `sudo apt install exiftool`
3. (optional) Prepare configuration (see `media-tool.example.yml` file and chapters bellow for more details). By default the application looks into following locations:
    - (preferable for a Linux env) `~/.config/media-tool/media-tool.yaml`
    - (preferable for a Windows env) `$HOME\.media-tool\media-tool.yaml`
    - (portable installation)`APP_DIR\conf\media-tool.yml`

## Usage

The application has a few different commands. Please use `media-tool -h` or `media-tool {cmd} -h` to get description and related arguments.

Each command has `--config` or `-c` arg to specify configuration file from non default location. May be useful if default (`~/.config/media-tool/media-tool.yaml`, `$HOME\.media-tool\media-tool.yaml` or `APP_DIR\conf\media-tool.yml`) location do not work well or if you need to temporally use different conf.

Each command has `-v` or `--verbose` arg which enable extra logging and may be useful for troubleshooting or initial learning phase.

And almost every command has `-d` or `--dry` arg which may be used preview changes without execute them.

And finally, all `import` commands work in two steps:

1. import files from device to temp directory
2. move files from temp folder to target one

If a media file cannot be processed (e.g. unexpected format or luck of disk space) these files will stay in temp directory. In case of any issues or incoplet operation please check your temp directory.

### Import GoPro Video

A `media-tool import gopro` command try to find connected GoPro camera and import files to specified directory.
If target dir was not specified, command takes it from config file.
It was tested with GoPro HERO8, however should also work with other models too.

### Import Photos from Camera or SD Card

A `media-tool import sdphotos` command try to find SD cart from DSLR cameras and import photos to specified directory.
If target dir was not specified, command takes it from config file.
It was tested with a few Nikon and Canon cameras, however should also work with everything what stores `jpeg`, `NEF` or `CR2`/`CR3` files.

### Import Video From Panasonic Camcorder

A `media-tool import camvideo` command try to find connected camcorder and import video into specified directory.
If target dir was not specified, command takes it from config file. WARNING Seems Panasonic cameras expose read only storage, this is why after successful import you have to manually remove files from camera.  
It was tested with Panasonic HC-V700 camera.

### Organize Files By Date

A `media-tool import local` command suppose to move video and image files from one local directory to another with creating date folders (e.g. `2020.01.02`).

### Correcting Photo and Video Dates

A `media-tool import fixDates` command will try to read date from file name and put it to the Exif and QuickTime metadata. The command will try to correct file creating date too. May be useful for files after post processing.

### Cleanup Image Names and Metadata

A `media-tool clean names` and `media-tool clean metadata` commands may be used to remove a `- Copy` and ` Copy` suffixes from filename and to wipe image metadata (e.g. wiping GPS data before publishing photos in Internet).

## Development

### How to Build

1. Install Go [v1.23.1 or later](https://go.dev/doc/install)
2. Install [Mage](https://github.com/magefile/mage). E.g. by `go install github.com/magefile/mage@latest &&
mage -init`
3. Use one of predefined tasks:


### Useful Commands

Dev routines:

* `mage -l` - show available tasks
* `mage goUpdateDeps` - update dependencies
* `mage build` - build application
* `mage buildClean` - cleanup project and build application
* `mage release` - build release packages
* `mage releaseVersion` - get or calculate new version ,prepare tag and build release packages


### TODOs

* Extract logging format to the config
* Parse Exiftool output. Warning: [minor] to debug
* Store image and videos formats to the config (mp4 and tsd)
* Coping speed and progress indicator
* try Exiftool -short -groupNames -if "$file:MIMEType=~/video/i" * for image and video
* Import data from SD/flash storage (Linux)
* Import data from MTP devices (Linux)
* deb packet?
* WhatApp and GPixel files handling
* Proper handling of unproper dates (a "1971 year for FAT32" issue)
* Handle more "Copy" naming patterns
