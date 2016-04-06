package smtpd

const (
	CodeSystemStatus                            = 211
	CodeHelpMessage                             = 214
	CodeServiceReady                            = 220
	CodeServiceClosing                          = 221
	CodeOK                                      = 250
	CodeUserNotLocal                            = 251
	CodeCannotVRFYUser                          = 252
	CodeStartMailInput                          = 354
	CodeServiceNotAvailable                     = 421
	CodeMailboxUnavailable                      = 450
	CodeLocalErrorInProcessing                  = 451
	CodeInsufficientSystemStorage               = 452
	CodeServerUnableToAccommodateParameters     = 455
	CodeSyntaxError                             = 500
	CodeSyntaxErrorInParametersOrArguments      = 501
	CodePermanentCommandNotImplemented          = 502
	CodePermanentBadSequenceOfCommands          = 503
	CodePermanentCommandParameterNotImplemented = 504
	CodePermanentMailboxUnavailable             = 550
	CodePermanentUserNotLocal                   = 551
	CodePermanentExceededStorageAllocation      = 552
	CodePermanentMailboxNameNotAllowed          = 553
	CodePermanentTransactionFailed              = 554
	CodePermanentMailRcptParameterError         = 555
)
