Overview
========

Switchboard is a link redirect server. It redirects incoming HTTP requests to
other URLs. It therefore functions as a personal link shortener.

Examples
========

With everything set up correctly and the example switchboard.conf:
- `http://s/wp` redirects to English-language Wikipedia's homepage.
- `http://s/wp/Switchboard` is equivalent to typing "Switchboard" into
  English-language Wikipedia's search box.

Requirements
============

- Go 1.1 or later.
- A local server on which port 80 is available (possibly `localhost`).
- The ability to add hostnames mapping to this server (via DNS or
  `/etc/hosts`).

Installation
============

Install the binary and configuration file:

- `go install github.com/nixprime/switchboard`. This builds a `bin/switchboard`
  somewhere in your $GOPATH, which you can move to wherever. (The included
  system startup scripts assume `switchboard` is located in `/usr/local/bin/`.)
- Create `switchboard.conf` somewhere. Use the included `switchboard.conf` as a
  guide. (By default, Switchboard assumes that `switchboard.conf` is located in
  `/etc/`.)

Enable the daemon to run as a non-root user:

- If using one of the system startup scripts, or running switchboard as
  non-root otherwise, allow switchboard to bind ports under 1024:

        sudo apt-get install libcap2-bin
        sudo setcap cap_net_bind_service+ep /usr/local/bin/switchboard

Enable automatic start at boot:

- If your system uses Upstart (e.g., Ubuntu prior to 15.04), and you want to
  use the included Upstart script to start Switchboard automatically, put
  `upstart/switchboard.conf` in `/etc/init/`. The daemon will now start on the
  next boot. To start it manually, run `service switchboard start` as root.

- If your system uses systemd (e.g., Arch Linux or Fedora, or Ubuntu 15.04 or
  later), and you want to use the included systemd script to start Switchboard
  automatically, put `systemd/switchboard.service` in
  `/usr/lib/systemd/system`. To enable the service on boot, run `systemctl
  enable switchboard` as root. To start the service manually, run `systemctl
  start switchboard` as root.

Configure hostname resolution:

- Configure either your host's local hostname resolution (e.g., the
  `/etc/hosts` file) or your local DNS zone, if applicable, to point whatever
  hostname(s) you want Switchboard to use to the IP address of the machine
  running Switchboard. How to configure DNS is beyond the scope of this
  document. A simple set of entries in `/etc/hosts` might be:

        127.0.0.1  s  # map `s` hostname to Switchboard on localhost
        127.0.0.1  m  # map `m` hostname to Switchboard on localhost

License
=======

Switchboard is provided under the MIT license:

> Copyright (C) 2013 Jamie Liu
>
> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the “Software”), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
