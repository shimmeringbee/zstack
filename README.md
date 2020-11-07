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

Z-Stack is a Zigbee Stack made available by Texas Instruments for use on their CC 2.4Ghz SOC processors. This 
library implements a Zigbee Network Processor that is capable of controlling a Z-Stack implementation, specifically it
supports the CC series of Zigbee sniffers flashed with the 
[zigbee2mqtt](https://www.zigbee2mqtt.io/getting_started/flashing_the_cc2531.html) Z-Stack coordinator firmware.

More information about Z-Stack is available from [Texas Instruments](https://www.ti.com/tool/Z-STACK) directly or from
[Z-Stack Developer's Guide](https://usermanual.wiki/Pdf/ZStack20Developers20Guide.1049398016/view).

[Another implementation](https://github.com/dyrkin/znp-go/) of a Z-Stack compatible ZNP exists for Golang, it did [hold no license for a period](https://github.com/dyrkin/zigbee-steward/issues/1)
and the author could not be contacted. This has been rectified, so it may be of interest you. This is a complete
reimplementation of the library, however it is likely there will be strong coincidences due to Golang standards.

## Supported Devices

The following chips and sticks are known to work, though it's likely others in the series will too:

* CC253X
  * Cheap Zigbee Sniffers from AliExpress - CC2531
* CC26X2R1
  * [Electrolama zig-a-sig-ah!](https://electrolama.com/projects/zig-a-zig-ah/) - CC2652R
   
Huge thanks to @Koenkk for his work in providing Z-Stack firmware for these chips. You can [grab the firmware from GitHub](https://github.com/Koenkk/Z-Stack-firmware/).

## Install

Add an import and most IDEs will `go get` automatically, if it doesn't `go build` will fetch.

```go
import "github.com/shimmeringbee/zstack"
```

## Usage

**This libraries API is unstable and should not yet be relied upon.**

### Open Serial Connection and Start ZStack

```go
/* Obtain a ReadWriter UART interface to CC253X */
serialPort :=

/* Construct node table, cache of network nodes. */
t := zstack.NewNodeTable()

/* Create a new ZStack struct. */
z := zstack.New(serialPort, t)

/* Generate random Zigbee network, on default channel (15) */
netCfg, _ := zigbee.GenerateNetworkConfiguration()

/* Obtain context for timeout of initialisation. */
ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
defer cancel()

/* Initialise ZStack and CC253X */)
err = z.Initialise(ctx, nc)
```

### Handle Events

**It is critical that this is handled until you wish to stop the Z-Stack instance.**

```go
for {
    ctx := context.Background()
    event, err := z.ReadEvent(ctx)

    if err != nil {
        return
    }

    switch e := event.(type) {
    case zigbee.NodeJoinEvent:
        log.Printf("join: %v\n", e.Node)
        go exploreDevice(z, e.Node)
    case zigbee.NodeLeaveEvent:
        log.Printf("leave: %v\n", e.Node)
    case zigbee.NodeUpdateEvent:
        log.Printf("update: %v\n", e.Node)
    case zigbee.NodeIncomingMessageEvent:
        log.Printf("message: %v\n", e)
    }
}
```

### Permit Joins

```go
err := z.PermitJoin(ctx, true)
```

### Deny Joins

```go
err := z.DenyJoin(ctx)
```

### Query Device For Details

```go
func exploreDevice(z *zstack.ZStack, node zigbee.Node) {
	log.Printf("node %v: querying", node.IEEEAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Minute)
	defer cancel()

	descriptor, err := z.QueryNodeDescription(ctx, node.IEEEAddress)

	if err != nil {
		log.Printf("failed to get node descriptor: %v", err)
		return
	}

	log.Printf("node %v: descriptor: %+v", node.IEEEAddress, descriptor)

	endpoints, err := z.QueryNodeEndpoints(ctx, node.IEEEAddress)

	if err != nil {
		log.Printf("failed to get node endpoints: %v", err)
		return
	}

	log.Printf("node %v: endpoints: %+v", node.IEEEAddress, endpoints)

	for _, endpoint := range endpoints {
		endpointDes, err := z.QueryNodeEndpointDescription(ctx, node.IEEEAddress, endpoint)

		if err != nil {
			log.Printf("failed to get node endpoint description: %v / %d", err, endpoint)
		} else {
			log.Printf("node %v: endpoint: %d desc: %+v", node.IEEEAddress, endpoint, endpointDes)
		}
	}
}
```

### Node Table Cache

`zstack` requires a `NodeTable` structure to cache a devices IEEE address to its Zibgee network address. A design 
decision for `zstack` was that all operations would reference the IEEE address. This cache must be persisted between 
program runs as the coordinator hardware does not retain this information between restarts.

```go
// Create new table
nodeTable := NewNodeTable()

// Dump current content
nodes := nodeTable.Nodes()

// Load previous content - this should be done before starting ZStack.
nodeTable.Load(nodes)
```

### ZCL

To handle ZCL messages you must handle `zigbee.NodeIncomingMessageEvent` messages and process the ZCL payload with the ZCL library, responses can be sent with `z.SendNodeMessage`.

## Maintainers

[@pwood](https://github.com/pwood)

## Contributing

Feel free to dive in! [Open an issue](https://github.com/shimmeringbee/zstack/issues/new) or submit PRs.

All Shimmering Bee projects follow the [Contributor Covenant](https://shimmeringbee.io/docs/code_of_conduct/) Code of Conduct.

## License

   Copyright 2019-2020 Shimmering Bee Contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.