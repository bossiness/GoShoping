#!/bin/bash -v

# set -x

cmd=$1
shift

## set config
### logs path
logpath=/home/jigangduan/workspace/micro/logs
env_cloud=YES
if [ ! -d $logpath ] ; then
	logpath=/Users/jigang.duan/WorkSpace/study/go/log
	if [ ! -d $logpath ] ; then
  	echo "$logpath : No such directory"
  	exit 1
	fi
	env_cloud=NO
fi
echo "logpath: $logpath"

### mongodb URL
if [ $env_cloud = YES ] ; then
	database_url=mongodb://root:GCCTS123@dds-bp18ebd1d16fc5b41515-pub.mongodb.rds.aliyuncs.com:3717,dds-bp18ebd1d16fc5b42738-pub.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-5412007
else
	database_url=localhost:27017
fi
echo "database_url: $database_url"


## Start Services
### Run Consul
start_consul() {
	echo 'Run Consul'
	if [ $env_cloud = YES ] ; then
		consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	else
		consul agent -dev 1>>${logpath}/consul.log 2>>${logpath}/consul.error.log &
	fi
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
    echo 'start taxons srv!'
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
	api=$1
  shift
  case $api in
	  auth)
    echo 'start merchant auth api!'
	  center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.merchant.api.auth --site_type=back 1>>${logpath}/api/merchant/info.log 2>>${logpath}/api/merchant/error.log &
	  ;;
	  shop)
    echo 'start merchant shop apis!'
	  merchant --register_ttl=30 --register_interval=15 shop 1>>${logpath}/api/merchant/info.log 2>>${logpath}/api/merchant/error.log &
		isSuccess
	  ;;
		taxons)
		echo 'start merchant taxons apis!'
	  merchant --register_ttl=30 --register_interval=15 taxons 1>>${logpath}/api/merchant/info.log 2>>${logpath}/api/merchant/error.log &
		isSuccess
	  ;;
	  *)
	  echo "run.sh start api merchant <auth|shop|taxons>"
	  exit
	  ;;
  esac
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
		taxons)
		echo 'start applet taxons apis!'
	  applet --register_ttl=30 --register_interval=15 taxons 1>>${logpath}/api/applet/info.log 2>>${logpath}/api/applet/error.log &
		isSuccess
	  ;;
	  *)
	  echo "run.sh start api applet <auth|shop|taxons>"
	  exit
	  ;;
  esac
}

isSuccess() {
	if [ $? -eq 0 ] ; then
		echo 'success!'
	else
		echo 'failure!'
	fi
}

start_micro() {
	sub=$1
  shift
  case $sub in
	  web)
		echo 'run micro web, at post: 3012'
	  micro --register_ttl=30 --register_interval=15 web --address=0.0.0.0:3012 1>>${logpath}/micro-web.log 2>>${logpath}/micro-web.err.log &
	  ;;
	  center_api)
		echo 'run micro api center, at post: 3001'
	  micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3001 --namespace=com.btdxcx.center.api   1>>${logpath}/micro-api.log 2>>${logpath}/micro-api.err.log &
	  ;;
		merchant_api)
		echo 'run micro api merchant, at post: 3002'
	  micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3002 --namespace=com.btdxcx.merchant.api 1>>${logpath}/micro-api.log 2>>${logpath}/micro-api.err.log &
	  ;;
		applet_api)
		echo 'run micro api applet, at post: 3003'
	  micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3003 --namespace=com.btdxcx.applet.api   1>>${logpath}/micro-api.log 2>>${logpath}/micro-api.err.log &
	  ;;
	  *)
	  echo "run.sh start micro <web|center_api|merchant_api|applet_api>"
	  exit
	  ;;
  esac
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
		consul)
	  start_consul
	  ;;
		micro)
	  start_micro $*
	  ;;
	  *)
	  echo "run.sh $0 <srv|api|consul|micro>"
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
