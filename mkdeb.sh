SRC1="controller"
SRC2="traycontroller-config/traycontroller-config"
DST="traycontroller-ubuntu_amd64.deb"
if [ -f $DST ];then
	if [ $SRC1 -nt $SRC2 ];then
		if [ $SRC1 -nt $DST ];then
			dpkg -b ./deb/ $DST
		else
			echo "nothing to do"
		fi
	else
		if [ $SRC2 -nt $DST ];then
			dpkg -b ./deb/ $DST
		else
			echo "nothing to do"
		fi 
	fi
else
	dpkg -b ./deb/ $DST
fi