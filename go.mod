module github.com/gar-id/Meta

go 1.22.2

replace github.com/gar-id/Meta/apps/agent v0.0.0 => ./apps/agent/

require (
	github.com/gar-id/Meta/apps/agent v0.0.0
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect
