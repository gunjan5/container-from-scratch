# ~/.xinitrc: Executed by slim login manager to startx X session.
# You can use F1 with Slim to change your window manager or configure
# it permanently with your personal applications.conf file.
#
. $HOME/.config/slitaz/applications.conf

case $1 in
	e17|enlightenment*)
		exec enlightenment_start ;;
	openbox|openbox-session|ob)
		exec openbox-session ;;
	dwm|dwm-session)
		exec dwm-session ;;
	fluxbox|startfluxbox)
		exec startfluxbox ;;
	awesome)
		exec awesome ;;
	pekwm)
		exec pekwm ;;
	karmen|karmen-session)
		exec karmen-session ;;
	jwm)
		lxpanel &
		exec jwm ;;
	xfce|xfce4|xfce4-session)
		xfce4-session ;;
	*)
		exec $WINDOW_MANAGER ;;
esac
