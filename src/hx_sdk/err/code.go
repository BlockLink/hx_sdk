package err

const (
	ErrAmount     = "40000"
	ErrDumpPubKey = "40001"

	ErrCodeInternalError     = "50000"
	ErrCodeNEOAmount         = "50001" // NEO数量必须为整数
	ErrCodeNEOAssetInvalid   = "50002" // 资产类型必须为NEO或NEOGAS
	ErrCodeNEOPrivKeyLen     = "50003" // NEO私钥长度错误(必须为32位)
	ErrCodeZeroAssetSend     = "50004" // 发送资产数量为0
	ErrCodeNEOInvalidWif     = "50005" // WIF或WIFC非法
	ErrCodeInvalidAPI        = "50006" // API地址非法
	ErrCodeMarshalReq        = "50007" // 序列化HTTP请求失败        注：有ErrMsg，系统函数调用错误
	ErrCodeCreateReq         = "50008" // 创建请求失败              注：有ErrMsg，系统函数调用错误
	ErrCodeDoRequest         = "50009" // 请求失败                  注：有ErrMsg，系统函数调用错误
	ErrCodeRespStatus        = "50010" // 请求返回码不是200          注： ErrMsg为HTTP response status
	ErrCodeReadRespBody      = "50011" // 读取response失败          注：有ErrMsg，系统函数ReadAll调用错误
	ErrCodeUnmarshalResp     = "50012" // 序列号response失败        注：有ErrMsg，系统函数Unmarshal调用错误
	ErrCodeUnmarshalRespData = "50013" // 序列号response data失败   注：有ErrMsg，为系统函数Unmarshal返回错误原因
	ErrCodeBuildNEOTxInput   = "50014" // 构建Tx Input失败          注: ErrMsg为详细错误原因
	ErrCodeBuildNEOTxAttrs   = "50015" // 构建Tx Attribute失败      注：ErrMsg为详细错误原因
	ErrCodeBuildNEOTxOutput  = "50016" // 构建Tx Output失败         注：ErrMsg为详细错误原因
	ErrCodeNEOSignTx         = "50017" // 构建NEO签名失败            注：ErrMsg为详细错误原因

	ErrCodeUnsupportCoin    = "60000" // unsupport coin type
	ErrAddressInvalid       = "60001"
	ErrAddressNotEqualBIP44 = "60002" // param from address should equal with bip44 m/44'/cointype'/account/index
	ErrBuildHCTx            = "60003" // build HC Tx failed
	ErrBroadcastHCTx        = "60004"
	ErrEstimateHCFee        = "60005"
	ErrGetHCAddress         = "60006"
	ErrDumpHCKey            = "60007"

	ErrBuildDcrTx     = "60053" // build Dcr Tx failed
	ErrBroadcastDcrTx = "60054"
	ErrEstimateDcrFee = "60055"
	ErrGetDcrAddress  = "60056"
	ErrDumpDcrKey     = "60057"

	ErrCodeCreateAdaWallet = "60100"
	ErrCodeGetAdaAddress   = "60101"
	ErrCodeAdaIndex        = "60102"
	ErrCodeAdaAddress      = "60103"
	ErrCodeBuildAdaTX      = "60104"
	ErrCodeBroadcastAda    = "60105"

	ErrGetPubkey      = "60200"
	ErrGetBTSAddress  = "60251"
	ErrDumpBTSKey     = "60252"
	ErrBTSTransaction = "60253"
	ErrBTSExportWif   = "60254"
	ErrBTSRefChainID  = "60255"
	ErrBTSResponse    = "60256"
	ErrBTSTransfer    = "60257"

	ErrInvalidAmount  = "70000"
	ErrPurchaseTicket = "70001"
	ErrCreateClient   = "70002"
	ErrMarshalJSON    = "70003"
	ErrUnmarshalJSON  = "70004"

	ErrGetHXAddress    = "60201"
	ErrDumpHXKey       = "60202"
	ErrHXTransaction   = "60203"
	ErrHXExportWif     = "60204"
	ErrHXRefChainID    = "60205"
	ErrHXInvalidAction = "60206"
	ErrHXUnknownAction = "60207"
	ErrGetHxRef        = "60208"
	ErrHxSignature     = "60209"

	ErrGetAEAddress    = "60301"
	ErrValidateAddress = "60302"
	ErrBuildAETX       = "60303"
	ErrTransferAE      = "60304"

	ErrGetTRXAddress      = "60401"
	ErrValidateTRXAddress = "60402"
	ErrSignTRX            = "60403"
	ErrDumpTRXPrivateKey  = "60404"

	ErrGetBTMAddress      = "60501"
	ErrValidateBTMAddress = "60502"
	ErrBuildBTMTX         = "60503"
	ErrDumpBTMPrivateKey  = "60504"

	ErrGetXWCAddress    = "60601"
	ErrDumpXWCKey       = "60602"
	ErrXWCTransaction   = "60603"
	ErrXWCExportWif     = "60604"
	ErrXWCRefChainID    = "60605"
	ErrXWCInvalidAction = "60606"
	ErrXWCUnknownAction = "60607"
	ErrGetXwcRef        = "60608"
	ErrXwcSignature     = "60609"
)
