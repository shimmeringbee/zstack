# Shimmering Bee: Z-Stack

[![license](https://img.shields.io/github/license/shimmeringbee/zstack.svg)](https://github.com/shimmeringbee/zstack/blob/master/LICENSE)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg)](https://github.com/RichardLitt/standard-readme)
[![Actions Status](https://github.com/shimmeringbee/zstack/workflows/test/badge.svg)](https://github.com/shimmeringbee/zstack/actions)

> Implementation of a ZNP and support code designed to interface with Texas Instruments Z-Stack, written in Go.

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Background

Z-Stack is a Zigbee Stack made available by Texas Instruments for use on their CC253X 2.4Ghz SOC processors. This 
library implements a Zigbee Network Processor that is capable of controlling a Z-Stack implementation, specifically it
supports the CC253X series of Zigbee sniffers flashed with the 
[zigbee2mqtt](https://www.zigbee2mqtt.io/getting_started/flashing_the_cc2531.html) Z-Stack coordinator firmware.

More information about Z-Stack is available from [Texas Instruments](https://www.ti.com/tool/Z-STACK) directly or from
[Z-Stack Developer's Guide](https://usermanual.wiki/Pdf/ZStack20Developers20Guide.1049398016/view).

[Another implementation](https://github.com/dyrkin/znp-go/) of a Z-Stack compatible ZNP exists for Golang, it did [hold no license for a period](https://github.com/dyrkin/zigbee-steward/issues/1)
and the author could not be contacted. This has been rectified, so it may be of interest you. This is a complete
reimplementation of the library, however it is likely there will be strong coincidences due to Golang standards.

## Install

Add an import and most IDEs will `go get` automatically, if it doesn't `go build` will fetch.

```go
import "github.com/shimmeringbee/zstack"
```

## Usage

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/shimmeringbee/zstack/issues/new) or submit PRs.

All Shimmering Bee projects follow the [Contributor Covenant](http://contributor-covenant.org/version/1/3/0/) Code of Conduct.

## License

   Copyright 2019 Peter Wood

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.