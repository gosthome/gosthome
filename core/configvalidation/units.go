package cv

type unitOfMeasurement interface {
	Domain() string
	UnitNames() []string
}

type unit[U any, P interface {
	*U
	unitOfMeasurement
}] struct {
	value float32
}

// frequency = float_with_unit("frequency", "(Hz|HZ|hz)?")
// resistance = float_with_unit("resistance", "(Ω|Ω|ohm|Ohm|OHM)?")
// current = float_with_unit("current", "(a|A|amp|Amp|amps|Amps|ampere|Ampere)?")
// voltage = float_with_unit("voltage", "(v|V|volt|Volts)?")
// distance = float_with_unit("distance", "(m)")
// framerate = float_with_unit("framerate", "(FPS|fps|Fps|FpS|Hz)")
// angle = float_with_unit("angle", "(°|deg)", optional_unit=True)

// decibel = float_with_unit("decibel", "(dB|dBm|db|dbm)", optional_unit=True)
// pressure = float_with_unit("pressure", "(bar|Bar)", optional_unit=True)

// _temperature_c = float_with_unit("temperature", "(°C|° C|°|C)?")
// _temperature_k = float_with_unit("temperature", "(° K|° K|K)?")
// _temperature_f = float_with_unit("temperature", "(°F|° F|F)?")
