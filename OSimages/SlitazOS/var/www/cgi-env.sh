#!/bin/sh
. /usr/bin/httpd_helper.sh
header

cat << EOT
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<title>CGI SHell Environment</title>
	<meta charset="utf-8" />
	<link rel="stylesheet" type="text/css" href="style.css" />
</head>
<body>

<!-- Header -->
<div id="header">
	<h1>CGI SHell Environment</h1>
</div>

<!-- Content -->
<div id="content">

<p>
	Welcome to the SliTaz web server CGI Shell environment. Let the power of
	SHell script meet the web! Here you can check HTTP info and try some
	requests. 
</p>
<p>
	Including /usr/bin/httpd_helper.sh in your scripts lets you
	use PHP-like syntax such as: \$(GET var)
</p>
<p>
	QUERY_STRING test: 
	<a href="$SCRIPT_NAME?var=value">$SCRIPT_NAME?var=value</a>
</p>

<h2>HTTP Info</h2>
<pre>
$(httpinfo)
</pre>

<!-- End content -->
</div>

</body>
</html>
EOT

exit 0
