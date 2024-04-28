#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

CORELIB_PKG := go-corelibs/go-sqlbuilder
VERSION_TAGS += MAIN
MAIN_MK_SUMMARY := ${CORELIB_PKG}
MAIN_MK_VERSION := v1.0.0

GOTESTS_SKIP += Example
GOTESTS_ARGV := .
CONVEY_EXCLUDE += integration_test

include CoreLibs.mk
