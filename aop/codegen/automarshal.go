package codegen

// AutoMarshal is the interface implemented by structs with weaver.AutoMarshal
// declarations.
type AutoMarshal interface {
	WeaverMarshal(enc *Encoder)
	WeaverUnmarshal(dec *Decoder)
}
