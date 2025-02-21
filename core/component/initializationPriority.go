package component

//go:generate go-enum

// ENUM(
// bus,// For communication buses like i2c/spi
// io,// For components that represent GPIO pins like PCF8573
// hardware,// For components that deal with hardware and are very important like GPIO switch
// data,// For components that import data from directly connected sensors like DHT.
// processor,// For components that use data from sensors like displays
// bluetooth,
// after_bluetooth,
// wifi,
// ethernet,
// before_connection,// For components that should be initialized after WiFi and before API is connected.
// after_wifi,// For components that should be initialized after WiFi is connected.
// after_connection,// For components that should be initialized after a data connection (API/MQTT) is connected.
// late,// For components that should be initialized at the very end of the setup process.
// )
type InitializationPriority int
