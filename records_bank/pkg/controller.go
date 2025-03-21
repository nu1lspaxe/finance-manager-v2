package pkg

import (
	"context"
	"records_bank/proto"
	"records_bank/utils"

	"github.com/bufbuild/protovalidate-go"
)

type BankRecordController struct {
	service *BankRecordService
	proto.UnimplementedBankRecordServiceServer
}

func NewBankRecordController(service *BankRecordService) *BankRecordController {
	return &BankRecordController{service: service}
}

func (c *BankRecordController) GetBankRecord(ctx context.Context, req *proto.GetBankRecordRequest) (*proto.BankRecordResponse, error) {
	bankRecord, err := c.service.GetBankRecord(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	protoBankRecord, err := utils.SqlcToProto(bankRecord)
	if err != nil {
		return nil, err
	}

	return &proto.BankRecordResponse{Record: protoBankRecord}, nil
}

func (c *BankRecordController) GetUserBankRecordsWithFilters(
	req *proto.GetUserBankRecordsWithFiltersRequest,
	stream proto.BankRecordService_GetUserBankRecordsWithFiltersServer,
) error {
	err := protovalidate.Validate(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(stream.Context(), utils.TIMEOUT_STREAM)
	defer cancel()

	bankRecords, err := c.service.GetBankRecordsWithFilters(
		ctx,
		req.UserId,
		req.AccountNumber,
		req.RecordType.String(),
		req.StartTime,
		req.EndTime,
	)

	if err != nil {
		return err
	}

	for _, bankRecord := range bankRecords {
		protoBankRecord, err := utils.SqlcToProto(&bankRecord)
		if err != nil {
			return err
		}
		if err := stream.Send(&proto.BankRecordResponse{Record: protoBankRecord}); err != nil {
			return err
		}
	}

	return nil
}
