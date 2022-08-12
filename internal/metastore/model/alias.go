package model

import pb "github.com/milvus-io/milvus/internal/proto/etcdpb"

type Alias struct {
	Name         string
	CollectionID int64
	CreatedTime  uint64
	State        pb.AliasState
}

func (a Alias) Available() bool {
	return a.State == pb.AliasState_AliasCreated
}

func (a Alias) Clone() *Alias {
	return &Alias{
		Name:         a.Name,
		CollectionID: a.CollectionID,
		CreatedTime:  a.CreatedTime,
		State:        a.State,
	}
}

func MarshalAliasModel(alias *Alias) *pb.AliasInfo {
	return &pb.AliasInfo{
		AliasName:    alias.Name,
		CollectionId: alias.CollectionID,
		CreatedTime:  alias.CreatedTime,
		State:        alias.State,
	}
}

func UnmarshalAliasModel(info *pb.AliasInfo) *Alias {
	return &Alias{
		Name:         info.GetAliasName(),
		CollectionID: info.GetCollectionId(),
		CreatedTime:  info.GetCreatedTime(),
		State:        info.GetState(),
	}
}
