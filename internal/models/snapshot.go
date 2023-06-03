package models

import (
	"fmt"
	"time"

	snapshot_v1 "github.com/f0rmul/sensor-control/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type State int64

const (
	NormalState State = iota
	WarningState
)

type Snapshot struct {
	ID          string    `bson:"snapshot_id,omitempty"`
	SensorID    string    `bson:"sensor_id,omitempty"`
	SensorName  string    `bson:"sensor_name,omitempty"`
	SensorState State     `bson:"state,omitempty"`
	Timestamp   time.Time `bson:"timestamp,omitempty"`
}

func NewFromProto(proto *snapshot_v1.Snapshot) *Snapshot {

	return &Snapshot{
		ID:          proto.SnapshotId,
		SensorID:    proto.SensorId,
		SensorName:  proto.SensorName,
		SensorState: State(proto.State),
		Timestamp:   proto.Timestamp.AsTime(),
	}
}

func (s *Snapshot) Stringify() string {
	return fmt.Sprintf("id: %s,sensor-id:%s,sensor-name:%s,state:%d,timestamp:%s",
		s.ID,
		s.SensorID,
		s.SensorName,
		s.SensorState,
		s.Timestamp.Format(time.RFC1123))
}

func (s *Snapshot) ToPoto() *snapshot_v1.Snapshot {
	return &snapshot_v1.Snapshot{
		SnapshotId: s.ID,
		SensorId:   s.SensorID,
		SensorName: s.SensorName,
		State:      snapshot_v1.Snapshot_SensorState(s.SensorState),
		Timestamp:  timestamppb.New(s.Timestamp),
	}
}
