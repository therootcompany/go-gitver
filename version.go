package main

// use recently generated version info as a fallback
// for when git isn't present (i.e. go run <url>)
func init() {
	GitRev = "9f05e2304ccd40ac8a6b6bdba176942b475e272f"
	GitVersion = "v1.1.0"
	GitTimestamp = "2019-06-21T00:01:09-06:00"
}
