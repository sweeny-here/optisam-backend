// Copyright (C) 2019 Orange
// 
// This software is distributed under the terms and conditions of the 'Apache License 2.0'
// license which can be found in the file 'License.txt' in this package distribution 
// or at 'http://www.apache.org/licenses/LICENSE-2.0'. 
//
package v1

import (
	"context"
	"errors"
	"optisam-backend/common/optisam/ctxmanage"
	"optisam-backend/common/optisam/logger"
	"optisam-backend/common/optisam/strcomp"
	v1 "optisam-backend/license-service/pkg/api/v1"
	repo "optisam-backend/license-service/pkg/repository/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *licenseServiceServer) CreateMetricOracleProcessorStandard(ctx context.Context, req *v1.CreateMetricOPS) (*v1.CreateMetricOPS, error) {
	userClaims, ok := ctxmanage.RetrieveClaims(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "cannot find claims in context")
	}
	if req.StartEqTypeId == "" {
		return nil, status.Error(codes.InvalidArgument, "start level is empty")
	}

	metrics, err := s.licenseRepo.ListMetrices(ctx, userClaims.Socpes)
	if err != nil && err != repo.ErrNoData {
		logger.Log.Error("service/v1 -CreateMetricSAGProcessorStandard - fetching metrics", zap.String("reason", err.Error()))
		return nil, status.Error(codes.Internal, "cannot fetch metrics")

	}

	if metricNameExistsAll(metrics, req.Name) != -1 {
		return nil, status.Error(codes.InvalidArgument, "metric name already exists")
	}

	eqTypes, err := s.licenseRepo.EquipmentTypes(ctx, userClaims.Socpes)
	if err != nil {
		logger.Log.Error("service/v1 -CreateMetricOracleProcessorStandard - fetching equipments", zap.String("reason", err.Error()))
		return nil, status.Error(codes.Internal, "cannot fetch equipment types")
	}

	parAncestors, err := parentHierarchy(eqTypes, req.StartEqTypeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "parent hierarchy doesnt exists")
	}

	baseLevelIdx, err := validateLevelsNew(parAncestors, 0, req.BaseEqTypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot find base level equipment type in parent hierarchy")
	}
	aggLevelIdx, err := validateLevelsNew(parAncestors, baseLevelIdx, req.AggerateLevelEqTypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot find aggregate level equipment type in parent hierarchy")
	}
	_, err = validateLevelsNew(parAncestors, aggLevelIdx, req.EndEqTypeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot find end level equipment type in parent hierarchy")
	}

	if err := validateAttributesOPS(parAncestors[baseLevelIdx].Attributes, req.NumCoreAttrId, req.NumCPUAttrId, req.CoreFactorAttrId); err != nil {
		return nil, err
	}

	met, err := s.licenseRepo.CreateMetricOPS(ctx, serverToRepoMetricOPS(req), userClaims.Socpes)
	if err != nil {
		logger.Log.Error("service/v1 - CreateMetricOracleProcessorStandard - fetching equipment", zap.String("reason", err.Error()))
		return nil, status.Error(codes.NotFound, "cannot create metric")
	}

	return repoToServerMetricOPS(met), nil

}

func parentHierarchy(eqTypes []*repo.EquipmentType, startID string) ([]*repo.EquipmentType, error) {
	equip, err := equipmentTypeExistsByID(startID, eqTypes)
	if err != nil {
		logger.Log.Error("service/v1 - parentHierarchy - fetching equipment type", zap.String("reason", err.Error()))
		return nil, status.Error(codes.NotFound, "cannot fetch equipment type with given Id")
	}
	ancestors := []*repo.EquipmentType{}
	ancestors = append(ancestors, equip)
	parID := equip.ParentID
	for parID != "" {
		equipAnc, err := equipmentTypeExistsByID(parID, eqTypes)
		if err != nil {
			logger.Log.Error("service/v1 - parentHierarchy - fetching equipment type", zap.String("reason", err.Error()))
			return nil, status.Error(codes.NotFound, "parent hierachy not found")
		}
		ancestors = append(ancestors, equipAnc)
		parID = equipAnc.ParentID
	}
	return ancestors, nil
}

func attributeExists(attributes []*repo.Attribute, attrID string) (*repo.Attribute, error) {
	for _, attr := range attributes {
		if attr.ID == attrID {
			return attr, nil
		}
	}
	return nil, status.Errorf(codes.Unknown, "attribute not exists")
}

func validateAttributesOPS(attr []*repo.Attribute, numCoreAttr string, numCPUAttr string, coreFactorAttr string) error {

	if numCoreAttr == "" {
		return status.Error(codes.InvalidArgument, "num of cores attribute is empty")
	}
	if numCPUAttr == "" {
		return status.Error(codes.InvalidArgument, "num of cpu attribute is empty")
	}
	if coreFactorAttr == "" {
		return status.Error(codes.InvalidArgument, "core factor attribute is empty")
	}

	numOfCores, err := attributeExists(attr, numCoreAttr)
	if err != nil {

		return status.Error(codes.InvalidArgument, "numofcores attribute doesnt exists")
	}
	if numOfCores.Type != repo.DataTypeInt && numOfCores.Type != repo.DataTypeFloat {
		return status.Error(codes.InvalidArgument, "numofcores attribute doesnt have valid data type")
	}

	numOfCPU, err := attributeExists(attr, numCPUAttr)
	if err != nil {

		return status.Error(codes.InvalidArgument, "numofcpu attribute doesnt exists")
	}
	if numOfCPU.Type != repo.DataTypeInt && numOfCPU.Type != repo.DataTypeFloat {
		return status.Error(codes.InvalidArgument, "numofcpu attribute doesnt have valid data type")
	}

	coreFactor, err := attributeExists(attr, coreFactorAttr)
	if err != nil {

		return status.Error(codes.InvalidArgument, "corefactor attribute doesnt exists")
	}

	if coreFactor.Type != repo.DataTypeInt && coreFactor.Type != repo.DataTypeFloat {
		return status.Error(codes.InvalidArgument, "corefactor attribute doesnt have valid data type")
	}
	return nil
}

