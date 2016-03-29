package database

type EventLog struct {}

func (n *EventLog) Event(eventName string) {
  //grpclog.Println(eventName)
}

// EventKv receives a notification when various events occur along with
// optional key/value data
func (n *EventLog) EventKv(eventName string, kvs map[string]string) {
  //grpclog.Println("%s ",eventName,kvs)
}

// EventErr receives a notification of an error if one occurs
func (n *EventLog) EventErr(eventName string, err error) error {
  //grpclog.Println("%s ",eventName,err)
  return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (n *EventLog) EventErrKv(eventName string, err error, kvs map[string]string) error {
  //grpclog.Println("%s ",eventName,err,kvs)
	return err
}

// Timing receives the time an event took to happen
func (n *EventLog) Timing(eventName string, nanoseconds int64) {
  //grpclog.Println("YES5")
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (n *EventLog) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
  //grpclog.Println("",eventName,nanoseconds,kvs)
}
