# KLD (Kernel Linker)

[![GoDoc](https://godoc.org/github.com/go-freebsd/kld?status.svg)](https://godoc.org/github.com/go-freebsd/kld)
[![Coverage 87%](https://img.shields.io/badge/coverage-87%25-green.svg)]()
[![FreeBSD 10.3](https://img.shields.io/badge/freebsd-10.3-green.svg)](https://www.freebsd.org/releases/10.3R/announce.html)
[![FreeBSD 11](https://img.shields.io/badge/freebsd-11-green.svg)](https://www.freebsd.org/releases/11.0R/announce.html)
[![FreeBSD HEAD](https://img.shields.io/badge/freebsd-HEAD-green.svg)](https://svnweb.freebsd.org/base/head/)

The FreeBSD operating system has a kernel that allows the user
to load modules after boot. This is done using the kernel linker.

This library enables easy access to the kernel linker to easily
load, unload kernel modules.
