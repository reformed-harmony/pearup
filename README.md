## pearup

[![Build Status](https://ci.quickmediasolutions.com/buildStatus/icon?job=pearup)](https://ci.quickmediasolutions.com/job/pearup/)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/reformed-harmony/pearup?status.svg)](https://pkg.go.dev/github.com/reformed-harmony/pearup)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

This web application provides a simplified interface for creating "pear-ups". The term "pear-up" refers to a group event in which each participant is paired with a member of the opposite gender to get to know through Messenger.

### Building pearup

To build the application, you will need the Go compiler installed. Once installed, execute the following commands:

    go generate
    go build

Assuming all goes well, you should end up with an executable in the current directory named `pearup` or `pearup.exe`. You can use the `--help` flag to learn which arguments should be passed to the application.
