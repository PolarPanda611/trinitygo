# trinitygo

[![Build Status](https://api.travis-ci.org/PolarPanda611/trinitygo.svg)](https://travis-ci.org/PolarPanda611/trinitygo)
[![Go Report Card](https://goreportcard.com/badge/github.com/PolarPanda611/trinitygo)](https://goreportcard.com/report/github.com/PolarPanda611/trinitygo)
[![GoDoc](https://godoc.org/github.com/PolarPanda611/trinitygo?status.svg)](https://godoc.org/github.com/PolarPanda611/trinitygo)
[![Release](https://img.shields.io/github/release/PolarPanda611/trinitygo.svg?style=flat-square)](https://github.com/PolarPanda611/trinitygo/releases)

golang restframework plugin with gin+gorm, fast and high scalable    

## Installation

```bash
$ go get github.com/PolarPanda611/trinitygo/trinitygo
$ trinitygo NewHttp [Your Project Name]
$ cd [Your Project Name] 
$ trinitygo NewCrud [Your Model Name]
$ swag init

// start your journey in Trinity
// you can check the demo under example folder 
```

## Overview

* Declarative router
* IOC container & Dependency Injection
* customize middleware
* customize runtime (Tracing Analysis, user authentication , event bus ...)
* support atomic request 
* support customize validator ( API permission , data validation ...)
* support URL query analyze (search , filter , order by , preload ...)
* integrate gorm
* integrate gin

## Getting Started

* Trinity Guides [wiki](https://github.com/PolarPanda611/trinitygo/wiki)

## Contributing

Feel free to create the Issue and PR . We need your help !