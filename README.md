



## additional info
ffprobe -v quiet VID_20200611_115538.MP4 -print_format json -show_entries stream_tags:format_tags

ffprobe -v quiet VID_20200611_115538.MP4 -print_format json -show_entries stream_tags:format_tags


ffprobe -v quiet VID_20200611_115538.MP4 -print_format json  -show_format
{
    "format": {
        "filename": "VID_20200611_115538.MP4",
        "nb_streams": 5,
        "nb_programs": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "50.730667",
        "size": "382249905",
        "bit_rate": "60279105",
        "probe_score": 100,
        "tags": {
            "major_brand": "mp41",
            "minor_version": "538120216",
            "compatible_brands": "mp41",
            "creation_time": "2020-06-11T11:55:38.000000Z",
            "firmware": "HD8.01.01.60.00"
        }
    }
}


Go
https://github.com/smith-30/go-ffprobe/blob/master/ffprobe.go
https://github.com/asticode/go-astiffprobe

https://github.com/dsoprea/go-exiftool
https://github.com/barasher/go-exiftool


exiftool
https://exiftool.org/#running
https://exiftool.org/faq.html#Q5

exiftool -short -groupNames 
exiftool.exe -short -groupNames

Exiftool –quicktime:createdate=2015:03:07 11:08:19 :\march_7_kiteboarding\FILE0020.mp4

exiftool -v2 -quicktime:CreateDate="2015:03:07 11:08:19" -quicktime:ModifyDate="2015:03:07 11:08:19" -quicktime:TrackCreateDate="2015:03:07 11:08:19" quicktime:TrackCreateDate="2015:03:07 11:08:19" quicktime:TrackModifiedDate="2015:03:07 11:08:19" quicktime:MediaCreateDate="2015:03:07 11:08:19" quicktime:MediaModifyDate="2015:03:07 11:08:19" FILE0020.mp4
exiftool -v2 -quicktime:CreateDate="2015:03:07 11:08:19" -quicktime:ModifyDate="2015:03:07 11:08:19" -quicktime:TrackCreateDate="2015:03:07 11:08:19" -quicktime:TrackCreateDate="2015:03:07 11:08:19" -quicktime:TrackModifyDate="2015:03:07 11:08:19" -quicktime:MediaCreateDate="2015:03:07 11:08:19" -quicktime:MediaModifyDate="2015:03:07 11:08:19" FILE0020.mp4


exiftool -v2 "-quicktime:CreateDate<filename" "-quicktime:ModifyDate<filename" "-quicktime:TrackCreateDate<filename" "-quicktime:TrackCreateDate<filename" "-quicktime:TrackModifyDate<filename" "-quicktime:MediaCreateDate<filename" "-quicktime:MediaModifyDate<filename" FILE0020.mp4

exiftool "-alldates<filename"


%Image::ExifTool::UserDefined::Shortcuts = (
    VideoDate => ['quicktime:CreateDate', 'quicktime:ModifyDate', 'quicktime:TrackCreateDate', 'quicktime:TrackCreateDate', 'quicktime:TrackModifyDate', 'quicktime:MediaCreateDate', 'quicktime:MediaModifyDate'],
    FileDate => ['File:FileModifyDate','File:FileAccessDate','File:FileCreateDate'],
);


exiftool -v2 "-videodate<filename"

exiftool -v2 -videodate="2010:03:04 18:51:58+02:00" - TZ does not work

exiftool -v2 "-videodate<filename" "-filedate<filename"

exiftool -v2 "-ImageDate<filename" "-VideoDate<filename" "-FileDate<filename"




exiftool "-FileName<CreateDate" -d "%Y%m%d_%H%M%S.%%e" DIR
Or a new directory can be specified by setting the value of the Directory tag. For example, the following command moves all images originally in directory "DIR" into a directory hierarchy organized by year/month/day:

exiftool "-Directory<DateTimeOriginal" -d "%Y.%m.%d" -d -ext jpg .



File name
https://exiftool.org/filename.html
https://exiftool.org/filename.html#ex6

FileName, Directory or TestName 


exiftool -d %Y%m%d_%H%M%%-c.%%e "-testname<CreateDate" tmp


exiftool -d  %Y.%m.%d/src/%Y%m%d_%H%M%S%-c.%%e "-testname<quicktime:CreateDate" 
exiftool -d %Y.%m.%d\src\VID_%Y%m%d_%H%M%S%-c.%%e "-testname<CreateDate"




exiftool "-FileName<DateTimeOriginal" -d %Y.%m.%d\%%f%%-c.%%e -ext jpg .
exiftool "-FileName<DateTimeOriginal" -d %Y.%m.%d\%%f%%-c.%%e -ext jpg .
exiftool "-FileName<CreateDateTime" -d %Y.%m.%d\%%f%%-c.%%e -ext jpg .





CO CLI
https://github.com/spf13/cobra

https://github.com/alecthomas/kong
https://github.com/alecthomas/kingpin



https://pkg.go.dev/github.com/spf13/viper?tab=overview


GO todos
https://taskfile.dev/#/installation
https://taskfile.dev/#/




Progress
-v0 


MTP



src: 'mtp:'
=== #0
0 - HERO8 BLACK (HERO8 BLACK)
s4#DCIM, ModTime:0, Size: 0, ChildCount: 1
#DCIM\100GOPRO, ModTime:0, Size: 0, ChildCount: 1
#DCIM\100GOPRO\GH010189.MP4, ModTime:0, Size: 159427, ChildCount: -1
#Get_started_with_GoPro.url, ModTime:0, Size: 139, ChildCount: -1


d:\dev\src\_tools\media-tool>media-tool import gopro
gopro called
src: 'mtp:'
Checking MTP#0
path separator: \
s4obj: &{{0 0 true} 0 s4 DEVICE  99ed0160-17ff-4c44-9d98-1d7a6f941921}
Processing : DCIM ( DCIM)
Processing : DCIM\100GOPRO (DCIM 100GOPRO)
Processing : DCIM\100GOPRO\GH010189.MP4 (DCIM\100GOPRO GH010189.MP4)
Processing : Get_started_with_GoPro.url ( Get_started_with_GoPro.url)
#DCIM, ModTime:0, Size: 0, ChildCount: 1
#DCIM\100GOPRO, ModTime:0, Size: 0, ChildCount: 1
#DCIM\100GOPRO\GH010189.MP4, ModTime:0, Size: 159427, ChildCount: -1
#Get_started_with_GoPro.url, ModTime:0, Size: 139, ChildCount: -1
0 - HERO8 BLACK (HERO8 BLACK)



TODO list
** Extract logging format to the config
* Remove sr from videcam
* Parse exiftool output. Warning: [minor] to debug
* Config
* Remember latest used folder
* Build script + prepare arhive
* Update version based on git info
* Store image and videos formats to the config (mp4 and tsd)
* Put tags istead non standard tag aliases
* Prepare Readme
* Copy speed indicator
* systemDir = "System Volume Information" //TODO take care about recycle bin
* exiftool -short -groupNames -if "$file:MIMEType=~/video/i" * for image and video