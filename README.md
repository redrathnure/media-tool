# Media Toolkit

TL;TR;
Tool to import video and photos from cameras like GoPro, Nikon DSLR or Panasonic camcoders. Plus additional features to correct file and EXIF/metadata dates. For Windows users only (at least now).

## Motivation

Primary goal was to simplify import from various digital cameras.   
Each device has "own way" how to be connected to PC/Laptop and how to name files. I need a tool which will abstragate me from all this differences and copy media files to my laptop with needed (for me) naming. Ideally doing this automatically without any additional questions.Also I would like to get rid of proprietary and/or outdated applications like HD Writer.   

So, my base workflow has following points:

1. I have a place for home/family video and separate place for photos. 
2. All media files organized by date (e.g. `2020.01.02` or `2020.01.02_Awesome_Event`). `YYYY.MM.DD` date format help me in searching, processing an arhivig activities. 
3. I came to ide to have unified file naming scheme with timestamp in file name. e.g. `VID_${TIMESTAMP}.mp4` and `IMG_{TIMESTAMP}.jpg`. Especially I dislike GoPro naming :) 
4. Sometimes I need to parse date from file name and put it into the embeded metadata (EXIF for photo and QuickTime attributes for video).

This application was designed according to these points. And as extra goal, I want to practice in GoLang programming a bit :)


*WARNING* You use it at your own risk and without any warranties. Author is not responsible for any kind of loss or damage of your data. It is strongly recommended to make backups. 

## Installation

Step 1. Installing Application

Build project, e.g. by executing `go build` command and copy result `media-tool.exe` file to some convinient location. If you prefer CLI interface it has sense to add this location to `$PATH`. Otherwise windows shortcuts may be good alternative (e.g. several shortcuts for different cases with well recognized icons). 

Step 2. Installing Exiftool

Application uses [ExifTool by Phil Harvey](https://exiftool.org/) to perform files and metadata manipulations. Please visit https://exiftool.org/ and download stable version. I use 12.01 for Windows. `exiftool.exe` should be placed into `APP_DIR\exiftool` dir OR into any `$PATH` location.

Step 3. (optional) Prepare Default Configuration

`$HOME\.media-tool\media-tool.yaml` OR `APP_DIR\conf\media-tool.yml` file contains configuration for default locations. May be ussefull if media files always should be imported to the same location(s). Please see `media-tool.example.yml` file and chapters bellow for more details.

## Commands And Common Usage Tips

Application provides several commands for differnt case. Please use `media-tool -h` or `media-tool cmd -h` to get description and related arguments.

Each command has `--config` or `-c` arg to specify configuration file from non default location. May be usefull if `$HOME\.media-tool\media-tool.yaml` OR `APP_DIR\conf\media-tool.yml` do not work well or if you just need to keep two or more configurations.

Each commang has `-v` or `--verbose` arg which enable extra logging and may be useful for troubleshooting or initial learning phase.

And almost every command has `-d` or `--dry` arg which may be used to simulate import wihtout any file removal.

And finally, all import commands work in two steps 1. is to copy files from device to temp folder and 2. come files from temp folder to target one. If media file cannot be processed (e.g. incorect/unexpected format or not enought disk space) these files will stay in temp directory. So pleasse check your temp directoy (and related configuration) in case of any issues.

### Import GoPro Video

`media-tool import gopro` command try to find connected GoPro camera and import files into specified directory.   
If target dir was not specified, command takes it from config file.   
Was tested with GoPro HERO8, however should also work with HERO 7 model.   

### Import Photos

`media-tool import sdphotos` command try to find SD cart from DSLR cameras and import photos into specified directory. 
If target dir was not specified, command takes it from config file.   
Was tested with Nikon D800, however should also work with everything what stores `jpeg` or `NEF` files in `DCIM` folder.   

### Import Video From Panasonic Camcorder

`media-tool import camvideo` command try to find connected camcoderand import video into specified directory. 
If target dir was not specified, command takes it from config file. WARNING Seems Panasonic cameras expose ReadOnly storage, this is why after successful import you have to manually remove files from camera.  
Was tested with Panasoic HC-V700 camera.

### Organize Files By Date

`media-tool import local` command suppose to move vide and image files from one local directory to another with creating date folders (e.g. `2020.01.02`).

### Correcting Photo and Video Dates

`media-tool import fixDates` command will try to read date from file name and put it to the EXIF or QuickTime metadata. Also command will try to set file creating date. May be useful for files after post processing.

## Technical Remarks

### How to Build

Install GO v1.14.4 or later and execute `go build` from root directy.

### TODOs

* Extract logging format to the config
* Parse exiftool output. Warning: [minor] to debug
* Build script + prepare installation package
* Update version based on git blame
* Store image and videos formats to the config (mp4 and tsd)
* Coping speed and progress indicator
* try exiftool -short -groupNames -if "$file:MIMEType=~/video/i" * for image and video
