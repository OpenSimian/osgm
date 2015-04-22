package core

import (
  "encoding/json"
  "github.com/satori/go.uuid"
)


type userRequest struct {
  MessageID int
  MessageType string
  Message json.RawMessage
}

func (ur *userRequest)load(msg []byte) {
  err := json.Unmarshal(msg, ur)
  if err != nil {
    ur.MessageType = err.Error()
  }
}

func(ur userRequest) readRegionID() (uuid.UUID, error) {
  type rid struct {
    RegionUUID uuid.UUID
  }
  r := rid{}
  err := json.Unmarshal(ur.Message, &r)
  if err != nil {
    return uuid.UUID{}, err
  }
  return r.RegionUUID, nil
}
