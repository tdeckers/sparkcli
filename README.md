[![Build Status](https://travis-ci.org/tdeckers/sparkcli.svg?branch=master)](https://travis-ci.org/tdeckers/sparkcli) [ ![Download](https://api.bintray.com/packages/tdeckers/sparkcli/sparkcli/images/download.svg) ](https://bintray.com/tdeckers/sparkcli/sparkcli/_latestVersion)

# Overview

sparkcli (say 'sparkly' :) is a command line interface to Cisco Spark.

# Setup - Configuration

**1. Obtain a Cisco Spark developer API key and secret**

Define a new Cisco Spark integration app here: [https://developer.ciscospark.com](https://developer.ciscospark.com/).  Fill in the fields as desired, with the exception of:
   
* App icon: Feel free to use `http://files.ducbase.com/spark.png` or use your own.
* Redirect Url: `http://files.ducbase.com/code.html`.
* Scopes: check all boxes.

You'll be provided with a `ClientID` and `ClientSecret`.  You'll need these for the 
   next step.

**2. Authorize**

Navige to the [Sparkcli authorization](http://files.ducbase.com/authorize.html) page 
and follow the instructions there.  Once you've obtained the authorization code
continue here.

**3. Configure**

Create a configuration file called `sparkcli.toml`.  This file is in 
[toml format](https://godoc.org/github.com/BurntSushi/toml).  Sparkcli will look for the 
file in these locations (in order):

* current working directory
* `/etc/sparkcli`
* users' home directory

Add the `ClientID`, `ClientSecret` and the `AuthCode` from the previous steps in the file:

    # cat ./sparkcli.toml
    ClientId = "C23d70022b9e6c4b348897daac846xf694e7f8ffa3cd38986c6974433def69784"
    ClientSecret = "dcca20a5b5cc89fbea1f2b3cd41x80248ff698277583bce69fa63923ef02dc64"
    AuthCode = "46cd20fe32936af96ecb385772896ff84208x14dc3b2f1aebf98cc5b5a02accd"

**4. Login**

Download a copy of sparkcli:  [ ![Download](https://api.bintray.com/packages/tdeckers/sparkcli/sparkcli/images/download.svg) ](https://bintray.com/tdeckers/sparkcli/sparkcli/_latestVersion).
Then run

    sparkcli login

This will update your configuration file with the neccesary tokens for Sparkcli
to authenticate against the Cisco Spark service.  If you use SparkCli frequent enough 
(once every 90 or so days at least), tokens will be refreshed and kept up to date 
as needed.

_**Note**: If Sparkcli gets confused and can't login for some reason, likely the easiest solution is
to remove followling fields - AuthCode, AccessToken, RefreshToken - from sparkcli.toml 
and restart from step 2 above._

# Usage

You'll notice that most commands have a short hand script which is listed below 
the long version

## Global arguments

    sparkcli -h

> Get help

    sparkcli -j=false ...

> Formats the results (if any) in a human readable format.  If this options is 
> set to true or not present the return value(s) as JSON.

## Rooms

List all rooms

    sparkcli rooms list
    sparkcli r l

> Lists all rooms you're subscribed too.

Create room

    sparkcli rooms create <name>
    sparkcli r c <name>

> Creates a room with the name specified.  The room name can include multiple words. 
> If -j=false, only the room id is printed so it can be assigned to a variable.

Get a specific room

    sparkcli rooms get <id>
    sparkcli r g <id>
    
    # using the default room
    sparkcli r g

> Gets details for the room.  If no id is provided, this command uses the default 
> room id if one is available in the config.  See how to set a default room in 
> the config later.

Delete a room

    sparkcli rooms delete <id>
    sparkcli r d <id>

> Deletes the room.

Set the default room

    sparkcli rooms default <id>
    
    # no short for default!
    sparkcli r default <id> 

> Saves a default room id to the config for use in other operations that support it.
> This won't check if the room actually exists. If no id is provided, this will 
> just diplay the saved room id.

## Messages

List messages

    sparkcli messages list <roomid>
    sparkcli m l <roomid>
    
    # to use the default room
    sparkcli m l

> List the messages for a given room.  If no room id is provided, the default room
> will be used if one exists.

Create message

    sparkcli messages create <roomid> <msg>
    sparkcli m c <roomid> <msg>
    
    # to post in the default room
    sparkcli m c - <msg>
    
> Creates a message is the specified room.  For posting to the default room, use
> a dash (-).

Get a message

    sparkcli messages get <id>
    sparkcli m g <id>
    
> Gets a messages' details.

Delete a message

    sparkcli messages delete <id>
    sparkcli m d <id>

> Deletes a messages.

## Other

Login

    sparkcli login

> Logs you into the Cisco Spark service, and stores access tokens on success.

# Development

See [Development](DEVELOPMENT.md)
