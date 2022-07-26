// Code generated by mockery v2.12.0. DO NOT EDIT.

package monitoring

import (
	http "net/http"
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// MetricsMock is an autogenerated mock type for the Metrics type
type MetricsMock struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields: networkName, networkID, chainID, oracleName, sender, feedName, feedPath, symbol, contractType, contractStatus, contractAddress, feedID
func (_m *MetricsMock) Cleanup(networkName string, networkID string, chainID string, oracleName string, sender string, feedName string, feedPath string, symbol string, contractType string, contractStatus string, contractAddress string, feedID string) {
	_m.Called(networkName, networkID, chainID, oracleName, sender, feedName, feedPath, symbol, contractType, contractStatus, contractAddress, feedID)
}

// HTTPHandler provides a mock function with given fields:
func (_m *MetricsMock) HTTPHandler() http.Handler {
	ret := _m.Called()

	var r0 http.Handler
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	return r0
}

// IncOffchainAggregatorAnswersTotal provides a mock function with given fields: contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) IncOffchainAggregatorAnswersTotal(contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetFeedContractLinkBalance provides a mock function with given fields: balance, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetFeedContractLinkBalance(balance float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(balance, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetFeedContractMetadata provides a mock function with given fields: chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName, symbol
func (_m *MetricsMock) SetFeedContractMetadata(chainID string, contractAddress string, feedID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string, symbol string) {
	_m.Called(chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName, symbol)
}

// SetFeedContractTransactionsFailed provides a mock function with given fields: numFailed, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetFeedContractTransactionsFailed(numFailed float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(numFailed, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetFeedContractTransactionsSucceeded provides a mock function with given fields: numSucceeded, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetFeedContractTransactionsSucceeded(numSucceeded float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(numSucceeded, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetHeadTrackerCurrentHead provides a mock function with given fields: blockNumber, networkName, chainID, networkID
func (_m *MetricsMock) SetHeadTrackerCurrentHead(blockNumber float64, networkName string, chainID string, networkID string) {
	_m.Called(blockNumber, networkName, chainID, networkID)
}

// SetLinkAvailableForPayment provides a mock function with given fields: amount, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetLinkAvailableForPayment(amount float64, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(amount, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetNodeMetadata provides a mock function with given fields: chainID, networkID, networkName, oracleName, sender
func (_m *MetricsMock) SetNodeMetadata(chainID string, networkID string, networkName string, oracleName string, sender string) {
	_m.Called(chainID, networkID, networkName, oracleName, sender)
}

// SetNodeTransactionsFailed provides a mock function with given fields: numFailed, oracleName, sender, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetNodeTransactionsFailed(numFailed float64, oracleName string, sender string, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(numFailed, oracleName, sender, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetNodeTransactionsSucceeded provides a mock function with given fields: numSucceeded, oracleName, sender, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetNodeTransactionsSucceeded(numSucceeded float64, oracleName string, sender string, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(numSucceeded, oracleName, sender, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorAnswerStalled provides a mock function with given fields: isSet, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorAnswerStalled(isSet bool, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(isSet, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorAnswers provides a mock function with given fields: answer, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorAnswers(answer float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(answer, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorAnswersRaw provides a mock function with given fields: answer, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorAnswersRaw(answer float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(answer, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorJuelsPerFeeCoin provides a mock function with given fields: juelsPerFeeCoin, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorJuelsPerFeeCoin(juelsPerFeeCoin float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(juelsPerFeeCoin, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorJuelsPerFeeCoinRaw provides a mock function with given fields: juelsPerFeeCoin, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorJuelsPerFeeCoinRaw(juelsPerFeeCoin float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(juelsPerFeeCoin, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorJuelsPerFeeCoinReceivedValues provides a mock function with given fields: value, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorJuelsPerFeeCoinReceivedValues(value float64, contractAddress string, feedID string, sender string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(value, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorRoundID provides a mock function with given fields: aggregatorRoundID, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorRoundID(aggregatorRoundID float64, contractAddress string, feedID string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(aggregatorRoundID, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// SetOffchainAggregatorSubmissionReceivedValues provides a mock function with given fields: value, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName
func (_m *MetricsMock) SetOffchainAggregatorSubmissionReceivedValues(value float64, contractAddress string, feedID string, sender string, chainID string, contractStatus string, contractType string, feedName string, feedPath string, networkID string, networkName string) {
	_m.Called(value, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName)
}

// NewMetricsMock creates a new instance of MetricsMock. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMetricsMock(t testing.TB) *MetricsMock {
	mock := &MetricsMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
