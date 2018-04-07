Run Service
====


## 数据库

### mongodb

database_url=mongodb://root:GCCTS123@dds-bp18ebd1d16fc5b41515-pub.mongodb.rds.aliyuncs.com:3717,dds-bp18ebd1d16fc5b42738-pub.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-5412007

## logs

logpath=/home/jigangduan/workspace/micro/logs

## Start Services

### Run Consul

```bash
consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
```

### Run Micro Web

```bash
micro web --address=0.0.0.0:3012 1>>${logpath}/web/info.log 2>>${logpath}/web/error.log &
```

### Run Micro API

总台 center

```bash
micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3001 --namespace=com.btdxcx.center.api 1>>${logpath}/api/center/info.log 2>>${logpath}/api/center/error.log &
```

店铺 merchant

```bash
micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3002 --namespace=com.btdxcx.merchant.api 1>>${logpath}/api/merchant/info.log 2>>${logpath}/api/merchant/error.log &
```

小程序 applet

```bash
micro --register_ttl=30 --register_interval=15 api --handler=proxy --address=0.0.0.0:3003 --namespace=com.btdxcx.applet.api 1>>${logpath}/api/applet/info.log 2>>${logpath}/api/applet/error.log &
```

### SRV

account-srv

```bash
account-srv --database_url=${database_url} 1>>${logpath}/srv/account-info.log 2>>${logpath}/srv/account-error.log 
```

jwtauth-srv

```bash
jwtauth-srv --database_url=${database_url} 1>>${logpath}/srv/jwtaut-info.log 2>>${logpath}/srv/jwtaut-error.log
```

shop-srv

```bash
shop-srv --database_url=${database_url} 1>>${logpath}/srv/shop-info.log 2>>${logpath}/srv/shop-error.log
```

### APIs

#### auth apis

- center

```bash
center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.center.api.auth --site_type=center 1>>${logpath}/api/center/auth-info.log 2>>${logpath}/api/center/auth-error.log
```

- merchant

```bash
center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.merchant.api.auth --site_type=back 1>>${logpath}/api/merchant/auth-info.log 2>>${logpath}/api/merchant/auth-error.log
```

- applet

```bash
center-api --register_ttl=30 --register_interval=15 auth --api_service=com.btdxcx.applet.api.auth --site_type=mini 1>>${logpath}/api/applet/auth-info.log 2>>${logpath}/api/applet/auth-error.log
```

#### shops apis

```bash
center-api --register_ttl=30 --register_interval=15 shop 1>>${logpath}/api/center/shop-info.log 2>>${logpath}/api/center/shop-error.log
```

#### taxons apis

- merchant

```bash
merchant --register_ttl=30 --register_interval=15 taxons 1>>${logpath}/api/merchant/taxons-info.log 2>>${logpath}/api/merchant/taxons-error.log
```

