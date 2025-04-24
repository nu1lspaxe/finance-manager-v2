package pkg

import (
	"context"
	"records/proto"
	"records/utils"

	"github.com/bufbuild/protovalidate-go"
)

type RecordController struct {
	service *RecordService
	proto.UnimplementedRecordServiceServer
}

func NewRecordController(service *RecordService) *RecordController {
	return &RecordController{service: service}
}

func (r *RecordController) CreateRecord(ctx context.Context, req *proto.CreateRecordRequest) (*proto.RecordResponse, error) {
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	record, err := r.service.CreateRecord(ctx, req.UserId, req.Amount, req.TransactionDate, req.RecordType.String(), req.Detail)
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
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	record, err := r.service.GetRecord(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	protoRecord, err := utils.SqlcToProto(record)
	if err != nil {
		return nil, err
	}
	return &proto.RecordResponse{Record: protoRecord}, nil
}

func (r *RecordController) GetUserRecordsWithFilters(
	req *proto.GetUserRecordsWithFiltersRequest,
	stream proto.RecordService_GetUserRecordsWithFiltersServer,
) error {
	err := protovalidate.Validate(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(stream.Context(), utils.TIMEOUT_STREAM)
	defer cancel()

	records, err := r.service.GetUserRecordsWithFilters(ctx, req.UserId, req.RecordType.String(), req.StartTime, req.EndTime)
	if err != nil {
		return err
	}

	for _, record := range records {
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
	err := protovalidate.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err = r.service.UpdateRecord(ctx, req.Id, req.Amount, req.TransactionDate, req.RecordType.String(), req.Detail)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateRecordResponse{Success: true}, nil
}

func (r *RecordController) DeleteRecord(ctx context.Context, req *proto.DeleteRecordRequest) (*proto.DeleteRecordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.TIMEOUT)
	defer cancel()

	err := r.service.DeleteRecord(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteRecordResponse{Success: true}, nil
}
