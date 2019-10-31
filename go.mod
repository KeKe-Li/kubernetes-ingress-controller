module kubernetes-ingress-controller

go 1.13

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.41.0
	golang.org/x/blog => github.com/golang/blog v0.0.0-20190708141629-e28c63452d36
	golang.org/x/build => github.com/golang/build v0.0.0-20190709001953-30c0e6b89ea0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190627132806-fd42eb6b336f
	golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190711165009-e47acb2ca7f9
	golang.org/x/net => github.com/golang/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190620143337-7c3f2128ad9b
	golang.org/x/review => github.com/golang/review v0.0.0-20190508204355-8102926ea734
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190710143415-6ec70d6a5542
	golang.org/x/talks => github.com/golang/talks v0.0.0-20190313194420-5ca518b26a55
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190711191110-9a621aea19f8
	golang.org/x/tour => github.com/golang/tour v0.0.0-20190611164551-1f1f3d2b450b
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.7.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.1
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190708153700-3bdd9d9f5532
	google.golang.org/grpc => github.com/grpc/grpc-go v1.22.0
	gopkg.in/jcmturner/gokrb5.v7 => github.com/jcmturner/gokrb5 v7.3.0+incompatible
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/chanxuehong/internal v0.0.0-20180430074813-85d6017afbc4 // indirect
	github.com/chanxuehong/log v0.0.0-20190928070219-69fffef55d8f
	github.com/chanxuehong/mns.aliyun.v20150606 v0.0.0-20180613022928-c2c6bc6a2603
	github.com/chanxuehong/rand v0.0.0-20180830053958-4b3aff17f488 // indirect
	github.com/chanxuehong/uuid v0.0.0-20180430073920-75ab5e2d8298 // indirect
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/prometheus/client_golang v1.1.0
	github.com/spf13/cast v1.3.0
	github.com/spf13/viper v1.4.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/apimachinery v0.0.0-20191030190112-bb31b70367b7 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20191030222137-2b95a09bc58d // indirect
)
