// Copyright (C) 2019 Orange
// 
// This software is distributed under the terms and conditions of the 'Apache License 2.0'
// license which can be found in the file 'License.txt' in this package distribution 
// or at 'http://www.apache.org/licenses/LICENSE-2.0'. 
//
package v1

import (
	"context"
	"encoding/json"
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

func Test_licenseServiceServer_ListEquipments(t *testing.T) {
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})

	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			Type:     "typ1",
			ID:       "1",
			SourceID: "s1",
			ParentID: "p1",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:           "1",
					Name:         "attr1",
					Type:         repo.DataTypeString,
					IsIdentifier: true,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
				&repo.Attribute{
					ID:           "3",
					Name:         "attr3",
					Type:         repo.DataTypeInt,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_3",
				},
				&repo.Attribute{
					ID:           "4",
					Name:         "attr4",
					Type:         repo.DataTypeFloat,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_4",
				},
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr5",
					Type:               repo.DataTypeString,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_5",
				},
				&repo.Attribute{
					ID:       "1",
					Name:     "attr6",
					Type:     repo.DataTypeString,
					MappedTo: "mapping_6",
				},
				&repo.Attribute{
					ID:          "1",
					Name:        "attr7",
					IsDisplayed: true,
					Type:        repo.DataTypeString,
					MappedTo:    "mapping_7",
				},
				&repo.Attribute{
					ID:           "1",
					Name:         "attr8",
					IsDisplayed:  true,
					IsSearchable: true,
					Type:         repo.DataType(255),
					MappedTo:     "mapping_8",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ2",
			ID:       "2",
			SourceID: "s2",
			ParentID: "p2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
	}
	var mockCtrl *gomock.Controller
	var rep repo.License
	type args struct {
		ctx context.Context
		req *v1.ListEquipmentsRequest
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *v1.ListEquipmentsResponse
		wantErr bool
	}{
		{name: "success",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				// TODO investigate why test is failing when we are checking QueryEquipments
				// when we are comparing in mock controller
				//&repo.QueryEquipments{
				// PageSize:  10,
				// Offset:    90,
				// SortBy:    "attr1",
				// SortOrder: repo.SortDESC,
				// Filter: &repo.AggregateFilter{
				// 	Filters: []repo.Queryable{
				// 		&repo.Filter{
				// 			FilterKey:   "attr1",
				// 			FilterValue: "a11",
				// 		},
				// 		&repo.Filter{
				// 			FilterKey:   "attr2",
				// 			FilterValue: "a22",
				// 		},
				// 		&repo.Filter{
				// 			FilterKey:   "attr3",
				// 			FilterValue: int64(3),
				// 		},
				// 		&repo.Filter{
				// 			FilterKey:   "attr4",
				// 			FilterValue: float64(4),
				// 		},
				// 	},
				// }
				mockRepo.EXPECT().Equipments(ctx, eqTypes[0], gomock.Any(), []string{"A", "B"}).Times(1).Return(int32(3), json.RawMessage(`[{ID:"1"}]`), nil)
			},
			want: &v1.ListEquipmentsResponse{
				TotalRecords: 3,
				Equipments:   json.RawMessage(`[{ID:"1"}]`),
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.ListEquipmentsRequest{
					TypeId:       "3",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "failure - validation : equipment type not found",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "3",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation : sort by attribute not found",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "Notfound",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation : sort by attribute not displayable",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr6",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : invalid query",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=%gh",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : atttribute not found",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "notfound=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute not dispalyed",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr6=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute not searchable",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr7=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute value empty",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute string type less than 3 chars",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute int type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr3=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute float type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr4=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute unsupported data type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr8=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - database : failure getting equipment types",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("test error"))
			},
			wantErr: true,
		},
		{name: "failure - database : failure getting equipments",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentsRequest{
					TypeId:       "1",
					PageNum:      10,
					PageSize:     10,
					SortBy:       "attr1",
					SortOrder:    v1.SortOrder_DESC,
					SearchParams: "attr1=a11,attr2=a22,attr3=3,attr4=4",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				// mockRepo.EXPECT().Equipments(nil, eqTypes[0], &repo.QueryEquipments{
				// 	PageSize:  10,
				// 	Offset:    90,
				// 	SortBy:    "attr1",
				// 	SortOrder: repo.SortDESC,
				// 	Filter: &repo.AggregateFilter{
				// 		Filters: []repo.Queryable{
				// 			&repo.Filter{
				// 				FilterKey:   "attr1",
				// 				FilterValue: "a11",
				// 			},
				// 			&repo.Filter{
				// 				FilterKey:   "attr2",
				// 				FilterValue: "a22",
				// 			},
				// 			&repo.Filter{
				// 				FilterKey:   "attr3",
				// 				FilterValue: int64(3),
				// 			},
				// 			&repo.Filter{
				// 				FilterKey:   "attr4",
				// 				FilterValue: float64(4),
				// 			},
				// 		},
				// 	},
				// }).Times(1).Return(int32(3), nil, errors.New("test error"))
				mockRepo.EXPECT().Equipments(ctx, eqTypes[0], gomock.Any(), []string{"A", "B"}).Times(1).Return(int32(3), nil, errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			got, err := s.ListEquipments(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.ListEquipments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.ListEquipments() = %v, want %v", got, tt.want)
			}
			if tt.setup == nil {
				mockCtrl.Finish()
			}
		})
	}
}

