package txmanager

import (
	"net/http"

	goCommonError "github.com/madevara24/go-common/errors"
)

const (
	ERR_CODE_DB_TRANSACTION_WRAPPER = "TX_01"
	ERR_CODE_ROLLBACK_TRANSACTION   = "TX_02"
	ERR_CODE_COMMIT_TRANSACTION     = "TX_03"
	ERR_CODE_TX_NOT_FOUND           = "TX_04"
)

var (
	ErrDbTransactionWrapper = goCommonError.NewErr(http.StatusInternalServerError, ERR_CODE_DB_TRANSACTION_WRAPPER, "error during transaction")
	ErrRollbackTransaction  = goCommonError.NewErr(http.StatusInternalServerError, ERR_CODE_ROLLBACK_TRANSACTION, "error on rollback transaction")
	ErrCommitTransaction    = goCommonError.NewErr(http.StatusInternalServerError, ERR_CODE_COMMIT_TRANSACTION, "error on commit transaction")
	ErrTxNotFound           = goCommonError.NewErr(http.StatusInternalServerError, ERR_CODE_TX_NOT_FOUND, "error TX is not found in context")
)
