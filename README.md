# Media Toolkit

TL;TR;
Tool to import video and photos from cameras like GoPro, Nikon DSLR or Panasonic camcorders. Plus additional features to manage file name and Exif/metadata routine. For Windows users only (at least now).

## Motivation

Primary goal was to simplify import from various digital cameras.

Each device has "own way" to be connected to PC/Laptop and file naming scheme. I need a tool that will hide all this complexity and copy media files to a laptop with predefined naming and without unnecessary questions. And I would get rid of proprietary (and outdated) applications like HD Writer.

The second group of requirements is related to the metadata handling and files organizing (file renaming, dates fixing, metadata wiping etc).

So, my base workflow has following points:

1. I have a place for home/family video and separate place for photos. 
2. All media files organized by date (e.g. `2020.01.02` or `2020.01.02_Awesome_Event`). `YYYY.MM.DD` date format help me in searching, processing an arhivig activities. 
3. I came to ide to have unified file naming scheme with timestamp in file name. e.g. `VID_${TIMESTAMP}.mp4` and `IMG_{TIMESTAMP}.jpg`. Especially I dislike GoPro naming :) 
4. Sometimes I need to parse date from file name and put it into the embedded metadata (Exif for photo and QuickTime attributes for video).

And in some cases I need to fix file names and metadata for certain files.

This application was designed to automatize these routines... and to practice in GoLang programming :)

*⚠️ WARNING* This application may perform destructive actions for media files (move, delete files or change metadata). You use it at your own risk and without any warranties. Author is not responsible for any kind of loss or damage of your data. It is strongly recommended to make data backups before any file operations.

## Installation

Step 1: put [`media-tool.exe`](https://github.com/redrathnure/media-tool/releases) to some folder (preferably in `$PATH` locations). 

Step 2. Install [ExifTool by Phil Harvey](https://Exiftool.org/) which is used to perform files and metadata manipulations. `Exiftool.exe` should be placed into `APP_DIR\Exiftool` dir OR into any `$PATH` location.

Step 3. (optional) Prepare Default Configuration. By default the application looks to `$HOME\.media-tool\media-tool.yaml` or `APP_DIR\conf\media-tool.yml` configuration file. Please see `media-tool.example.yml` file and chapters bellow for more details.

## Usage

The application has a few different commands. Please use `media-tool -h` or `media-tool {cmd} -h` to get description and related arguments.

Each command has `--config` or `-c` arg to specify configuration file from non default location. May be useful if default `$HOME\.media-tool\media-tool.yaml` OR `APP_DIR\conf\media-tool.yml` locations do not work well or if you need to keep a few different configurations.

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
If target dir was not specified, command takes it from config file. WARNING Seems Panasonic cameras expose ReadOnly storage, this is why after successful import you have to manually remove files from camera.  
It was tested with Panasoic HC-V700 camera.

### Organize Files By Date

A `media-tool import local` command suppose to move video and image files from one local directory to another with creating date folders (e.g. `2020.01.02`).

### Correcting Photo and Video Dates

A `media-tool import fixDates` command will try to read date from file name and put it to the Exif and QuickTime metadata. The command will try to correct file creating date too. May be useful for files after post processing.

### Cleanup Image Names and Metadata

A `media-tool clean names` and `media-tool clean metadata` commands may be used to remove `- Copy` filename suffixes and to wipe image metadata (e.g. wiping GPS data before publishing photos in Internet).

## Development

### How to Build

1. Install Go [v1.23.1 or later](https://go.dev/doc/install)
2. Install [Mage](https://github.com/magefile/mage). E.g. by `mkdir %GOPATH%\src && cd %GOPATH%\src && git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go`
3. Use one of predefined tasks:
    * `mage -l` - show available tasks
    * `mage releasePkg` - prepare release package
    * `mage reBuild` - build application

### TODOs

* Extract logging format to the config
* Parse Exiftool output. Warning: [minor] to debug
* Build script + prepare installation package
* Update version based on git blame
* Store image and videos formats to the config (mp4 and tsd)
* Coping speed and progress indicator
* try Exiftool -short -groupNames -if "$file:MIMEType=~/video/i" * for image and video
