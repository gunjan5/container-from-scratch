# ~/.xsession: Start X window session manually on your system (startx).
#
. $HOME/.config/slitaz/applications.conf
Xorg &
#xterm &
#xpad &
exec $WINDOW_MANAGER
