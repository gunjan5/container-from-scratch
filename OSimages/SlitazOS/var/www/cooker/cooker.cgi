#!/bin/sh
#
# SliTaz Cooker CGI/web interface.
#
echo "Content-Type: text/html"
echo ""

[ -f "/etc/slitaz/cook.conf" ] && . /etc/slitaz/cook.conf
[ -f "cook.conf" ] && . ./cook.conf

# The same wok as cook.
wok="$WOK"

# Cooker DB files.
activity="$CACHE/activity"
commits="$CACHE/commits"
cooklist="$CACHE/cooklist"
cookorder="$CACHE/cookorder"
command="$CACHE/command"
blocked="$CACHE/blocked"
broken="$CACHE/broken"
cooknotes="$CACHE/cooknotes"
wokrev="$CACHE/wokrev"

# We're not logged and want time zone to display correct server date.
export TZ=$(cat /etc/TZ)

#
# Functions
#

# RSS feed generator
if [ "$QUERY_STRING" == "rss" ]; then
	pubdate=$(date "+%a, %d %b %Y %X")
	cat << EOT
<?xml version="1.0" encoding="utf-8" ?>
<rss version="2.0">
<channel>
	<title>SliTaz Cooker</title>
	<description>The SliTaz packages cooker feed</description>
	<link>$COOKER_URL</link>
	<lastBuildDate>$pubdate GMT</lastBuildDate>
	<pubDate>$pubdate GMT</pubDate>
