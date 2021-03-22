# mp3tag

A command line utility for modifying mp3 files and setting IDv3 tags.

## Examples

`mp3tag` has four main functions.

1.   View the metadata in an MP3 file. If no specific attributes are requested,
     all attributes will be displayed.

```
mp3tag show 07-_The_Roman_Washington.mp3
========
Album:  The History Of Rome
Artist: Mike Duncan
Title:  07- The Roman Washington
Year:   2007
Genre:  Podcast
```

2.   Find files with matching attributes. To find all the `mp3` files without
     an artist set.

```
mp3tag find 'artist:""' data/*
```

3.   Update the metadata in an MP3 file. Files can be updated one at a time, or
     many at a time.

```
mp3tag update --title="The Roman Washington" 07-_The_Roman_Washington.mp3
mp3tag update --genre=Podcast *.mp3
```

4.   Rename the MP3 file based on the attributes in the metadata. A rename
     update will take the attributes of the file that is about to be renamed,
     so it is possible to update all the files in a directory with one command.

```
mp3tag update --rename="{artist} - {album} - {title}.mp3" The\ History\ of\ Rome/*.mp3
```


## Development

```
go mod tidy
go generate ./...
go build .
go test ./...
```
