module webrtc-monitor

go 1.12

require (
	github.com/deepch/av v0.0.0-20160612005306-c437a98c9300
	github.com/djwackey/dorsvr v0.0.0-20180829130634-3a045ea76ec0
	github.com/djwackey/gitea v0.0.0-20170413062720-8b0b6b776461 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-chi/render v1.0.1
	github.com/go-xorm/core v0.6.3 // indirect
	github.com/go-xorm/xorm v0.7.5 // indirect
	github.com/pion/webrtc/v2 v2.0.27
	github.com/sirupsen/logrus v1.4.2
)

replace (
	github.com/go-xorm/builder v0.3.4 => xorm.io/builde v0.3.4
	github.com/go-xorm/core v0.6.3 => xorm.io/core v0.6.3
)
