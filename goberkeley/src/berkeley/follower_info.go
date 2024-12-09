package berkeley

import (
	"fmt"
	"time"
)

// FollowerInfo encapsula información sobre un nodo seguidor en el sistema.
type FollowerInfo struct {
	Name              string        // Nombre del seguidor
	Address           string        // Dirección del seguidor en formato "host:puerto"
	State             FollowerState // Estado actual del seguidor
	CurrentTime       int64         // T0: Hora inicial del líder
	FollowerTime      int64         // Hora local del seguidor en milisegundos desde la época UNIX
	CommunicationTime int64         // Tiempo de comunicación entre líder y seguidor en milisegundos
	TripTime          int64         // Tiempo de ida y vuelta (triptime) estimado
	DiffTime          int64         // Diferencia de tiempo calculada entre líder y seguidor
	Delta             int64         // Diferencia global de tiempo (delta) aplicada al seguidor
}

// NewFollowerInfo crea un nuevo objeto FollowerInfo con los valores proporcionados.
func NewFollowerInfo(address, name string, followerTime, currentTimer, communicationTime, nowTime int64) *FollowerInfo {
	tripTime := (nowTime - currentTimer) / 2
	diffTime := (followerTime + tripTime) - nowTime

	return &FollowerInfo{
		Name:              name,
		Address:           address,
		FollowerTime:      followerTime,
		CurrentTime:       currentTimer,
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

// GetDiffTime devuelve la diferencia de tiempo entre el líder y el seguidor.
func (f *FollowerInfo) GetDiffTime() int64 {
	return f.DiffTime
}

// GetLocalTime devuelve la hora local del seguidor.
func (f *FollowerInfo) GetFollowerTime() int64 {
	return f.FollowerTime
}

// CurrentTime devuelve la hora T0 del líder.
func (f *FollowerInfo) GetCurrentTime() int64 {
	return f.CurrentTime
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
	// Convertir DateFollower de int64 (timestamp en milisegundos) a time.Time
	timestampFollower := time.Unix(0, f.GetFollowerTime()*int64(time.Millisecond))
	timestampLeader := time.Unix(0, f.GetCurrentTime()*int64(time.Millisecond))
	// Usar el formato adecuado para la fecha
	return fmt.Sprintf("Nombre: %s, Estado: %s, Hora local del Seguidor: %d, Fecha: %s, "+
		"Hora inicial T0 del líder: %d, Fecha: %s,  Tiempo de comunicación: %d ms, TripTime: %d ms, "+
		"Diferencia de tiempo: %d ms, Diferencia global (delta): %d ms, Dirección: %s",
		f.Name, f.State, f.FollowerTime, timestampFollower.Format("2006-01-02 15:04:05"), // Formato para la fecha
		f.CurrentTime, timestampLeader, f.CommunicationTime, f.TripTime, f.DiffTime, f.Delta, f.Address)
}
