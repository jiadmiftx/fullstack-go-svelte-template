package gin_utils

import (
	"reflect"

	"sekolah-madrasah/pkg/paginate_utils"
)

type DataPaginateI interface {
	GetMessage() string
	GetData() interface{}
	SetData(interface{})
	GetPaginate() *paginate_utils.PaginateData
}

type MessageResponse struct {
	Message string `json:"message"`
}

type DataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DataWithPaginateResponse struct {
	DataResponse
	Paginate *paginate_utils.PaginateData `json:"paginate,omitempty"`
}

func MakeDataPaginateResponse(data DataWithPaginateResponse) DataWithPaginateResponse {
	if data.Data == nil || reflect.ValueOf(data.Data).IsZero() {
		data.Data = []string{}
		return data
	}
	return data
}

func MakePaginateResponse(data DataPaginateI) interface{} {
	if data.GetData() == nil || reflect.ValueOf(data.GetData()).IsZero() {
		data.SetData([]string{})
		return data
	}
	return data
}
