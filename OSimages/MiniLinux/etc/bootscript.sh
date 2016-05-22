#!/bin/sh

dmesg -n 1
mount -t devtmpfs none /dev
mount -t proc none /proc
mount -t sysfs none /sys

for DEVICE in /sys/class/net/* ; do
  ip link set ${DEVICE##*/} up
  [ ${DEVICE##*/} != lo ] && udhcpc -b -i ${DEVICE##*/} -s /etc/rc.dhcp
done

