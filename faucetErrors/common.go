package faucetErrors

var INVALID_REQUEST_ERROR = NewBaseError(400, "invalid request")

var TIME_ERROR = NewBaseError(401, "You can claim only once an hour")