EOT
	for rss in $(ls -lt $FEEDS/*.xml | head -n 12)
	do
		cat $rss
	done
	cat << EOT
</channel>
</rss>
EOT
	exit 0
fi

# Put some colors in log and DB files.
syntax_highlighter() {
	case $1 in
		log)
			sed -e 's#OK$#<span class="span-ok">OK</span>#g' \
				-e 's#yes$#<span class="span-ok">yes</span>#g' \
				-e 's#no$#<span class="span-no">no</span>#g' \
				-e 's#error$#<span class="span-red">error</span>#g' \
				-e 's#ERROR:#<span class="span-red">ERROR:</span>#g' \
				-e 's#WARNING:#<span class="span-red">WARNING:</span>#g' \
				-e s"#^Executing:\([^']*\).#<span class='sh-val'>\0</span>#"g \
				-e s"#^====\([^']*\).#<span class='span-line'>\0</span>#"g \
				-e s"#^[a-zA-Z0-9]\([^']*\) :: #<span class='span-sky'>\0</span>#"g \
				-e s"#ftp://\([^']*\).*#<a href='\0'>\0</a>#"g	\
				-e s"#http://\([^']*\).*#<a href='\0'>\0</a>#"g ;;
		receipt)
			sed -e s'|&|\&amp;|g' -e 's|<|\&lt;|g' -e 's|>|\&gt;|'g \
				-e s"#^\#\([^']*\)#<span class='sh-comment'>\0</span>#"g \
				-e s"#\"\([^']*\)\"#<span class='sh-val'>\0</span>#"g ;;
		diff)
			sed -e 's|&|\&amp;|g' -e 's|<|\&lt;|g' -e 's|>|\&gt;|g' \
				-e s"#^-\([^']*\).#<span class='span-red'>\0</span>#"g \
				-e s"#^+\([^']*\).#<span class='span-ok'>\0</span>#"g \
				-e s"#@@\([^']*\)@@#<span class='span-sky'>@@\1@@</span>#"g ;;
		activity)
			sed s"#^\([^']* : \)#<span class='log-date'>\0</span>#"g ;;
	esac
}

# Latest build pkgs.
list_packages() {
	cd $PKGS
	ls -1t *.tazpkg | head -20 | \
	while read file
	do
		echo -n $(stat -c '%y' $PKGS/$file | cut -d . -f 1 | sed s/:[0-9]*$//)
		echo " : $file"
	done
}

# xHTML header. Pages can be customized with a separated html.header file.
if [ -f "header.html" ]; then
	cat header.html
else
	cat << EOT
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<title>SliTaz Cooker</title>
	<meta charset="utf-8" />
	<link rel="shortcut icon" href="favicon.ico" />
	<link rel="stylesheet" type="text/css" href="style.css" />
</head>
<body>

<div id="header">
	<div id="logo"></div>
	<h1><a href="cooker.cgi">SliTaz Cooker</a></h1>
</div>

<!-- Content -->
<div id="content">
EOT
fi

#
# Load requested page
#

case "${QUERY_STRING}" in
	pkg=*)
		pkg=${QUERY_STRING#pkg=}
		log=$LOGS/$pkg.log
		echo "<h2>Package: $pkg</h2>"

		# Package info.
		echo '<div id="info">'
		if [ -f "$wok/$pkg/receipt" ]; then
			echo "<a href='cooker.cgi?receipt=$pkg'>receipt</a>"
			unset WEB_SITE
			. $wok/$pkg/receipt
			[ -n "$WEB_SITE" ] && # busybox wget -s $WEB_SITE &&
			echo "<a href='$WEB_SITE'>home</a>"
			echo "<a href='cooker.cgi?files=$pkg'>files</a>"
		else
			echo "No package named: $pkg"
		fi
		echo '</div>'

		# Check for a log file and display summary if it exists.
		if [ -f "$log" ]; then
			if grep -q "cook:$pkg$" $command; then
				echo "<pre>The Cooker is currently building: $pkg</pre>"
			fi
			if fgrep -q "Summary for:" $LOGS/$pkg.log; then
				echo "<h3>Cook summary</h3>"
				echo '<pre>'
				grep -A 8 "^Summary for:" $LOGS/$pkg.log | sed /^$/d | \
					syntax_highlighter log
				echo '</pre>'
			fi
			if fgrep -q "Debug information" $LOGS/$pkg.log; then
				echo "<h3>Cook failed</h3>"
				echo '<pre>'
				grep -A 8 "^Debug information" $LOGS/$pkg.log | sed /^$/d | \
						syntax_highlighter log
				echo '</pre>'
			fi
			echo "<h3>Cook log</h3>"
			echo '<pre>'
			cat $log | syntax_highlighter log
			echo '</pre>'
		else
			echo "<pre>No log: $pkg</pre>"
		fi ;;
	file=*)
		# Dont allow all files on the system for security reasons.
		file=${QUERY_STRING#file=}
		case "$file" in
			activity|cooknotes|cooklist)
				[ "$file" == "cooklist" ] && \
					nb="- Packages: $(cat $cooklist | wc -l)"
				echo "<h2>DB: $file $nb</h2>"
				echo '<pre>'
				tac $CACHE/$file | syntax_highlighter activity
				echo '</pre>' ;;
			broken)
				nb=$(cat $broken | wc -l)
				echo "<h2>DB: broken - Packages: $nb</h2>"
				echo '<pre>'
				cat $CACHE/$file | sort | \
					sed s"#^[^']*#<a href='cooker.cgi?pkg=\0'>\0</a>#"g
				echo '</pre>' ;;
			*.diff)
				diff=$CACHE/$file
				echo "<h2>Diff for: ${file%.diff}</h2>"
				[ "$file" == "installed.diff" ] && echo \
					"<p>This is the latest diff between installed packages \
					and installed build dependencies to cook.</p>"
				echo '<pre>'
				cat $diff | syntax_highlighter diff
				echo '</pre>' ;;
			*.log)
				log=$LOGS/$file
				name=$(basename $log)
				echo "<h2>Log for: ${name%.log}</h2>"
				if [ -f "$log" ]; then
					if fgrep -q "Summary" $log; then
						echo '<pre>'
						grep -A 20 "^Summary" $log | sed /^$/d | \
							syntax_highlighter log
						echo '</pre>'
					fi
					echo '<pre>'
					cat $log | syntax_highlighter log
					echo '</pre>'
				else
					echo "<pre>No log file: $log</pre>"
				fi ;;
		esac ;;
	receipt=*)
		pkg=${QUERY_STRING#receipt=}
		echo "<h2>Receipt for: $pkg</h2>"
		if [ -f "$wok/$pkg/receipt" ]; then
			echo '<pre>'
			cat $wok/$pkg/receipt | syntax_highlighter receipt
			echo '</pre>'
		else
			echo "<pre>No receipt for: $pkg</pre>"
		fi ;;
	files=*)
		pkg=${QUERY_STRING#files=}
		echo "<h2>Installed files by: $pkg</h2>"
		dir=$(ls -d $WOK/$pkg/taz/$pkg-*)
		if [ -d "$dir/fs" ]; then
			echo '<pre>'
			find $dir/fs -not -type d | xargs ls -ld | \
				sed "s|$dir/fs||" | syntax_highlighter log
			echo '</pre>'
		else
			echo "<pre>No files list for: $pkg</pre>"
		fi ;;
	*)
		# Main page with summary.
		inwok=$(ls $WOK | wc -l)
		cooked=$(ls $PKGS/*.tazpkg | wc -l)
		unbuilt=$(($inwok - $cooked))
		pct=0
		[ $inwok -gt 0 ] && pct=$(( ($cooked * 100) / $inwok ))
		cat << EOT
<div style="float: right;">
	<form method="get" action="$SCRIPT_NAME">
		Package:
		<input type="text" name="pkg" />
	</form>
</div>

<h2>Summary</h2>

<pre>
Running command  : $([ -s "$command" ] && cat $command || echo "Not running")
Wok revision     : <a href="$WOK_URL">$(cat $wokrev)</a>
Commits to cook  : $(cat $commits | wc -l)
Current cooklist : $(cat $cooklist | wc -l)
Broken packages  : $(cat $broken | wc -l)
Blocked packages : $(cat $blocked | wc -l)
</pre>

<p>
	Packages: $inwok in the wok - $cooked cooked - $unbuilt unbuilt -
	Server date: $(date '+%Y-%m-%d %H:%M')
</p>
<div class="pctbar">
	<div class="pct" style="width: ${pct}%;">${pct}%</div>
</div>

<p>
	Latest:
	<a href="cooker.cgi?file=cookorder.log">cookorder.log</a>
	<a href="cooker.cgi?file=commits.log">commits.log</a>
	<a href="cooker.cgi?file=installed.diff">installed.diff</a>
	- Architecture $ARCH:
	<a href="cooker.cgi?pkg=slitaz-toolchain">toolchain.log</a>
</p>

<a name="activity"></a>
<h2>Activity</h2>
<pre>
$(tac $CACHE/activity | head -n 12 | syntax_highlighter activity)
</pre>
<a class="button" href="cooker.cgi?file=activity">More activity</a>

<a name="cooknotes"></a>
<h2>Cooknotes</h2>
<pre>
$(tac $cooknotes | head -n 12 | syntax_highlighter activity)
</pre>
<a class="button" href="cooker.cgi?file=cooknotes">More notes</a>

<a name="commits"></a>
<h2>Commits</h2>
<pre>
$(cat $commits)
</pre>

<a name="cooklist"></a>
<h2>Cooklist</h2>
<pre>
$(cat $cooklist | head -n 20)
</pre>
<a class="button" href="cooker.cgi?file=cooklist">Full cooklist</a>

<a name="broken"></a>
<h2>Broken</h2>
<pre>
$(cat $broken | head -n 20 | sed s"#^[^']*#<a href='cooker.cgi?pkg=\0'>\0</a>#"g)
</pre>
<a class="button" href="cooker.cgi?file=broken">All broken packages</a>

<a name="blocked"></a>
<h2>Blocked</h2>
<pre>
$(cat $blocked | sed s"#^[^']*#<a href='cooker.cgi?pkg=\0'>\0</a>#"g)
</pre>

<a name="lastcook"></a>
<h2>Latest cook</h2>
<pre>
$(list_packages | sed s"#^\([^']*\).* : #<span class='log-date'>\0</span>#"g)
</pre>
EOT
	;;
esac

# Close xHTML page
cat << EOT
</div>

<div id="footer">
	<a href="http://www.slitaz.org/">SliTaz Website</a>
	<a href="cooker.cgi">Cooker</a>
	<a href="http://hg.slitaz.org/cookutils/raw-file/tip/doc/cookutils.en.html">
		Documentation</a>
</div>

</body>
</html>
EOT

exit 0
