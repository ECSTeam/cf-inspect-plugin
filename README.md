# Cloud Foundry Get Events CLI Plugin

Cloud Foundry plugin to view inspect deployed apps.

## Install

```
$ go get github.com/ECSTeam/cf-inspect-plugin
$ cf install-plugin $GOPATH/bin/cf-inspect-plugin
```

## Motivation 

_Why do I need this plugin?_ 

A client needed a list containing the apps, their guids, 
the org name, its guid, and space name, and its guid as a list. 

## Usage

```
 $> cf inspect --help
NAME:
   inspect - Get list of applications (by akoranne@ecsteam.com)

Usage: cf inspect [--json]
              --json : list output in json format (default is csv)
```

## Access 

The `inspect` plugin will show apps for the orgs and spaces that the current user has access to.

## Sample Output

```
```

## Uninstall

```
 $> cf uninstall-plugin inspect
```

