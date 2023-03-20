# comPass - compromised password checks
download/check compromised password hashes from haveibeenpwned.com


### Download
Download inspired by https://github.com/HaveIBeenPwned/PwnedPasswordsDownloader  
Uses range search api endpoint (https://haveibeenpwned.com/API/v3#SearchingPwnedPasswordsByRange) to get hashes for all (`0x00000-0xFFFFF`) ranges.

> A range search typically returns approximately 800 hash suffixes, although this number will differ depending on the hash prefix being searched for and will increase as more passwords are added. There are 1,048,576 different hash prefixes between 00000 and FFFFF (16^5) and every single one will return HTTP 200; there is no circumstance in which the API should return HTTP 404.

```bash
xh https://api.pwnedpasswords.com/range/21BD1                                                                                dev-cookie/search 
HTTP/2.0 200 OK                                                                                                                                    
access-control-allow-origin: *
age: 366950
arr-disable-session-affinity: True
cache-control: public, max-age=2678400
cf-cache-status: HIT
cf-ray: 7a3102c8aa5ebaff-MXP
content-type: text/plain
date: Sun, 05 Mar 2023 11:27:38 GMT
etag: W/"0x8DB19B10CBC28C1"
expires: Wed, 05 Apr 2023 11:27:38 GMT
last-modified: Tue, 28 Feb 2023 17:27:21 GMT
request-context: appId=cid-v1:639b3d62-d78b-45f0-8442-2b7f52b50c2e
server: cloudflare
strict-transport-security: max-age=31536000; includeSubDomains; preload
vary: Accept-Encoding
x-content-type-options: nosniff

0018A45C4D1DEF81644B54AB7F969B88D65:3
00D4F6E8FA6EECAD2A3AA415EEC418D38EC:3
011053FD0102E94D6AE2F8B83D76FAF94F6:1
012A7CA357541F0AC487871FEEC1891C49C:3
0136E006E24E7D152139815FB0FC6A50B15:4

```

`If-None-Match` can be used to check if range changed since last request. 
```bash 
xh https://api.pwnedpasswords.com/range/21BD1 If-None-Match:'W/"0x8DB19B10CBC28C1"'
HTTP/2.0 304 Not Modified
access-control-allow-origin: *
arr-disable-session-affinity: True
cache-control: public, max-age=2678400
cf-cache-status: HIT
cf-ray: 7a3a25aa2cdb0f7e-MXP
date: Sun, 05 Mar 2023 11:51:33 GMT
etag: "0x8DB19B10CBC28C1"
expires: Wed, 05 Apr 2023 11:51:33 GMT
last-modified: Tue, 28 Feb 2023 17:27:21 GMT
request-context: appId=cid-v1:639b3d62-d78b-45f0-8442-2b7f52b50c2e
server: cloudflare
set-cookie: __cf_bm=P2IWNHQFrbYPD24XuumU7erqcKhHIUoPLHL5xHUD8V8-1678017093-0-Af3wDk1topj7VSLSfQsCDtu+WPMIWhjzi5sV3WlrHCgoQWr/YkHGCnW57gbQ5feRL+4fSqAhxDC9etuNVeFmWyc=; path=/; expires=Sun, 05-Mar-23 12:21:33 GMT; domain=.pwnedpasswords.com; HttpOnly; Secure; SameSite=None
strict-transport-security: max-age=31536000; includeSubDomains; preload
vary: Accept-Encoding
x-content-type-options: nosniff

```

## todo 
 - bloom build
 - bloom check 
 - store config dir in e.g. `~/.config/compass/config.yml`