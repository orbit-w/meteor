package mailbox

type (
	SuspendMailbox struct{}
	ResumeMailbox  struct{}
)

var (
	gSuspendMailbox = new(SuspendMailbox)
	gResumeMailbox  = new(ResumeMailbox)
)
