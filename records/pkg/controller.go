package pkg

import (
	"context"
	"records/proto"
	"records/utils"
)

type RecordController struct {
	service *RecordService
	proto.UnimplementedRecordServiceServer
}

func NewRecordController(service *RecordService) *RecordController {
	return &RecordController{service: service}
}

func (r *RecordController) CreateRecord(ctx context.Context, req *proto.CreateRecordRequest) (*proto.RecordResponse, error) {
	record, err := r.service.CreateRecord(req.UserId, req.Amount, req.TransactionDate, req.RecordType, req.Detail)
	if err != nil {
		return nil, err
	}
	protoRecord, err := utils.SqlcToProto(record)
	if err != nil {
		return nil, err
	}
	return &proto.RecordResponse{Record: protoRecord}, nil
}

func (r *RecordController) GetRecord(ctx context.Context, req *proto.GetRecordRequest) (*proto.RecordResponse, error) {
	record, err := r.service.GetRecord(req.Id)
	if err != nil {
		return nil, err
	}
	protoRecord, err := utils.SqlcToProto(record)
	if err != nil {
		return nil, err
	}
	return &proto.RecordResponse{Record: protoRecord}, nil
}

func (r *RecordController) GetUserRecordsWithFilters(req *proto.GetUserRecordsWithFiltersRequest, stream proto.RecordService_GetUserRecordsWithFiltersServer) error {
	records, err := r.service.GetUserRecordsWithFilters(req.UserId, req.RecordType, req.StartTime, req.EndTime)
	if err != nil {
		return err
	}
	for _, record := range *records {
		protoRecord, err := utils.SqlcToProto(&record)
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.RecordResponse{Record: protoRecord}); err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordController) UpdateRecord(ctx context.Context, req *proto.UpdateRecordRequest) (*proto.UpdateRecordResponse, error) {
	err := r.service.UpdateRecord(req.Id, req.Amount, req.TransactionDate, req.RecordType, req.Detail)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateRecordResponse{Success: true}, nil
}

func (r *RecordController) DeleteRecord(ctx context.Context, req *proto.DeleteRecordRequest) (*proto.DeleteRecordResponse, error) {
	err := r.service.DeleteRecord(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteRecordResponse{Success: true}, nil
}
