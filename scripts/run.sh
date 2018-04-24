#!/bin/bash

# Just a script to run the demo

cmd=$1
shift

## set config
### logs path
logpath=/home/jigangduan/workspace/micro/logs
### mongodb URL
database_url=mongodb://root:GCCTS123@dds-bp18ebd1d16fc5b41515-pub.mongodb.rds.aliyuncs.com:3717,dds-bp18ebd1d16fc5b42738-pub.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-5412007

## Start Services
### Run Consul
start_consul() {
  consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
}
stop_consul() {
  pkill consul
}

### SRV
start_srv() {
  srv=$1
  case $srv in
	  account)
	  account-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
	  jwtauth)
	  jwtauth-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
    shop)
	  shop-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
	  *)
	  echo "$0 <account|jwtauth|shop>"
	  exit
	  ;;
  esac
}

### API
start_api() {
  echo 'start api cmd'
}

start() {
  sub=$1
  shift
  case $sub in
	  srv)
	  start_srv $*
	  ;;
	  api)
	  start_api
	  ;;
	  *)
	  echo "$0 <srv|api>"
	  exit
	  ;;
  esac
}

stop() {
  echo 'stop cmd'
}

case $cmd in
	start)
	start $*
	;;
	stop)
	stop
	;;
	restart)
	stop
	start
	;;
	*)
	echo "$0 <start|stop|restart> {subcmd}"
	exit
	;;
esac