func metricNameExistsOPS(metrics []*repo.MetricOPS, name string) int {
	for i, met := range metrics {
		if strcomp.CompareStrings(met.Name, name) {
			return i
		}
	}
	return -1
}

func serverToRepoMetricOPS(met *v1.CreateMetricOPS) *repo.MetricOPS {
	return &repo.MetricOPS{
		ID:                    met.ID,
		Name:                  met.Name,
		NumCoreAttrID:         met.NumCoreAttrId,
		NumCPUAttrID:          met.NumCPUAttrId,
		CoreFactorAttrID:      met.CoreFactorAttrId,
		StartEqTypeID:         met.StartEqTypeId,
		BaseEqTypeID:          met.BaseEqTypeId,
		AggerateLevelEqTypeID: met.AggerateLevelEqTypeId,
		EndEqTypeID:           met.EndEqTypeId,
	}
}

func repoToServerMetricOPS(met *repo.MetricOPS) *v1.CreateMetricOPS {
	return &v1.CreateMetricOPS{

		ID:                    met.ID,
		Name:                  met.Name,
		NumCoreAttrId:         met.NumCoreAttrID,
		NumCPUAttrId:          met.NumCPUAttrID,
		CoreFactorAttrId:      met.CoreFactorAttrID,
		StartEqTypeId:         met.StartEqTypeID,
		BaseEqTypeId:          met.BaseEqTypeID,
		AggerateLevelEqTypeId: met.AggerateLevelEqTypeID,
		EndEqTypeId:           met.EndEqTypeID,
	}
}

func validateLevelsNew(levels []*repo.EquipmentType, startIdx int, base string) (int, error) {
	for i := startIdx; i < len(levels); i++ {
		if levels[i].ID == base {
			return i, nil
		}
	}
	return -1, errors.New("Not found")
}
func (s *licenseServiceServer) computedLicensesOPS(ctx context.Context, eqTypes []*repo.EquipmentType, ID string, met string) (uint64, error) {
	userClaims, _ := ctxmanage.RetrieveClaims(ctx)
	metrics, err := s.licenseRepo.ListMetricOPS(ctx, userClaims.Socpes)
	if err != nil && err != repo.ErrNoData {
		return 0, status.Error(codes.Internal, "cannot fetch metric OPS")

	}
	ind := 0
	if ind = metricNameExistsOPS(metrics, met); ind == -1 {
		return 0, status.Error(codes.Internal, "metric name doesnot exists")
	}
	parTree, err := parentHierarchy(eqTypes, metrics[ind].StartEqTypeID)
	if err != nil {
		return 0, status.Error(codes.Internal, "cannot fetch equipment types")

	}

	baseLevelIdx, err := validateLevelsNew(parTree, 0, metrics[ind].BaseEqTypeID)
	if err != nil {
		return 0, status.Error(codes.Internal, "cannot find base level equipment type in parent hierarchy")
	}

	aggLevelIdx, err := validateLevelsNew(parTree, baseLevelIdx, metrics[ind].AggerateLevelEqTypeID)
	if err != nil {
		return 0, status.Error(codes.Internal, "cannot find aggregate level equipment type in parent hierarchy")
	}

	endLevelIdx, err := validateLevelsNew(parTree, aggLevelIdx, metrics[ind].EndEqTypeID)
	if err != nil {
		return 0, status.Error(codes.Internal, "cannot find end level equipment type in parent hierarchy")
	}

	numOfCores, err := attributeExists(parTree[baseLevelIdx].Attributes, metrics[ind].NumCoreAttrID)
	if err != nil {
		return 0, status.Error(codes.Internal, "numofcores attribute doesnt exits")

	}
	numOfCPU, err := attributeExists(parTree[baseLevelIdx].Attributes, metrics[ind].NumCPUAttrID)
	if err != nil {
		return 0, status.Error(codes.Internal, "numofcpu attribute doesnt exits")

	}
	coreFactor, err := attributeExists(parTree[baseLevelIdx].Attributes, metrics[ind].CoreFactorAttrID)
	if err != nil {
		return 0, status.Error(codes.Internal, "coreFactor attribute doesnt exits")

	}
	mat := &repo.MetricOPSComputed{
		EqTypeTree:     parTree[:endLevelIdx+1],
		BaseType:       parTree[baseLevelIdx],
		AggregateLevel: parTree[aggLevelIdx],
		NumCoresAttr:   numOfCores,
		NumCPUAttr:     numOfCPU,
		CoreFactorAttr: coreFactor,
	}
	computedLicenses, err := s.licenseRepo.MetricOPSComputedLicenses(ctx, ID, mat, userClaims.Socpes)
	if err != nil {
		return 0, status.Error(codes.Internal, "cannot compute licenses for metric OPS")

	}

	return computedLicenses, nil
}
