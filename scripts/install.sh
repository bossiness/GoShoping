#!/bin/bash

cmd=$1
shift

srv() {
    echo 'install srv'
}

api() {
    echo 'install api'
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
	echo "$0 api <center|merchant|applet> {subcmd}"
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