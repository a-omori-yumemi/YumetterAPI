package data_source_wrapper

type DataSourceWrapper interface {
	Get() interface{}
}

type DataSourceMaker interface {
	NewDataSourceWrapper(func() (interface{}, error)) DataSourceWrapper
}
