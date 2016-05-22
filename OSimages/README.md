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

### NOTE: 
- Some Linux distributions have initial ram disk compressed in a different format, even though file extension is `.gz`, so you can't use `gunzip` on them.
- Slitaz, for example, has it's initial RAM file `rootfs4.gz` compressed with LZMA
- Best way to find that out is to use the `file` command: 
```bash
$ file rootfs4.gz
rootfs4.gz: LZMA compressed data, streamed
```
- To unzip it, use `unlzma rootfs4.gz -S .gz`