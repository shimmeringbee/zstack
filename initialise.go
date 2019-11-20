package zstack

import (
	"context"
	"time"
)

func (z *ZStack) Initialise(ctx context.Context) error {
	if err := Retry(ctx, 5*time.Second, 3, func(ctx context.Context) error {
		return z.resetAdapter(ctx, Soft)
	}); err != nil {
		return err
	}

	// Reset (SOFT)

	// Perform Configuration and State Reset - NVRAM - ZCD_NV_STARTUP_OPTION (0x0003) 0x03 (Clear State, Clear Config)

	// Reset (SOFT)

	// Set Logical Type as COORDINATOR - NVRAM - ZCD_NV_LOGICAL_TYPE (0x0087) 0x00 (Coordinator) (01 = Router, 02 = End Device)

	// Reset (SOFT)

	// Enable Network Security - NVRAM - ZCD_NV_SECURITY_MODE (0x0064) 0x01 (Enable Security)

	// Enable distribution of network keys - NVRAM - ZCD_NV_PRECFGKEYS_ENABLE (0x0063) 0x01 (Use precfg keys)

	// Set Network Key - NVRAM - ZCD_NV_PRECFGKEY (0x0062) [16]byte (Set network key)

	// Set ZDO Direct Callback - NVRAM - ZCD_NV_ZDO_DIRECT_CB (0x008f) 0x01 (True, Don't write in ZDO_MSG_CB_INCOMING)

	// Set Channels - NVRAM - ZCD_NV_CHANLIST (0x0084) Bitmap (Setting 1 statically uses that, multiple scan least busy)

	// Set PAN ID - NVRAM - ZCD_NV_PANID (0x0083) [2]byte (Set PAN ID)

	// Set Extended PAN ID - NVRAM - ZCD_NV_EXTPANID (0x002d) [8]byte (Set extended PAN ID)

	// Set Enable TC Link Key - NVRAM - ZCD_NV_USE_DEFAULT_TCLK (0x006d) 0x01 (Enable TC Link Key)

	// Set TC Link Key - NVRAM - ZCD_NV_TCLK_TABLE_START (0x0101) [20]byte (Default TC Link, more complex than just key structure)

	return nil
}
