#!/bin/sh
# mktazdevs.sh: Make device files for SliTaz GNU/Linux
# 2007/10/02 - pankso@slitaz.org
#

# Script functions.
status()
{
	local CHECK=$?
	echo -en "\\033[70G[ "
	if [ $CHECK = 0 ]; then
		echo -en "\\033[1;33mOK"
	else
		echo -en "\\033[1;31mFailed"
	fi
	echo -e "\\033[0;39m ]"
}

# We do our work in the dev/ directory.
if [ -z "$1" ] ; then
	echo "usage: `basename $0` path/to/dev"
	exit 1
fi

# Script start.
echo -n "Moving to $1... "
cd $1
status

# Make useful directories.
echo -n "Starting to build directories... "
mkdir pts input net usb shm
status

# Script start.
#
echo -n "Starting to build devices... "

# Input devs.
#
mknod input/event0 c 13 64
mknod input/event1 c 13 65
mknod input/event2 c 13 66
mknod input/mouse0 c 13 32
mknod input/mice c 13 63
mknod input/ts0 c 254 0

# Miscellaneous one-of-a-kind stuff.
#
mknod logibm c 10 0
mknod psaux c 10 1
mknod inportbm c 10 2
mknod atibm c 10 3
mknod console c 5 1
mknod full c 1 7
mknod kmem c 1 2
mknod mem c 1 1
mknod null c 1 3
mknod port c 1 4
mknod random c 1 8
mknod urandom c 1 9
mknod zero c 1 5
mknod rtc c 10 135
mknod sr0 b 11 0
mknod sr1 b 11 1
mknod agpgart c 10 175
mknod ttyS0 c 4 64
mknod audio c 14 4
mknod beep c 10 128
mknod ptmx c 5 2
mknod nvram c 10 144
ln -s /proc/kcore core
# DSP
mknod -m 0666 dsp c 14 3
# PPP dev.
mknod ppp c 108 0

# net/tun device.
#
mknod net/tun c 10 200

# Framebuffer devs.
#
mknod fb0 c 29 0
mknod fb1 c 29 32
mknod fb2 c 29 64
mknod fb3 c 29 96
mknod fb4 c 29 128
mknod fb5 c 29 160
mknod fb6 c 29 192

# usb/hiddev.
#
mknod usb/hiddev0 c 180 96
mknod usb/hiddev1 c 180 97
mknod usb/hiddev2 c 180 98
mknod usb/hiddev3 c 180 99
mknod usb/hiddev4 c 180 100
mknod usb/hiddev5 c 180 101
mknod usb/hiddev6 c 180 102

# IDE HD devs.
# With a few conceivable partitions, you can do
# more of them yourself as you need 'em.
#

# hda devs.
#
mknod hda b 3 0
mknod hda1 b 3 1
mknod hda2 b 3 2
mknod hda3 b 3 3
mknod hda4 b 3 4
mknod hda5 b 3 5
mknod hda6 b 3 6
mknod hda7 b 3 7
mknod hda8 b 3 8
mknod hda9 b 3 9

# hdb devs.
#
mknod hdb b 3 64
mknod hdb1 b 3 65
mknod hdb2 b 3 66
mknod hdb3 b 3 67
mknod hdb4 b 3 68
mknod hdb5 b 3 69
mknod hdb6 b 3 70
mknod hdb7 b 3 71
mknod hdb8 b 3 72
mknod hdb9 b 3 73

# hdc and hdd.
#
mknod hdc b 22 0
mknod hdd b 22 64

# sda devs.
#
mknod sda  b 8 0
mknod sda1 b 8 1
mknod sda2 b 8 2
mknod sda3 b 8 3
mknod sda4 b 8 4
mknod sda5 b 8 5
mknod sda6 b 8 6
mknod sda7 b 8 7
mknod sda8 b 8 8
mknod sda9 b 8 9
ln -s sda1 flash

# sdb devs.
#
mknod sdb b 8 16
mknod sdb1 b 8 17
mknod sdb2 b 8 18
mknod sdb3 b 8 19
mknod sdb4 b 8 20
mknod sdb5 b 8 21
mknod sdb6 b 8 22
mknod sdb7 b 8 23
mknod sdb8 b 8 24
mknod sdb9 b 9 25

# Floppy device.
#
mknod fd0 b 2 0

# loop devs.
#
for i in `seq 0 7`; do
	mknod loop$i b 7 $i
done

# ram devs.
#
for i in `seq 0 7`; do
	mknod ram$i b 1 $i
done
ln -s ram1 ram

# tty devs.
#
mknod tty c 5 0
for i in `seq 0 7`; do
	mknod tty$i c 4 $i
done

# Virtual console screen devs.
#
for i in `seq 0 7`; do
	mknod vcs$i c 7 $i
done
ln -s vcs0 vcs

# Virtual console screen w/ attributes devs.
#
for i in `seq 0 7`; do
	mknod vcsa$i c 7 $(($i + 128))
done
ln -s vcsa0 vcsa

status

# Symlinks.
#
ln -snf /proc/self/fd fd
ln -snf /proc/self/fd/0 stdin
ln -snf /proc/self/fd/1 stdout
ln -snf /proc/self/fd/2 stderr

# Changes permissions.
#
echo -n "Changing permissions on devices... "
chmod 0666 ptmx
chmod 0666 null
chmod 0622 console
chmod 0666 tty*
status

# Script end.
echo -n "All devices are built..."
status

