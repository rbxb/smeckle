# smeckle

smeckle is used to convert graphics assets from the game *Vainglory* to standard formats.

 - [Installation](#Installation)
 - [convertmodel](#convertmodel)
 - [converttexture](#converttexture)
 - [compare](#compare)

## Installation

Install the command line tools using Go:

```shell
$ go get github.com/rbxb/smeckle
$ go install github.com/rbxb/smeckle/cmd/convertmodel
$ go install github.com/rbxb/smeckle/cmd/converttexture
$ go install github.com/rbxb/smeckle/cmd/convertcompare
```

Or you can download precompiled binaries for Windows:

[windows amd64](windows_amd64.zip)

## convertmodel

convertmodel converts models from SEMC's format to a wavefront (.obj) file.

`-source` The directory or file that you want to convert. Defaults to `./source`.  
`-ex` The save directory. Defaults to `./ex`.  
`-threads` The maximum number of threads to use. Defaults to `8`.  

```shell
$ convertmodel -source <source directory> -ex <save directory>
```

## converttexture

converttexture converts textures from SEMC's format to PNGs.  
Each texture has a color texture and a specular map.

`-source` The directory or file that you want to convert. Defaults to `./source`.  
`-ex` The save directory. Defaults to `./ex`.  
`-threads` The maximum number of threads to use. Defaults to `8`.  

```shell
$ converttexture -source <source directory> -ex <save directory>
```

## compare

compare is a general tool that compares two similar file directories. Any new files or files that have changed in the second directory are copied to a save location.

`-a` The first directory. Defaults to `./a`.  
`-b` The second directory. Defaults to `./b`.  
`-diff` The save directory. Defaults to `./diff`.  

```shell
$ compare -a <old update source> -b <new update source> -diff <save directory>
```