#!/bin/bash

cmd=$1
shift

## set config
### logs path
logpath=/home/jigangduan/workspace/micro/logs
if [ ! -d $logpath ] ; then
  echo "$logpath : No such directory"
  exit 1
fi
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
  shift
  case $srv in
	  account)
    echo 'start account srv!'
	  account-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
	  jwtauth)
    echo 'start jwtauth srv!'
	  jwtauth-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
    shop)
    echo 'start shop srv!'
	  shop-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
		taxons)
    echo 'start jwtauth srv!'
	  taxons-srv --database_url=${database_url} 1>>${logpath}/srv/info.log 2>>${logpath}/srv/error.log &
	  ;;
	  *)
	  echo "run.sh start $0 <account|jwtauth|shop|taxons>"
	  exit
	  ;;
  esac
}

### APIs
start_api() {
  endpoint=$1
  shift
  case $endpoint in
	  center)
	  center_apis $*
	  ;;
	  merchant)
	  merchant_apis $*
	  ;;
    applet)
	  applet_apis $*
	  ;;
	  *)
	  echo "run.sh start $0 <center|merchant|applet>"
	  exit
	  ;;
  esac
}

center_apis() {
  api=$1
  shift
  case $api in
	  auth)
    echo 'start center auth api!'
	  center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.center.api.auth --site_type=center 1>>${logpath}/api/center/info.log 2>>${logpath}/api/center/error.log &
	  ;;
	  shop)
    echo 'start center shop apis!'
	  center-api --register_ttl=30 --register_interval=15 shop 1>>${logpath}/api/center/info.log 2>>${logpath}/api/center/error.log &
	  ;;
	  *)
	  echo "run.sh start api center <auth|shop>"
	  exit
	  ;;
  esac
}

merchant_apis() {
  echo 'start merchant apis'
}

applet_apis() {
  api=$1
  shift
  case $api in
	  auth)
    echo 'start applet auth api!'
	  center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.applet.api.auth --site_type=mini 1>>${logpath}/api/applet/info.log 2>>${logpath}/api/applet/error.log &
	  ;;
	  shop)
    echo 'start applet shop apis!'
	  applet --register_ttl=30 --register_interval=15 shop 1>>${logpath}/api/applet/info.log 2>>${logpath}/api/applet/error.log &
		isSuccess
	  ;;
		echo 'start applet taxons apis!'
	  applet --register_ttl=30 --register_interval=15 taxons 1>>${logpath}/api/applet/info.log 2>>${logpath}/api/applet/error.log &
		isSuccess
	  ;;
	  *)
	  echo "run.sh start api center <auth|shop|taxons>"
	  exit
	  ;;
  esac
}

isSuccess() {
	if [ $? -eq 0 ] ; then
		echo 'start center shop apis success!'
	else
		echo 'start center shop apis failure!'
	fi
}

start() {
  sub=$1
  shift
  case $sub in
	  srv)
	  start_srv $*
	  ;;
	  api)
	  start_api $*
	  ;;
	  *)
	  echo "run.sh $0 <srv|api>"
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
