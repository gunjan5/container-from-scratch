#!/bin/sh
# /etc/init.d/system.sh - SliTaz hardware configuration.
#
# This script configures the sound card and screen. Tazhw is used earlier
# at boot time to autoconfigure PCI and USB devices. It also configures
# system language, keyboard and TZ in live mode.
#
. /etc/init.d/rc.functions

# Parse cmdline args for boot options (See also rcS and bootopts.sh).
XARG=""
for opt in $(cat /proc/cmdline)
do
	case $opt in
		console=*)
			sed -i "s/tty1/${opt#console=}/g;/^tty[2-9]::/d" \
				/etc/inittab ;;
		sound=*)
			DRIVER=${opt#sound=} ;;
		xarg=*)
			XARG="$XARG ${opt#xarg=}" ;;
		screen=text)
				SCREEN=text
				# Disable X.
				echo -n "Disabling X login manager: slim..."
				. /etc/rcS.conf
				RUN_DAEMONS=$(echo $RUN_DAEMONS | sed s/' slim'/''/)
				sed -i s/"RUN_DAEMONS.*"/"RUN_DAEMONS=\"$RUN_DAEMONS\"/" /etc/rcS.conf
				status ;;
		screen=*)
			SCREEN=${opt#screen=} ;;
		*)
			continue ;;
	esac
done

# Sound configuration stuff. First check if sound=no and remove all
# sound Kernel modules.
if [ -n "$DRIVER" ]; then
	case "$DRIVER" in
	no)
		echo -n "Removing all sound kernel modules..."
		rm -rf /lib/modules/`uname -r`/kernel/sound
		status
		echo -n "Removing all sound packages..."
		for i in $(grep -l '^DEPENDS=.*alsa-lib' /var/lib/tazpkg/installed/*/receipt) ; do
			pkg=${i#/var/lib/tazpkg/installed/}
			echo 'y' | tazpkg remove ${pkg%/*} > /dev/null
		done
		for i in alsa-lib mhwaveedit asunder libcddb ; do
			echo 'y' | tazpkg remove $i > /dev/null
		done
		status ;;
	noconf)
		echo "Sound configuration was disabled from cmdline..." ;;
	*)
		if [ -x /usr/sbin/soundconf ]; then
			echo "Using sound kernel module $DRIVER..."
			/usr/sbin/soundconf -M $DRIVER
		fi ;;
	esac
# Sound card may already be detected by PCI-detect.
elif [ -d /proc/asound ]; then
	# Restore sound config for installed system.
	if [ -s /var/lib/alsa/asound.state ]; then
		echo -n "Restoring last alsa configuration..."
		alsactl restore
		status
	else
		/usr/sbin/setmixer
	fi
	# Start soundconf to config driver and load module for Live mode
	# if not yet detected.
	/usr/bin/amixer >/dev/null || /usr/sbin/soundconf
else
	echo "Unable to configure sound card."
fi

# Start TazPanel
[ -x /usr/bin/tazpanel ] && tazpanel start

# Auto recharge packages list (after network connection of course)
[ "$RECHARGE_PACKAGES_LIST" == "yes" ] && tazpkg recharge &

# Locale config. Do a gui config for both lang/keymap.
echo "Checking if /etc/locale.conf exists... "
if [ ! -s "/etc/locale.conf" ]; then
	if [ "$SCREEN" != "text" ] && [ -x /usr/bin/Xorg ]; then
		echo "Starting TazBox configuration..."
		DISPLAY=:1 tazbox boot
	else
		tazlocale
	fi
else
	lang=$(cat /etc/locale.conf | fgrep LANG | cut -d "=" -f 2)
	echo -n "Locale configuration: $lang" && status
fi

# Keymap config.
if [ -s "/etc/keymap.conf" ]; then
	KEYMAP=$(cat /etc/keymap.conf)
	echo "Keymap configuration: $KEYMAP"
	if [ -x /bin/loadkeys ]; then
		loadkeys $KEYMAP
	else
		loadkmap < /usr/share/kmap/$KEYMAP.kmap
	fi
else
	tazkeymap
fi

# Timezone config. Set timezone using the keymap config for fr, be, fr_CH
# and ca with Montreal.
if [ ! -s "/etc/TZ" ]; then
	map=$(cat /etc/keymap.conf)
	case "$map" in
		fr-latin1|be-latin1)
			echo "Europe/Paris" > /etc/TZ ;;
		fr_CH-latin1|de_CH-latin1)
			echo "Europe/Zurich" > /etc/TZ ;;
		cf)
			echo "America/Montreal" > /etc/TZ ;;
		*)
			echo "UTC" > /etc/TZ ;;
	esac
fi

# Xorg auto configuration.
if [ "$SCREEN" != "text" -a ! -s /etc/X11/xorg.conf -a -x /usr/bin/Xorg ]; then
	echo "Configuring Xorg..."
	# $HOME is not yet set.
	HOME=/root
	sed -i 's|/usr/bin/Xvesa|/usr/bin/Xorg|' /etc/slim.conf
	sed -i s/"^xserver_arguments"/'\#xserver_arguments'/ /etc/slim.conf
	tazx config-xorg 2>/var/log/xorg.configure.log
fi

# Start X sesssion as soon as possible in Live/frugal mode. HD install
# can use FAST_BOOT_X which starts X beforehand. In live mode we need
# keymap config for Xorg configuration and a working Xorg config.
#if [ "$SCREEN" != "text" ] && [ -x /usr/bin/slim ]; then
	#if fgrep -q root=/dev/null /proc/cmdline; then
		#/etc/init.d/slim start
	#fi
#fi

# Firefox hack to get the right locale.
if fgrep -q "fr_" /etc/locale.conf; then
	# But is the fox installed ?
	if [ -f "/var/lib/tazpkg/installed/firefox/receipt" ]; then
		. /var/lib/tazpkg/installed/firefox/receipt
		sed -i 's/en-US/fr/' /etc/firefox/pref/firefox-l10n.js
	fi
fi
