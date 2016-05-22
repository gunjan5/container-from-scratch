#!/bin/sh

[ ! -d ..$QUERY_STRING ] && echo "HTTP/1.1 404 Not Found" || cat <<EOT  
Content-type: text/html

<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<title>Index of $QUERY_STRING</title>
	<meta charset="utf-8" />
	<link rel="stylesheet" type="text/css" href="/style.css" />
</head>

<!-- Header -->
<div id="header">
	<h1>Index of $QUERY_STRING</h1>
</div>

<!-- Content -->
<div id="content">
<body>
	<ul>
$({ [ "$QUERY_STRING" != "/" ] && echo "../"; ls -p ..$QUERY_STRING; } | \
  sed 's|.*|		<li><a href="&">&</a></li>|')
	</ul>
</div>

<!-- Footer -->
<div id="footer">
    Copyright &copy; $(date +%Y) <a href="http://www.slitaz.org/">SliTaz GNU/Linux</a>
</div>

</body>
</html>
EOT
