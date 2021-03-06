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
	"optisam-backend/common/optisam/token/claims"
	v1 "optisam-backend/license-service/pkg/api/v1"
	repo "optisam-backend/license-service/pkg/repository/v1"
	"optisam-backend/license-service/pkg/repository/v1/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_licenseServiceServer_CreateMetricIBMPvuStandard(t *testing.T) {
	var mockCtrl *gomock.Controller
	var rep repo.License
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})

	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			ID:       "e2",
			ParentID: "e3",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:   "a1",
					Type: repo.DataTypeInt,
				},
				&repo.Attribute{
					ID:   "a2",
					Type: repo.DataTypeInt,
				},
				&repo.Attribute{
					ID:   "a3",
					Type: repo.DataTypeInt,
				},
			},
		},
	}

	type args struct {
		ctx context.Context
		req *v1.CreateMetricIPS
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.CreateMetricIPS
		setup   func()
		wantErr bool
	}{
		{name: "SUCCESS",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				mockRepo.EXPECT().CreateMetricIPS(ctx, &repo.MetricIPS{
					Name:             "IPS",
					NumCoreAttrID:    "a1",
					CoreFactorAttrID: "a3",
					BaseEqTypeID:     "e2",
				}, []string{"A", "B"}).Times(1).Return(&repo.MetricIPS{
					ID:               "IPS",
					Name:             "IPS",
					NumCoreAttrID:    "a1",
					CoreFactorAttrID: "a3",
					BaseEqTypeID:     "e2",
				}, nil)
			},
			want: &v1.CreateMetricIPS{
				ID:               "IPS",
				Name:             "IPS",
				NumCoreAttrId:    "a1",
				CoreFactorAttrId: "a3",
				BaseEqTypeId:     "e2",
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch metric types",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment types",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot create metric",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				mockRepo.EXPECT().CreateMetricIPS(ctx, &repo.MetricIPS{
					Name:             "IPS",
					NumCoreAttrID:    "a1",
					CoreFactorAttrID: "a3",
					BaseEqTypeID:     "e2",
				}, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - metric name already exists",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "IPS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - metric name already exists - case insensitive",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "ips",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "IPS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot find base level equipment type",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return([]*repo.EquipmentType{
					&repo.EquipmentType{
						ID:       "e1",
						ParentID: "e2",
					},
					&repo.EquipmentType{
						ID:       "e3",
						ParentID: "e4",
					},
					&repo.EquipmentType{
						ID: "e4",
					},
				}, nil)

			},
			wantErr: true,
		},
		{name: "FAILURE - num of cores attribute is empty",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - core factor attribute is empty",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - numofcores attribute doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a4",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - numofcores attribute doesnt have valid data type",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return([]*repo.EquipmentType{
					&repo.EquipmentType{
						ID:       "e2",
						ParentID: "e3",
						Attributes: []*repo.Attribute{
							&repo.Attribute{
								ID:   "a1",
								Type: repo.DataTypeString,
							},
							&repo.Attribute{
								ID:   "a2",
								Type: repo.DataTypeInt,
							},
							&repo.Attribute{
								ID:   "a3",
								Type: repo.DataTypeInt,
							},
						},
					},
				}, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - core factor attribute doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a4",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - core factor attribute doesnt have valid data type",
			args: args{
				ctx: ctx,
				req: &v1.CreateMetricIPS{
					Name:             "IPS",
					NumCoreAttrId:    "a1",
					CoreFactorAttrId: "a3",
					BaseEqTypeId:     "e2",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().ListMetrices(ctx, []string{"A", "B"}).Times(1).Return([]*repo.Metric{
					&repo.Metric{
						Name: "ONS",
					},
					&repo.Metric{
						Name: "WS",
					},
				}, nil)
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return([]*repo.EquipmentType{
					&repo.EquipmentType{
						ID:       "e2",
						ParentID: "e3",
						Attributes: []*repo.Attribute{
							&repo.Attribute{
								ID:   "a1",
								Type: repo.DataTypeInt,
							},
							&repo.Attribute{
								ID:   "a2",
								Type: repo.DataTypeInt,
							},
							&repo.Attribute{
								ID:   "a3",
								Type: repo.DataTypeString,
							},
						},
					},
				}, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			if tt.setup == nil {
				defer mockCtrl.Finish()
			}
			got, err := s.CreateMetricIBMPvuStandard(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.CreateMetricIBMPvuStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.CreateMetricIBMPvuStandard() = %v, want %v", got, tt.want)
			}
		})
	}
}
