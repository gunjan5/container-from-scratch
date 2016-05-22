# Images 

## How to build an image

- Download the iso image for linux distro you want (`Core` Linux in this example)
- Mount it to your host (`sudo mount Core-current.iso tmp`)
- Copy the cpio zipped file to a directory `<img_name>` (`core.gz` in this case)
- Unzip it ( `gunzip core.gz`)
- `sudo cpio -i < core`
-  Do an `ls` to make sure you see the familiar Linux file system

OR 

- Single command: `zcat core.gz | sudo cpio -i`
