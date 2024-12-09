package berkeley

import (
	"fmt"
	"time"
)

<<<<<<< HEAD
=======
// FollowerState representa el estado actual de un seguidor.
type FollowerState string

const (
	// Estados posibles para el seguidor
	RequestNotSent FollowerState = "REQUEST_NOT_SENT"
	Responded      FollowerState = "RESPONDED"
	Failed         FollowerState = "FAILED"
)

>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
// FollowerInfo encapsula información sobre un nodo seguidor en el sistema.
type FollowerInfo struct {
	Name              string        // Nombre del seguidor
	Address           string        // Dirección del seguidor en formato "host:puerto"
	State             FollowerState // Estado actual del seguidor
	DateFollower      time.Time     // Fecha y hora local del seguidor
	LocalTime         int64         // Hora local del seguidor en milisegundos desde la época UNIX
	CommunicationTime int64         // Tiempo de comunicación entre líder y seguidor en milisegundos
	TripTime          int64         // Tiempo de ida y vuelta (triptime) estimado
	DiffTime          int64         // Diferencia de tiempo calculada entre líder y seguidor
	Delta             int64         // Diferencia global de tiempo (delta) aplicada al seguidor
}

// NewFollowerInfo crea un nuevo objeto FollowerInfo con los valores proporcionados.
func NewFollowerInfo(address, name string, localTime, communicationTime, nowTime int64) *FollowerInfo {
	tripTime := communicationTime / 2
	diffTime := (localTime + tripTime) - nowTime

	return &FollowerInfo{
		Name:              name,
		Address:           address,
		LocalTime:         localTime,
<<<<<<< HEAD
=======
		DateFollower:      time.UnixMilli(localTime),
>>>>>>> 4df5b3de8e63625e208a4e5d62139f2f2b1cb612
		CommunicationTime: communicationTime,
		TripTime:          tripTime,
		DiffTime:          diffTime,
		State:             RequestNotSent, // Estado inicial por defecto
	}
}

// GetAddress devuelve la dirección del seguidor.
func (f *FollowerInfo) GetAddress() string {
	return f.Address
}

// GetName devuelve el nombre del seguidor.
func (f *FollowerInfo) GetName() string {
	return f.Name
}

// GetDateFollower devuelve la fecha y hora local del seguidor.
func (f *FollowerInfo) GetDateFollower() time.Time {
	return f.DateFollower
}

// GetDiffTime devuelve la diferencia de tiempo entre el líder y el seguidor.
func (f *FollowerInfo) GetDiffTime() int64 {
	return f.DiffTime
}

// GetLocalTime devuelve la hora local del seguidor.
func (f *FollowerInfo) GetLocalTime() int64 {
	return f.LocalTime
}

// GetCommunicationTime devuelve el tiempo de comunicación entre líder y seguidor.
func (f *FollowerInfo) GetCommunicationTime() int64 {
	return f.CommunicationTime
}

// GetDelta devuelve la diferencia global de tiempo (delta).
func (f *FollowerInfo) GetDelta() int64 {
	return f.Delta
}

// SetDelta establece la diferencia global de tiempo (delta) para el seguidor.
func (f *FollowerInfo) SetDelta(delta int64) {
	f.Delta = delta
}

// GetState devuelve el estado actual del seguidor.
func (f *FollowerInfo) GetState() FollowerState {
	return f.State
}

// SetState establece el estado del seguidor.
func (f *FollowerInfo) SetState(state FollowerState) {
	f.State = state
}

// String genera una representación en cadena del objeto FollowerInfo.
func (f *FollowerInfo) String() string {
	return fmt.Sprintf("Nombre: %s, Estado: %s, Hora local: %d, Fecha: %s, "+
		"Tiempo de comunicación: %d ms, TripTime: %d ms, Diferencia de tiempo: %d ms, "+
		"Diferencia global (delta): %d ms, Dirección: %s",
		f.Name, f.State, f.LocalTime, f.DateFollower, f.CommunicationTime,
		f.TripTime, f.DiffTime, f.Delta, f.Address)
}
