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
	taxons_srv
	;;
	member)
	member_srv
	;;
	*)
	echo "$0 srv <account|jwtauth|order|product|shop|taxons|member>"
	exit
	;;
	esac
}

account_srv() {
	if [ -f "${GOPATH}/bin/account-srv" ] ; then
		echo 'clean account srv program.'
		rm -fv ${GOPATH}/bin/account-srv
	fi
	go install -x btdxcx.com/micro/account-srv
}

jwtauth_srv() {
	if [ -f "${GOPATH}/bin/jwtauth-srv" ] ; then
		echo 'clean jwtauth srv program.'
		rm -fv ${GOPATH}/bin/jwtauth-srv
	fi
	go install -x btdxcx.com/micro/jwtauth-srv
}

order_srv() {
	if [ -f "${GOPATH}/bin/order-srv" ] ; then
		echo 'clean order srv program.'
		rm -fv ${GOPATH}/bin/order-srv
	fi
	go install -x btdxcx.com/micro/order-srv
}

product_srv() {
	if [ -f "${GOPATH}/bin/product-srv" ] ; then
		echo 'clean product srv program.'
		rm -fv ${GOPATH}/bin/product-srv
	fi
	go install -x btdxcx.com/micro/product-srv
}

shop_srv() {
	if [ -f "${GOPATH}/bin/shop-srv" ] ; then
		echo 'clean shop srv program.'
		rm -fv ${GOPATH}/bin/shop-srv
	fi
	go install -x btdxcx.com/micro/shop-srv
}

taxons_srv() {
	if [ -f "${GOPATH}/bin/taxons-srv" ] ; then
		echo 'clean taxons srv program.'
		rm -fv ${GOPATH}/bin/taxons-srv
	fi
	go install -x btdxcx.com/micro/taxons-srv
}

member_srv() {
	if [ -f "${GOPATH}/bin/member-srv" ] ; then
		echo 'clean member srv program.'
		rm -fv ${GOPATH}/bin/member-srv
	fi
	go install -x btdxcx.com/micro/member-srv
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
	if [ -f "${GOPATH}/bin/applet" ] ; then
		echo 'clean applet api program.'
		rm -fv ${GOPATH}/bin/applet
	fi
    go install -x btdxcx.com/applet
}

merchant() {
	if [ -f "${GOPATH}/bin/merchant" ] ; then
		echo 'clean merchant api program.'
		rm -fv ${GOPATH}/bin/merchant
	fi
  go install -x btdxcx.com/merchant
}

center() {
	if [ -f "${GOPATH}/bin/center-api" ] ; then
		echo 'clean center api program.'
		rm -fv ${GOPATH}/bin/center-api
	fi
  go install -x btdxcx.com/center/center-api
}


web () {
	if [ -f "${GOPATH}/bin/shop-web" ] ; then
		echo 'clean shop web program.'
		rm -fv ${GOPATH}/bin/cshop-web
	fi
	go install -x btdxcx.com/shop/shop-web
}

case $cmd in
	srv)
	srv $*
	;;
	api)
	api $*
	;;
	web)
	web
	;;
	*)
	echo "$0 <srv|api|web> {subcmd}"
	exit
	;;
esac