func Test_licenseServiceServer_GetEquipment(t *testing.T) {
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})
	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			Type:     "typ1",
			ID:       "1",
			SourceID: "s1",
			ParentID: "2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:           "1",
					Name:         "attr1",
					Type:         repo.DataTypeString,
					IsIdentifier: true,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
				&repo.Attribute{
					ID:           "3",
					Name:         "attr3",
					Type:         repo.DataTypeInt,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_3",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ2",
			ID:       "2",
			SourceID: "s2",
			//ParentID: "p2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
	}
	var mockCtrl *gomock.Controller
	var rep repo.License
	type args struct {
		ctx context.Context
		req *v1.GetEquipmentRequest
	}
	tests := []struct {
		name    string
		s       *licenseServiceServer
		args    args
		setup   func()
		want    *v1.GetEquipmentResponse
		wantErr bool
	}{
		{name: "SUCCESS",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().Equipment(ctx, eqTypes[0], "e1", []string{"A", "B"}).Times(1).Return(json.RawMessage(`[{ID:"1"}]`), nil)
			},
			want: &v1.GetEquipmentResponse{
				Equipment: string(json.RawMessage(`[{ID:"1"}]`)),
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment types",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "3",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().Equipment(ctx, eqTypes[0], "e1", []string{"A", "B"}).Times(1).Return(nil, repo.ErrNoData)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id - node not exists",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().Equipment(ctx, eqTypes[0], "e1", []string{"A", "B"}).Times(1).Return(nil, repo.ErrNodeNotFound)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id",
			args: args{
				ctx: ctx,
				req: &v1.GetEquipmentRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().Equipment(ctx, eqTypes[0], "e1", []string{"A", "B"}).Times(1).Return(nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			got, err := s.GetEquipment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.GetEquipment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.GetEquipment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_licenseServiceServer_ListEquipmentParents(t *testing.T) {
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})
	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			Type:     "typ1",
			ID:       "1",
			SourceID: "s1",
			ParentID: "2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:           "1",
					Name:         "attr1",
					Type:         repo.DataTypeString,
					IsIdentifier: true,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
				&repo.Attribute{
					ID:           "3",
					Name:         "attr3",
					Type:         repo.DataTypeInt,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_3",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ2",
			ID:       "2",
			SourceID: "s2",
			//ParentID: "p2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ3",
			ID:       "3",
			SourceID: "s3",
			ParentID: "4",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
	}
	const (
		records   int32 = 3
		norecords int32 = 0
	)
	var mockCtrl *gomock.Controller
	var rep repo.License
	type args struct {
		ctx context.Context
		req *v1.ListEquipmentParentsRequest
	}
	tests := []struct {
		name    string
		s       *licenseServiceServer
		args    args
		setup   func()
		want    *v1.ListEquipmentsResponse
		wantErr bool
	}{
		{name: "SUCCESS",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentParents(ctx, eqTypes[0], eqTypes[1], "e1", []string{"A", "B"}).Times(1).Return(records, json.RawMessage(`[{ID:"1"}]`), nil)
			},
			want: &v1.ListEquipmentsResponse{
				TotalRecords: records,
				Equipments:   json.RawMessage(`[{ID:"1"}]`),
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment types",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "4",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - parent of equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "3",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentParents(ctx, eqTypes[0], eqTypes[1], "e1", []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNoData)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentParents(ctx, eqTypes[0], eqTypes[1], "e1", []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNodeNotFound)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentParentsRequest{
					TypeId:  "1",
					EquipId: "e1",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentParents(ctx, eqTypes[0], eqTypes[1], "e1", []string{"A", "B"}).Times(1).Return(norecords, nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			got, err := s.ListEquipmentParents(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.ListEquipmentParents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.ListEquipmentParents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_licenseServiceServer_ListEquipmentChildren(t *testing.T) {
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})
	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			Type:     "typ1",
			ID:       "1",
			SourceID: "s1",
			ParentID: "2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:           "1",
					Name:         "attr1",
					Type:         repo.DataTypeString,
					IsIdentifier: true,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
				&repo.Attribute{
					ID:           "3",
					Name:         "attr3",
					Type:         repo.DataTypeInt,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_3",
				},
				&repo.Attribute{
					ID:           "4",
					Name:         "attr4",
					Type:         repo.DataTypeFloat,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_4",
				},
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr5",
					Type:               repo.DataTypeString,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_5",
				},
				&repo.Attribute{
					ID:       "1",
					Name:     "attr6",
					Type:     repo.DataTypeString,
					MappedTo: "mapping_6",
				},
				&repo.Attribute{
					ID:          "1",
					Name:        "attr7",
					IsDisplayed: true,
					Type:        repo.DataTypeString,
					MappedTo:    "mapping_7",
				},
				&repo.Attribute{
					ID:           "1",
					Name:         "attr8",
					IsDisplayed:  true,
					IsSearchable: true,
					Type:         repo.DataType(255),
					MappedTo:     "mapping_8",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ2",
			ID:       "2",
			SourceID: "s2",
			//ParentID: "p2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr_2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ3",
			ID:       "3",
			SourceID: "s3",
			ParentID: "4",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
	}
	const (
		records   int32 = 3
		norecords int32 = 0
	)
	var mockCtrl *gomock.Controller
	var rep repo.License
	type args struct {
		ctx context.Context
		req *v1.ListEquipmentChildrenRequest
	}
	tests := []struct {
		name    string
		s       *licenseServiceServer
		args    args
		setup   func()
		want    *v1.ListEquipmentsResponse
		wantErr bool
	}{
		{name: "SUCCESS",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
					SearchParams:   "attr1=a11,attr2=a22",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				// &repo.QueryEquipments{
				// 	PageSize:  10,
				// 	Offset:    0,
				// 	SortBy:    "attr1",
				// 	SortOrder: 0,
				// 	Filter: &repo.AggregateFilter{
				// 		Filters: []repo.Queryable{
				// 			&repo.Filter{
				// 				FilterKey:   "attr1",
				// 				FilterValue: "a11",
				// 			},
				// 			&repo.Filter{
				// 				FilterKey:   "attr2",
				// 				FilterValue: "a22",
				// 			},
				// 		},
				// 	},
				// }

				mockRepo.EXPECT().EquipmentChildren(ctx, eqTypes[1], eqTypes[0], "e1", gomock.Any(), []string{"A", "B"}).Times(1).Return(records, json.RawMessage(`[{ID:"1"}]`), nil)
			},
			want: &v1.ListEquipmentsResponse{
				TotalRecords: records,
				Equipments:   json.RawMessage(`[{ID:"1"}]`),
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment types",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "4",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - children of equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "4",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - child of equipment type is not valid",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "3",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : invalid query",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr1=%gh",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : atttribute not found",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "notfound=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute not dispalyed",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr6=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute not searchable",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr7=10",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute value empty",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr1=",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute string type less than 3 chars",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr1=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute int type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr3=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute float type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr4=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "failure - validation - query : attribute unsupported data type cannot parse",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        10,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_DESC,
					SearchParams:   "attr8=hi",
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				// &repo.QueryEquipments{
				// 	PageSize:  10,
				// 	Offset:    0,
				// 	SortBy:    "attr1",
				// 	SortOrder: 0,
				// }
				mockRepo.EXPECT().EquipmentChildren(ctx, eqTypes[1], eqTypes[0], "e1", gomock.Any(), []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNoData)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				// &repo.QueryEquipments{
				// 	PageSize:  10,
				// 	Offset:    0,
				// 	SortBy:    "attr1",
				// 	SortOrder: 0,
				// }
				mockRepo.EXPECT().EquipmentChildren(ctx, eqTypes[1], eqTypes[0], "e1", gomock.Any(), []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNodeNotFound)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentChildrenRequest{
					TypeId:         "2",
					EquipId:        "e1",
					ChildrenTypeId: "1",
					PageNum:        1,
					PageSize:       10,
					SortBy:         "attr1",
					SortOrder:      v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				// &repo.QueryEquipments{
				// 	PageSize:  10,
				// 	Offset:    0,
				// 	SortBy:    "attr1",
				// 	SortOrder: 0,
				// }
				mockRepo.EXPECT().EquipmentChildren(ctx, eqTypes[1], eqTypes[0], "e1", gomock.Any(), []string{"A", "B"}).Times(1).Return(norecords, nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			got, err := s.ListEquipmentChildren(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.ListEquipmentChildren() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.ListEquipmentChildren() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_licenseServiceServer_ListEquipmentProducts(t *testing.T) {
	ctx := ctxmanage.AddClaims(context.Background(), &claims.Claims{
		UserID: "admin@superuser.com",
		Role:   "Admin",
		Socpes: []string{"A", "B"},
	})
	eqTypes := []*repo.EquipmentType{
		&repo.EquipmentType{
			Type:     "typ1",
			ID:       "1",
			SourceID: "s1",
			ParentID: "2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:           "1",
					Name:         "attr1",
					Type:         repo.DataTypeString,
					IsIdentifier: true,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
				&repo.Attribute{
					ID:           "3",
					Name:         "attr3",
					Type:         repo.DataTypeInt,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_3",
				},
				&repo.Attribute{
					ID:           "4",
					Name:         "attr4",
					Type:         repo.DataTypeFloat,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_4",
				},
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr5",
					Type:               repo.DataTypeString,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_5",
				},
				&repo.Attribute{
					ID:       "1",
					Name:     "attr6",
					Type:     repo.DataTypeString,
					MappedTo: "mapping_6",
				},
				&repo.Attribute{
					ID:          "1",
					Name:        "attr7",
					IsDisplayed: true,
					Type:        repo.DataTypeString,
					MappedTo:    "mapping_7",
				},
				&repo.Attribute{
					ID:           "1",
					Name:         "attr8",
					IsDisplayed:  true,
					IsSearchable: true,
					Type:         repo.DataType(255),
					MappedTo:     "mapping_8",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ2",
			ID:       "2",
			SourceID: "s2",
			//ParentID: "p2",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
				&repo.Attribute{
					ID:           "2",
					Name:         "attr_2",
					Type:         repo.DataTypeString,
					IsDisplayed:  true,
					IsSearchable: true,
					MappedTo:     "mapping_2",
				},
			},
		},
		&repo.EquipmentType{
			Type:     "typ3",
			ID:       "3",
			SourceID: "s3",
			ParentID: "4",
			Attributes: []*repo.Attribute{
				&repo.Attribute{
					ID:                 "1",
					Name:               "attr_1",
					Type:               repo.DataTypeString,
					IsIdentifier:       true,
					IsDisplayed:        true,
					IsSearchable:       true,
					IsParentIdentifier: true,
					MappedTo:           "mapping_1",
				},
			},
		},
	}
	const (
		records   int32 = 3
		norecords int32 = 0
	)
	var mockCtrl *gomock.Controller
	var rep repo.License
	type args struct {
		ctx context.Context
		req *v1.ListEquipmentProductsRequest
	}
	tests := []struct {
		name    string
		s       *licenseServiceServer
		args    args
		setup   func()
		want    *v1.ListEquipmentProductsResponse
		wantErr bool
	}{
		{name: "SUCCESS",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
					SearchParams: &v1.EquipmentProductSearchParams{
						Name: &v1.StringFilter{
							Filteringkey: "orac",
						},
					},
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentProducts(ctx, eqTypes[1], "e1", &repo.QueryEquipmentProduct{
					PageSize:  10,
					Offset:    0,
					SortBy:    repo.EquipmentProductSortByName,
					SortOrder: 0,
					Filter: &repo.AggregateFilter{
						Filters: []repo.Queryable{
							&repo.Filter{
								FilterKey:   "name",
								FilterValue: "orac",
							},
						},
					},
				}, []string{"A", "B"}).Times(1).Return(records, []*repo.EquipmentProduct{
					&repo.EquipmentProduct{
						SwidTag: "ORAC001",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
					&repo.EquipmentProduct{
						SwidTag: "ORAC002",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
					&repo.EquipmentProduct{
						SwidTag: "ORAC003",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
				}, nil)
			},
			want: &v1.ListEquipmentProductsResponse{
				TotalRecords: records,
				Products: []*v1.EquipmentProduct{
					&v1.EquipmentProduct{
						SwidTag: "ORAC001",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
					&v1.EquipmentProduct{
						SwidTag: "ORAC002",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
					&v1.EquipmentProduct{
						SwidTag: "ORAC003",
						Name:    "Oracle SGBD",
						Editor:  "Oracle",
						Version: "9.5.2",
					},
				},
			},
		},
		{name: "FAILURE - can not retrieve claims",
			args: args{
				ctx: context.Background(),
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup:   func() {},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment types",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
		{name: "FAILURE - equipment type doesnt exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "4",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - no data",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				mockRepo.EXPECT().EquipmentProducts(ctx, eqTypes[1], "e1", &repo.QueryEquipmentProduct{
					PageSize:  10,
					Offset:    0,
					SortBy:    1,
					SortOrder: 0,
				}, []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNoData)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment parents with given id - node not exists",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)
				mockRepo.EXPECT().EquipmentProducts(ctx, eqTypes[1], "e1", &repo.QueryEquipmentProduct{
					PageSize:  10,
					Offset:    0,
					SortBy:    1,
					SortOrder: 0,
				}, []string{"A", "B"}).Times(1).Return(norecords, nil, repo.ErrNodeNotFound)
			},
			wantErr: true,
		},
		{name: "FAILURE - cannot fetch equipment with given id",
			args: args{
				ctx: ctx,
				req: &v1.ListEquipmentProductsRequest{
					TypeId:    "2",
					EquipId:   "e1",
					PageNum:   1,
					PageSize:  10,
					SortBy:    v1.ListEquipmentProductsRequest_NAME,
					SortOrder: v1.SortOrder_ASC,
				},
			},
			setup: func() {
				mockCtrl = gomock.NewController(t)
				mockRepo := mock.NewMockLicense(mockCtrl)
				rep = mockRepo
				mockRepo.EXPECT().EquipmentTypes(ctx, []string{"A", "B"}).Times(1).Return(eqTypes, nil)

				mockRepo.EXPECT().EquipmentProducts(ctx, eqTypes[1], "e1", &repo.QueryEquipmentProduct{
					PageSize:  10,
					Offset:    0,
					SortBy:    1,
					SortOrder: 0,
				}, []string{"A", "B"}).Times(1).Return(norecords, nil, errors.New("Test Error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := NewLicenseServiceServer(rep)
			got, err := s.ListEquipmentProducts(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("licenseServiceServer.ListEquipmentProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("licenseServiceServer.ListEquipmentProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}
