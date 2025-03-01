<p align="center">
  <img width="500" src="https://user-images.githubusercontent.com/58973699/145937515-c8438450-4aa8-4653-88bf-f0a5e91e163d.png">  
</p>

<p align="center">
  <a href="LICENSE.md" target="_blank">
    <img src="https://badgen.net/badge/license/MIT/blue" alt="">
  </a>
  <a href="https://github.com/hnimminh/shield/releases" target="_blank">
    <img src="https://badgen.net/github/tag/hnimminh/shield" alt="">
  </a>
</p>

<p align="center">
  <br>
  <strong>Firewall Controller Daemon</strong>
  <br>
  <code>Simple Daemon implementation with Go, receiving REST API & PUBSUB firewall config command and execution. </code>
  <br><br>
</p>

## Get Start
```
env GOOS=linux GOARCH=amd64 go build -o shield cmd/shield.go
env GOOS=darwin GOARCH=arm64 go build -o shield cmd/shield.go

shield --redisurl tcp://:@127.0.0.1:6379/0
```


## Disclaimer
* It is an experiental based project for a proof of work, *but usable by brave souls who don't mind some unstable API between updates*
* It is not finished yet, and the work is still in progress.
## License
[MIT](./LICENSE)



