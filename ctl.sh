#!/bin/bash

exe=dailynews
pid=$(ps -ef|grep ${exe}|grep -v "grep $exe"|grep -v "ctl_${exe}.sh"|awk '{print $2}')


export ENV_APPPATH=/root/dailynews
case $1 in
	start)
        if [[ "$pid" != "" ]]; then
            echo "the exe $exe is run, pid-$pid"
            exit
        fi
	  nohup /root/dailynews/${exe} > /root/dailynews/log.txt 2>&1 &
        ;;
	
	stop)
        if [[ "$pid" == "" ]]; then
            echo "the exe $exe is already stop"
            exit
        fi

        kill $pid
	;;
	
	status)

        if [[ "$pid" == "" ]]; then
            echo "the exe $exe is stop"
        else 
            echo "the exe $exe is run, pip-$pid"
        fi

        exit
	;;
esac
