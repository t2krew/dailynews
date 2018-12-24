#!/bin/bash

exe=dailynews
pid=$(ps -ef|grep ${exe}|grep -v "grep $exe"|grep -v "ctl_${exe}.sh"|awk '{print $2}')

case $1 in
	start)
        if [[ "$pid" != "" ]]; then
            echo "the exe $exe is run, pid-$pid"
            exit
        fi
            nohup ./${exe}  >./${exe}_demon.log 2>&1 &
        ;;
	stop)
        if [[ "$pid" == "" ]]; then
            echo "the exe $exe is already stop"
            exit
        fi

        kill ${pid}
		;;
	restart|force-reload)
		$0 stop
		sleep 1
		$0 start
		;;

    status)

       if [[ "$pid" == "" ]]; then
            echo "the exe $exe is stop"
       else
            echo "the exe $exe is run, pip-$pid"
       fi

       exit
	;;
	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
		;;
esac
