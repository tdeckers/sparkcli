[![Build Status](https://travis-ci.org/tdeckers/sparkcli.svg?branch=master)](https://travis-ci.org/tdeckers/sparkcli)

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

Download a copy of sparkcli [here](https://bintray.com/tdeckers/sparkcli/sparkcli#files).
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

Download 

For help:

    sparkcli -h

Examples

    # List rooms
    sparkcli room list
    
    # Send message to a room
    sparkcli message create <roomid> <msg>

# Development

See [Development](DEVELOPMENT.md)
