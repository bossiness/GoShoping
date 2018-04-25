#!/bin/bash

cmd=$1
shift

srv() {
	subcmd=$1
    shift
    case $subcmd in
	account)
	account_srv
	;;
	jwtauth)
	jwtauth_srv
	;;
    order)
	order_srv
	;;
	product)
	product_srv
	;;
	shop)
	shop_srv
	;;
	taxons)
	taxons_str
	;;
	*)
	echo "$0 srv <account|jwtauth|order|product|shop|taxons>"
	exit
	;;
	esac
}

account_srv() {
	go install -x btdxcx.com/micro/account-srv
}

jwtauth_srv() {
	go install -x btdxcx.com/micro/jwtauth-srv
}

order_srv() {
	go install -x btdxcx.com/micro/order-srv
}

product_srv() {
	go install -x btdxcx.com/micro/product-srv
}

shop_srv() {
	go install -x btdxcx.com/micro/shop-srv
}

taxons_str() {
	go install -x btdxcx.com/micro/taxons-srv
}

api() {
    subcmd=$1
    shift
    case $subcmd in
	center)
	center
	;;
	merchant)
	merchant
	;;
    applet)
	applet
	;;
	*)
	echo "$0 api <center|merchant|applet>"
	exit
	;;
	esac
}

applet() {
    go install -x btdxcx.com/applet
}

merchant() {
    go install -x btdxcx.com/merchant
}

center() {
    go install -x btdxcx.com/center/center-api
}

case $cmd in
	srv)
	srv $*
	;;
	api)
	api $*
	;;
	*)
	echo "$0 <srv|api> {subcmd}"
	exit
	;;
esac