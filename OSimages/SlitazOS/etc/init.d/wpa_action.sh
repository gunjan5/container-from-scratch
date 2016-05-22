#!/bin/sh
INTERFACE=$1

# For a dynamic IP with DHCP.
if [ "$2" = "CONNECTED" ]; then
     [ -f /var/run/udhcpc.$INTERFACE.pid] && killall udhcpc
	/sbin/udhcpc -b -i $INTERFACE -p /var/run/udhcpc.$INTERFACE.pid
elif [ "$2" = "DISCONNECTED" ]; then
   	/sbin/udhcpc -b -i $INTERFACE -p /var/run/udhcpc.$INTERFACE.pid
fi

