# Sleep

## Usage

You can send a request to [http://sleep.hui.lu](http://hui.lu), and Sleep server will wait a moment to response.


### Rand

`curl http://sleep.hui.lu -i`


### Specified time

* `curl http://sleep.hui.lu/2`
* `curl http://sleep.hui.lu -H 'X-Sleep: 2'`
