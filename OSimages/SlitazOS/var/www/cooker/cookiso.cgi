#!/bin/sh
#
# SliTaz Cookiso CGI/web interface.
#
echo "Content-Type: text/html"
echo ""

[ -f "/etc/slitaz/cook.conf" ] && . /etc/slitaz/cook.conf
[ -f "cook.conf" ] && . ./cook.conf

# Cookiso DB files.
cache="$CACHE/cookiso"
iso="$SLITAZ/iso"
activity="$cache/activity"
command="$cache/command"
rollog="$cache/rolling.log"
synclog="$cache/rsync.log"

#
# Functions
#

# Put some colors in log and DB files.
syntax_highlighter() {
	case $1 in
		log)
			sed -e 's#OK#<span class="span-ok">OK</span>#g' \
				-e 's#Failed#<span class="span-red">Failed</span>#g' \
				-e 's|\(Filesystem size:\).*G\([0-9\.]*M\) *$|\1 \2|' \
				-e 's|.\[1m|<b>|' -e 's|.\[0m|</b>|' -e 's|.\[[0-9Gm;]*||g' ;;
		activity)
			sed s"#^\([^']* : \)#<span class='log-date'>\0</span>#"g ;;
	esac
}

# Latest build pkgs.
list_isos() {
	cd $iso
	ls -1t *.iso | head -6 | \
	while read file
	do
		echo -n $(stat -c '%y' $file | cut -d . -f 1 | sed s/:[0-9]*$//)
		echo " : $file"
	done
}

# xHTML header. Pages can be customized with a separate html.header file.
if [ -f "header.html" ]; then
	cat header.html | sed s'/Cooker/ISO Cooker/'
else
	cat << EOT
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<title>SliTaz ISO Cooker</title>
	<meta charset="utf-8" />
	<link rel="shortcut icon" href="favicon.ico" />
	<link rel="stylesheet" type="text/css" href="style.css" />
</head>
<body>

<div id="header">
	<div id="logo"></div>
	<h1><a href="cookiso.cgi">SliTaz ISO Cooker</a></h1>
</div>

<!-- Content -->
<div id="content">
EOT
fi

#
# Load requested page
#

case "${QUERY_STRING}" in
	distro=*)
		distro=${QUERY_STRING#distro=}
		ver=${distro%-core-4in1}
		log=$iso/slitaz-$ver.log
		. $SLITAZ/flavors/${distro#*-}/receipt
		echo "<h2>Distro: $distro</h2>"
		echo "<p>Description: $SHORT_DESC</p>"
		echo '<h3>Summary</h3>'
		echo '<pre>'
		fgrep "Build time" $log
		fgrep "Build date" $log
		fgrep "Packages" $log
		fgrep "Rootfs size" $log
		fgrep "ISO image size" $log
		echo '</pre>'
		echo '<h3>Cookiso log</h3>'
		echo '<pre>'
		cat $log | syntax_highlighter log
		echo '</pre>' ;;
	*)
		# Main page with summary.
		echo -n "Running command  : "
		if [ -f "$command" ]; then
			cat $command
		else
			echo "Not running"
		fi
		cat << EOT
<h2>Activity</h2>
<pre>
$(tac $activity | head -n 12 | syntax_highlighter activity)
</pre>

<h2>Latest ISO</h2>
<pre>
$(list_isos | syntax_highlighter activity)
</pre>
EOT
		# Rolling Bot log.
		if [ -f "$rollog" ]; then
			echo "<h2>Rolling log</h2>"
			echo '<pre>'
			cat $rollog
			echo '</pre>'
		fi
		# Rsync log.
		if [ -f "$synclog" ]; then
			echo "<h2>Rsync log</h2>"
			echo '<pre>'
			awk '{
	if (/\/s/) h=$0; 
	else {
		if (h!="") print h;
		h="";
		print;
	}
}'< $synclog
			echo '</pre>'
		fi ;;
esac

# Close xHTML page
cat << EOT
</div>

<div id="footer">
	<a href="http://www.slitaz.org/">SliTaz Website</a>
	<a href="cookiso.cgi">Cookiso</a>
	<a href="http://hg.slitaz.org/cookutils/raw-file/tip/doc/cookutils.en.html">
		Documentation</a>
</div>

</body>
</html>
EOT

exit 0
